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

package client

import (
	"errors"
	"fmt"
	"github.com/airy/client/logger"
	"github.com/airy/constants"
	"github.com/airy/interfaces"
	"github.com/airy/pb"
	"github.com/airy/push"
	"github.com/airy/util"
	"math"
	"math/rand"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type pendingRequest struct {
	packet *pb.Packet
	sentAt time.Time
}

// Client struct
type Client struct {
	conn             net.Conn
	pid              uint32
	nextID           uint32
	heatBeatInterval int
	Connected        bool
	Serializer       interfaces.Serializer
	Codec            interfaces.PacketCodec
	Compress         interfaces.Compress
	IncomingMsgChan  chan *pb.Packet
	SendMsgChan      chan *pb.Packet
	pendingChan      chan bool
	pendingRequests  map[int32]*pendingRequest
	pendingReqMutex  sync.Mutex
	requestTimeout   time.Duration
	HandshakeData    *pb.Packet
	HeartBeatData    *pb.Packet
	closeChan        chan struct{}
	player           *Player
}

// MsgChannel return the incoming message channel
func (c *Client) MsgChannel() chan *pb.Packet {
	return c.IncomingMsgChan
}

// ConnectedStatus return the connection status
func (c *Client) ConnectedStatus() bool {
	return c.Connected
}

// NewClient returns a new client
func NewClient(heartBeatInterval int, packetCodec interfaces.PacketCodec, compress interfaces.Compress, serializer interfaces.Serializer, requestTimeout ...time.Duration) *Client {
	reqTimeout := 30 * time.Second
	if len(requestTimeout) > 0 {
		reqTimeout = requestTimeout[0]
	}
	return &Client{
		heatBeatInterval: heartBeatInterval,
		Connected:        false,
		Codec:            packetCodec,
		Compress:         compress,
		Serializer:       serializer,
		IncomingMsgChan:  make(chan *pb.Packet, 100),
		SendMsgChan:      make(chan *pb.Packet, 100),
		pendingChan:      make(chan bool, 30),
		pendingRequests:  make(map[int32]*pendingRequest),
		requestTimeout:   reqTimeout,
		HandshakeData:    new(pb.Packet),
		HeartBeatData: &pb.Packet{
			Type: pb.Type_Heartbeat_,
		},
	}
}

func (c *Client) Bind(p *Player, spaceID string) error {
	mapData, err := util.ReadMapData(spaceID)
	if err != nil {
		return err
	}
	r := rand.Intn(len(mapData.Birth))
	p.X, p.Z = mapData.Birth[r].X, mapData.Birth[r].Z
	p.SpaceID = spaceID
	handshake := &pb.Handshake{
		PID:     p.PID,
		SpaceID: spaceID,
	}
	byt, _ := c.Serializer.Marshal(handshake)
	c.HandshakeData = &pb.Packet{
		Type:    pb.Type_Handshake_,
		PID:     c.pid,
		SpaceID: spaceID,
		Data:    byt,
	}
	c.player = p
	c.pid = p.PID
	return nil
}

func (c *Client) handleHandshakeResponse() error {
	byt, err := c.Codec.Decode(c.conn)
	if err != nil {
		logger.Log.Errorf("handleHandshakeResponse codec decode error : %s", err.Error())
		return err
	}
	if c.Compress != nil {
		byt, _ = c.Compress.Decode(byt)
	}
	p := new(pb.Packet)
	err = c.Serializer.Unmarshal(byt, p)
	if err != nil {
		logger.Log.Errorf("handleHandshakeResponse serializer error")
		return err
	}
	if p.Type == pb.Type_HandshakeAck_ && p.ErrorCode == 0 {
		logger.Log.Infof("receive server handshake ack,client connected status is true")
		c.Connected = true
		// send heartbeat
		go c.sendHeartbeats(c.heatBeatInterval)
		// send request msg
		go c.handleSendMsg()
		// read push msg
		go c.handleReadMsg()
		// hand incoming msg
		go c.handleInComingMsg()
		// hand request
		go c.pendingRequestsReaper()
		return nil
	}
	return errors.New("handle handshake response error")
}

// pendingRequestsReaper delete timedout requests
func (c *Client) pendingRequestsReaper() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			toDelete := make([]*pendingRequest, 0, 5)
			c.pendingReqMutex.Lock()
			for _, v := range c.pendingRequests {
				if time.Now().Sub(v.sentAt) > c.requestTimeout {
					toDelete = append(toDelete, v)
				}
			}
			for _, pendingReq := range toDelete {
				logger.Log.Debugf("timeout ,delete request message id = %d", pendingReq.packet.RequestID)
				p := &pb.Packet{
					Type:        pb.Type_Response_,
					RequestID:   pendingReq.packet.RequestID,
					RequestCode: pendingReq.packet.RequestCode,
					ErrorCode:   500,
				}
				delete(c.pendingRequests, pendingReq.packet.RequestID)
				<-c.pendingChan
				c.IncomingMsgChan <- p
			}
			c.pendingReqMutex.Unlock()
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) sendHeartbeats(interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Second)
	defer func() {
		t.Stop()
		c.Disconnect()
	}()
	for {
		select {
		case <-t.C:
			c.SendMsgChan <- c.HeartBeatData
		case <-c.closeChan:
			return
		}
	}
}

