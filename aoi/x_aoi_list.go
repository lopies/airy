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

type xAOIList struct {
	aoidist float32
	head    *player.Player
	tail    *player.Player
}

// newXAOIList initializes X aoi list
func newXAOIList(aoidist float32) *xAOIList {
	return &xAOIList{aoidist: aoidist}
}

// Insert add a player node
func (sl *xAOIList) Insert(aoi *player.Player) {
	insertCoord := aoi.X
	if sl.head != nil {
		p := sl.head
		for p != nil && p.X < insertCoord {
			p = p.XNext
		}
		// now, p == nil or p.coord >= insertCoord
		if p == nil { // if p == nil, insert xzaoi at the end of list
			tail := sl.tail
			tail.XNext = aoi
			aoi.XPrev = tail
			sl.tail = aoi
		} else { // otherwise, p >= xzaoi, insert xzaoi before p
			prev := p.XPrev
			aoi.XNext = p
			p.XPrev = aoi
			aoi.XPrev = prev

			if prev != nil {
				prev.XNext = aoi
			} else { // p is the head, so xzaoi should be the new head
				sl.head = aoi
			}
		}
	} else {
		sl.head = aoi
		sl.tail = aoi
	}
}

// Remove delete a player node
func (sl *xAOIList) Remove(p *player.Player) {
	prev := p.XPrev
	next := p.XNext
	if prev != nil {
		prev.XNext = next
		p.XPrev = nil
	} else {
		sl.head = next
	}
	if next != nil {
		next.XPrev = prev
		p.XNext = nil
	} else {
		sl.tail = prev
	}
}

// Move a player node x coordinate
func (sl *xAOIList) Move(p *player.Player, oldCoord float32) {
	coord := p.X
	if coord > oldCoord {
		// moving to next ...
		next := p.XNext
		if next == nil || next.X >= coord {
			// no need to adjust in list
			return
		}
		prev := p.XPrev
		//fmt.Println(1, prev, next, prev == nil || prev.XNext == xzaoi)
		if prev != nil {
			prev.XNext = next // remove xzaoi from list
		} else {
			sl.head = next // xzaoi is the head, trim it
		}
		next.XPrev = prev

		//fmt.Println(2, prev, next, prev == nil || prev.XNext == next)
		prev, next = next, next.XNext
		for next != nil && next.X < coord {
			prev, next = next, next.XNext
			//fmt.Println(2, prev, next, prev == nil || prev.XNext == next)
		}
		//fmt.Println(3, prev, next)
		// no we have prev.X < coord && (next == nil || next.X >= coord), so insert between prev and next
		prev.XNext = p
		p.XPrev = prev
		if next != nil {
			next.XPrev = p
		} else {
			sl.tail = p
		}
		p.XNext = next

		//fmt.Println(4)
	} else {
		// moving to prev ...
		prev := p.XPrev
		if prev == nil || prev.X <= coord {
			// no need to adjust in list
			return
		}

		next := p.XNext
		if next != nil {
			next.XPrev = prev
		} else {
			sl.tail = prev // xzaoi is the head, trim it
		}
		prev.XNext = next // remove xzaoi from list

		next, prev = prev, prev.XPrev
		for prev != nil && prev.X > coord {
			next, prev = prev, prev.XPrev
		}
		// no we have next.X > coord && (prev == nil || prev.X <= coord), so insert between prev and next
		next.XPrev = p
		p.XNext = next
		if prev != nil {
			prev.XNext = p
		} else {
			sl.head = p
		}
		p.XPrev = prev
	}
}

// markNeighbors update neighbors
func (sl *xAOIList) markNeighbors(p *player.Player) {
	prev := p.XPrev
	minXCoord := p.X - sl.aoidist
	minZCoord := p.Z - sl.aoidist
	maxZCoord := p.Z + sl.aoidist
	for prev != nil && prev.X >= minXCoord {
		if section(prev.Z, minZCoord, maxZCoord) {
			prev.MarkVal += 1
			p.Neighbors[prev] = struct{}{}
			prev.Neighbors[p] = struct{}{}

			p.MarkNeighbors = append(p.MarkNeighbors, prev)
		}
		prev = prev.XPrev
	}

	next := p.XNext
	maxXCoord := p.X + sl.aoidist
	for next != nil && next.X <= maxXCoord {
		if section(next.Z, minZCoord, maxZCoord) {
			next.MarkVal += 1
			p.Neighbors[next] = struct{}{}
			next.Neighbors[p] = struct{}{}

			p.MarkNeighbors = append(p.MarkNeighbors, next)
		}
		next = next.XNext
	}
}

// clearMark clear mark
func (sl *xAOIList) clearMark(p *player.Player) {
	for _, player := range p.MarkNeighbors {
		player.MarkVal = 0
	}
	p.MarkNeighbors = p.MarkNeighbors[0:0]
}

// section if value Coordinate interval[min,max],return true
func section(value, min, max float32) bool {
	if min > value {
		return false
	}
	if max < value {
		return false
	}
	return true
}
