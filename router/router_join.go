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
)

type RouteJoin struct {
	BaseRouter
}

func (r *RouteJoin) Handle(ctx context.Context, request *pb.Packet) {
	space := actx.GetSpaceFromCtx(ctx)
	session := actx.GetSessionFromCtx(ctx)
	player := session.Player()
	serializer := r.serializer

	data := new(pb.Enter)
	err := serializer.Unmarshal(request.Data, data)
	if err != nil {
		logger.Errorf("serializer fail,the byte array cannot be parsed into its corresponding structure")
	}
	logger.Debugf("logic server receive request message[ENTER],pid = %d,name = %s,x = %f,z = %f", data.PID, data.Name, data.X, data.Z)

	aoiMgr := space.AOIManager()
	succ := aoiMgr.Enter(player, data.X, data.Z)
	if !succ {
		return
	}
	myComeIntoView := message.BuildComeIntoView(player, serializer)
	m := aoiMgr.Neighbors(player)
	myPacket := make([]*pb.Packet, 0, len(m))
	pushMyComeIntoView := make(map[string][]uint32)
	for other, _ := range m {
		otherComeIntoView := message.BuildComeIntoView(other, serializer)
		myPacket = append(myPacket, otherComeIntoView)

		pids := pushMyComeIntoView[other.GateServerID]
		if pids == nil {
			pids = make([]uint32, 0, 32)
		}
		pids = append(pids, other.PID)
		pushMyComeIntoView[other.GateServerID] = pids
	}

	for serverID, pids := range pushMyComeIntoView {
		space.SendMessage(&pb.ServerCombination{
			ServerID:    serverID,
			Combination: &pb.Combination{PIDs: pids, Packets: []*pb.Packet{myComeIntoView}},
		})
	}
	if len(myPacket) > 0 {
		space.SendMessage(&pb.ServerCombination{
			ServerID:    player.GateServerID,
			Combination: &pb.Combination{PIDs: []uint32{player.PID}, Packets: myPacket},
		})
	}
}
