package app

import (
	"github.com/MQEnergy/gin-framework/app/controller/base"
	"github.com/MQEnergy/gin-framework/pkg/response"
	"github.com/MQEnergy/gin-framework/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

var tokenContainer = make(map[string]struct{})

type TokenController struct {
	base.Controller
}

var Token = TokenController{}

// Create 生成token
func (c *TokenController) Create(ctx *gin.Context) {
	token := util.GetToken()
	response.ResponseJson(ctx, http.StatusOK, response.Success, "", token)
}

//// View token解析
//func (c *TokenController) View(ctx *gin.Context) {
//	token := ctx.GetHeader(global.Cfg.Jwt.TokenKey)
//	if token == "" {
//		response.UnauthorizedException(ctx, "")
//		return
//	}
//	flag := strings.Contains(token, "Bearer")
//	if !flag {
//		response.UnauthorizedException(ctx, "")
//		return
//	}
//	token = strings.TrimSpace(strings.TrimLeft(token, "Bearer"))
//	jwtTokenArr, err := auth.ParseJwtToken(token, global.Cfg.Jwt.Secret)
//	if err != nil {
//		response.UnauthorizedException(ctx, "")
//		return
//	}
//	response.ResponseJson(ctx, http.StatusOK, response.Success, "", jwtTokenArr)
//}
