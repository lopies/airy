// Copyright (c) TFG Co and AIRY. All Rights Reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package component

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/interfaces"
	"github.com/airy/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
	"strings"
	"sync"
	"time"
)

var deleteCallbacks = make([]func(sv *config.Server), 0)

type EtcdServiceDiscovery struct {
	cli                    *clientv3.Client
	syncServersInterval    time.Duration
	heartbeatTTL           time.Duration
	logHeartbeat           bool
	lastHeartbeatTime      time.Time
	leaseID                clientv3.LeaseID
	mapByTypeLock          sync.RWMutex
	serverMapByType        map[string]map[string]*config.Server
	serverMapByID          sync.Map
	etcdEndpoints          []string
	etcdUser               string
	etcdPass               string
	etcdPrefix             string
	etcdDialTimeout        time.Duration
	running                bool
	server                 *config.Server
	stopChan               chan bool
	stopLeaseChan          chan bool
	lastSyncTime           time.Time
	listeners              []interfaces.SDListener
	revokeTimeout          time.Duration
	grantLeaseTimeout      time.Duration
	grantLeaseMaxRetries   int
	grantLeaseInterval     time.Duration
	shutdownDelay          time.Duration //Sleep for a short while to ensure shutdown has propagated
	appDieChan             chan struct{}
	serverTypesBlacklist   []string
	syncServersParallelism int
	syncServersRunning     chan bool
	BaseComponent
}

func NewETCD() *EtcdServiceDiscovery {
	e := new(EtcdServiceDiscovery)
	e.SetType(constants.ETCDComponent)
	return e
}

// Init starts the service discovery client
func (sd *EtcdServiceDiscovery) Init(conf *config.AiryConfig) {
	if len(sd.listeners) <= 0 {
		logger.Fatalf("please add grpc component in etcd component before first")
	}
	sd.running = false
	sd.server = conf.Server
	sd.serverMapByType = make(map[string]map[string]*config.Server)
	//sd.listeners = make([]interfaces.SDListener, 0) not allowed
	sd.stopChan = make(chan bool)
	sd.stopLeaseChan = make(chan bool)
	sd.appDieChan = conf.Stop.DieChan
	sd.syncServersRunning = make(chan bool)

	sd.configure(conf)
}

func (sd *EtcdServiceDiscovery) AfterInit() {
	sd.running = true
	var err error

	if sd.cli == nil {
		err = sd.InitETCDClient()
		if err != nil {
			logger.Errorf("init etcd client error,err=%s", err.Error())
			panic(fmt.Sprintf("init etcd client error,err=%s", err.Error()))
		}
	} else {
		sd.cli.KV = namespace.NewKV(sd.cli.KV, sd.etcdPrefix)
		sd.cli.Watcher = namespace.NewWatcher(sd.cli.Watcher, sd.etcdPrefix)
		sd.cli.Lease = namespace.NewLease(sd.cli.Lease, sd.etcdPrefix)
	}
	go sd.watchEtcdChanges()

	if err = sd.bootstrap(); err != nil {
		logger.Errorf("etcd bootstrap error,err=%s", err.Error())
		panic(fmt.Sprintf("etcd bootstrap error,err=%s", err.Error()))
	}

	// update modules
	syncServersTicker := time.NewTicker(sd.syncServersInterval)
	go func() {
		for sd.running {
			select {
			case <-syncServersTicker.C:
				err := sd.SyncServers(false)
				if err != nil {
					logger.Errorf("error resyncing modules: %s", err.Error())
				}
			case <-sd.stopChan:
				return
			}
		}
	}()
}

