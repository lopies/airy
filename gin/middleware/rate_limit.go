package middleware

import (
	"github.com/MQEnergy/gin-framework/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"time"
)

func RateLimitMiddleware(fillInterval time.Duration, cap, quantum int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(fillInterval, cap, quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.InternalServerException(c, "您点击太快了，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
