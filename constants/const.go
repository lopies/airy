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

package constants

type ComponentType string

const (
	TCPComponent        ComponentType = "tcp"
	ETCDComponent       ComponentType = "etcd"
	GRPCClientComponent ComponentType = "grpc_client"
	GRPCServerComponent ComponentType = "grpc_server"
	MongoComponent      ComponentType = "mongo"
	MysqlComponent      ComponentType = "mysql"
	RedisComponent      ComponentType = "redis"
)

type ModuleType string

const (
	GateModule  ModuleType = "gate"
	GMSModule   ModuleType = "gms"
	LogicModule ModuleType = "logic"
)

// SessionCtxKey is the context key where the session will be set
var SessionCtxKey = "session"

// SpaceCtxKey is the context key where the space will be set
var SpaceCtxKey = "space"

// GRPCHostKey is the key for grpc host on server metadata
var GRPCHostKey = "grpcHost"

// GRPCExternalHostKey is the key for grpc external host on server metadata
var GRPCExternalHostKey = "grpc-external-host"

// GRPCPortKey is the key for grpc port on server metadata
var GRPCPortKey = "grpcPort"

var AcceptorPort = "acceptorPort"

// GRPCExternalPortKey is the key for grpc external port on server metadata
var GRPCExternalPortKey = "grpc-external-port"

// RegionKey is the key to save the region server is on
var RegionKey = "region"
