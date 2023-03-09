package common

import (
	"errors"
	"github.com/MQEnergy/gin-framework/global"
	"github.com/MQEnergy/gin-framework/pkg/auth"
	"github.com/MQEnergy/gin-framework/pkg/lib"
	"github.com/MQEnergy/gin-framework/pkg/util"
	"github.com/MQEnergy/gin-framework/types"
)

type UserService struct{}

var User = UserService{}

// Login 登录操作
func (s UserService) Login(requestParams types.LoginRequest) (interface{}, error) {
	var userInfo types.User
	if err := global.DB.Where("phone = ? and password = ?", requestParams.Phone, requestParams.Password).First(&userInfo).Error; err != nil {
		return userInfo, errors.New("未查找到用户")
	}
	jwtToken, err := auth.GenerateJwtToken(global.Cfg.Jwt.Secret, global.Cfg.Jwt.TokenExpire, userInfo, global.Cfg.Jwt.TokenIssuer)
	if err != nil {
		return "", errors.New("token生成失败")
	}
	token := util.GetToken()
	url := lib.GetURL()
	return types.LoginResponse{
		URL:    url,
		Bearer: jwtToken,
		Token:  token,
	}, nil
}
