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
	"github.com/airy/component"
	"github.com/airy/config"
	"github.com/airy/constants"
	actx "github.com/airy/context"
	"github.com/airy/logger"
	"github.com/airy/mgr"
	"github.com/airy/pb"
	"github.com/airy/session"
	"github.com/airy/space"
	"github.com/golang/protobuf/ptypes/empty"
	"runtime"
	"sync/atomic"
	"time"
)

type Logic struct {
	BaseModule
	spaceMgr              *mgr.SpaceMgr
	addHeartBeatChan      chan *session.Session
	delHeartBeatChan      chan *session.Session
	serverCombinationChan chan *pb.ServerCombination
	doRequestCount        int32
}

func NewLogic() *Logic {
	l := new(Logic)
	l.SetType(constants.LogicModule)
	return l
}

func (l *Logic) Init(config *config.AiryConfig) {
	l.BaseModule.Init(config)
	l.spaceMgr = mgr.NewSpaceMgr()

	l.addHeartBeatChan = make(chan *session.Session, 1<<10)
	l.delHeartBeatChan = make(chan *session.Session, 1<<10)
	l.serverCombinationChan = make(chan *pb.ServerCombination, 1<<10)
}

func (l *Logic) AfterInit() {
	session.BindCallBack(l.sessionTimeoutCallBack)
	space.BindCallBack(l.spaceCreateCallBack, l.spaceCloseCallBack)
	t := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-t.C:
				logger.Infof("推送的消息通道大小为：%d,1秒处理完的消息：%d\n", len(l.serverCombinationChan), l.doRequestCount)
				l.doRequestCount = 0
			}
		}
	}()
}

func (l *Logic) BeforeShutdown() {
	close(l.addHeartBeatChan)
	close(l.delHeartBeatChan)
	close(l.serverCombinationChan)
	l.clearSpace()
}

func (l *Logic) heartbeatCheck() {
	for {
		select {
		case s, ok := <-l.addHeartBeatChan:
			if !ok {
				return
			}
			l.Timer().AddCountTimer(s.ID(), l.HeartBeat(), -1, func(...any) {
				now := time.Now()
				deadline := now.Add(-2 * l.HeartBeat()).Unix()
				if atomic.LoadInt64(s.LastAt()) <= deadline {
					logger.Errorf("Session heartbeat timeout,pid = %d, LastTime=%d, Deadline=%d,Now = %d", s.PID(), atomic.LoadInt64(s.LastAt()), deadline, now.Unix())
					l.delHeartBeatChan <- s
				}
			})
			logger.Debugf("adding the heartbeat detection succeeded,pid = %d,id = %s", s.PID(), s.ID())
		case s, ok := <-l.delHeartBeatChan:
			if !ok {
				return
			}
			s.Timeout()
			l.Timer().RemoveTimer(s.ID())
		}
	}
}

func (l *Logic) Run(stop chan struct{}) {
	go l.heartbeatCheck()
	go l.combination()
}

func (l *Logic) AddComponents(comps ...component.Component) {
	l.BaseModule.AddComponents(comps...)
}

// Request method that the logical server, data packet return
func (l *Logic) Request(context context.Context, req *pb.Packet) (*empty.Empty, error) {
	sess := session.SessionByPID(req.PID)
	if sess != nil {
		sess.Space().WriteRequestBuffer(req)
		return &empty.Empty{}, nil
	}
	spaceID := req.SpaceID
	if spaceID != "" {
		s := l.spaceMgr.SpaceByID(spaceID)
		if s == nil {
			if validateCreateSpaceParam(req) {
				logger.Debugf("not found spaceID[%s] from space manager,creating space...", spaceID)
				newSpace, err := space.NewSpace(spaceID, l.CreateAOIManager())
				if err != nil {
					logger.Errorf("request create space error : %s", err.Error())
					return &empty.Empty{}, nil
				}
				newSpace.WriteRequestBuffer(req)
			}
		} else {
			s.WriteRequestBuffer(req)
		}

	}
	return &empty.Empty{}, nil
}

