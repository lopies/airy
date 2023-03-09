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

type GrID struct {
	GID       int             // Grid ID
	MinX      int             // Grid left boundary coordinates
	MaxX      int             // Grid right boundary coordinates
	MinY      int             // Grid up boundary coordinates
	MaxY      int             // Grid down boundary coordinates
	playerIDs map[uint32]bool // Player id collection
}

// NewGrID initialize
func NewGrID(gID, minX, maxX, minY, maxY int) *GrID {
	return &GrID{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[uint32]bool),
	}
}

// add Adds a player to the grid
func (g *GrID) add(pid uint32) {
	g.playerIDs[pid] = true
}

// remove delete a player from the grid
func (g *GrID) remove(pid uint32) {
	delete(g.playerIDs, pid)
}

// getPlayerIDs gets all the players in the current grid
func (g *GrID) getPlayerIDs() (pids []uint32) {
	for k, _ := range g.playerIDs {
		pids = append(pids, k)
	}
	return
}
