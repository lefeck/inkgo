package utils

import (
	"github.com/gin-gonic/gin"
	"inkgo/model"
)

const UserContextKey = "user"

func UserFromContext(c *gin.Context) (*model.User, bool) {
	val, exists := c.Get(UserContextKey)
	if !exists {
		return nil, false
	}
	user, ok := val.(*model.User)
	if !ok {
		return nil, false
	}
	return user, true
}

func SetUserInContext(c *gin.Context, user *model.User) {
	if c == nil || user == nil {
		return
	}
	c.Set(UserContextKey, user)
}

func IsAdmin(user *model.User) bool {
	if user == nil {
		return false
	}
	return user.Role == model.RoleAdmin
}