func (sd *EtcdServiceDiscovery) configure(config *config.AiryConfig) {
	c := config.EtcdServiceDiscovery
	sd.etcdEndpoints = c.Endpoints
	sd.etcdUser = c.User
	sd.etcdPass = c.Pass
	sd.etcdDialTimeout = c.DialTimeout
	sd.etcdPrefix = c.Prefix
	sd.heartbeatTTL = c.HeartbeatTTL
	sd.logHeartbeat = c.HeartbeatLog
	sd.syncServersInterval = c.SyncServersInterval
	sd.revokeTimeout = c.RevokeTimeout
	sd.grantLeaseTimeout = c.GrantLeaseTimeout
	sd.grantLeaseMaxRetries = c.GrantLeaseMaxRetries
	sd.grantLeaseInterval = c.GrantLeaseRetryInterval
	sd.shutdownDelay = c.ShutdownDelay
	sd.serverTypesBlacklist = c.ServerTypesBlacklist
	sd.syncServersParallelism = c.SyncServersParallelism
}

// BeforeShutdown executes before shutting down and will remove the server from the list
func (sd *EtcdServiceDiscovery) BeforeShutdown() {
	sd.revoke()
	time.Sleep(sd.shutdownDelay) // Sleep for a short while to ensure shutdown has propagated
}

// Shutdown executes on shutdown and will clean etcd
func (sd *EtcdServiceDiscovery) Shutdown() error {
	sd.running = false
	close(sd.stopChan)
	return nil
}

func (sd *EtcdServiceDiscovery) OnDeleteServerBind(f func(sv *config.Server)) {
	deleteCallbacks = append(deleteCallbacks, f)
}

func (sd *EtcdServiceDiscovery) watchLeaseChan(c <-chan *clientv3.LeaseKeepAliveResponse) {
	failedGrantLeaseAttempts := 0
	for {
		select {
		case <-sd.stopChan:
			return
		case <-sd.stopLeaseChan:
			return
		case leaseKeepAliveResponse, ok := <-c:
			if !ok {
				logger.Error("ETCD lease KeepAlive died, retrying in 10 seconds")
				time.Sleep(10000 * time.Millisecond)
			}
			if leaseKeepAliveResponse != nil {
				if sd.logHeartbeat {
					logger.Debugf("sd: etcd lease %x renewed", leaseKeepAliveResponse.ID)
				}
				failedGrantLeaseAttempts = 0
				continue
			}
			logger.Warnf("sd: error renewing etcd lease, reconfiguring")
			for {
				err := sd.renewLease()
				if err != nil {
					failedGrantLeaseAttempts = failedGrantLeaseAttempts + 1
					if err == constants.ErrEtcdGrantLeaseTimeout {
						logger.Warnf("sd: timed out trying to grant etcd lease")
						if sd.appDieChan != nil {
							sd.appDieChan <- struct{}{}
						}
						return
					}
					if failedGrantLeaseAttempts >= sd.grantLeaseMaxRetries {
						logger.Warnf("sd: exceeded max attempts to renew etcd lease")
						if sd.appDieChan != nil {
							sd.appDieChan <- struct{}{}
						}
						return
					}
					logger.Warnf("sd: error granting etcd lease, will retry in %d seconds", uint64(sd.grantLeaseInterval.Seconds()))
					time.Sleep(sd.grantLeaseInterval)
					continue
				}
				return
			}
		}
	}
}

// renewLease reestablishes connection with etcd
func (sd *EtcdServiceDiscovery) renewLease() error {
	c := make(chan error)
	go func() {
		defer close(c)
		logger.Infof("waiting for etcd lease")
		err := sd.grantLease()
		if err != nil {
			c <- err
			return
		}
		err = sd.bootstrapServer(sd.server)
		c <- err
	}()
	select {
	case err := <-c:
		return err
	case <-time.After(sd.grantLeaseTimeout):
		return constants.ErrEtcdGrantLeaseTimeout
	}
}

