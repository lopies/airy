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
	"github.com/airy/interfaces"
	"github.com/airy/player"
	"math"
)

// WorldManager is an implementation of no limit
type WorldManager struct {
	players map[uint32]*player.Player // The collection of players currently online
	num     int                       // aoi people number
}

// NewGridAOIManager initializes an AOI region
func NewWorldManager() interfaces.AOIManager {
	aoiMgr := &WorldManager{}
	return aoiMgr
}

func (g *WorldManager) Neighbors(p *player.Player) map[*player.Player]struct{} {
	players := make(map[*player.Player]struct{}, len(g.players))
	for _, pl := range g.players {
		players[pl] = struct{}{}
	}
	return players
}

// Enter is called when Entity enters Space
func (g *WorldManager) Enter(p *player.Player, x, z float32) bool {
	if _, ok := g.players[p.PID]; ok {
		return false
	}
	// Add player to the World Manager
	g.players[p.PID] = p
	g.num++
	// Add player to AOI network planning
	p.X, p.Z = x, z
	return true
}

// Leave is called when Entity leaves Space
func (g *WorldManager) Leave(p *player.Player) {
	// Remove the current player
	delete(g.players, p.PID)
	g.num--
}

// Move is called when Entity moves in Space
func (g *WorldManager) Move(p *player.Player, x, z float32) {
	p.X = x
	p.Z = z
}

// AllPlayers get information on all players
func (g *WorldManager) AllPlayers() []*player.Player {
	players := make([]*player.Player, 0, math.MaxInt8)
	for _, v := range g.players {
		players = append(players, v)
	}
	return players
}

// Num return aoi people number
func (g *WorldManager) Num() int {
	return g.num
}

// String print aoi info
func (g *WorldManager) String() {
}
