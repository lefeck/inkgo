package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"inkgo/authentication"
	"inkgo/authentication/oauth"
	"inkgo/common/request"
	"inkgo/model"
	"inkgo/service"
	"inkgo/utils"
	"net/http"
	"strconv"
)

type AuthController struct {
	userService  service.UserService
	jwtService   *authentication.JWT
	oauthManager *oauth.OAuthManager
	authService  service.AuthService
}

func NewAuthController(userService service.UserService, jwtService *authentication.JWT, oauthManager *oauth.OAuthManager, auth service.AuthService) Controller {
	return &AuthController{
		userService:  userService,
		jwtService:   jwtService,
		oauthManager: oauthManager,
		authService:  auth,
	}
}

type AuthUser struct {
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	AuthType   string `json:"auth_type"` // OAuth 登录类型
	AuthCode   string `json:"auth_code"` // OAuth 授权码
}

// Login 用户登录
func (auth *AuthController) Login(c *gin.Context) {
	var loginReq request.LoginRequest
	// 解析请求体中的用户登录信息
	if err := c.ShouldBind(&loginReq); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	var user model.User

	if !oauth.IsEmptyAuthType(loginReq.AuthType) && loginReq.Identifier == "" {
		// 如果是 OAuth 登录
		provider, err := auth.oauthManager.GetAuthProvider(loginReq.AuthCode)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err)
			return
		}
		// 获取 OAuth 提供商的访问令牌
		authToken, err := provider.GetAccessToken(loginReq.AuthCode)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, err)
			return
		}
		// 获取用户信息
		userInfo, err := provider.GetUserInfo(authToken)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, err)
			return
		}

		// 创建 OAuth 用户
		user, err := auth.userService.CreateOAuthUser(userInfo.User())
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, err)
			return
		}
		// 生成 JWT 令牌
		token, err := auth.jwtService.GenerateToken(user)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, errors.New("生成token失败"))
			return
		}

		utils.Success(c, gin.H{
			"id":    user.ID,
			"name":  user.UserName,
			"role":  user.Role,
			"token": token,
		})

		return
	} else if oauth.IsEmptyAuthType(loginReq.AuthType) && loginReq.Identifier != "" {
		// 如果是普通登录, 使用用户名或邮箱或手机号登录
		user := loginReq.GetUser()
		if user, err := auth.userService.LoginByPassword(loginReq.Identifier, loginReq.Password, user); err != nil {
			utils.Error(c, http.StatusUnauthorized, err)

		} else {
			token, err := auth.jwtService.GenerateToken(user)
			if err != nil {
				utils.Error(c, http.StatusInternalServerError, errors.New("生成token失败"))
				return
			}
			utils.Success(c, gin.H{
				"id":    user.ID,
				"name":  user.UserName,
				"role":  user.Role,
				"token": token,
			})
		}
	}

	// 设置 Cookie 中的 JWT 令牌
	c.SetCookie("jwt_token", "", 0, "/", "localhost", false, true)
	// 设置 JWT 令牌到响应头
	c.Header("Authorization", fmt.Sprintf("Bearer %s", c.GetHeader("Authorization")))
	// 设置用户 ID 到上下文
	c.Set("user_id", user.ID)

}

// 注册新用户
func (auth *AuthController) Register(c *gin.Context) {
	var user request.RegisterRequest
	if err := c.ShouldBind(&user); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	// 创建新用户
	newUser, err := auth.userService.Create(user.GetUser())
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, newUser)
}

// Logout 用户登出
func (auth *AuthController) Logout(c *gin.Context) {
	// 从请求头中获取 JWT 令牌
	token := c.GetHeader("Authorization")
	if token == "" {
		utils.Error(c, http.StatusBadRequest, errors.New("缺少 JWT 令牌"))
		return
	}

	// 清除用户的 JWT 令牌
	if err := auth.jwtService.RevokeToken(token); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	// 清除用户的登录状态
	c.Set("user_id", nil)
	c.Set("claims", nil)
	// 清除 Cookie 中的 JWT 令牌
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	// 返回成功响应
	utils.Success(c, nil)
}

func (auth *AuthController) ResetPasswordByToken(c *gin.Context) {
	var resetReq request.ResetPasswordRequest
	if err := c.ShouldBind(&resetReq); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	token, err := c.Cookie("reset_token")
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	email, err := auth.authService.GetEmailByResetToken(token)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	// 通过邮箱查看用户
	user, err := auth.userService.FindByEmail(email)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	if err := auth.userService.UpdatePassword(strconv.Itoa(int(user.ID)), resetReq.Password); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	//_ = u.emailService.DeleteToken(token)
	c.SetCookie("reset_token", "", -1, "/", "localhost", false, true)

	utils.Success(c, gin.H{"message": "密码重置成功"})

}

// SendCode 发送邮箱验证码
func (auth *AuthController) SendResetCode(c *gin.Context) {
	var email request.EmailRequest

	if err := c.ShouldBind(&email); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	// 检查邮箱是否存在
	user, err := auth.userService.FindByEmail(email.To)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	user.Email = email.To

	if err := auth.authService.SendEmailVerificationCode(email.From, user.Email, email.AuthCode); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, gin.H{"message": "验证码已发送"})
}

// VerifyCode 验证邮箱验证码
func (auth *AuthController) VerifyResetCode(c *gin.Context) {
	var email request.VerifyCodeRequest

	if err := c.ShouldBind(&email); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	token, err := auth.authService.VerifyEmailCode(email.Email, email.Code)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	// ✅ 写入 Cookie，防止前端手动传
	c.SetCookie("reset_token", token, 300, "/", "127.0.0.1", false, true)

	utils.Success(c, gin.H{"message": "验证码验证成功", "reset_token": token})
}

// 注销用户,  清除用户相关的所有数据
func (auth *AuthController) Revoke(c *gin.Context) {
	// 清除用户的登录状态
	c.Set("user_id", nil)
	c.Set("claims", nil)
	// 清除 Cookie 中的 JWT 令牌
	c.SetCookie("jwt_token", "", -1, "/", "localhost", false, true)
	// 返回成功响应
	utils.Success(c, nil)
}

func (auth *AuthController) RegisterRoute(api *gin.RouterGroup) {
	api.POST("/auth/token", auth.Login)
	api.POST("/auth/user", auth.Register)
	api.POST("/auth/logout", auth.Logout)
	api.POST("/auth/send", auth.SendResetCode)                  // 发送验证码
	api.POST("/auth/verify", auth.VerifyResetCode)              // 验证验证码
	api.POST("/user/reset_password", auth.ResetPasswordByToken) // 重置密码
}

func (auth *AuthController) Name() string {
	return "auth"
}
