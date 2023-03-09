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

package space

import (
	"github.com/airy/aoi"
	"github.com/airy/interfaces"
	"github.com/airy/logger"
	"github.com/airy/pb"
	"github.com/airy/util"
	"time"
)

var spaceCloseCallBack func(space *Space) error
var spaceCreateCallBack func(space *Space) error

type Space struct {
	close                 bool
	id                    string
	goMap                 [][]bool
	mapWidth              int
	mapHeight             int
	aoiManager            interfaces.AOIManager
	closeChan             chan struct{}
	requestChan           chan *pb.Packet
	serverCombinationChan chan *pb.ServerCombination
}

func NewSpace(id string, aoiManager interfaces.AOIManager) (*Space, error) {
	//load map
	mapData, err := util.ReadMapData(id)
	if err != nil {
		return nil, err
	}
	mapWidth, mapHeight := mapData.Width, mapData.Height
	goMap := mapData.Map
	var am interfaces.AOIManager
	if _, ok := aoiManager.(*aoi.XZListAOIManager); ok {
		am = aoi.NewXZListAOIManager()
	} else {
		am = aoi.NewGridAOIManager(0, mapWidth-1, mapWidth, 0, mapHeight-1, mapHeight)
	}
	space := &Space{
		id:          id,
		goMap:       goMap,
		mapWidth:    mapWidth,
		mapHeight:   mapHeight,
		aoiManager:  am,
		closeChan:   make(chan struct{}),
		requestChan: make(chan *pb.Packet, 1<<12),
	}
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		defer func() {
			ticker.Stop()
		}()
		for {
			select {
			case <-ticker.C:
				logger.Infof("space = %s players num:%d,requestchan length:%d", space.id, space.aoiManager.Num(), len(space.requestChan))
			case <-space.closeChan:
				return
			}
		}
	}()
	spaceCreateCallBack(space)
	return space, nil
}

func BindCallBack(cb1, cb2 func(space *Space) error) {
	spaceCreateCallBack = cb1
	spaceCloseCallBack = cb2
}

func (s *Space) BindChan(serverCombinationChan chan *pb.ServerCombination) {
	s.serverCombinationChan = serverCombinationChan
}

func (s *Space) SendMessage(msg *pb.ServerCombination) {
	s.ServerCombinationChan() <- msg
}

func (s *Space) ServerCombinationChan() chan *pb.ServerCombination {
	return s.serverCombinationChan
}

func (s *Space) ID() string {
	return s.id
}

func (s *Space) AOIManager() interfaces.AOIManager {
	return s.aoiManager
}

func (s *Space) WriteRequestBuffer(c *pb.Packet) {
	if !s.close {
		s.requestChan <- c
	}
}

func (s *Space) ReadRequestBuffer() chan *pb.Packet {
	return s.requestChan
}

func (s *Space) Close() {
	if s.close {
		return
	}
	s.close = true
	spaceCloseCallBack(s)
	close(s.closeChan)
}

func (s *Space) IsClose() bool {
	return s.close
}
