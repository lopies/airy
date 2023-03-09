package app

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/MQEnergy/gin-framework/global"
	"github.com/MQEnergy/gin-framework/pkg/auth"
	"strconv"
)

type TokenPayload struct {
	UserId int64 `json:"id"`
}

// ParseUserByToken 通过token解析用户
func ParseUserByToken(token string) (TokenPayload, error) {
	user := TokenPayload{}
	if token == "" {
		return user, errors.New("token 为空")
	}
	jwtPayload, err := auth.ParseJwtToken(token, global.Cfg.Jwt.Secret)
	if err != nil {
		return user, err
	}
	byteSlice, err := json.Marshal(jwtPayload.User)
	if err != nil {
		return user, err
	}
	if err = json.Unmarshal(byteSlice, &user); err != nil {
		return user, err
	}
	if user.UserId == 0 {
		return user, errors.New("非法登录")
	}
	_, err = global.Redis.Get(context.Background(), global.Cfg.Redis.LoginPrefix+strconv.FormatInt(user.UserId, 10)).Result()
	if err != nil {
		return TokenPayload{}, errors.New("会话过期，请重新登录")
	}
	return user, nil
}
