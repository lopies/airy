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

import "errors"

var (
	ErrUnsupportedLength        = errors.New("unsupported length")
	ErrTooLessLength            = errors.New("too less length")
	ErrIncompletePacket         = errors.New("incomplete packet")
	ErrUnexpectedEOF            = errors.New("unexpected EOF")
	ErrWrongValueType           = errors.New("types that cannot be converted")
	ErrWrongMessageType         = errors.New("message type is wrong")
	ErrIllegalPID               = errors.New("illegal pid")
	ErrPIDAlreadyBound          = errors.New("agent is already bound to an pid")
	ErrExecCallBack             = errors.New("exec call back fail")
	ErrEtcdGrantLeaseTimeout    = errors.New("timed out waiting for etcd lease grant")
	ErrInvalidSpanCarrier       = errors.New("tracing: invalid span carrier")
	ErrMetricNotKnown           = errors.New("the provided metric does not exist")
	ErrNoBindingStorageModule   = errors.New("for sending remote pushes or using unique session module while using cluster you need to pass it a BindingStorage")
	ErrNoConnectionToServer     = errors.New("rpc client has no connection to the chosen server")
	ErrNoServerWithID           = errors.New("can't find any server with the provided ID")
	ErrNoServerWithType         = errors.New("can't find any server with the server type")
	ErrNoServersAvailableOfType = errors.New("no modules available of this type")
	ErrNotImplemented           = errors.New("method not implemented")
	ErrPlayerOnline             = errors.New("player is online")
	ErrPlayerOffline            = errors.New("player is offline")
	ErrPlayerBindSpace          = errors.New("player bind space error")
	ErrPlayerUnBindSpace        = errors.New("player unbound space error")
	ErrCreateDockerContainer    = errors.New("create a docker container fail")
	ErrStartDockerContainer     = errors.New("start a docker container fail")
	ErrStopDockerContainer      = errors.New("stop a docker container fail")
	ErrConnectionClosed         = errors.New("client connection closed")
	ErrInvalidCertificates      = errors.New("certificates must be exactly two")
	ErrWriteConn                = errors.New("failed to write in conn")
	ErrInvalidStatus            = errors.New("invalid agent status")
	ErrNotFoundServer           = errors.New("not found server")
	ErrCreateSpace              = errors.New("failed to create space from logic server")
	ErrRequest                  = errors.New("failed to hand request")
	ErrNotify                   = errors.New("failed to hand notify")
	ErrInvalidToken             = errors.New("token is invalid")
	ErrTimeInterval             = errors.New("the timer delay cannot be less than 0")
	ErrAgentHandshake           = errors.New("agent handshake error")
)
