package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 第一种方式:
// 主要是给OPTION组件做返回的
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}

// 第二种方式:
//func CORSMiddlewares() gin.HandlerFunc {
//	return cors.New(cors.Config{
//		AllowOriginFunc: func(origin string) bool {
//			return true
//		},
//		AllowMethods:     []string{"POST, GET, OPTIONS, DELETE, PATCH, PUT"},
//		AllowHeaders:     []string{"Origin", "Authorization", "Content-Length", "Content-Type"},
//		ExposeHeaders:    []string{"Content-Length"},
//		AllowCredentials: true,
//		MaxAge:           12 * time.Hour,
//		AllowWebSockets:  true,
//	})
//}

//func Addheader() cors.Config {
//	cors := cors.Config{}
//	cors.AddAllowHeaders("AccessToken", "X-CSRF-Token", "Token", "x-token")
//	return cors
//}
