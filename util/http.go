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
	"github.com/airy/logger"
	"io/ioutil"
	"net/http"
	"net/url"
)

type httpResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type loginResponse struct {
	URL    string `json:"url"`
	Bearer string `json:"bearer"`
	Token  string `json:"token"`
}

func ValidToken(token string) bool {
	resp, err := http.Post("http://127.0.0.1:9527/validate/token/"+token, "application/json", nil)
	if err != nil {
		logger.Errorf("Failed to post request[/validate/token]: %s", err.Error())
		return false
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to post read response[/validate/token]: %s", err.Error())
		return false
	}
	response := new(httpResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		logger.Errorf("Failed to unmarshal response post request[/validate/token]: %s", err.Error())
		return false
	}
	if response.Data.(bool) {
		return true
	}
	return false
}

// 13888888888 123456
func Login(phone string, passwd string) (*loginResponse, error) {
	param := url.Values{}
	param.Add("phone", phone)
	param.Add("password", passwd)
	resp, err := http.PostForm("http://127.0.0.1:9527/login", param)
	if err != nil {
		logger.Errorf("Failed to post request[/login]: %s", err.Error())
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to post read response[/login]: %s", err.Error())
		return nil, err
	}
	response := new(httpResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		logger.Errorf("Failed to unmarshal response post request[/login]: %s", err.Error())
		return nil, err
	}
	res := new(loginResponse)
	byt, _ := json.Marshal(response.Data)
	json.Unmarshal(byt, res)
	return res, nil
}
