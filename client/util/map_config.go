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

package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Birth struct {
	X float32 `json:"x"`
	Z float32 `json:"z"`
}

type MapData struct {
	Map    [][]bool `json:"map"`
	Birth  []*Birth `json:"birth"`
	Width  int      `json:"width"`
	Height int      `json:"height"`
}

func ReadMapData(id string) (*MapData, error) {
	file, err := os.Open(fmt.Sprintf("../map_data/%s.json", id))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("map[%s] data not found,%s", id, err.Error()))
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		panic(fmt.Errorf(err.Error()))
	}
	mapData := new(MapData)
	err = json.Unmarshal(data, mapData)
	if err != nil {
		return nil, err
	}
	return mapData, nil
}
