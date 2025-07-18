package controller

import (
	"errors"
	"fmt"
	"inkgo/common/request"
	"inkgo/utils"

	"github.com/gin-gonic/gin"
	"inkgo/service"
	"net/http"
	"strconv"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userSerivce service.UserService) *UserController {
	return &UserController{
		userService: userSerivce,
	}
}

// List 列出所有用户
func (u *UserController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	uid := c.Param("id")

	currentUser, ok := utils.UserFromContext(c)
	if !ok || currentUser == nil {
		utils.Error(c, http.StatusBadRequest, errors.New("未登录或 token 无效"))
		return
	}

	id, err := strconv.Atoi(uid)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errors.New("无效的用户ID"))
		return
	}

	// 非管理员只能获取自己
	if uint(id) != currentUser.ID && !utils.IsAdmin(currentUser) {
		utils.Error(c, http.StatusUnauthorized, errors.New("无权查看其他用户信息"))
		return
	}

	users, total, err := u.userService.List(pageSize, page)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessWithPage(c, users, total, page, pageSize)
}

func (u *UserController) GetUserByID(c *gin.Context) {
	uid := c.Param("id")
	currentUser, ok := utils.UserFromContext(c)
	if !ok || currentUser == nil {
		utils.Error(c, http.StatusBadRequest, errors.New("未登录或 token 无效"))
		return
	}
	id, err := strconv.Atoi(uid)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errors.New("无效的用户ID"))
		return
	}

	// 非管理员只能获取自己
	if uint(id) != currentUser.ID && !utils.IsAdmin(currentUser) {
		utils.Error(c, http.StatusUnauthorized, errors.New("无权查看其他用户信息"))
		return
	}

	user, err := u.userService.GetUserByID(uid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, user)
}

// GetUserByName 获取用户信息
func (u *UserController) GetUserByName(c *gin.Context) {
	// 从URL参数中获取用户名
	name := c.Param("user_name")
	currentUser, ok := utils.UserFromContext(c)
	if !ok || currentUser == nil {
		utils.Error(c, http.StatusBadRequest, errors.New("未登录或 token 无效"))
		return
	}

	// 非管理员只能获取自己
	fmt.Println(currentUser.UserName)
	if name != currentUser.UserName && !utils.IsAdmin(currentUser) {
		utils.Error(c, http.StatusUnauthorized, errors.New("无权查看其他用户信息"))
		return
	}

	user, err := u.userService.FindByUserName(name)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, user)
}

// Update 更新用户信息
func (u *UserController) Update(c *gin.Context) {
	uid := c.Param("id")

	currentUser, ok := utils.UserFromContext(c)
	if !ok || currentUser == nil {
		utils.Error(c, http.StatusBadRequest, errors.New("未登录或 token 无效"))
		return
	}
	id, err := strconv.Atoi(uid)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errors.New("无效的用户ID"))
		return
	}

	// 非管理员只能获取自己
	if uint(id) != currentUser.ID && !utils.IsAdmin(currentUser) {
		utils.Error(c, http.StatusUnauthorized, errors.New("无权查看其他用户信息"))
		return
	}

	user := &request.UpdateUser{}
	if err := c.ShouldBind(user); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	updatedUser, err := u.userService.Update(user.GetUser(uint(id)))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, updatedUser)
}

// Delete 删除用户
func (u *UserController) Delete(c *gin.Context) {
	uid := c.Param("id")
	currentUser, ok := utils.UserFromContext(c)
	if !ok || currentUser == nil {
		utils.Error(c, http.StatusBadRequest, errors.New("未登录或 token 无效"))
		return
	}
	id, err := strconv.Atoi(uid)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, errors.New("无效的用户ID"))
		return
	}

	// 非管理员只能获取自己
	if uint(id) != currentUser.ID && !utils.IsAdmin(currentUser) {
		utils.Error(c, http.StatusUnauthorized, errors.New("无权限操作其他用户信息"))
		return
	}

	err = u.userService.Delete(uid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, nil)
}

func (u *UserController) Name() string {
	return "users"
}

func (u *UserController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/users", u.List) // 列出所有用户
	api.GET("/user/:id", u.GetUserByID)
	api.DELETE("/user/:id", u.Delete)
	api.PUT("/user/:id", u.Update)
	api.GET("/user/name/:user_name", u.GetUserByName)

}

//func (u *UserController) Export(c *gin.Context) {
//	var users []model.User
//	fileName := "test.xlsx"
//	headerName := []string{"ID", "Name", "Remark", "Status"}
//	users, err := u.userService.FindAll(users)
//	if err != nil {
//		common.ResponseFailed(c, http.StatusInternalServerError, err)
//		return
//	}
//	err = u.userService.Export(&users, headerName, fileName, c)
//	if err != nil {
//		common.ResponseFailed(c, http.StatusInternalServerError, err)
//		return
//	}
//	common.ResponseSuccess(c, nil)
//}
