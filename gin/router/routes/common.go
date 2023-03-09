package routes

import (
	"github.com/MQEnergy/gin-framework/app/controller/common"
	"github.com/gin-gonic/gin"
)

func InitCommonGroup(r *gin.RouterGroup) (router gin.IRoutes) {
	commonGroup := r.Group("")
	{
		// ping
		commonGroup.GET("/ping", common.Common.Ping)
		//TOKEN校验
		commonGroup.POST("/validate/token/:token", common.Common.ValidateToken)
		//登陆
		commonGroup.POST("/login", common.User.Login)
	}
	return commonGroup
}
