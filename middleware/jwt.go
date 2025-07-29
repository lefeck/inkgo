package middleware

import (
	"fmt"
	"github.com/casbin/casbin/v2"
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
		// 1. 从请求头中的 Authorization 认证字段, 获取token
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
		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		// 4. 检查用户是否存在
		user, err := userRepo.GetUserByID(claims.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		// 检查用户状态
		if user.Status != "active" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User is frozen or deleted"})
			return
		}
		// 5. 将用户信息存入上下文
		//c.Set("user_id", user.ID)
		//c.Set("user_name", user.UserName)
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

func CasbinMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {

		sub, ok := utils.UserFromContext(c)

		if !ok || sub == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		obj := c.FullPath()
		act := c.Request.Method
		fmt.Println(sub.Role, obj, act)

		ok, err := e.Enforce(sub.Role, obj, act)
		// 还是在这里会失败?为啥
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "无权限"})
			return
		}
		c.Next()
	}
}

//func RequireCasbinPermission(resource string ) gin.HandlerFunc  {
//	return func(c *gin.Context) {
//		currentUser, ok  := utils.UserFromContext(c)
//		if !ok || currentUser == nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
//			return
//		}
//		resourceKey := fmt.Sprintf("%s:*", resource,)
//		if currentUser.Role == model. {
//			resourceKey = fmt.Sprintf("%s:own", resource)
//		}
//		allowed, err := db.Enforcer.Enforce(string(currentUser.Role), resourceKey, c.Request.Method)
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
//			return
//		}
//		if !allowed {
//			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "无权访问该资源"})
//			return
//		}
//		c.Next()
//	}
//}
