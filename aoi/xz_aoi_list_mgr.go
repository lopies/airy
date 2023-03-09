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

package aoi

import (
	"fmt"
	"github.com/airy/interfaces"
	"github.com/airy/logger"
	"github.com/airy/player"
	"math"
)

// XZListAOIManager is an implementation of AOICalculator using XZ lists
type XZListAOIManager struct {
	xSweepList *xAOIList
	zSweepList *zAOIList
	players    map[uint32]*player.Player
	num        int
}

// NewXZListAOIManager creates a new XZListAOIManager
func NewXZListAOIManager() interfaces.AOIManager {
	return &XZListAOIManager{
		xSweepList: newXAOIList(dist),
		zSweepList: newZAOIList(dist),
		players:    map[uint32]*player.Player{},
	}
}

// Neighbors return players in a nearby field of view
func (xzm *XZListAOIManager) Neighbors(p *player.Player) map[*player.Player]struct{} {
	neighbors := make(map[*player.Player]struct{}, len(p.Neighbors))
	for player, _ := range p.Neighbors {
		neighbors[player] = struct{}{}
	}
	return neighbors
}

// Enter is called when Entity enters Space
func (xzm *XZListAOIManager) Enter(p *player.Player, x, z float32) bool {
	if _, ok := xzm.players[p.PID]; ok {
		return false
	}
	p.X, p.Z = x, z
	xzm.xSweepList.Insert(p)
	xzm.zSweepList.Insert(p)
	xzm.adjust(p)
	xzm.num++
	xzm.players[p.PID] = p
	logger.Debugf("one player joins,pid = %d", p.PID)
	return true
}

// Leave is called when Entity leaves Space
func (xzm *XZListAOIManager) Leave(p *player.Player) {
	xzm.xSweepList.Remove(p)
	xzm.zSweepList.Remove(p)
	xzm.adjust(p)

	delete(xzm.players, p.PID)
	xzm.num--
	logger.Debugf("one player leaves,pid = %d", p.PID)
}

// Move is called when Entity moves in Space
func (xzm *XZListAOIManager) Move(p *player.Player, x, z float32) {
	oldX := p.X
	oldZ := p.Z
	p.X, p.Z = x, z
	if oldX != x {
		xzm.xSweepList.Move(p, oldX)
	}
	if oldZ != z {
		xzm.zSweepList.Move(p, oldZ)
	}
	xzm.adjust(p)
	logger.Debugf("one player moves,pid = %d,x = %f,z = %f", p.PID, p.X, p.Z)
}

// AllPlayers return all players
func (xzm *XZListAOIManager) AllPlayers() []*player.Player {
	//创建返回的player集合切片
	players := make([]*player.Player, 0, math.MaxInt8)
	//添加切片
	for _, v := range xzm.players {
		players = append(players, v)
	}
	//返回
	return players
}

// Num return aoi people number
func (xzm *XZListAOIManager) Num() int {
	return xzm.num
}

// adjust is called by Entity to adjust neighbors
func (xzm *XZListAOIManager) adjust(p *player.Player) {
	xzm.xSweepList.markNeighbors(p)
	for neighbor := range p.Neighbors {
		if neighbor.MarkVal == 0 {
			delete(p.Neighbors, neighbor)
			delete(neighbor.Neighbors, p)
		}
	}
	xzm.xSweepList.clearMark(p)
}

// String print aoi info
func (xzm *XZListAOIManager) String() {
	fmt.Print("X : ")
	for node := xzm.xSweepList.head; node != nil; node = node.XNext {
		fmt.Printf("--->(%f,%f)", node.X, node.Z)
	}
	fmt.Printf("\nY : ")
	for node := xzm.zSweepList.head; node != nil; node = node.ZNext {
		fmt.Printf("--->(%f,%f)", node.Z, node.X)
	}
	fmt.Println()
}
