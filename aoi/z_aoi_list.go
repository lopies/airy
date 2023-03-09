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

import "github.com/airy/player"

type zAOIList struct {
	aoidist float32
	head    *player.Player
	tail    *player.Player
}

// newZAOIList initializes Z aoi list
func newZAOIList(aoidist float32) *zAOIList {
	return &zAOIList{aoidist: aoidist}
}

// Insert add a player node
func (sl *zAOIList) Insert(p *player.Player) {
	insertCoord := p.Z
	if sl.head != nil {
		p1 := sl.head
		for p1 != nil && p1.Z < insertCoord {
			p1 = p1.ZNext
		}
		// now, p == nil or p.coord >= insertCoord
		if p1 == nil { // if p == nil, insert xzaoi at the end of list
			tail := sl.tail
			tail.ZNext = p
			p.ZPrev = tail
			sl.tail = p
		} else { // otherwise, p >= xzaoi, insert xzaoi before p
			prev := p1.ZPrev
			p.ZNext = p1
			p1.ZPrev = p
			p.ZPrev = prev

			if prev != nil {
				prev.ZNext = p
			} else { // p is the head, so xzaoi should be the new head
				sl.head = p
			}
		}
	} else {
		sl.head = p
		sl.tail = p
	}
}

// Remove delete a player node
func (sl *zAOIList) Remove(p *player.Player) {
	prev := p.ZPrev
	next := p.ZNext
	if prev != nil {
		prev.ZNext = next
		p.ZPrev = nil
	} else {
		sl.head = next
	}
	if next != nil {
		next.ZPrev = prev
		p.ZNext = nil
	} else {
		sl.tail = prev
	}
}

// Move a player node z coordinate
func (sl *zAOIList) Move(p *player.Player, oldCoord float32) {
	coord := p.Z
	if coord > oldCoord {
		// moving to next ...
		next := p.ZNext
		if next == nil || next.Z >= coord {
			// no need to adjust in list
			return
		}
		prev := p.ZPrev
		//fmt.Println(1, prev, next, prev == nil || prev.yNext == xzaoi)
		if prev != nil {
			prev.ZNext = next // remove xzaoi from list
		} else {
			sl.head = next // xzaoi is the head, trim it
		}
		next.ZPrev = prev

		//fmt.Println(2, prev, next, prev == nil || prev.yNext == next)
		prev, next = next, next.ZNext
		for next != nil && next.Z < coord {
			prev, next = next, next.ZNext
			//fmt.Println(2, prev, next, prev == nil || prev.yNext == next)
		}
		//fmt.Println(3, prev, next)
		// no we have prev.X < coord && (next == nil || next.X >= coord), so insert between prev and next
		prev.ZNext = p
		p.ZPrev = prev
		if next != nil {
			next.ZPrev = p
		} else {
			sl.tail = p
		}
		p.ZNext = next

		//fmt.Println(4)
	} else {
		// moving to prev ...
		prev := p.ZPrev
		if prev == nil || prev.Z <= coord {
			// no need to adjust in list
			return
		}

		next := p.ZNext
		if next != nil {
			next.ZPrev = prev
		} else {
			sl.tail = prev // xzaoi is the head, trim it
		}
		prev.ZNext = next // remove xzaoi from list

		next, prev = prev, prev.ZPrev
		for prev != nil && prev.Z > coord {
			next, prev = prev, prev.ZPrev
		}
		// no we have next.X > coord && (prev == nil || prev.X <= coord), so insert between prev and next
		next.ZPrev = p
		p.ZNext = next
		if prev != nil {
			prev.ZNext = p
		} else {
			sl.head = p
		}
		p.ZPrev = prev
	}
}
