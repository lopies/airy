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

package interfaces

import "github.com/airy/player"

type AOIManager interface {
	// Enter is called when Entity enters Space
	Enter(p *player.Player, x, z float32) bool
	// Leave is called when Entity leaves Space
	Leave(p *player.Player)
	// Move is called when Entity moves in Space
	Move(p *player.Player, x, z float32)
	// Neighbors return players in a nearby field of view
	Neighbors(*player.Player) map[*player.Player]struct{}
	// AllPlayers return all players
	AllPlayers() []*player.Player
	// Num return aoi people number
	Num() int
	// String print aoi info
	String()
}
