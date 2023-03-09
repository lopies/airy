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

package player

type Player struct {
	PID           uint32
	X             float32
	Y             float32
	Z             float32
	V             float32 // Rotate 0-360 degrees
	MarkVal       int
	GateServerID  string
	Name          string
	Neighbors     map[*Player]struct{}
	MarkNeighbors []*Player
	XPrev, XNext  *Player
	ZPrev, ZNext  *Player
}

// NewPlayer initializes a player
func NewPlayer(pid uint32, gateServerID string) *Player {
	u := &Player{
		PID:           pid,
		GateServerID:  gateServerID,
		Neighbors:     map[*Player]struct{}{},
		MarkNeighbors: make([]*Player, 0, 30),
	}
	return u
}

// SetName set a player name
func (p *Player) SetName(name string) *Player {
	p.Name = name
	return p
}

// SetGateServerID set gate server where the player resides
func (p *Player) SetGateServerID(gateServerID string) *Player {
	p.GateServerID = gateServerID
	return p
}