func (sd *EtcdServiceDiscovery) grantLease() error {
	// grab lease
	l, err := sd.cli.Grant(context.TODO(), int64(sd.heartbeatTTL.Seconds()))
	if err != nil {
		return err
	}
	sd.leaseID = l.ID
	logger.Debugf("sd: got leaseID: %x", l.ID)
	// this will keep alive forever, when channel c is closed
	// it means we probably have to rebootstrap the lease
	c, err := sd.cli.KeepAlive(context.TODO(), sd.leaseID)
	if err != nil {
		return err
	}
	// need to receive here as per etcd docs
	<-c
	go sd.watchLeaseChan(c)
	return nil
}

func (sd *EtcdServiceDiscovery) addServerIntoEtcd(server *config.Server) error {
	_, err := sd.cli.Put(
		context.TODO(),
		getKey(server.ID, string(server.Type)),
		server.String(),
		clientv3.WithLease(sd.leaseID),
	)
	return err
}

func (sd *EtcdServiceDiscovery) bootstrapServer(server *config.Server) error {
	if err := sd.addServerIntoEtcd(server); err != nil {
		return err
	}

	sd.SyncServers(true)
	return nil
}

// AddListener adds a listener to etcd service discovery
func (sd *EtcdServiceDiscovery) AddListener(listener interfaces.SDListener) {
	sd.listeners = append(sd.listeners, listener)
	return
}

func (sd *EtcdServiceDiscovery) notifyListeners(act interfaces.Action, sv *config.Server) {
	for _, l := range sd.listeners {
		if act == interfaces.DEL {
			l.RemoveServer(sv)
		} else if act == interfaces.ADD {
			l.AddServer(sv)
		}
	}
}

func (sd *EtcdServiceDiscovery) writeLockScope(f func()) {
	sd.mapByTypeLock.Lock()
	defer sd.mapByTypeLock.Unlock()
	f()
}

func (sd *EtcdServiceDiscovery) readLockScope(f func()) {
	sd.mapByTypeLock.RLock()
	defer sd.mapByTypeLock.RUnlock()
	f()
}

func (sd *EtcdServiceDiscovery) deleteServer(serverID string) {
	if actual, ok := sd.serverMapByID.Load(serverID); ok {
		sv := actual.(*config.Server)
		sd.serverMapByID.Delete(sv.ID)
		sd.writeLockScope(func() {
			if svMap, ok := sd.serverMapByType[sv.Type]; ok {
				delete(svMap, sv.ID)
				for _, callback := range deleteCallbacks {
					callback(sv)
				}
			}
		})
		sd.notifyListeners(interfaces.DEL, sv)
	}
}

func (sd *EtcdServiceDiscovery) deleteLocalInvalidServers(actualServers []string) {
	sd.serverMapByID.Range(func(key interface{}, value interface{}) bool {
		k := key.(string)
		if !sliceContainsString(actualServers, k) {
			logger.Warnf("deleting invalid local server %s", k)
			sd.deleteServer(k)
		}
		return true
	})
}

func getKey(serverID, serverType string) string {
	return fmt.Sprintf("modules/%s/%s", serverType, serverID)
}

func getServerFromEtcd(cli *clientv3.Client, serverType, serverID string) (*config.Server, error) {
	svKey := getKey(serverID, serverType)
	svEInfo, err := cli.Get(context.TODO(), svKey)
	if err != nil {
		return nil, fmt.Errorf("error getting server: %s from etcd, error: %s", svKey, err.Error())
	}
	if len(svEInfo.Kvs) == 0 {
		return nil, fmt.Errorf("didn't found server: %s in etcd", svKey)
	}
	return parseServer(svEInfo.Kvs[0].Value)
}

// GetServersByType returns a slice with all the modules of a certain type
func (sd *EtcdServiceDiscovery) GetServersByType(typ string) (map[string]*config.Server, error) {
	sd.mapByTypeLock.RLock()
	defer sd.mapByTypeLock.RUnlock()
	if m, ok := sd.serverMapByType[typ]; ok && len(m) > 0 {
		// Create a new map to avoid concurrent read and write access to the
		// map, this also prevents accidental changes to the list of modules
		// kept by the service discovery.
		ret := make(map[string]*config.Server, len(sd.serverMapByType[typ]))
		for k, v := range sd.serverMapByType[typ] {
			ret[k] = v
		}
		return ret, nil
	}
	return nil, constants.ErrNoServersAvailableOfType
}

