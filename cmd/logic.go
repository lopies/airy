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

package main

import (
	"github.com/airy/app"
	"github.com/airy/component"
	"github.com/airy/constants"
	"github.com/airy/module"
	"github.com/airy/router"
	"github.com/airy/serializer/protobuf"
)

func main() {
	m := module.NewLogic()
	//reset grpc server port
	app.WithGRPCPort(9400)
	//reset serializer
	app.WithSerializer(protobuf.NewSerializer())
	//app.WithCompress(compress.NewSnappyCompress())
	//register router
	app.WithRouter(constants.HeartBeatRequest, &router.RouteHeartBeat{})
	app.WithRouter(constants.JoinRequest, &router.RouteJoin{})
	app.WithRouter(constants.MoveRequest, &router.RouteMove{})
	app.WithRouter(constants.LeaveRequest, &router.RouteLeave{})
	app.WithRouter(constants.DelayRequest, &router.RouteDelay{})
	app.WithRouter(constants.GlobalChatRequest, &router.RouteGlobalChat{})
	//add component...
	comp := []component.Component{
		component.NewGRPCClient(),
		component.NewGRPCServer(),
		component.NewRedis(),
		component.NewETCD(),
	}
	//run module
	app.RunModuleComponent(m, comp...)
}
