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
	"github.com/airy/constants"
	actx "github.com/airy/context"
	"github.com/airy/interfaces"
	"github.com/airy/logger"
	"github.com/airy/pb"
)

type BaseRouter struct {
	serializer interfaces.Serializer
}

// Init initializes serializer
func (br *BaseRouter) Init(serializer interfaces.Serializer) {
	br.serializer = serializer
}

// Handle Logical execution
func (br *BaseRouter) Handle(ctx context.Context, request *pb.Packet) {}

// Post Logical post-execution
func (br *BaseRouter) Post(ctx context.Context, request *pb.Packet) {
	if request.RequestID != 0 {
		sess := actx.GetSessionFromCtx(ctx)
		p := constants.NewResponse(constants.Success, request.RequestID, pb.Type_Response_)
		p.RequestCode = request.RequestCode

		actx.GetSpaceFromCtx(ctx).SendMessage(&pb.ServerCombination{
			ServerID:    actx.GetSessionFromCtx(ctx).Player().GateServerID,
			Combination: &pb.Combination{PIDs: []uint32{sess.Player().PID}, Packets: []*pb.Packet{p}},
		})
		logger.Debugf("push response packet,packet = %+v", p)
	}
}
