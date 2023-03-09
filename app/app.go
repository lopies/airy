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

package app

import (
	"fmt"
	"github.com/airy/component"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/interfaces"
	"github.com/airy/logger"
	"github.com/airy/module"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	opts        = make([]option, 0, 10)
	appName     = "airy"
	appVersion  = "v1.0.0"
	appUsage    = "lightweight distributed scalable game server"
	authorName  = "Franco Zhou"
	authorEmail = "1140282958@qq.com"
	logo        = `
 ______  ______       ______  __  ______  __  __    
/\  ___\/\  __ \     /\  __ \/\ \/\  == \/\ \_\ \   
\ \ \__ \ \ \/\ \    \ \  __ \ \ \ \  __<\ \____ \  
 \ \_____\ \_____\    \ \_\ \_\ \_\ \_\ \_\/\_____\ 
  \/_____/\/_____/     \/_/\/_/\/_/\/_/ /_/\/_____/ 
`
	appInfo = &cli.App{
		Name:    appName,
		Version: appVersion,
		Usage:   appUsage,
		Authors: []*cli.Author{
			{
				Name:  authorName,
				Email: authorEmail,
			},
		},
		Action: func(context *cli.Context) error {
			fmt.Println(fmt.Sprintf("\u001B[34m%s\u001B[0m", logo))
			return nil
		},
	}
)

// coverConfig execute set parameter
func coverConfig(mod module.Module, opts ...option) {
	for _, opt := range opts {
		opt.apply(mod)
	}
}

// RunModuleComponent execute components and modules
func RunModuleComponent(mod module.Module, comp ...component.Component) {
	compMap := make(map[constants.ComponentType]component.Component)
	for _, c := range comp {
		if _, ok := compMap[c.Type()]; ok {
			panic(fmt.Sprintf("duplicate components were registered,component type = %s", c.Type()))
		}
		compMap[c.Type()] = c
	}
	var conf *config.AiryConfig
	switch mod.Type() {
	case constants.GateModule:
		conf = config.NewDefaultAiryConfig(string(mod.Type()))
		if c, ok := compMap[constants.GRPCServerComponent]; ok {
			c.(*component.GRPCServer).SetAiryGateServer(mod.(*module.Gate))
		}
	case constants.LogicModule:
		conf = config.NewDefaultAiryConfig(string(mod.Type()))
		if c, ok := compMap[constants.GRPCServerComponent]; ok {
			c.(*component.GRPCServer).SetAiryLogicServer(mod.(*module.Logic))
		}
	}
	err := module.RegisterModule(mod)
	if err != nil {
		logger.Fatalf("server register module error : %s", err.Error())
		panic(err.Error())
	}
	module.StartModules(conf)
	coverConfig(mod, opts...)
	for _, c := range comp {
		err = component.RegisterComponent(c)
		if err != nil {
			panic(err.Error())
		}
	}
	mod.AddComponents(comp...)
	component.StartComponents(conf)
	mod.Run(conf.DieChan)
	logger.Infof("%s server run successfully,listen signal waiting for exit....", mod.Type())
	_ = appInfo.Run(os.Args)
	signal.Notify(conf.SgChan, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case <-conf.DieChan:
		logger.Infof("the app will shutdown in a few seconds")
	case s := <-conf.SgChan:
		logger.Infof("got signal: ", s, ", shutting down...")
		close(conf.DieChan)
	}
	module.ShutdownModules()
	component.ShutdownComponents()
	logger.Infof("server is stopping...")
}

type option interface {
	apply(module.Module)
}

type FuncOption struct {
	f func(module.Module)
}

func (fo FuncOption) apply(mod module.Module) {
	fo.f(mod)
}

// WithGRPCPort set grpc port
func WithGRPCPort(port int) {
	opts = append(opts, FuncOption{
		f: func(mod module.Module) {
			mod.SetGRPCServerPort(port)
		},
	})
}

// WithTCPPort set tcp port
func WithTCPPort(port int) {
	opts = append(opts, FuncOption{
		f: func(mod module.Module) {
			mod.SetTCPPort(port)
		},
	})
}

// WithSerializer set protocol json | protobuf
func WithSerializer(serializer interfaces.Serializer) {
	opts = append(opts, FuncOption{
		f: func(mod module.Module) {
			mod.SetSerializer(serializer)
		},
	})
}

// WithCodec set codec
func WithCodec(codec interfaces.PacketCodec) {
	opts = append(opts, FuncOption{
		f: func(mod module.Module) {
			mod.SetCodec(codec)
		},
	})
}

// WithHeartBeat set heartbeat interval
func WithHeartBeat(d time.Duration) {
	opts = append(opts, FuncOption{
		f: func(mod module.Module) {
			mod.SetHeartBeat(d)
		},
	})
}

// WithCompress set compress
func WithCompress(compress interfaces.Compress) {
	opts = append(opts, FuncOption{
		f: func(mod module.Module) {
			mod.SetCompress(compress)
		},
	})
}

// WithRouter register router
func WithRouter(routerCode int32, router interfaces.Router) {
	opts = append(opts, FuncOption{
		f: func(mod module.Module) {
			mod.AddRouter(routerCode, router)
		},
	})
}

// WithRouters register routers
func WithRouters(routers map[int32]interfaces.Router) {
	opts = append(opts, FuncOption{
		f: func(mod module.Module) {
			mod.AddRouters(routers)
		},
	})
}

// WithTimer set time wheel config
func WithTimer(interval time.Duration, slotNum int) {
	opts = append(opts, FuncOption{
		f: func(mod module.Module) {
			mod.SetTimer(interval, slotNum)
		},
	})
}

// WithAuthorization whether verification is required, if true,need run authorization server
func WithAuthorization(authorization bool) {
	opts = append(opts, FuncOption{
		f: func(mod module.Module) {
			mod.SetAuthorization(authorization)
		},
	})
}