// ConnectToTCP  connects to the server at addr, for now the only supported protocol is tcp
func (c *Client) ConnectToTCP(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("client start err, exit!", err)
		return err
	}
	c.conn = conn
	if err = c.handleHandshake(); err != nil {
		return err
	}
	c.closeChan = make(chan struct{})
	return nil
}

func (c *Client) handleHandshake() error {
	c.sendHandshakeRequest()
	if err := c.handleHandshakeResponse(); err != nil {
		return err
	}
	return nil
}

// SendRequest sends a request to the server
func (c *Client) SendRequest(packetType pb.Type, requestCode int32, data []byte, ack bool) {
	p := &pb.Packet{
		Type:        packetType,
		RequestCode: requestCode,
		PID:         c.pid,
		SpaceID:     c.player.SpaceID,
		Data:        data,
	}
	if ack {
		msgID := atomic.AddUint32(&c.nextID, 1)
		if msgID == math.MaxInt32 {
			c.nextID = 0
			msgID = 0
		}
		p.RequestID = int32(msgID)
		c.pendingChan <- true
		c.pendingReqMutex.Lock()
		if _, ok := c.pendingRequests[p.RequestID]; !ok {
			c.pendingRequests[p.RequestID] = &pendingRequest{
				packet: p,
				sentAt: time.Now(),
			}
			logger.Log.Debugf("send request,id = %d", p.RequestID)
		}
		c.pendingReqMutex.Unlock()
	}
	c.SendMsgChan <- p
}

// Disconnect disconnects the client
func (c *Client) Disconnect() {
	if c.Connected {
		c.Connected = false
		close(c.closeChan)
		c.conn.Close()
	}
}

func (c *Client) write(p *pb.Packet) error {
	byt, _ := c.Serializer.Marshal(p)
	if c.Compress != nil {
		byt = c.Compress.Encode(byt)
	}
	byt, _ = c.Codec.Encode(byt)
	_, err := c.conn.Write(byt)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) read() error {
	byt, err := c.Codec.Decode(c.conn)
	if err != nil {
		logger.Log.Errorf("codec decode error : %s", err.Error())
		return err
	}
	if c.Compress != nil {
		byt, _ = c.Compress.Decode(byt)
	}
	p := new(pb.Packet)
	err = c.Serializer.Unmarshal(byt, p)
	if err != nil {
		logger.Log.Errorf("serializer error")
		return err
	}
	c.IncomingMsgChan <- p
	return nil
}

func (c *Client) sendHandshakeRequest() {
	err := c.write(c.HandshakeData)
	if err != nil {
		logger.Log.Errorf("conn write data error : %s", err.Error())
		return
	}
}

func (c *Client) handleReadMsg() {
	for {
		err := c.read()
		if err != nil {
			logger.Log.Errorf("read message error : %+v", err)
			return
		}
	}
}

