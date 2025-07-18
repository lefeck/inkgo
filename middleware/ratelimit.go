package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"inkgo/config"
	"net/http"
)

func RateLimitMiddleware(conf *config.RateLimitsConfigs) gin.HandlerFunc {
	bucket := ratelimit.NewBucketWithQuantum(conf.FillInterval, conf.Cap, conf.Quantum)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.String(http.StatusForbidden, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}
