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

package client

import (
	"github.com/airy/client/behaviortree"
	"github.com/airy/client/config"
	"github.com/airy/client/logger"
	"github.com/airy/client/util"
	"github.com/airy/interfaces"
	"time"
)

type Player struct {
	Map       *util.MapData
	EnterTime time.Time
	Connected bool
	PID       uint32
	X, Z      float32
	Name      string
	SpaceID   string

	BehaviorTree *behaviortree.NodeBase
	Client       *Client
	Order        int
}

func NewPlayer(addr string, order int, heartBeatInterval int, packetCodec interfaces.PacketCodec, compress interfaces.Compress, serializer interfaces.Serializer, m *util.MapData) *Player {
	pid := uint32(config.Config.PIDPre + order)
	cli := NewClient(heartBeatInterval, packetCodec, compress, serializer)
	p := &Player{
		Map:    m,
		PID:    pid,
		Name:   util.RandomName(),
		Client: cli,
		Order:  order,
	}

	var err error
	err = cli.Bind(p, config.Config.SpaceId)
	if err != nil {
		return nil
	}
	err = cli.ConnectToTCP(addr)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil
	}
	tree := behaviortree.NewTree()
	if tree == nil {
		logger.Log.Error("create player behavior tree failed")
		return nil
	}
	p.BehaviorTree = tree

	return p
}

func (p *Player) Start() {
	go p.behaviorTreeRun()
}

func (p *Player) behaviorTreeRun() {
	ticker := time.NewTicker(time.Millisecond * 150)
	for {
		select {
		case <-ticker.C:
			p.BehaviorTree.Run(p)
		}
	}
}
