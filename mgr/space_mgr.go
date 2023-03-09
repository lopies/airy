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

package mgr

import (
	"github.com/airy/constants"
	"github.com/airy/logger"
	"github.com/airy/space"
	"sync"
)

type SpaceMgr struct {
	Spaces sync.Map
}

func NewSpaceMgr() *SpaceMgr {
	return &SpaceMgr{}
}

func (sm *SpaceMgr) AddSpace(s *space.Space) error{
	_,load := sm.Spaces.LoadOrStore(s.ID(), s)
	if load {
		logger.Errorf("add space = %s to manager error",s.ID())
		return constants.ErrCreateSpace
	}
	logger.Infof("add space to spaceMgr,spaceID = %s", s.ID())
	return nil
}

func (sm *SpaceMgr) DelSpace(s *space.Space) {
	sm.Spaces.Delete(s.ID())
	logger.Infof("remove space[%s] from spaceMgr",s.ID())
}

func (sm *SpaceMgr) SpaceByID(id string) *space.Space {
	v, ok := sm.Spaces.Load(id)
	if ok {
		return v.(*space.Space)
	}
	return nil
}

func (sm *SpaceMgr) Clear() {
	sm.Spaces.Range(func(key, value any) bool {
		sm.Spaces.Delete(key)
		return true
	})
}
