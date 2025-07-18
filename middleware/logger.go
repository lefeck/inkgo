package middleware

import (
	"inkgo/global"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	hostname, _ = os.Hostname()
)

// GinLogger 接收gin框架默认的日志
func LoggerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		global.Log.Info(path,
			zap.String("hostname", hostname),
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
