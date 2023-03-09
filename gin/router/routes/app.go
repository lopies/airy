package routes

import (
	"github.com/MQEnergy/gin-framework/app/controller/app"
	"github.com/gin-gonic/gin"
)

func InitAppGroup(r *gin.RouterGroup) gin.IRoutes {
	appGroup := r.Group("app")
	{
		//重新获取TOKEN
		appGroup.GET("/token", app.Token.Create)
		// 获取用户列表
		appGroup.GET("/user/index", app.User.Index)
		// 获取用户列表
		appGroup.GET("/user/list", app.User.List)
		// 上传附件
		appGroup.POST("/attachment/upload", app.Attachment.Upload)
	}
	return appGroup
}
