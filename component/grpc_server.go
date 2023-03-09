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

package component

import (
	"fmt"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/interfaces"
	"github.com/airy/logger"
	"github.com/airy/mgr"
	"github.com/airy/pb"
	"google.golang.org/grpc"
	"net"
	"strconv"
)

// GRPCServer rpc server struct
type GRPCServer struct {
	server          *config.Server
	port            int
	grpcSv          *grpc.Server
	airyGateServer  pb.AiryGateServer
	airyLogicServer pb.AiryLogicServer
	configured      bool
	BaseComponent
}

func NewGRPCServer() *GRPCServer {
	g := new(GRPCServer)
	g.SetType(constants.GRPCServerComponent)
	return g
}

// Init inits grpc rpc server
func (gs *GRPCServer) Init(conf *config.AiryConfig) {
	if !gs.configured {
		panic("the grpc server is not configured,")
	}
	gs.server = conf.Server
	var grpcStrategy interfaces.InitGRPCStrategy
	switch conf.Server.Type {
	case string(constants.GateModule):
		grpcStrategy = &GateGRPC{
			conf: conf,
			gs:   gs,
		}
	case string(constants.LogicModule):
		grpcStrategy = &LogicGRPC{
			conf: conf,
			gs:   gs,
		}
	}
	grpcStrategy.SetPort()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", gs.port))
	if err != nil {
		logger.Errorf("grpc server init error,err=%s", err.Error())
		panic(fmt.Sprintf("grpc server init error,err=%s", err.Error()))
	}
	gs.grpcSv = grpc.NewServer(grpc.MaxRecvMsgSize(conf.GRPCServer.MaxRecvMessage))
	grpcStrategy.RegisterGRPCServer()
	go gs.grpcSv.Serve(lis)
	ip, err := ipv4()
	if err != nil {
		logger.Fatalf("ip net not found,error : %s", err.Error())
		panic(err)
	}

	if conf.Server.Metadata != nil {
		conf.Server.Metadata[constants.GRPCHostKey] = ip
		conf.Server.Metadata[constants.GRPCPortKey] = strconv.Itoa(gs.port)
	} else {
		conf.Server.Metadata = map[string]string{
			constants.GRPCHostKey: ip,
			constants.GRPCPortKey: strconv.Itoa(gs.port),
		}
	}

	logger.Infof("grpc %s is running on port :%d", conf.Server.Type, gs.port)
}

func (gs *GRPCServer) SetAiryGateServer(server pb.AiryGateServer) *GRPCServer {
	gs.airyGateServer = server
	gs.configured = true
	return gs
}

func (gs *GRPCServer) SetAiryLogicServer(server pb.AiryLogicServer) *GRPCServer {
	gs.airyLogicServer = server
	gs.configured = true
	return gs
}

// Shutdown stops grpc rpc server
func (gs *GRPCServer) Shutdown() error {
	// graceful: stops the server from accepting new connections and RPCs and
	// blocks until all the pending RPCs are finished.
	// source: https://godoc.org/google.golang.org/grpc#Server.GracefulStop
	gs.grpcSv.GracefulStop()
	return nil
}

type GateGRPC struct {
	conf *config.AiryConfig
	gs   *GRPCServer
}

func (g *GateGRPC) SetPort() {
	g.gs.port = mgr.GetPort(g.conf.Port.GatePort)
	g.conf.GatePort = g.gs.port
}

func (g *GateGRPC) RegisterGRPCServer() {
	pb.RegisterAiryGateServer(g.gs.grpcSv, g.gs.airyGateServer)
}

type LogicGRPC struct {
	conf *config.AiryConfig
	gs   *GRPCServer
}

func (l *LogicGRPC) SetPort() {
	l.gs.port = mgr.GetPort(l.conf.Port.LogicPort)
	l.conf.LogicPort = l.gs.port
}

func (l *LogicGRPC) RegisterGRPCServer() {
	pb.RegisterAiryLogicServer(l.gs.grpcSv, l.gs.airyLogicServer)
}

func ipv4() (ipv4Address string, err error) {
	addrs, err := net.InterfaceAddrs()
	for _, addr := range addrs {
		ipNet, isIpNet := addr.(*net.IPNet)
		if isIpNet && !ipNet.IP.IsLoopback() {
			ip := ipNet.IP.To4()
			if ip != nil {
				return ip.String(), nil
			}
		}
	}
	return
}
