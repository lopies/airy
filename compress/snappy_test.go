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

package compress

import (
	"fmt"
	"testing"
)

func TestSnappy(t *testing.T) {
	data := []byte{10, 24, 54, 50, 98, 98, 101, 54, 52, 55, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 100, 18, 10, 50, 53, 49, 51, 48, 50, 48, 57, 48, 48, 26, 24, 54, 50, 98, 98, 101, 54, 51, 97, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 57, 10, 24, 54, 50, 98, 98, 101, 54, 52, 55, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 100, 18, 10, 50, 53, 49, 51, 48, 50, 48, 57, 48, 48, 26, 24, 54, 50, 98, 98, 101, 54, 51, 97, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 57, 10, 24, 54, 50, 98, 98, 101, 54, 52, 55, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 100, 18, 10, 50, 53, 49, 51, 48, 50, 48, 57, 48, 48, 26, 24, 54, 50, 98, 98, 101, 54, 51, 97, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 57, 10, 24, 54, 50, 98, 98, 101, 54, 52, 55, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 100, 18, 10, 50, 53, 49, 51, 48, 50, 48, 57, 48, 48, 26, 24, 54, 50, 98, 98, 101, 54, 51, 97, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 57}
	fmt.Println("src data=", data, "src len=", len(data))
	snappy := NewSnappyCompress()
	encodeData := snappy.Encode(data)
	fmt.Println("encode data=", encodeData, "encode len=", len(encodeData))
	decodeData, err := snappy.Decode(encodeData)
	if err != nil {
		fmt.Println("err===")
	}
	fmt.Println("decode data=", decodeData, "decode len=", len(decodeData))
}

func BenchmarkWithSnappy(b *testing.B) {
	data := []byte{10, 24, 54, 50, 98, 98, 101, 54, 52, 55, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 100, 18, 10, 50, 53, 49, 51, 48, 50, 48, 57, 48, 48, 26, 24, 54, 50, 98, 98, 101, 54, 51, 97, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 57, 10, 24, 54, 50, 98, 98, 101, 54, 52, 55, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 100, 18, 10, 50, 53, 49, 51, 48, 50, 48, 57, 48, 48, 26, 24, 54, 50, 98, 98, 101, 54, 51, 97, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 57, 10, 24, 54, 50, 98, 98, 101, 54, 52, 55, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 100, 18, 10, 50, 53, 49, 51, 48, 50, 48, 57, 48, 48, 26, 24, 54, 50, 98, 98, 101, 54, 51, 97, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 57, 10, 24, 54, 50, 98, 98, 101, 54, 52, 55, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 100, 18, 10, 50, 53, 49, 51, 48, 50, 48, 57, 48, 48, 26, 24, 54, 50, 98, 98, 101, 54, 51, 97, 99, 57, 49, 52, 55, 52, 57, 99, 57, 99, 50, 101, 51, 57, 48, 57}
	snappy := NewSnappyCompress()
	for n := 0; n < b.N; n++ {
		encodeData := snappy.Encode(data)
		//fmt.Println("encode data=", encodeData, "encode len=", len(encodeData))
		decodeData, err := snappy.Decode(encodeData)
		if err != nil {
			fmt.Println("err===", decodeData)
		}
	}
}
