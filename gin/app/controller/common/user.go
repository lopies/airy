package common

import (
	"github.com/MQEnergy/gin-framework/app/controller/base"
	"github.com/MQEnergy/gin-framework/app/service/common"
	"github.com/MQEnergy/gin-framework/pkg/response"
	"github.com/MQEnergy/gin-framework/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserController struct {
	*base.Controller
}

var User = &UserController{}

// Login 用户登录
func (c *UserController) Login(ctx *gin.Context) {
	var requestParams types.LoginRequest
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.BadRequestException(ctx, err.Error())
		return
	}
	res, err := common.User.Login(requestParams)
	if err != nil {
		response.BadRequestException(ctx, "")
		return
	}
	response.ResponseJson(ctx, http.StatusOK, response.Success, "", res)
}
