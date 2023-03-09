// Copyright (c) Airy Author. All Rights Reserved.
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

package module

import (
	"context"
	"github.com/airy/aoi"
	"github.com/airy/codec"
	"github.com/airy/component"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/interfaces"
	"github.com/airy/logger"
	"github.com/airy/pb"
	"github.com/airy/serializer/protobuf"
	"github.com/airy/timer"
	"time"
)

type BaseModule struct {
	//
	conf       *config.AiryConfig
	typ        constants.ModuleType
	serializer interfaces.Serializer
	codec      interfaces.PacketCodec
	compress   interfaces.Compress
	timer      interfaces.Timer
	aoiManager interfaces.AOIManager
	// logic server router
	routers map[int32]interfaces.Router

	redis         *component.Redis
	mysql         *component.Mysql
	mongo         *component.Mongo
	grpcClient    *component.GRPCClient
	grpcServer    *component.GRPCServer
	etcdDiscovery *component.EtcdServiceDiscovery
	tcpAcceptor   *component.TCPAcceptor
}

// Init runs initialization
func (bm *BaseModule) Init(config *config.AiryConfig) {
	bm.serializer = protobuf.NewSerializer()
	bm.routers = make(map[int32]interfaces.Router)

	bm.conf = config

	bm.redis = new(component.Redis)
	bm.mysql = new(component.Mysql)
	bm.mongo = new(component.Mongo)
	bm.grpcClient = new(component.GRPCClient)
	bm.grpcServer = new(component.GRPCServer)
	bm.etcdDiscovery = new(component.EtcdServiceDiscovery)
	bm.tcpAcceptor = new(component.TCPAcceptor)
}

// AfterInit runs after initialization
func (bm *BaseModule) AfterInit() {}

// BeforeShutdown runs before shutdown
func (bm *BaseModule) BeforeShutdown() {
	if bm.timer != nil {
		bm.timer.Stop()
	}
}

// Shutdown runs module stop
func (bm *BaseModule) Shutdown() error {
	return nil
}

func (bm *BaseModule) Run() {}

// Type return module type
func (bm *BaseModule) Type() constants.ModuleType {
	return bm.typ
}

func (bm *BaseModule) SetType(typ constants.ModuleType) {
	bm.typ = typ
}

// AddComponent run module add components
func (bm *BaseModule) AddComponents(comps ...component.Component) {
	for _, c := range comps {
		switch c.Type() {
		case constants.RedisComponent:
			bm.redis = c.(*component.Redis)
		case constants.MongoComponent:
			bm.mongo = c.(*component.Mongo)
		case constants.MysqlComponent:
			bm.mysql = c.(*component.Mysql)
		case constants.TCPComponent:
			bm.tcpAcceptor = c.(*component.TCPAcceptor)
		case constants.ETCDComponent:
			bm.etcdDiscovery = c.(*component.EtcdServiceDiscovery)
			bm.etcdDiscovery.AddListener(bm.grpcClient)
		case constants.GRPCClientComponent:
			bm.grpcClient = c.(*component.GRPCClient)
		case constants.GRPCServerComponent:
			bm.grpcServer = c.(*component.GRPCServer)
		default:
			logger.Fatalf("unsupported component type : %s", c.Type())
		}
	}
}

func (bm *BaseModule) SetSerializer(serializer interfaces.Serializer) {
	bm.serializer = serializer
}

func (bm *BaseModule) Serializer() interfaces.Serializer {
	return bm.serializer
}

func (bm *BaseModule) SetTimer(interval time.Duration, slotNum int) {
	t := timer.NewTimeWheel(interval, slotNum)
	bm.timer = t
}

func (bm *BaseModule) Timer() interfaces.Timer {
	if bm.timer == nil {
		t := timer.NewTimeWheel(time.Second, 3600)
		bm.timer = t
		// 启动时间轮
		t.Start()
	}
	return bm.timer
}

func (bm *BaseModule) SetCodec(codec interfaces.PacketCodec) {
	bm.codec = codec
}

