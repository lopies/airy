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

type RouteLeave struct {
	BaseRouter
}

func (r *RouteLeave) Handle(ctx context.Context, request *pb.Packet) {
	logger.Debugf("logic server receive request message[Exit],pid = %d", request.PID)
	space := actx.GetSpaceFromCtx(ctx)
	session := actx.GetSessionFromCtx(ctx)
	serializer := r.serializer
	player := session.Player()

	aoiMgr := space.AOIManager()
	aoiMgr.Leave(player)
	exit := aoiMgr.Neighbors(player)

	myExit := message.BuildExit(player, serializer)
	pushMyExit := make(map[string][]uint32)
	for other, _ := range exit {
		pids := pushMyExit[other.GateServerID]
		if pids == nil {
			pids = make([]uint32, 0, 32)
		}
		pids = append(pids, other.PID)
		pushMyExit[other.GateServerID] = pids
	}
	for serverID, pids := range pushMyExit {
		space.SendMessage(&pb.ServerCombination{
			ServerID:    serverID,
			Combination: &pb.Combination{PIDs: pids, Packets: []*pb.Packet{myExit}},
		})
	}
}
