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
	"time"
)

type Condition struct{}

func (c *Condition) EnterSpaceCondition(p *client.Player, childSLice []*behaviortree.NodeBase) bool {
	if p.Connected {
		return false
	}
	if p.EnterTime.Add(time.Second * 10).After(time.Now()) {
		return false
	}
	return true
}

func (c *Condition) MoveCondition(p *client.Player, childSLice []*behaviortree.NodeBase) bool {
	if !p.Connected {
		return false
	}
	return true
}

func (c *Condition) DelayCondition(p *client.Player, childSLice []*behaviortree.NodeBase) bool {
	if !p.Connected {
		logger.Log.Warnf("[DelayCondition]player[%d] status is not connected", p.PID)
		return false
	}
	return true
}

func (c *Condition) GlobalChatCondition(p *client.Player, childSLice []*behaviortree.NodeBase) bool {
	if !p.Connected {
		logger.Log.Warnf("[GlobalChatCondition]player[%d] status is not connected", p.PID)
		return false
	}
	return true
}