// GetServers returns a slice with all the modules
func (sd *EtcdServiceDiscovery) GetServers() []*config.Server {
	ret := make([]*config.Server, 0)
	sd.serverMapByID.Range(func(k, v interface{}) bool {
		ret = append(ret, v.(*config.Server))
		return true
	})
	return ret
}

// bootstrap run and sync etcd server
func (sd *EtcdServiceDiscovery) bootstrap() error {
	if err := sd.grantLease(); err != nil {
		return err
	}
	if err := sd.bootstrapServer(sd.server); err != nil {
		return err
	}

	return nil
}

// GetServer returns a server given it's id
func (sd *EtcdServiceDiscovery) GetServer(id string) (*config.Server, error) {
	if sv, ok := sd.serverMapByID.Load(id); ok {
		return sv.(*config.Server), nil
	}
	return nil, constants.ErrNoServerWithID
}

// InitETCDClient initializes etcd client
func (sd *EtcdServiceDiscovery) InitETCDClient() error {
	logger.Debugf("Initializing ETCD client")
	var cli *clientv3.Client
	var err error
	c := clientv3.Config{
		Endpoints:   sd.etcdEndpoints,
		DialTimeout: sd.etcdDialTimeout,
	}
	if sd.etcdUser != "" && sd.etcdPass != "" {
		c.Username = sd.etcdUser
		c.Password = sd.etcdPass
	}
	cli, err = clientv3.New(c)
	if err != nil {
		logger.Errorf("error initializing etcd client: %s", err.Error())
		return err
	}
	sd.cli = cli

	// namespaced etcd :)
	sd.cli.KV = namespace.NewKV(sd.cli.KV, sd.etcdPrefix)
	sd.cli.Watcher = namespace.NewWatcher(sd.cli.Watcher, sd.etcdPrefix)
	sd.cli.Lease = namespace.NewLease(sd.cli.Lease, sd.etcdPrefix)
	return nil
}

// parseEtcdKey
func parseEtcdKey(key string) (string, string, error) {
	splittedServer := strings.Split(key, "/")
	if len(splittedServer) != 3 {
		return "", "", fmt.Errorf("error parsing etcd key %s (server name can't contain /)", key)
	}
	svType := splittedServer[1]
	svID := splittedServer[2]
	return svType, svID, nil
}

// parseServer return server
func parseServer(value []byte) (*config.Server, error) {
	var sv *config.Server
	err := json.Unmarshal(value, &sv)
	if err != nil {
		logger.Warnf("failed to load server %s, error: %s", sv, err.Error())
		return nil, err
	}
	return sv, nil
}

func (sd *EtcdServiceDiscovery) printServers() {
	sd.mapByTypeLock.RLock()
	defer sd.mapByTypeLock.RUnlock()
	for k, v := range sd.serverMapByType {
		logger.Debugf("type: %s, modules: %+v", k, v)
	}
}

// Struct that encapsulates a parallel/concurrent etcd get
// it spawns goroutines and receives work requests through a channel
type parallelGetterWork struct {
	serverType string
	serverID   string
	payload    []byte
}

type parallelGetter struct {
	cli         *clientv3.Client
	numWorkers  int
	wg          *sync.WaitGroup
	resultMutex sync.Mutex
	result      *[]*config.Server
	workChan    chan parallelGetterWork
}

// newParallelGetter init parallel and start
func newParallelGetter(cli *clientv3.Client, numWorkers int) parallelGetter {
	if numWorkers <= 0 {
		numWorkers = 10
	}
	p := parallelGetter{
		cli:        cli,
		numWorkers: numWorkers,
		workChan:   make(chan parallelGetterWork),
		wg:         new(sync.WaitGroup),
		result:     new([]*config.Server),
	}
	p.start()
	return p
}

