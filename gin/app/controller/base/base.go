package base

import (
	"errors"
	"github.com/MQEnergy/gin-framework/pkg/validator"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

var Base = Controller{}

// ValidateReqParams 验证请求参数
func (c *Controller) ValidateReqParams(ctx *gin.Context, requestParams interface{}) error {
	var err error
	if ctx.ContentType() != "application/json" {
		err = ctx.Bind(requestParams)
	} else {
		err = ctx.BindJSON(requestParams)
	}
	if err != nil {
		translate := validator.Translate(err)
		return errors.New(translate[0])
	}
	return nil
}
