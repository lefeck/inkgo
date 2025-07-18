package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"inkgo/authentication"
	"inkgo/repository"
	"inkgo/utils"
	"net/http"
	"strings"
)

// JWTAuth 中间件
func AuthenticationMiddleware(jwt *authentication.JWT, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取请求头中的 Authorization 认证字段信息
		token, _ := getTokenFromAuthorizationHeader(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}
		// 2. 检查token是否被注销
		if revoked, err := jwt.IsTokenRevoked(token); revoked || err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Token has been revoked"})
			return
		}
		// 3. 解析token
		user, err := jwt.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		}
		utils.SetUserInContext(c, user)
		c.Next()
	}
}

func getTokenFromAuthorizationHeader(c *gin.Context) (string, error) {
	// HTTP Bearer： https://www.cnblogs.com/qtiger/p/14868110.html#autoid-3-4-0
	// 在 http 请求头当中添加 Authorization: Bearer (token) 字段完成验证
	// 1. 获取请求头中的 Authorization 认证字段信息
	auth := c.Request.Header.Get("Authorization")
	if auth == "" {
		return "", nil
	}
	token := strings.Fields(auth)
	if len(token) != 2 || strings.ToLower(token[0]) != "bearer" || token[1] == "" {
		return "", fmt.Errorf("Authorization header invaild")
	}
	return token[1], nil
}
