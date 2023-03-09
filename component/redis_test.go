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

package component

import (
	"fmt"
	"github.com/airy/config"
	"testing"
)

func TestStringStore(t *testing.T) {
	r := NewRedis()
	r.Init(config.NewDefaultAiryConfig(""))
	err := r.StoreSpaceIDLogicID("space1", "logic1")
	if err != nil {
		fmt.Println("err : ", err.Error())
	}
}

func TestStringGet(t *testing.T) {
	r := NewRedis()
	r.Init(config.NewDefaultAiryConfig(""))
	serverID, err := r.GetLogicIDBySpaceID("space1")
	if err != nil {
		fmt.Printf("err : %s", err.Error())
	}
	fmt.Println("serverID = ", serverID)
}

func TestStringDel(t *testing.T) {
	r := NewRedis()
	r.Init(config.NewDefaultAiryConfig(""))
	err := r.DelSpaceIDLogicID("space1", "logic1")
	if err != nil {
		fmt.Printf("err : %s", err.Error())
	}
}

func TestSAdd(t *testing.T) {
	r := NewRedis()
	r.Init(config.NewDefaultAiryConfig(""))
	r.StoreLogicIDSpaceID("server1", "space1")
	r.StoreLogicIDSpaceID("server1", "space2")
	r.StoreLogicIDSpaceID("server1", "space3")
}

func TestSAddDEL(t *testing.T) {
	r := NewRedis()
	r.Init(config.NewDefaultAiryConfig(""))
	r.DelLogicIDSpaceID("server1")
}
