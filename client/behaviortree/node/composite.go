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
)

type Composite struct{}

/*
Selector select a node，true if one returns true
*/
func (c *Composite) Selector(p *client.Player, childSlice []*behaviortree.NodeBase) bool {
	//logger.Log.Infof("%s do Selector", player.BaseData.Name)
	if childSlice != nil {
		childSize := len(childSlice)
		for i := 0; i < childSize; i++ {
			result := childSlice[i].Run(p)
			if !result {
				return false
			}
		}
	}
	return true
}

/*
Sequence sequential nodes, all return success is true, otherwise false
*/
func (c *Composite) Sequence(p *client.Player, childSlice []*behaviortree.NodeBase) bool {
	//logger.Log.Infof("%s do Sequence", player.BaseData.Name)
	if childSlice != nil {
		childSize := len(childSlice)
		for i := 0; i < childSize; i++ {
			result := childSlice[i].Run(p)
			if !result {
				return false
			}
		}
		return true
	}
	return false
}
