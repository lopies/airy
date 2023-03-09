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
	"fmt"
	"github.com/airy/client/behaviortree"
	"github.com/airy/client/behaviortree/node"
	"github.com/airy/client/client"
	"github.com/airy/client/config"
	"github.com/airy/client/logger"
	"github.com/airy/client/util"
	codec2 "github.com/airy/codec"
	"github.com/airy/serializer/protobuf"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

func main() {
	parseConfig()
	acceptorAddress := addr()
	logger.InitLogger(config.Config.LogLevel)
	behaviortree.InitBehaviorTreeXmlConfig()
	registerNode()
	mapData, err := loadSpaceMapInfo()
	if err != nil {
		panic(err.Error())
	}
	heartBeatInterval := 5
	codec := codec2.NewLengthFieldBasedFrameCodec()
	//compress := compress2.NewSnappyCompress()
	seria := protobuf.NewSerializer()
	for i := 1; i <= config.Config.PlayerNum; i++ {
		p := client.NewPlayer(acceptorAddress, i, heartBeatInterval, codec, nil, seria, mapData)
		if p != nil {
			p.Start()
			randSleep := rand.Intn(100) //1秒内随机睡眠
			logger.Log.Infof("player[%s] create successfully", p.PID)
			time.Sleep(time.Duration(randSleep) * time.Millisecond)
		} else {
			logger.Log.Errorf("player[%d] create successfully", p.PID)
			i--
		}
	}
	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	select {
	case <-sg:
		logger.Log.Infof("shut down!!!")
	}
}

func parseConfig() {
	file, err := os.Open("./config.yaml")
	if err != nil {
		panic(fmt.Errorf("player config err: %v", err.Error()))
	}
	defer file.Close()
	robotConfigData, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	config.Config = &config.ClientConfig{}
	err = yaml.Unmarshal(robotConfigData, config.Config)
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
}
func addr() (addr string) {
	switch config.Config.Env {
	case "dev":
		return config.Dev
	case "beta":
		return config.Beta
	case "release":
		return config.Release
	case "local":
		return config.Local
	default:
		panic("env config error")
	}
}

func registerNode() {
	behaviortree.TypeMap["Action"] = reflect.TypeOf(node.Action{})
	behaviortree.TypeMap["Composite"] = reflect.TypeOf(node.Composite{})
	behaviortree.TypeMap["Condition"] = reflect.TypeOf(node.Condition{})
}

//初始化加载地图数据
func loadSpaceMapInfo() (*util.MapData, error) {
	if config.Config.SpaceId == "" {
		panic("spaceID is empty")
	}
	mapData, err := util.ReadMapData(config.Config.SpaceId)
	if err != nil {
		return nil, err
	}
	return mapData, nil
}
