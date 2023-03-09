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

package module

import (
	"context"
	"errors"
	"fmt"
	"github.com/airy/agent"
	"github.com/airy/component"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/logger"
	"github.com/airy/pb"
	"github.com/airy/util"
	"github.com/golang/protobuf/ptypes/empty"
	"runtime"
	"strconv"
	"sync/atomic"
	"time"
)

type Gate struct {
	BaseModule
	addHeartBeatChan chan *agent.Agent
	delHeartBeatChan chan *agent.Agent
	combinationsChan chan *pb.Combinations
}

func NewGate() *Gate {
	g := new(Gate)
	g.SetType(constants.GateModule)
	return g
}

func (g *Gate) Init(config *config.AiryConfig) {
	g.BaseModule.Init(config)
	g.addHeartBeatChan = make(chan *agent.Agent, 1<<14)
	g.delHeartBeatChan = make(chan *agent.Agent, 1<<14)
	g.combinationsChan = make(chan *pb.Combinations, 1<<10)
}

func (g *Gate) AfterInit() {
	agent.BindCallBack(g.agentCloseCallBack, g.agentHandshakeCallBack, g.agentHeartBeatCallBack)
	t := time.NewTicker(time.Second * 5)
	go func() {
		for {
			select {
			case <-t.C:
				logger.Infof("发送的消息通道大小为:%d,\n", len(g.combinationsChan))
			}
		}
	}()
}

func (g *Gate) heartbeatCheck() {
	for {
		select {
		case a, ok := <-g.addHeartBeatChan:
			if !ok {
				return
			}
			g.Timer().AddCountTimer(a.UniqueID(), g.HeartBeat(), -1, func(...any) {
				deadline := time.Now().Add(-2 * g.HeartBeat()).Unix()
				if atomic.LoadInt64(a.LastAt()) <= deadline {
					logger.Debugf("Session heartbeat timeout, LastTime=%d, Deadline=%d", atomic.LoadInt64(a.LastAt()), deadline)
					g.delHeartBeatChan <- a
				}
			})
			logger.Debugf("adding the heartbeat detection succeeded,pid = %d,uniqueID = %s", a.PID(), a.UniqueID())
		case a, ok := <-g.delHeartBeatChan:
			if !ok {
				return
			}
			a.Shutdown()
			g.Timer().RemoveTimer(a.UniqueID())
			logger.Debugf("removing the heartbeat detection succeeded,pid = %d,uniqueID = %s", a.PID(), a.UniqueID())
		}
	}
}

func (g *Gate) Run(stop chan struct{}) {
	// push to client
	go g.write()
	tcpAcceptor := g.AcceptorComponent()
	if tcpAcceptor != nil {
		go g.heartbeatCheck()
		go func() {
			for {
				select {
				case conn := <-tcpAcceptor.GetConnChan():
					// create a client agent
					logger.Debugf("receive client connection")
					a := agent.NewAgent(conn, g.BaseModule.Serializer(), g.codec)
					g.addHeartBeatChan <- a
					go g.read(a)
				case <-stop:
					return
				}
			}
		}()
	}
}

func (g *Gate) AddComponents(comps ...component.Component) {
	g.BaseModule.AddComponents(comps...)
}

func (g *Gate) BeforeShutdown() {
	close(g.addHeartBeatChan)
	close(g.delHeartBeatChan)
	g.clearPlayer()
}

func (g *Gate) PushToUsers(context context.Context, req *pb.Combinations) (*empty.Empty, error) {
	g.combinationsChan <- req
	return &empty.Empty{}, nil
}

func (g *Gate) parse(a *agent.Agent, dm []byte) error {
	p := new(pb.Packet)
	err := g.Serializer().Unmarshal(dm, p)
	if err != nil {
		logger.Errorf("Failed to Unmarshal in conn: %s", err.Error())
		return err
	}
	switch p.Type {
	case pb.Type_Handshake_:
		var err error
		logger.Debugf("receive client handshake packet,pid = %d,uniqueID = %s", a.PID(), a.UniqueID())
		if a.Status() > agent.Pending {
			return errors.New("agent status error")
		}
		handshake := new(pb.Handshake)
		err = g.Serializer().Unmarshal(p.Data, handshake)
		if err != nil {
			logger.Errorf("Failed to Unmarshal in handshake: %s", err.Error())
			return err
		}
		logger.Debugf("receive handshake msg,token = %s,spaceID = %s,pid =%d,uniqueID = %s", handshake.Token, handshake.SpaceID, handshake.PID, a.UniqueID())
		if g.Authorization() {
			if !util.ValidToken(handshake.Token) {
				logger.Errorf("Failed to valid token: %s", err.Error())
				return constants.ErrInvalidToken
			}
		}
		if !a.Handshake(handshake.PID, handshake.SpaceID) {
			a.Write(constants.NewResponse(constants.ErrHandshake, p.RequestID, pb.Type_HandshakeAck_))
			return errors.New(fmt.Sprintf("the player[%d] is online", handshake.PID))
		}
		a.Write(constants.NewResponse(constants.Success, p.RequestID, pb.Type_HandshakeAck_))
	case pb.Type_Heartbeat_:
		a.SetLastAt(time.Now().Unix())
		logger.Debugf("receive client heartbeat packet,pid = %d,uniqueID = %s", a.PID(), a.UniqueID())
		return nil
	default:
		if a.Status() < agent.Working {
			logger.Errorf("invalid agent status,pid = %d,uniqueID = %s", a.PID(), a.UniqueID())
			return constants.ErrInvalidStatus
		}
		logger.Debugf("receive client request packet,pid = %d,requestCode = %d,uniqueID = %s", a.PID(), p.RequestCode, a.UniqueID())
		if p.Type == pb.Type_Join {
			p.SpaceID = a.SpaceID()
			p.SourceServerID = g.Server().ID
		}
		p.PID = a.PID()
		err = g.GRPCClientComponent().Request(context.Background(), a.LogicServerID(), p)
		if err != nil {
			logger.Errorf("send grpc request error : %s", err.Error())
		}
	}
	return nil
}

