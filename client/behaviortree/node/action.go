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

package node

import (
	"github.com/airy/client/behaviortree"
	"github.com/airy/client/client"
	"github.com/airy/client/logger"
	"github.com/airy/constants"
	"github.com/airy/pb"
	"math/rand"
	"time"
)

type Action struct {
}

func (a *Action) EnterSpace(p *client.Player, childSlice []*behaviortree.NodeBase) bool {
	if !p.Client.Connected {
		logger.Log.Errorf("[EnterSpace]the connection not established")
		return false
	}
	result := verifyCondition(p, childSlice)
	if !result {
		return true
	}
	logger.Log.Debugf("Exec EnterSpace")
	msg := &pb.Enter{
		PID:  p.PID,
		Name: p.Name,
		X:    p.X,
		Z:    p.Z,
	}
	//logger.RobotLog.Errorf("%s image=%s", r.BaseData.Name, msg.I)
	data, _ := p.Client.Serializer.Marshal(msg)
	p.Client.SendRequest(pb.Type_Join, constants.JoinRequest, data, true)
	p.EnterTime = time.Now()
	return true
}

func (a *Action) Move(p *client.Player, childSlice []*behaviortree.NodeBase) bool {
	if !p.Client.Connected {
		logger.Log.Errorf("[Move]the connection not established")
		return false
	}
	result := verifyCondition(p, childSlice)
	if !result {
		return true
	}
	for {
		r := rand.Intn(4)
		switch r {
		case 0:
			p.Z = p.Z - 1
		case 1:
			p.Z = p.Z + 1
		case 2:
			p.X = p.X - 1
		case 3:
			p.X = p.X + 1
		}
		if p.X <= float32(p.Map.Width) && p.X >= 0 && p.Z <= float32(p.Map.Height) && p.Z >= 0 {
			break
		}
	}
	logger.Log.Debugf("Exec Move,X = %f,Z = %f", p.X, p.Z)
	msg := &pb.Move{
		X: p.X,
		Z: p.Z,
	}
	//logger.RobotLog.Errorf("%s image=%s", r.BaseData.Name, msg.I)
	data, _ := p.Client.Serializer.Marshal(msg)
	p.Client.SendRequest(pb.Type_Request_, constants.MoveRequest, data, false)
	return true
}

func (a *Action) Delay(p *client.Player, childSlice []*behaviortree.NodeBase) bool {
	if !p.Client.Connected {
		logger.Log.Errorf("[Delay]the connection not established")
		return false
	}
	result := verifyCondition(p, childSlice)
	if !result {
		return true
	}
	logger.Log.Debugf("Exec Delay")
	msg := &pb.Delay{
		StartTime: time.Now().UnixMilli(),
	}
	data, _ := p.Client.Serializer.Marshal(msg)
	p.Client.SendRequest(pb.Type_Request_, constants.DelayRequest, data, false)
	return true
}

func (a *Action) GlobalChat(p *client.Player, childSlice []*behaviortree.NodeBase) bool {
	if !p.Client.Connected {
		logger.Log.Errorf("[GlobalChat]the connection not established")
		return false
	}
	result := verifyCondition(p, childSlice)
	if !result {
		return true
	}
	logger.Log.Debugf("Exec GlobalChat")
	msg := &pb.CommonChat{
		Content: "hello world",
	}
	//logger.RobotLog.Errorf("%s image=%s", r.BaseData.Name, msg.I)
	data, _ := p.Client.Serializer.Marshal(msg)
	p.Client.SendRequest(pb.Type_Request_, constants.GlobalChatRequest, data, false)
	return true
}

// verifyCondition verify all condition nodes
func verifyCondition(p *client.Player, conditions []*behaviortree.NodeBase) bool {
	if conditions != nil && len(conditions) > 0 {
		childSize := len(conditions)
		//Condition都为true则执行
		for i := 0; i < childSize; i++ {
			childNode := conditions[i]
			if childNode.Name != "Condition" {
				continue
			}
			result := childNode.Run(p)
			if !result {
				return false
			}
		}
	}
	return true
}
