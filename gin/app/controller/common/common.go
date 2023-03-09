package common

import (
	"fmt"
	"github.com/MQEnergy/gin-framework/app/controller/base"
	"github.com/MQEnergy/gin-framework/pkg/response"
	"github.com/MQEnergy/gin-framework/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommonController struct {
	*base.Controller
}

var Common = &CommonController{}

// Ping 心跳
func (c *CommonController) Ping(ctx *gin.Context) {
	response.ResponseJson(ctx, http.StatusOK, response.Success, "Pong!", "")
}

func (c *CommonController) ValidateToken(ctx *gin.Context) {
	token := ctx.Param("token")
	fmt.Println("收到校验token请求，token = ", token)
	success := util.ValidateToken(token)
	response.ResponseJson(ctx, http.StatusOK, response.Success, "", success)
}
