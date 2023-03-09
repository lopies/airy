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

package component

import "C"
import (
	"fmt"
	"github.com/airy/config"
	"github.com/airy/constants"
	"github.com/airy/logger"
	"github.com/garyburd/redigo/redis"
)

type Redis struct {
	client *redis.Pool
	BaseComponent
}

func NewRedis() *Redis {
	r := new(Redis)
	r.SetType(constants.RedisComponent)
	return r
}

func (rs *Redis) Init(config *config.AiryConfig) {
	pool := &redis.Pool{
		MaxIdle: config.Redis.MaxIdle,
		//连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		MaxActive: config.Redis.MaxActive,
		//连接关闭时间
		IdleTimeout: config.Redis.IdleTimeout,
		Wait:        config.Redis.Wait,
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			return redis.Dial("tcp", config.Redis.Path, redis.DialDatabase(config.Redis.Index),
				redis.DialPassword(config.Redis.Password),
				redis.DialConnectTimeout(config.Redis.ConnectTimeout))
		},
	}
	conn := pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		logger.Errorf("redis client init fail,err=%s", err.Error())
		panic(fmt.Sprintf("redis client init fail,err=%s", err.Error()))
	}
	logger.Infof("redis component init success")
	rs.client = pool
}

func (rs *Redis) Shutdown() error {
	return rs.client.Close()
}

const (
	SPACEID_LOGICID = "space_logic:"
	PLAYERID_GATEID = "player_gate:"
	LOGICID_SPACEID = "logic_space:"
)

func (rs *Redis) StoreSpaceIDLogicID(spaceId, serverID string) error {
	conn := rs.client.Get()
	defer conn.Close()
	_, err := conn.Do("SET", SPACEID_LOGICID+spaceId, serverID)
	if err != nil {
		logger.Errorf("redis client put key[%s] fail,err=%s", SPACEID_LOGICID+spaceId, err.Error())
		return err
	}
	return nil
}

func (rs *Redis) DelSpaceIDLogicID(spaceId, serverID string) error {
	conn := rs.client.Get()
	defer conn.Close()
	conn.Send("MULTI")
	conn.Send("DEL", SPACEID_LOGICID+spaceId)
	conn.Send("SREM", LOGICID_SPACEID+serverID, spaceId)
	_, err := conn.Do("EXEC")
	if err != nil {
		logger.Errorf("redis client del key[%s] fail,err=%s", SPACEID_LOGICID+spaceId, err.Error())
		return err
	}
	return nil
}

func (rs *Redis) GetLogicIDBySpaceID(spaceID string) (string, error) {
	conn := rs.client.Get()
	defer conn.Close()
	serverID, err := redis.String(conn.Do("GET", SPACEID_LOGICID+spaceID))
	if err != nil {
		//logger.Errorf("redis client get key[%s] fail,err=%s", SPACEID_LOGICID+spaceID, err.Error())
		return "", err
	}
	return serverID, nil
}

func (rs *Redis) StorePlayerIDGateID(playerID, serverID string) error {
	conn := rs.client.Get()
	defer conn.Close()
	_, err := conn.Do("SET", PLAYERID_GATEID+playerID, serverID)
	if err != nil {
		logger.Errorf("redis client put key[%s] fail,err=%s", PLAYERID_GATEID+playerID, err.Error())
		return err
	}
	return nil
}

func (rs *Redis) DelPlayerIDGateID(playerID string) error {
	conn := rs.client.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", PLAYERID_GATEID+playerID)
	if err != nil {
		logger.Errorf("redis client del key[%s] fail,err=%s", PLAYERID_GATEID+playerID, err.Error())
		return err
	}
	return nil
}

func (rs *Redis) GetGateIDByPlayerID(playerID string) string {
	conn := rs.client.Get()
	defer conn.Close()
	serverID, err := redis.String(conn.Do("GET", PLAYERID_GATEID+playerID))
	if err != nil {
		return ""
	}
	return serverID
}

func (rs *Redis) StoreLogicIDSpaceID(serverID, spaceID string) error {
	conn := rs.client.Get()
	defer conn.Close()
	_, err := conn.Do("SADD", LOGICID_SPACEID+serverID, spaceID)
	if err != nil {
		logger.Errorf("redis client sadd key[%s] fail,err=%s", LOGICID_SPACEID+serverID, err.Error())
		return err
	}
	return nil
}

func (rs *Redis) DelLogicIDSpaceID(serverID string) error {
	conn := rs.client.Get()
	defer conn.Close()
	spaceIDs, err := redis.Strings(conn.Do("SMEMBERS", LOGICID_SPACEID+serverID))
	if err != nil {
		logger.Errorf("redis client scard key[%s] fail,err=%s", LOGICID_SPACEID+serverID, err.Error())
		return err
	}
	conn.Send("MULTI")
	for _, spaceID := range spaceIDs {
		conn.Send("DEL", SPACEID_LOGICID+spaceID)
	}
	conn.Send("DEL", LOGICID_SPACEID+serverID)
	_, err = conn.Do("EXEC")
	if err != nil {
		return err
	}
	return nil
}