// start exec work chan
func (p *parallelGetter) start() {
	for i := 0; i < p.numWorkers; i++ {
		go func() {
			for work := range p.workChan {
				logger.Debugf("loading info from missing server: %s/%s", work.serverType, work.serverID)
				var sv *config.Server
				var err error
				if work.payload == nil {
					sv, err = getServerFromEtcd(p.cli, work.serverType, work.serverID)
				} else {
					sv, err = parseServer(work.payload)
				}
				if err != nil {
					logger.Errorf("Error parsing server from etcd: %s, error: %s", work.serverID, err.Error())
					p.wg.Done()
					continue
				}

				p.resultMutex.Lock()
				*p.result = append(*p.result, sv)
				p.resultMutex.Unlock()

				p.wg.Done()
			}
		}()
	}
}

// waitAndGetResult
func (p *parallelGetter) waitAndGetResult() []*config.Server {
	p.wg.Wait()
	close(p.workChan)
	return *p.result
}

// addWorkWithPayload add a job to work chan with payload
func (p *parallelGetter) addWorkWithPayload(serverType, serverID string, payload []byte) {
	p.wg.Add(1)
	p.workChan <- parallelGetterWork{
		serverType: serverType,
		serverID:   serverID,
		payload:    payload,
	}
}

// addWork add a job to work chan
func (p *parallelGetter) addWork(serverType, serverID string) {
	p.wg.Add(1)
	p.workChan <- parallelGetterWork{
		serverType: serverType,
		serverID:   serverID,
	}
}

// SyncServers gets all modules from etcd
func (sd *EtcdServiceDiscovery) SyncServers(firstSync bool) error {
	sd.syncServersRunning <- true
	defer func() {
		sd.syncServersRunning <- false
	}()
	//start := time.Now()
	var kvs *clientv3.GetResponse
	var err error
	if firstSync {
		kvs, err = sd.cli.Get(
			context.TODO(),
			"modules/",
			clientv3.WithPrefix(),
		)
	} else {
		kvs, err = sd.cli.Get(
			context.TODO(),
			"modules/",
			clientv3.WithPrefix(),
			clientv3.WithKeysOnly(),
		)
	}
	if err != nil {
		logger.Errorf("Error querying etcd server: %s", err.Error())
		return err
	}

	// delete invalid modules (local ones that are not in etcd)
	var allIds = make([]string, 0)

	// Spawn worker goroutines that will work in parallel
	parallelGetter := newParallelGetter(sd.cli, sd.syncServersParallelism)

	for _, kv := range kvs.Kvs {
		svType, svID, err := parseEtcdKey(string(kv.Key))
		if err != nil {
			logger.Warnf("failed to parse etcd key %s, error: %s", kv.Key, err.Error())
			continue
		}

		// Check whether the server type is blacklisted or not
		if sd.isServerTypeBlacklisted(svType) && svID != sd.server.ID {
			logger.Debugf("ignoring blacklisted server type '%s'", svType)
			continue
		}

		allIds = append(allIds, svID)

		if _, ok := sd.serverMapByID.Load(svID); !ok {
			// Add new work to the channel
			if firstSync {
				parallelGetter.addWorkWithPayload(svType, svID, kv.Value)
			} else {
				parallelGetter.addWork(svType, svID)
			}
		}
	}

	// Wait until all goroutines are finished
	servers := parallelGetter.waitAndGetResult()

	for _, server := range servers {
		logger.Debugf("adding server %s", server)
		sd.addServer(server)
	}

	sd.deleteLocalInvalidServers(allIds)

	//sd.printServers()
	sd.lastSyncTime = time.Now()
	//elapsed := time.Since(start)
	//logger.Debugf("SyncServers took : %s to run", elapsed)
	return nil
}