func (bm *BaseModule) Codec() interfaces.PacketCodec {
	if bm.codec == nil {
		return codec.NewLengthFieldBasedFrameCodec()
	}
	return bm.codec
}

func (bm *BaseModule) SetTCPPort(port int) {
	bm.conf.Port.TCPPort = port
}

func (bm *BaseModule) SetGRPCServerPort(port int) {
	switch bm.conf.Server.Type {
	case string(constants.GateModule):
		bm.conf.Port.GatePort = port
	case string(constants.GMSModule):
		bm.conf.Port.GMSPort = port
	case string(constants.LogicModule):
		bm.conf.Port.LogicPort = port
	}
}

func (bm *BaseModule) Server() *config.Server {
	return bm.conf.Server
}

func (bm *BaseModule) SetHeartBeat(d time.Duration) {
	bm.conf.Heartbeat.Interval = d
}

func (bm *BaseModule) HeartBeat() time.Duration {
	return bm.conf.Heartbeat.Interval
}

func (bm *BaseModule) GPRCServerComponent() *component.GRPCServer {
	return bm.grpcServer
}

func (bm *BaseModule) GRPCClientComponent() *component.GRPCClient {
	return bm.grpcClient
}

func (bm *BaseModule) RedisComponent() *component.Redis {
	return bm.redis
}

func (bm *BaseModule) MongoComponent() *component.Mongo {
	return bm.mongo
}

func (bm *BaseModule) MysqlComponent() *component.Mysql {
	return bm.mysql
}

func (bm *BaseModule) AcceptorComponent() *component.TCPAcceptor {
	return bm.tcpAcceptor
}

func (bm *BaseModule) ETCDComponent() *component.EtcdServiceDiscovery {
	return bm.etcdDiscovery
}

func (bm *BaseModule) SetCompress(compress interfaces.Compress) {
	bm.compress = compress
}

func (bm *BaseModule) Compress() interfaces.Compress {
	return bm.compress
}

func (bm *BaseModule) AddRouter(routerCode int32, router interfaces.Router) {
	if bm.conf.Server.Type == string(constants.LogicModule) {
		if _, ok := bm.routers[routerCode]; ok {
			logger.Warnf("route_code = %s has been registered")
			return
		}
		router.Init(bm.serializer)
		bm.routers[routerCode] = router
		logger.Infof("register router successfully")
		return
	}
	logger.Warnf("service type does not belong to logic")
}

func (bm *BaseModule) AddRouters(routers map[int32]interfaces.Router) {
	if bm.conf.Server.Type == string(constants.LogicModule) {
		if len(bm.routers) != 0 {
			logger.Warnf("the previous router configuration was overwritten")
		}
		for _, router := range routers {
			router.Init(bm.serializer)
		}
		bm.routers = routers
		logger.Infof("register router successfully")
		return
	}
	logger.Warnf("service type does not belong to logic")
}

func (bm *BaseModule) HandRoute(ctx context.Context, request *pb.Packet) {
	requestCode := request.RequestCode
	h := bm.routers[requestCode]
	if h != nil {
		h.Handle(ctx, request)
		h.Post(ctx, request)
		return
	}
	logger.Warnf("not found request router: %d", requestCode)
}

func (bm *BaseModule) SetAuthorization(authorization bool) {
	bm.conf.Server.Authorization = authorization
}

func (bm *BaseModule) Authorization() bool {
	return bm.conf.Server.Authorization
}

func (bm *BaseModule) SetAOIManager(manager interfaces.AOIManager) {
	bm.aoiManager = manager
}

func (bm *BaseModule) CreateAOIManager() interfaces.AOIManager {
	if _, ok := bm.aoiManager.(*aoi.GridAOIManager); ok {
		return new(aoi.GridAOIManager)
	}
	return new(aoi.XZListAOIManager)
}

func (bm *BaseModule) SysStorage() interfaces.Storage {
	return bm.redis
}
