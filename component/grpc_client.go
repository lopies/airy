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
	"fmt"
	"github.com/airy/config"
	"github.com/airy/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"sync"
	"time"

	"github.com/airy/constants"
	"github.com/airy/logger"
	"google.golang.org/grpc"
)

// GRPCClient rpc client struct
type GRPCClient struct {
	gateServerCount  int
	logicServerCount int
	clientMap        sync.Map
	dialTimeout      time.Duration
	lazy             bool
	reqTimeout       time.Duration
	server           *config.Server
	BaseComponent
}

type grpcClient struct {
	gateCli  pb.AiryGateClient
	logicCli pb.AiryLogicClient
	address    string
	serverType string
	conn       *grpc.ClientConn
	connected  bool
	lock       sync.Mutex
}

func NewGRPCClient() *GRPCClient {
	g := new(GRPCClient)
	g.SetType(constants.GRPCClientComponent)
	return g
}

// Init inits grpc rpc client
func (gs *GRPCClient) Init(config *config.AiryConfig) {
	gs.server = config.Server
	gs.dialTimeout = config.GRPCClient.DialTimeout
	gs.lazy = config.GRPCClient.LazyConnection
	gs.reqTimeout = config.GRPCClient.RequestTimeout
}

func (gs *GRPCClient) GateServerCount() int {
	return gs.gateServerCount
}

func (gs *GRPCClient) LogicServerCount() int {
	return gs.logicServerCount
}

func (gs *GRPCClient) PushToUsers(ctx context.Context, serverID string, req *pb.Combinations) {
	logger.Debugf("PushToUsers ,svID = %s", serverID)
	c, ok := gs.clientMap.Load(serverID)
	if !ok {
		logger.Errorf("PushToUsers not found serverID : %s", serverID)
		return
	}
	ctxT, done := context.WithTimeout(ctx, gs.reqTimeout)
	defer done()
	_, err := c.(*grpcClient).pushToUsers(ctxT, req)
	if err != nil {
		logger.Errorf("grpc call pushToUsers error,err = %s", err.Error())
		return
	}
	return
}

//Request gate request
func (gs *GRPCClient) Request(ctx context.Context, serverID string, req *pb.Packet) error {
	c, ok := gs.clientMap.Load(serverID)
	if !ok {
		logger.Errorf("grpc request ,not found server id : %s", serverID)
		return constants.ErrNotFoundServer
	}
	ctxT, done := context.WithTimeout(ctx, gs.reqTimeout)
	defer done()
	_, err := c.(*grpcClient).request(ctxT, req)
	if err != nil {
		logger.Errorf("grpc call request interface error,code = %d,err = %s", req.RequestCode, err.Error())
		return err
	}
	return nil
}

// AddServer is called when a new server is discovered
func (gs *GRPCClient) AddServer(sv *config.Server) {
	var host, port, portKey string
	var ok bool

	host, portKey = gs.getServerHost(sv)
	if host == "" {
		logger.Errorf("server %s has no grpcHost specified in metadata", sv.ID)
		return
	}

	if port, ok = sv.Metadata[portKey]; !ok {
		logger.Errorf("server %s has no %s specified in metadata", sv.ID, portKey)
		return
	}

	address := fmt.Sprintf("%s:%s", host, port)
	client := &grpcClient{address: address, serverType: sv.Type}
	if !gs.lazy {
		if err := client.connect(); err != nil {
			logger.Errorf("unable to connect to server %s at %s: %v", sv.ID, address, err)
		}
	}
	switch sv.Type {
	case string(constants.GateModule):
		gs.gateServerCount++
	case string(constants.LogicModule):
		gs.logicServerCount++
	}
	gs.clientMap.Store(sv.ID, client)
	logger.Debugf("added server %s at %s,type = %s", sv.ID, address, sv.Type)
}

// RemoveServer is called when a server is removed
func (gs *GRPCClient) RemoveServer(sv *config.Server) {
	if c, ok := gs.clientMap.Load(sv.ID); ok {
		c.(*grpcClient).disconnect()
		gs.clientMap.Delete(sv.ID)

		switch sv.Type {
		case string(constants.GateModule):
			gs.gateServerCount--
		case string(constants.LogicModule):
			gs.logicServerCount--
		}
		logger.Debugf("removed server %s,type = %s", sv.ID, sv.Type)
	}
}

func (gs *GRPCClient) getServerHost(sv *config.Server) (host, portKey string) {
	var (
		serverRegion, hasRegion   = sv.Metadata[constants.RegionKey]
		externalHost, hasExternal = sv.Metadata[constants.GRPCExternalHostKey]
		internalHost, _           = sv.Metadata[constants.GRPCHostKey]
	)

	hasRegion = hasRegion && serverRegion != ""
	hasExternal = hasExternal && externalHost != ""

	if !hasRegion {
		if hasExternal {
			logger.Warnf("server %s has no region specified in metadata, using external host", sv.ID)
			return externalHost, constants.GRPCExternalPortKey
		}

		logger.Debugf("server %s has no region nor external host specified in metadata, using internal host", sv.ID)
		return internalHost, constants.GRPCPortKey
	}

	logger.Infof("server %s is in other region, using external host", sv.ID)
	return externalHost, constants.GRPCExternalPortKey
}

func (gc *grpcClient) connect() error {
	gc.lock.Lock()
	defer gc.lock.Unlock()
	if gc.connected {
		return nil
	}

	conn, err := grpc.Dial(
		gc.address,
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(2 << 26)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(2 << 26)),
		grpc.WithInsecure(),
	)
	if err != nil {
		return err
	}
	switch gc.serverType {
	case string(constants.GateModule):
		gc.gateCli = pb.NewAiryGateClient(conn)
	case string(constants.LogicModule):
		gc.logicCli = pb.NewAiryLogicClient(conn)
	default:
		logger.Errorf("listen grpc server type error,type = %s", gc.serverType)
	}
	gc.conn = conn
	gc.connected = true
	return nil
}

func (gc *grpcClient) disconnect() {
	gc.lock.Lock()
	if gc.connected {
		gc.conn.Close()
		gc.connected = false
	}
	gc.lock.Unlock()
}

func (gc *grpcClient) request(ctx context.Context, req *pb.Packet) (*empty.Empty, error) {
	if !gc.connected {
		if err := gc.connect(); err != nil {
			return nil, err
		}
	}
	return gc.logicCli.Request(ctx, req)
}

func (gc *grpcClient) pushToUsers(ctx context.Context, req *pb.Combinations) (*empty.Empty, error) {
	if !gc.connected {
		if err := gc.connect(); err != nil {
			return nil, err
		}
	}
	return gc.gateCli.PushToUsers(ctx, req)
}
