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

package codec

import (
	"encoding/binary"
	"fmt"
	"github.com/airy/constants"
	"github.com/airy/interfaces"
	"github.com/airy/logger"
	"io"
	"io/ioutil"
	"net"
)

func NewLengthFieldBasedFrameCodec() interfaces.PacketCodec {
	return &LengthFieldBasedFrameCodec{encoderConfig: encoderConfig{
		ByteOrder:         binary.BigEndian,
		LengthFieldLength: 2,
		LengthAdjustment:  0,

		LengthIncludesLengthFieldLength: true,
	}, decoderConfig: decoderConfig{
		ByteOrder:           binary.BigEndian,
		LengthFieldOffset:   0,
		LengthFieldLength:   2,
		LengthAdjustment:    -2,
		InitialBytesToStrip: 2,
	}}
}

// LengthFieldBasedFrameCodec is the refactoring from
// It encodes/decodes frames into/from TCP stream with value of the length field in the message.
type LengthFieldBasedFrameCodec struct {
	encoderConfig encoderConfig
	decoderConfig decoderConfig
}

// EncoderConfig config for encoder.
type encoderConfig struct {
	// ByteOrder is the ByteOrder of the length field.
	ByteOrder binary.ByteOrder
	// LengthFieldLength is the length of the length field.
	LengthFieldLength int
	// LengthAdjustment is the compensation value to add to the value of the length field
	LengthAdjustment int
	// LengthIncludesLengthFieldLength is true, the length of the prepended length field is added to the value of
	// the prepended length field
	LengthIncludesLengthFieldLength bool
}

// DecoderConfig config for decoder.
type decoderConfig struct {
	// ByteOrder is the ByteOrder of the length field.
	ByteOrder binary.ByteOrder
	// LengthFieldOffset is the offset of the length field
	LengthFieldOffset int64
	// LengthFieldLength is the length of the length field
	LengthFieldLength int
	// LengthAdjustment is the compensation value to add to the value of the length field
	LengthAdjustment int64
	// InitialBytesToStrip is the number of first bytes to strip out from the decoded frame
	InitialBytesToStrip int
}

// Encode ...
func (cc *LengthFieldBasedFrameCodec) Encode(buf []byte) (out []byte, err error) {
	length := len(buf) + cc.encoderConfig.LengthAdjustment
	if cc.encoderConfig.LengthIncludesLengthFieldLength {
		length += cc.encoderConfig.LengthFieldLength
	}

	if length < 0 {
		return nil, constants.ErrTooLessLength
	}

	switch cc.encoderConfig.LengthFieldLength {
	case 1:
		if length >= 256 {
			return nil, fmt.Errorf("length does not fit into a byte: %d", length)
		}
		out = []byte{byte(length)}
	case 2:
		if length >= 65536 {
			return nil, fmt.Errorf("length does not fit into a short integer: %d", length)
		}
		out = make([]byte, 2)
		cc.encoderConfig.ByteOrder.PutUint16(out, uint16(length))
	case 3:
		if length >= 16777216 {
			return nil, fmt.Errorf("length does not fit into a medium integer: %d", length)
		}
		out = writeUint24(cc.encoderConfig.ByteOrder, length)
	case 4:
		out = make([]byte, 4)
		cc.encoderConfig.ByteOrder.PutUint32(out, uint32(length))
	case 8:
		out = make([]byte, 8)
		cc.encoderConfig.ByteOrder.PutUint64(out, uint64(length))
	default:
		return nil, constants.ErrUnsupportedLength
	}

	out = append(out, buf...)
	return
}

// Decode ...
func (cc *LengthFieldBasedFrameCodec) Decode(c net.Conn) ([]byte, error) {
	var (
		header []byte
		err    error
	)

	if cc.decoderConfig.LengthFieldOffset > 0 {
		header, err = peak(c, cc.decoderConfig.LengthFieldOffset)
		if err != nil {
			return nil, err
		}
	}
	lenBuf, frameLength, err := cc.getUnadjustedFrameLength(c)
	if err != nil {
		if err == constants.ErrUnsupportedLength {
			return nil, err
		}
		return nil, constants.ErrUnexpectedEOF
	}

	// real message length
	msgLength := int64(frameLength) + cc.decoderConfig.LengthAdjustment
	msg, _ := peak(c, msgLength)
	packetLen := len(header) + len(lenBuf) + int(msgLength)
	fullMessage := make([]byte, packetLen)
	copy(fullMessage, header)
	copy(fullMessage[len(header):], lenBuf)
	copy(fullMessage[len(header)+len(lenBuf):], msg)
	return fullMessage[cc.decoderConfig.InitialBytesToStrip:], nil
}

func (cc *LengthFieldBasedFrameCodec) getUnadjustedFrameLength(c net.Conn) ([]byte, uint64, error) {
	switch cc.decoderConfig.LengthFieldLength {
	case 1:
		b, err := peak(c, 1)
		if err != nil {
			return nil, 0, constants.ErrUnexpectedEOF
		}
		return b, uint64(b[0]), nil
	case 2:
		lenBuf, err := peak(c, 2)
		if err != nil {
			return nil, 0, constants.ErrUnexpectedEOF
		}
		return lenBuf, uint64(cc.decoderConfig.ByteOrder.Uint16(lenBuf)), nil
	case 3:
		lenBuf, err := peak(c, 3)
		if err != nil {
			return nil, 0, constants.ErrUnexpectedEOF
		}
		return lenBuf, readUint24(cc.decoderConfig.ByteOrder, lenBuf), nil
	case 4:
		lenBuf, err := peak(c, 4)
		if err != nil {
			return nil, 0, constants.ErrUnexpectedEOF
		}
		return lenBuf, uint64(cc.decoderConfig.ByteOrder.Uint32(lenBuf)), nil
	case 8:
		lenBuf, err := peak(c, 8)
		if err != nil {
			return nil, 0, constants.ErrUnexpectedEOF
		}
		return lenBuf, cc.decoderConfig.ByteOrder.Uint64(lenBuf), nil
	default:
		return nil, 0, constants.ErrUnsupportedLength
	}
}

func peak(c net.Conn, n int64) (buf []byte, err error) {
	buf, err = ioutil.ReadAll(io.LimitReader(c, n))
	if err != nil {
		logger.Errorf("connection peak error : %s", err.Error())
		return nil, err
	}
	if len(buf) == 0 {
		return nil, constants.ErrUnexpectedEOF
	}
	return
}

func readUint24(byteOrder binary.ByteOrder, b []byte) uint64 {
	_ = b[2]
	if byteOrder == binary.LittleEndian {
		return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16
	}
	return uint64(b[2]) | uint64(b[1])<<8 | uint64(b[0])<<16
}

func writeUint24(byteOrder binary.ByteOrder, v int) []byte {
	b := make([]byte, 3)
	if byteOrder == binary.LittleEndian {
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
	} else {
		b[2] = byte(v)
		b[1] = byte(v >> 8)
		b[0] = byte(v >> 16)
	}
	return b
}
