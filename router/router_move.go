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

package router

import (
	"context"
	message "github.com/airy/build_message"
	actx "github.com/airy/context"
	"github.com/airy/logger"
	"github.com/airy/pb"
	"github.com/airy/player"
)

type RouteMove struct {
	BaseRouter
}

func (r *RouteMove) Handle(ctx context.Context, request *pb.Packet) {
	space := actx.GetSpaceFromCtx(ctx)
	session := actx.GetSessionFromCtx(ctx)
	p := session.Player()
	serializer := r.serializer

	data := new(pb.Move)
	err := serializer.Unmarshal(request.Data, data)
	if err != nil {
		logger.Errorf("serializer fail,the byte array cannot be parsed into its corresponding structure")
	}
	logger.Debugf("logic server receive request message[MOVE],pid = %d,x = %f,z = %f", request.PID, data.X, data.Z)

	aoiMgr := space.AOIManager()
	moveBefore := aoiMgr.Neighbors(p)
	aoiMgr.Move(p, data.X, data.Z)
	moveAfter := aoiMgr.Neighbors(p)

	outOfView := make(map[*player.Player]struct{})
	comeIntoView := make(map[*player.Player]struct{})

	myMove := message.BuildMove(p, serializer)
	myComeIntoView := message.BuildComeIntoView(p, serializer)
	myOutOfView := message.BuildOutOfView(p, serializer)

	pushMyPacket := make([]*pb.Packet, 0, 32)

	pushOthersMyMove := make(map[string][]uint32)
	pushOthersMyComeIntoView := make(map[string][]uint32)
	pushOthersMyOutOfView := make(map[string][]uint32)

	for before, _ := range moveBefore {
		if _, ok := moveAfter[before]; ok {
			delete(moveBefore, before)
			delete(moveAfter, before)

			pids := pushOthersMyMove[before.GateServerID]
			if pids == nil {
				pids = make([]uint32, 0, 32)
			}
			pushOthersMyMove[before.GateServerID] = append(pids, before.PID)
		}
	}
	outOfView = moveBefore
	comeIntoView = moveAfter

	for other, _ := range outOfView {
		otherOutOfView := message.BuildOutOfView(other, serializer)
		pushMyPacket = append(pushMyPacket, otherOutOfView)

		pids := pushOthersMyOutOfView[other.GateServerID]
		if pids == nil {
			pids = make([]uint32, 0, 32)
		}
		pushOthersMyMove[other.GateServerID] = append(pids, other.PID)
	}

	for other, _ := range comeIntoView {
		otherComeIntoView := message.BuildComeIntoView(other, serializer)
		pushMyPacket = append(pushMyPacket, otherComeIntoView)

		pids := pushOthersMyComeIntoView[other.GateServerID]
		if pids == nil {
			pids = make([]uint32, 0, 32)
		}
		pushOthersMyMove[other.GateServerID] = append(pids, other.PID)
	}

	if len(pushMyPacket) > 0 {
		space.SendMessage(&pb.ServerCombination{
			ServerID:    p.GateServerID,
			Combination: &pb.Combination{PIDs: []uint32{p.PID}, Packets: pushMyPacket},
		})
	}

	for serverID, pids := range pushOthersMyMove {
		space.SendMessage(&pb.ServerCombination{
			ServerID:    serverID,
			Combination: &pb.Combination{PIDs: pids, Packets: []*pb.Packet{myMove}},
		})
	}

	for serverID, pids := range pushOthersMyComeIntoView {
		space.SendMessage(&pb.ServerCombination{
			ServerID:    serverID,
			Combination: &pb.Combination{PIDs: pids, Packets: []*pb.Packet{myComeIntoView}},
		})
	}

	for serverID, pids := range pushOthersMyOutOfView {
		space.SendMessage(&pb.ServerCombination{
			ServerID:    serverID,
			Combination: &pb.Combination{PIDs: pids, Packets: []*pb.Packet{myOutOfView}},
		})
	}
}