// revoke prevents Airy from crashing when etcd is not available
func (sd *EtcdServiceDiscovery) revoke() error {
	close(sd.stopLeaseChan)
	c := make(chan error)
	defer close(c)
	go func() {
		logger.Debugf("waiting for etcd revoke")
		_, err := sd.cli.Revoke(context.TODO(), sd.leaseID)
		c <- err
		logger.Debugf("finished waiting for etcd revoke")
	}()
	select {
	case err := <-c:
		return err // completed normally
	case <-time.After(sd.revokeTimeout):
		logger.Warnf("timed out waiting for etcd revoke")
		return nil // timed out
	}
}

// addServer add server to memory if a new server connect etcd
func (sd *EtcdServiceDiscovery) addServer(sv *config.Server) {
	if _, loaded := sd.serverMapByID.LoadOrStore(sv.ID, sv); !loaded {
		sd.writeLockScope(func() {
			mapSvByType, ok := sd.serverMapByType[sv.Type]
			if !ok {
				mapSvByType = make(map[string]*config.Server)
				sd.serverMapByType[sv.Type] = mapSvByType
			}
			mapSvByType[sv.ID] = sv
		})
		if sv.ID != sd.server.ID {
			sd.notifyListeners(interfaces.ADD, sv)
		}
	}
}

// watchEtcdChanges listen etcd prefix key change
func (sd *EtcdServiceDiscovery) watchEtcdChanges() {
	w := sd.cli.Watch(context.Background(), "modules/", clientv3.WithPrefix())
	failedWatchAttempts := 0
	go func(chn clientv3.WatchChan) {
		for sd.running {
			select {
			// Block here if SyncServers() is running and consume the watcher channel after it's finished, to avoid conflicts
			case syncServersState := <-sd.syncServersRunning:
				for syncServersState {
					syncServersState = <-sd.syncServersRunning
				}
			case wResp, ok := <-chn:
				if wResp.Err() != nil {
					logger.Warnf("etcd watcher response error: %s", wResp.Err())
					time.Sleep(100 * time.Millisecond)
				}
				if !ok {
					logger.Error("etcd watcher died, retrying to watch in 1 second")
					failedWatchAttempts++
					time.Sleep(1000 * time.Millisecond)
					if failedWatchAttempts > 10 {
						if err := sd.InitETCDClient(); err != nil {
							failedWatchAttempts = 0
							continue
						}
						chn = sd.cli.Watch(context.Background(), "modules/", clientv3.WithPrefix())
						failedWatchAttempts = 0
					}
					continue
				}
				failedWatchAttempts = 0
				for _, ev := range wResp.Events {
					svType, svID, err := parseEtcdKey(string(ev.Kv.Key))
					if err != nil {
						logger.Warnf("failed to parse key from etcd: %s", ev.Kv.Key)
						continue
					}

					if sd.isServerTypeBlacklisted(svType) && sd.server.ID != svID {
						continue
					}

					switch ev.Type {
					case clientv3.EventTypePut:
						var sv *config.Server
						var err error
						if sv, err = parseServer(ev.Kv.Value); err != nil {
							logger.Errorf("Failed to parse server from etcd: %v", err)
							continue
						}

						sd.addServer(sv)
						logger.Debugf("server %s added by watcher", ev.Kv.Key)
						sd.printServers()
					case clientv3.EventTypeDelete:
						sd.deleteServer(svID)
						logger.Debugf("server %s deleted by watcher", svID)
						sd.printServers()
					}
				}
			case <-sd.stopChan:
				return
			}

		}
	}(w)
}

// isServerTypeBlacklisted return true if svType in black list
func (sd *EtcdServiceDiscovery) isServerTypeBlacklisted(svType string) bool {
	for _, blacklistedSv := range sd.serverTypesBlacklist {
		if blacklistedSv == svType {
			return true
		}
	}
	return false
}

// sliceContainsString returns true if a slice contains the string
func sliceContainsString(slice []string, str string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}
	return false
}
