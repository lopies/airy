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
	"github.com/airy/codec"
	"github.com/airy/component"
	"github.com/airy/module"
	"github.com/airy/serializer/protobuf"
	"time"
)

func main() {
	m := module.NewGate()
	app.WithGRPCPort(9000)
	app.WithCodec(codec.NewLengthFieldBasedFrameCodec())
	app.WithTCPPort(8800)
	app.WithSerializer(protobuf.NewSerializer())
	app.WithHeartBeat(5 * time.Minute)
	//app.WithCompress(compress.NewSnappyCompress())

	comp := []component.Component{component.NewTCPAcceptor(),
		component.NewGRPCClient(),
		component.NewGRPCServer(),
		component.NewRedis(),
		component.NewETCD(),
	}
	app.RunModuleComponent(m, comp...)
}