func (g *Gate) write() {
	for {
		select {
		case msg, ok := <-g.combinationsChan:
			if !ok {
				return
			}
			for _, combination := range msg.Combinations {
				if len(combination.PIDs) > 0 {
					packet := combination.Packets[0]
					for _, pid := range combination.PIDs {
						a := agent.GetAgentByPID(pid)
						if a != nil {
							a.Write(packet)
						}
					}
				} else {
					a := agent.GetAgentByPID(combination.PIDs[0])
					packets := combination.Packets
					a.Write(packets...)
				}
			}
		}
	}
}

func (g *Gate) read(a *agent.Agent) {
	defer func() {
		err := recover()
		if err != nil {
			buf := make([]byte, 1<<10)
			runtime.Stack(buf, true)
			logger.Errorf("gate loop func occur  a serious error occurs,err : %s", string(buf))
		}
		g.delHeartBeatChan <- a
	}()

	codec := g.Codec()
	for {
		var err error
		dm, err := a.Conn().ReadFrom(codec)
		if err != nil {
			logger.Errorf("Failed to recv in conn: %s", err.Error())
			return
		}
		err = g.parse(a, dm)
		if err != nil {
			logger.Errorf("Failed to parse message: err = %s", err.Error())
			return
		}
	}
}

func (g *Gate) clearPlayer() {
	f := func(key, value any) bool {
		_ = g.SysStorage().DelPlayerIDGateID(strconv.Itoa(int(key.(uint32))))
		return true
	}
	agent.Scan(f)
}

func (g *Gate) agentHandshakeCallBack(a *agent.Agent) error {
	gateID := g.SysStorage().GetGateIDByPlayerID(strconv.Itoa(int(a.PID())))
	if gateID != "" {
		return constants.ErrAgentHandshake
	}
	var serverID string
	var err error
	serverID, err = g.SysStorage().GetLogicIDBySpaceID(a.SpaceID())
	if err != nil {
		serverID = g.selectLogicServerID()
		logger.Debugf("select logic server id : %s", serverID)
		_ = g.SysStorage().StoreSpaceIDLogicID(a.SpaceID(), serverID)
		_ = g.SysStorage().StoreLogicIDSpaceID(serverID, a.SpaceID())
	}
	_ = a.SetLogicServerID(serverID)
	_ = a.SetWorkingStatus()
	_ = g.SysStorage().StorePlayerIDGateID(strconv.Itoa(int(a.PID())), g.Server().ID)
	return nil
}

func (g *Gate) agentCloseCallBack(a *agent.Agent) error {
	if a.Status() == agent.Working {
		_ = g.GRPCClientComponent().Request(context.Background(), a.LogicServerID(), &pb.Packet{
			Type:        pb.Type_Leave,
			PID:         a.PID(),
			RequestCode: constants.LeaveRequest,
		})
	}
	_ = g.SysStorage().DelPlayerIDGateID(strconv.Itoa(int(a.PID())))
	return nil
}

func (g *Gate) agentHeartBeatCallBack(a *agent.Agent) error {
	if a.Status() == agent.Working {
		_ = g.GRPCClientComponent().Request(context.Background(), a.LogicServerID(), &pb.Packet{
			Type:        pb.Type_Heartbeat_,
			RequestCode: constants.JoinRequest,
			PID:         a.PID(),
		})
	}
	return nil
}

func (g *Gate) selectLogicServerID() string {
	servers, err := g.ETCDComponent().GetServersByType(string(constants.LogicModule))
	if err != nil {
		logger.Errorf("no modules available of this type")
		return ""
	}
	//TODO : Load balancing strategy
	for _, server := range servers {
		return server.ID
	}
	return ""
}
