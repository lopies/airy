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

type RouteGlobalChat struct {
	BaseRouter
}

func (r *RouteGlobalChat) Handle(ctx context.Context, request *pb.Packet) {
	space := actx.GetSpaceFromCtx(ctx)
	session := actx.GetSessionFromCtx(ctx)

	serializer := r.serializer
	data := new(pb.CommonChat)
	err := serializer.Unmarshal(request.Data, data)
	if err != nil {
		logger.Errorf("serializer fail,the byte array cannot be parsed into its corresponding structure")
	}
	logger.Debugf("logic server receive request message[RouteGlobalChat],pid = %d", request.PID)

	player := session.Player()
	myGlobalChat := message.BuildGlobalChat(player, serializer, 1, data.Content)
	aoiMgr := space.AOIManager()
	m := aoiMgr.AllPlayers()
	pushMyGlobalChat := make(map[string][]uint32)
	for _, other := range m {
		pids := pushMyGlobalChat[other.GateServerID]
		if pids == nil {
			pids = make([]uint32, 0, 32)
		}
		pids = append(pids, other.PID)
		pushMyGlobalChat[other.GateServerID] = pids
	}

	for serverID, pids := range pushMyGlobalChat {
		space.SendMessage(&pb.ServerCombination{
			ServerID:    serverID,
			Combination: &pb.Combination{PIDs: pids, Packets: []*pb.Packet{myGlobalChat}},
		})
	}
}