func (l *Logic) loop(space *space.Space) {
	defer func() {
		err := recover()
		if err != nil {
			buf := make([]byte, 1<<10)
			runtime.Stack(buf, true)
			logger.Errorf("logic loop func occur a serious error occurs,err : %s", string(buf))
		}
	}()
	for {
		select {
		case request := <-space.ReadRequestBuffer():
			if space.IsClose() {
				return
			}
			logger.Debugf("space = %s receive a rpc request,requestCode = %d", space.ID(), request.RequestCode)
			sess := session.SessionByPID(request.PID)
			if request.Type == pb.Type_Join {
				if sess == nil {
					logger.Debugf("not found session,pid = %d", request.PID)
					sess = session.NewSession(request.PID, request.SourceServerID, space)
					logger.Debugf("create session successfully,pid = %d", request.PID)
					l.addHeartBeatChan <- sess
				} else {
					sess = nil
				}
			}
			if sess == nil {
				break
			}
			sess.SetLastAt(time.Now().Unix())
			l.HandRoute(actx.AddSessionSpaceToCtx(context.Background(), space, sess), request)
			if request.Type == pb.Type_Leave {
				sess.Close()
				l.delHeartBeatChan <- sess
				if space.AOIManager().Num() == 0 {
					space.Close()
					return
				}
			}
			l.doRequestCount++
		}
	}
}

func (l *Logic) spaceCreateCallBack(s *space.Space) error{
	s.BindChan(l.serverCombinationChan)
	err := l.spaceMgr.AddSpace(s)
	if err != nil {
		logger.Errorf(err.Error())
		return err
	}
	go l.loop(s)
	return nil
}

func (l *Logic) spaceCloseCallBack(s *space.Space) error{
	l.spaceMgr.DelSpace(s)
	_ = l.SysStorage().DelSpaceIDLogicID(s.ID(), l.Server().ID)
	return nil
}

func (l *Logic) sessionTimeoutCallBack(s *session.Session) {
	_, _ = l.Request(context.Background(), &pb.Packet{
		Type:        pb.Type_Leave,
		PID:         s.PID(),
		RequestCode: constants.LeaveRequest,
	})
	logger.Debugf("removing the heartbeat detection succeeded,pid = %d,uniqueID = %s", s.PID(), s.ID())
}

func (l *Logic) clearSpace() {
	_ = l.SysStorage().DelLogicIDSpaceID(l.Server().ID)
}

func (l *Logic) combination() {
	pushMessages := make([]*pb.ServerCombination, 0, 1<<4)
	ticker := time.NewTicker(time.Millisecond)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case pmc, ok := <-l.serverCombinationChan:
			if !ok {
				return
			}
			pushMessages = append(pushMessages, pmc)
		case <-ticker.C:
			length := l.GRPCClientComponent().GateServerCount()
			if length == 0 {
				break
			}
			messageMap := make(map[string][]*pb.Combination, length)
			for i := 0; i < len(pushMessages); i++ {
				combinations := messageMap[pushMessages[i].ServerID]
				if combinations == nil {
					combinations = make([]*pb.Combination, 0, len(pushMessages)/length)
				}
				messageMap[pushMessages[i].ServerID] = append(combinations, pushMessages[i].Combination)
			}
			for serverID, combinations := range messageMap {
				l.GRPCClientComponent().PushToUsers(context.Background(), serverID, &pb.Combinations{Combinations: combinations})
				logger.Debugf("push grpc data, : %+v", &pb.Combinations{Combinations: combinations})
			}
			pushMessages = pushMessages[0:0]
		}
	}
}

func validateCreateSpaceParam(req *pb.Packet) bool {
	if req.Type == pb.Type_Join {
		return true
	}
	return false
}
