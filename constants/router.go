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

package constants

// 推送的CODE列表
const (
	OnHandshake    int32 = 8000
	OnComeIntoView int32 = 8001
	OnOutOfView    int32 = 8002
	OnMove         int32 = 8003
	OnLeave        int32 = 8004
	OnDelay        int32 = 8005
	OnGlobalChat   int32 = 8006
)

//请求的CODE列表
const (
	HeartBeatRequest  int32 = 1000
	JoinRequest       int32 = 1001
	MoveRequest       int32 = 1002
	LeaveRequest      int32 = 1003
	DelayRequest      int32 = 1004
	GlobalChatRequest int32 = 1005
)
