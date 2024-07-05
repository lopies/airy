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

// GridAOIManager is an implementation of AOICalculator using grid
type GridAOIManager struct {
	MinX    int                       // Grid left boundary coordinates
	MaxX    int                       // Grid right boundary coordinates
	CntsX   int                       // The number of grids in the x direction
	MinY    int                       // Grid up boundary coordinates
	MaxY    int                       // Grid down boundary coordinates
	CntsY   int                       // The number of grids in the y direction
	grIDs   map[int]*GrID             // Grid in current area, key= grid ID, value= grid object
	players map[uint32]*player.Player // The collection of players currently online
	num     int                       // aoi people number
}

// NewGridAOIManager initializes an AOI region
func NewGridAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) interfaces.AOIManager {
	aoiMgr := &GridAOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grIDs: make(map[int]*GrID),
	}

	// initialize all grids in the AOI region
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// grid id：ID = IDy *nx + IDx
			gID := y*cntsX + x
			// initializes a grid in the map of AOI, where key is the ID of the current cell
			aoiMgr.grIDs[gID] = NewGrID(gID,
				aoiMgr.MinX+x*aoiMgr.grIDWIDth(),
				aoiMgr.MinX+(x+1)*aoiMgr.grIDWIDth(),
				aoiMgr.MinY+y*aoiMgr.grIDLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.grIDLength())
		}
	}
	return aoiMgr
}

func (g *GridAOIManager) Neighbors(p *player.Player) map[*player.Player]struct{} {
	// Get all Pids for the current AOI region
	pids := g.getPIDsByPos(p.X, p.Z)
	// Add the Player corresponding to all Pids to the Player slice
	players := make(map[*player.Player]struct{}, len(pids))
	for _, pid := range pids {
		players[g.getPlayerByPID(pid)] = struct{}{}
	}
	return players
}

// Enter is called when Entity enters Space
func (g *GridAOIManager) Enter(p *player.Player, x, z float32) bool {
	if _, ok := g.players[p.PID]; ok {
		return false
	}
	// Add player to the World Manager
	g.players[p.PID] = p
	g.num++
	// Add player to AOI network planning
	p.X, p.Z = x, z
	g.addToGrIDByPos(p.PID, x, z)
	return true
}

// Leave is called when Entity leaves Space
func (g *GridAOIManager) Leave(p *player.Player) {
	// Remove the current player from AOI
	g.removeFromGrIDByPos(p.PID, p.X, p.Z)
	g.removePlayerByPID(p.PID)
}

// Move is called when Entity moves in Space
func (g *GridAOIManager) Move(p *player.Player, x, z float32) {
	oldGID := g.getGIDByPos(p.X, p.Z)
	newGID := g.getGIDByPos(x, z)
	if oldGID != newGID {
		// Remove pid from an old grid
		g.removePIDFromGrID(p.PID, oldGID)
		// Add pid to a new grid
		g.addPIDToGrID(p.PID, newGID)
	}
	p.X, p.Z = x, z
}

// AllPlayers get information on all players
func (g *GridAOIManager) AllPlayers() []*player.Player {
	players := make([]*player.Player, 0, math.MaxInt8)
	for _, v := range g.players {
		players = append(players, v)
	}
	return players
}

// Num return aoi people number
func (g *GridAOIManager) Num() int {
	return g.num
}

// String print aoi info
func (g *GridAOIManager) String() {
}

// grIDWIDth get the width of each grid in the x direction
func (g *GridAOIManager) grIDWIDth() int {
	return (g.MaxX - g.MinX) / g.CntsX
}

// grIDLength get the height of each grid in the x direction
func (g *GridAOIManager) grIDLength() int {
	return (g.MaxY - g.MinY) / g.CntsY
}

// getSurroundGrIDsByGID according to the gID of the grid, the information of the current surrounding nine grids is obtained
func (g *GridAOIManager) getSurroundGrIDsByGID(gid int) (grIDs []*GrID) {
	if _, ok := g.grIDs[gid]; !ok {
		return
	}
	grIDs = append(grIDs, g.grIDs[gid])
	x, y := gid%g.CntsX, gid/g.CntsX
	surroundGID := make([]int, 0)
	// new eight direction vector: top left: (1, 1), left in: (1, 0), bottom left: (1, 1), upper: (0, 1), middle: (0, 1), the upper right (1, 1)
	// Middle right: (1, 0), bottom right: (1, 1), respectively write the direction vectors of the 8 directions into the component array of x and y in order
	dx := []int{-1, -1, -1, 0, 0, 1, 1, 1}
	dy := []int{-1, 0, 1, -1, 1, -1, 0, 1}

	// According to the 8 direction vectors, the relative coordinates of the surrounding points are obtained, the coordinates that do not cross the boundary are picked out, and the coordinates are converted to gID
	for i := 0; i < 8; i++ {
		newX := x + dx[i]
		newY := y + dy[i]

		if newX >= 0 && newX < g.CntsX && newY >= 0 && newY < g.CntsY {
			surroundGID = append(surroundGID, newY*g.CntsX+newX)
		}
	}

	// The grid information is obtained based on the gID that does not cross the boundary
	for _, gid := range surroundGID {
		grIDs = append(grIDs, g.grIDs[gid])
	}
	return
}

// getPIDsByPos you get all the PlayerIDs in the surrounding nine grids
func (g *GridAOIManager) getPIDsByPos(x, y float32) (playerIDs []uint32) {
	// According to the horizontal and vertical coordinates to get the current coordinates belong to which grid ID
	gid := g.getGIDByPos(x, y)

	// According to the ID of the grid to get the information of the surrounding nine grids
	grIDs := g.getSurroundGrIDsByGID(gid)
	for _, v := range grIDs {
		playerIDs = append(playerIDs, v.getPlayerIDs()...)
		//fmt.Printf("===> grID ID : %d, pIDs : %v  ====", v.GID, v.GetPlyerIDs())
	}
	return
}

// getGIDByPos obtain the corresponding grid ID by the horizontal and vertical coordinates
func (g *GridAOIManager) getGIDByPos(x, y float32) int {
	gx := (int(x) - g.MinX) / g.grIDWIDth()
	gy := (int(y) - g.MinY) / g.grIDLength()
	return gy*g.CntsX + gx
}

// addPIDToGrID add a PlayerID to a grid
func (g *GridAOIManager) addPIDToGrID(pid uint32, gid int) {
	g.grIDs[gid].add(pid)
}

// addToGrIDByPos add a Player to a grid by horizontal and vertical coordinates
func (g *GridAOIManager) addToGrIDByPos(pid uint32, x, z float32) {
	gID := g.getGIDByPos(x, z)
	grID := g.grIDs[gID]
	grID.add(pid)
}

// removeFromGrIDByPos remove a Player from the corresponding grid by horizontal and vertical coordinates
func (g *GridAOIManager) removeFromGrIDByPos(pid uint32, x, z float32) {
	gid := g.getGIDByPos(x, z)
	grID := g.grIDs[gid]
	grID.remove(pid)
}

// removePIDFromGrID removes a PlayerID in a grid
func (g *GridAOIManager) removePIDFromGrID(pid uint32, gid int) {
	g.grIDs[gid].remove(pid)
}

// removePlayerByPID remove a player from the player information table
func (g *GridAOIManager) removePlayerByPID(pid uint32) {
	delete(g.players, pid)
	g.num--
}

// getPlayerByPID get the corresponding player information through the player ID
func (g *GridAOIManager) getPlayerByPID(pid uint32) *player.Player {
	return g.players[pid]
}