func (c *Client) handleSendMsg() {
	for {
		select {
		case msg := <-c.SendMsgChan:
			err := c.write(msg)
			if err != nil {
				logger.Log.Errorf("conn write data error : %s", err.Error())
				return
			}
		case <-c.closeChan:
			return
		}
	}
}

func (c *Client) handleInComingMsg() {
	for {
		select {
		case p := <-c.IncomingMsgChan:
			switch p.Type {
			case pb.Type_Push_:
				logger.Log.Debug("client receive a push packet")
				c.handPushMsg(p)
			case pb.Type_Kick_:
				logger.Log.Debug("client receive a kick packet")
				c.Disconnect()
			case pb.Type_Response_:
				logger.Log.Debug("client receive a response packet")
				c.pendingReqMutex.Lock()
				if _, ok := c.pendingRequests[p.RequestID]; ok {
					delete(c.pendingRequests, p.RequestID)
					<-c.pendingChan
					logger.Log.Debugf("receive server response,packet = %+v", p)
				} else {
					logger.Log.Errorf("not receive server response,packet = %+v", p)
					c.pendingReqMutex.Unlock()
					continue // TODO:do not process msg for already timedout request
				}
				c.pendingReqMutex.Unlock()

				if p.RequestCode == constants.JoinRequest && p.ErrorCode == 0 {
					c.player.Connected = true
					logger.Log.Infof("player status is connected")
				}

			}
		case <-c.closeChan:
			return
		}

	}
}

func (c *Client) handPushMsg(p *pb.Packet) {
	if p.ErrorCode != 0 {
		logger.Log.Errorf("error code = %d", p.ErrorCode)
		return
	}
	switch p.RouteCode {
	case constants.OnHandshake:
		logger.Log.Debugf("receive server OnHandshake push data")
	case constants.OnComeIntoView:
		data := new(push.OnComeIntoView)
		err := c.Serializer.Unmarshal(p.Data, data)
		if err != nil {
			logger.Log.Errorf("serializer unmarshal decode err : %s", err.Error())
		}
		logger.Log.Debugf("one player come into my view,pid = %d,name = %s,x = %f,z = %f", data.PID, data.Name, data.X, data.Z)
	case constants.OnOutOfView:
		data := new(push.OnOutOfView)
		err := c.Serializer.Unmarshal(p.Data, data)
		if err != nil {
			logger.Log.Errorf("serializer unmarshal decode err : %s", err.Error())
		}
		logger.Log.Debugf("one player out of my view,pid = %d", data.PID)
	case constants.OnMove:
		data := new(push.OnMove)
		err := c.Serializer.Unmarshal(p.Data, data)
		if err != nil {
			logger.Log.Errorf("serializer unmarshal decode err : %s", err.Error())
		}
		logger.Log.Debugf("one player move,pid = %d,x = %f,z = %f", data.PID, data.X, data.Z)
	case constants.OnLeave:
		data := new(push.OnExit)
		err := c.Serializer.Unmarshal(p.Data, data)
		if err != nil {
			logger.Log.Errorf("serializer unmarshal decode err : %s", err.Error())
		}
		logger.Log.Debugf("one player exit,pid = %d", data.PID)
	case constants.OnDelay:
		data := new(push.OnDelay)
		err := c.Serializer.Unmarshal(p.Data, data)
		if err != nil {
			logger.Log.Errorf("serializer unmarshal decode err : %s", err.Error())
		}
		logger.Log.Infof("delay : %d ms,pid = %d", time.Now().UnixMilli()-data.StartTime, c.pid)
		if c.pid == 80017 || c.pid == 800117 {
			logger.Log.Errorf("delay : %d ms", time.Now().UnixMilli()-data.StartTime)
		}
	case constants.OnGlobalChat:
		data := new(push.OnCommonChat)
		err := c.Serializer.Unmarshal(p.Data, data)
		if err != nil {
			logger.Log.Errorf("serializer unmarshal decode err : %s", err.Error())
		}
		logger.Log.Infof("one player[%d] send a global message,content = %s", data.Pid, data.Content)
	}
}
