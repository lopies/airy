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

package timer

import (
	"container/list"
	"github.com/airy/constants"
	"time"
)

type TimeWheel struct {
	interval          time.Duration  // How often does the pointer move forward one space
	ticker            *time.Ticker   // system ticker
	slots             []*list.List   // Time wheel slots
	timer             map[string]int // key: unique identifier of the timer value: slot where the timer resides. The slot is used to delete the timer. Concurrent read/write operations are not performed
	currentPos        int            // Which slot the current pointer points to
	slotNum           int            // Number of slots
	addTaskChannel    chan *Task     // Adds a task channel
	removeTaskChannel chan string    // Deletes the task channel
	stopChannel       chan struct{}  // Stops the timer channel
}

// NewTimeWheel Create a new timeWheel
func NewTimeWheel(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:          interval,
		slots:             make([]*list.List, slotNum),
		timer:             make(map[string]int),
		currentPos:        0,
		slotNum:           slotNum,
		addTaskChannel:    make(chan *Task),
		removeTaskChannel: make(chan string),
		stopChannel:       make(chan struct{}),
	}
	tw.initSlots()
	return tw
}

// Initializes slots, each pointing to a bidirectional linked list
func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

// Start time wheel
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
}

// Stop time wheel
func (tw *TimeWheel) Stop() {
	tw.stopChannel <- struct{}{}
}

// AddAfterTimer Add a delay execution task after interval time
// key: task unique tag
// interval: delay time
// f: task func
// args: task parameter
func (tw *TimeWheel) AddAfterTimer(key string, interval time.Duration, f func(...any), args ...any) error {
	if interval < 0 {
		return constants.ErrTimeInterval
	}
	task := newTask(key, interval, 1, f, args)
	tw.addTaskChannel <- task
	return nil
}

// AddCountTimer Add a delay execution task
// key: task unique tag
// interval: delay time
// count: exec times,if count = -1,indicates that there is no limit on the number of times
// f: task func
// args: task parameter
func (tw *TimeWheel) AddCountTimer(key string, interval time.Duration, count int, f func(...any), args ...any) error {
	if interval < 0 {
		return constants.ErrTimeInterval
	}
	task := newTask(key, interval, count, f, args)
	tw.addTaskChannel <- task
	return nil
}

// RemoveTimer the key of deleting a timer is the unique identifier of the timer passed when the timer is added
func (tw *TimeWheel) RemoveTimer(key string) {
	if key == "" {
		return
	}
	tw.removeTaskChannel <- key
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChannel:
			tw.addTask(task)
		case key := <-tw.removeTaskChannel:
			tw.removeTask(key)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAndRunTask(l)
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

// scanAndRunTask scan the expired timer in the linked list and execute the callback function
func (tw *TimeWheel) scanAndRunTask(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}
		go task.job(task.args)
		task.counter--
		next := e.Next()
		l.Remove(e)
		if task.key != "" {
			delete(tw.timer, task.key)
		}
		if task.counter != 0 {
			tw.addTask(task)
		}
		e = next
	}
}

// addTask Adds a task to the linked list
func (tw *TimeWheel) addTask(task *Task) {
	pos, circle := tw.getPositionAndCircle(task.interval)
	task.circle = circle
	if task.key != "" {
		if _, ok := tw.timer[task.key]; ok {
			return
		}
		tw.timer[task.key] = pos
	}
	tw.slots[pos].PushBack(task)

}

// getPositionAndCircle Gets the position of the timer in the slot and the number of turns the time wheel needs to turn
func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	forward := delaySeconds / intervalSeconds
	circle = forward / tw.slotNum
	pos = (tw.currentPos + forward) % tw.slotNum

	return
}

// removeTask Deletes a task from a linked list
func (tw *TimeWheel) removeTask(key string) {
	// gets the slot where the timer is located
	position, ok := tw.timer[key]
	if !ok {
		return
	}
	// gets the linked list to which the slot points
	l := tw.slots[position]
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.key == key {
			delete(tw.timer, task.key)
			l.Remove(e)
		}
		e = e.Next()
	}
}
