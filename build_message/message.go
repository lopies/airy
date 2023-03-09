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

package message

import (
	"github.com/airy/constants"
	"github.com/airy/interfaces"
	"github.com/airy/pb"
	"github.com/airy/player"
	"github.com/airy/push"
	"time"
)

/**
BuildXXX Used to construct data packets to send client
*/

func BuildComeIntoView(p *player.Player, serializer interfaces.Serializer) *pb.Packet {
	data, _ := serializer.Marshal(&push.OnComeIntoView{
		PID:  p.PID,
		Name: p.Name,
		X:    p.X,
		Z:    p.Z,
	})
	return &pb.Packet{
		Type:      pb.Type_Push_,
		RouteCode: constants.OnComeIntoView,
		Data:      data,
	}
}

func BuildOutOfView(p *player.Player, serializer interfaces.Serializer) *pb.Packet {
	data, _ := serializer.Marshal(&push.OnComeIntoView{
		PID:  p.PID,
		Name: p.Name,
		X:    p.X,
		Z:    p.Z,
	})
	return &pb.Packet{
		Type:      pb.Type_Push_,
		RouteCode: constants.OnOutOfView,
		Data:      data,
	}
}

func BuildMove(p *player.Player, serializer interfaces.Serializer) *pb.Packet {
	data, _ := serializer.Marshal(&push.OnMove{
		PID: p.PID,
		X:   p.X,
		Z:   p.Z,
	})
	return &pb.Packet{
		Type:      pb.Type_Push_,
		RouteCode: constants.OnMove,
		Data:      data,
	}
}

func BuildExit(p *player.Player, serializer interfaces.Serializer) *pb.Packet {
	data, _ := serializer.Marshal(&push.OnExit{
		PID: p.PID,
	})
	return &pb.Packet{
		Type:      pb.Type_Push_,
		RouteCode: constants.OnLeave,
		Data:      data,
	}
}

func BuildGlobalChat(p *player.Player, serializer interfaces.Serializer, tp int32, content string) *pb.Packet {
	data, _ := serializer.Marshal(&push.OnCommonChat{
		Tp:      tp,
		Pid:     p.PID,
		Time:    time.Now().UnixMilli(),
		Content: content,
	})
	return &pb.Packet{
		Type:      pb.Type_Push_,
		RouteCode: constants.OnGlobalChat,
		Data:      data,
	}
}

func BuildDelay(serializer interfaces.Serializer, startTime int64) *pb.Packet {
	data, _ := serializer.Marshal(&push.OnDelay{
		StartTime: startTime,
	})
	return &pb.Packet{
		Type:      pb.Type_Push_,
		RouteCode: constants.OnDelay,
		Data:      data,
	}
}
