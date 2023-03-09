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
	"github.com/airy/component"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/interfaces"
	"github.com/airy/pb"
	"time"
)

type Module interface {
	Init(config *config.AiryConfig)
	AfterInit()
	BeforeShutdown()
	Shutdown() error
	Run(chan struct{})
	Type() constants.ModuleType
	SetType(typ constants.ModuleType)
	AddComponents(c ...component.Component)
	Server() *config.Server
	SetSerializer(interfaces.Serializer)
	Serializer() interfaces.Serializer
	SetCodec(interfaces.PacketCodec)
	Codec() interfaces.PacketCodec
	SetTCPPort(int)
	SetGRPCServerPort(port int)
	SetHeartBeat(d time.Duration)
	SetCompress(compress interfaces.Compress)
	Compress() interfaces.Compress
	HeartBeat() time.Duration
	GPRCServerComponent() *component.GRPCServer
	GRPCClientComponent() *component.GRPCClient
	RedisComponent() *component.Redis
	MongoComponent() *component.Mongo
	MysqlComponent() *component.Mysql
	AcceptorComponent() *component.TCPAcceptor
	ETCDComponent() *component.EtcdServiceDiscovery
	SetTimer(interval time.Duration, slotNum int)
	Timer() interfaces.Timer
	HandRoute(ctx context.Context, request *pb.Packet)
	AddRouter(routerCode int32, router interfaces.Router)
	AddRouters(map[int32]interfaces.Router)
	SetAuthorization(authorization bool)
	Authorization() bool
	SetAOIManager(manager interfaces.AOIManager)
	CreateAOIManager() interfaces.AOIManager
	SysStorage() interfaces.Storage
}
