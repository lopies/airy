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

package agent

import (
	"github.com/airy/constants"
	"github.com/airy/interfaces"
	"github.com/airy/logger"
	"github.com/airy/pb"
	"github.com/airy/util"
	"sync"
	"time"
)

var (
	agentByPID        sync.Map
	closeCallBack     func(*Agent) error
	handshakeCallBack func(*Agent) error
	heartbeatCallBack func(*Agent) error
)

type Status int8

const (
	_       Status = iota
	Pending        // Pending status
	Working        // Working status
	Closed         // Closed status
)

// Agent for storing player information
type Agent struct {
	pid           uint32
	spaceID       string
	uniqueID      string
	conn          interfaces.PlayerConn
	lastAt        int64
	status        Status
	serializer    interfaces.Serializer
	codec         interfaces.PacketCodec
	logicServerID string
	closeChan     chan struct{}
}

// NewAgent initialize player information
func NewAgent(conn interfaces.PlayerConn, serialize interfaces.Serializer, codec interfaces.PacketCodec) *Agent {
	a := &Agent{
		uniqueID:   util.UUID(),
		conn:       conn,
		lastAt:     time.Now().Unix(),
		status:     Pending,
		serializer: serialize,
		codec:      codec,
		closeChan:  make(chan struct{}),
	}
	return a
}

func BindCallBack(closeCb, handshakeCb, heartbeatCb func(agent *Agent) error) {
	closeCallBack = closeCb
	handshakeCallBack = handshakeCb
	heartbeatCallBack = heartbeatCb
}

// Scan injected function is executed while server close
func Scan(f func(key, value any) bool) {
	agentByPID.Range(f)
}

// Handshake set agent pid and return true if success
func (a *Agent) Handshake(pid uint32, spaceId string) bool {
	a.pid = pid
	a.spaceID = spaceId
	err := handshakeCallBack(a)
	if err != nil {
		logger.Errorf(err.Error())
		a.pid = 0
		a.spaceID = ""
		return false
	}
	agentByPID.Store(pid, a)
	return true
}

// PID return agent pid only if the handshake is successful
func (a *Agent) PID() uint32 {
	return a.pid
}

// UniqueID return player unique id
func (a *Agent) UniqueID() string {
	return a.uniqueID
}

// LastAt return player last execution time
func (a *Agent) LastAt() *int64 {
	return &a.lastAt
}

// SetLastAt set player last execution time
func (a *Agent) SetLastAt(lastAt int64) {
	a.lastAt = lastAt
	_ = heartbeatCallBack(a)
}

// SetWorkingStatus status (working)
func (a *Agent) SetWorkingStatus() error {
	if a.status >= Working {
		return constants.ErrInvalidStatus
	}
	a.status = Working
	return nil
}

// SetClosedStatus status (closed)
func (a *Agent) SetClosedStatus() error {
	a.status = Closed
	return nil
}

// SetPendingStatus status (pending)
func (a *Agent) SetPendingStatus() error {
	a.status = Pending
	a.spaceID = ""
	a.logicServerID = ""
	return nil
}

// Status return player status
func (a *Agent) Status() Status {
	return a.status
}

// SetSpaceID set space to which the player belongs
func (a *Agent) SetSpaceID(spaceID string) {
	a.spaceID = spaceID
}

// SpaceID return to which the player belongs
func (a *Agent) SpaceID() string {
	return a.spaceID
}

// SetLogicServerID set logic server to which the player belongs
func (a *Agent) SetLogicServerID(serverID string) error {
	logger.Debugf("agent set logic server id = %s,pid = %d,uniqueID = %s", serverID, a.PID(), a.UniqueID())
	if serverID == "" {
		return constants.ErrNotFoundServer
	}
	a.logicServerID = serverID
	return nil
}

// LogicServerID return logic server to which the player belongs
func (a *Agent) LogicServerID() string {
	return a.logicServerID
}

// GetAgentByPID get agent by pid
func GetAgentByPID(pid uint32) *Agent {
	au, ok := agentByPID.Load(pid)
	if !ok {
		return nil
	}
	return au.(*Agent)
}

// Shutdown agent end
func (a *Agent) Shutdown() {
	defer func() {
		_ = recover()
	}()
	_ = closeCallBack(a)
	agentByPID.Delete(a.pid)
	_ = a.SetClosedStatus()
	close(a.closeChan)
	a.conn.Close()
}

// Write send packet to client
func (a *Agent) Write(packets ...*pb.Packet) {
	if a.status == Working {
		var msg []byte
		for _, packet := range packets {
			logger.Debugf("receiver[%d], push packet to client :%+v", a.PID(), packet)
			byt, _ := a.serializer.Marshal(packet)
			bytEnc, _ := a.codec.Encode(byt)
			msg = append(msg, bytEnc...)
		}
		if err := a.Conn().WriteTo(msg); err != nil {
			logger.Errorf("Failed to write in conn: err = %s", err.Error())
			return
		}
	}
}

// ReadCloseChan return agent close channel
func (a *Agent) ReadCloseChan() chan struct{} {
	return a.closeChan
}

// Conn return connection
func (a *Agent) Conn() interfaces.PlayerConn {
	return a.conn
}
