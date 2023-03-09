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

package behaviortree

import (
	"github.com/airy/client/logger"
	"reflect"
)

var TypeMap = make(map[string]reflect.Type)

//type Node interface {}

type NodeBase struct {
	Name           string      `xml:"Name,attr"`
	Method         string      `xml:"Method,attr"`
	ChildNodeSlice []*NodeBase `xml:"Node"`

	MethodValue reflect.Value
}

func (n *NodeBase) Init() {
	logger.Log.Debugf("%s Init", n.Name)
	t := TypeMap[n.Name]
	if t != nil {
		realNode := reflect.New(t)
		n.MethodValue = realNode.MethodByName(n.Method)
	}
	if n.ChildNodeSlice == nil {
		return
	}
	childSize := len(n.ChildNodeSlice)
	for i := 0; i < childSize; i++ {
		n.ChildNodeSlice[i].Init()
	}
}

func (n *NodeBase) Run(player interface{}) bool {
	//logger.Log.Errorf("method point:%d,player point:%d", &runMethod, player)
	if !n.MethodValue.IsValid() {
		logger.Log.Errorf("Method -%s- not exist in %s", n.Method, n.Name)
		return false
	}
	if player == nil {
		return false
	}
	callValue := []reflect.Value{reflect.ValueOf(player), reflect.ValueOf(n.ChildNodeSlice)}
	ret := n.MethodValue.Call(callValue)
	return ret[0].Bool()
}
