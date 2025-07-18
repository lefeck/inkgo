package request

import (
	"github.com/asaskevich/govalidator"
	"gorm.io/gorm"
	"inkgo/model"
	"unicode"
)

type LoginRequest struct {
	//Identifier string `json:"identifier" binding:"required"`
	//Password   string `json:"password" form:"password" binding:"required,min=6,max=20"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
	AuthType   string `json:"auth_type"`
	AuthCode   string `json:"auth_code"`
}

type IdentifierType string

const (
	ByUserName IdentifierType = "username"
	ByEmail    IdentifierType = "email"
	ByMobile   IdentifierType = "mobile"
)

func Detect(identifier string) IdentifierType {
	switch {
	case govalidator.IsEmail(identifier):
		return ByEmail
	case isMobile(identifier):
		return ByMobile
	default:
		return ByUserName
	}
}

func isMobile(s string) bool {
	if len(s) != 11 {
		return false
	}
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func (l *LoginRequest) GetUser() *model.User {
	user := &model.User{Password: l.Password}
	switch Detect(l.Identifier) {
	case ByEmail:
		user.Email = l.Identifier
	case ByMobile:
		user.Mobile = l.Identifier
	default:
		user.UserName = l.Identifier
	}
	return user
}

// Register 用于用户注册
type RegisterRequest struct {
	UserName        string `json:"user_name" form:"user_name" binding:"required"`
	Email           string `json:"email" form:"email" binding:"required"`                                        // 邮箱格式校验
	Mobile          string `json:"mobile" form:"mobile" binding:"required,len=11"`                               // 手机号（可选，长度11位）
	Password        string `json:"password" form:"password" binding:"required,min=6,max=20"`                     //密码长度6-20位
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=Password"` // 确认密码，必须与密码相同
}

func (r *RegisterRequest) GetUser() *model.User {
	return &model.User{
		UserName: r.UserName,
		Password: r.Password,
		Email:    r.Email,
		Mobile:   r.Mobile,
	}
}

type ResetPasswordRequest struct {
	Password        string `json:"password" form:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=Password"`
}

type UpdateUser struct {
	UserName string `json:"user_name" form:"user_name"` // 用户名（唯一）
	Email    string `json:"email" form:"email"`         // 邮箱格式校验
	Avatar   string `json:"avatar" form:"avatar"`       // 头像
	Intro    string `json:"intro" form:"intro"`         // 简介/签名
}

func (u *UpdateUser) GetUser(uid uint) *model.User {
	return &model.User{
		Model: gorm.Model{
			ID: uid,
		},
		UserName: u.UserName,
		Email:    u.Email,
		Avatar:   u.Avatar,
		Intro:    u.Intro,
	}
}

// Email
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}

type EmailRequest struct {
	From     string `json:"from" binding:"required"`      // 发件人
	To       string `json:"to" binding:"required"`        // 收件人
	AuthCode string `json:"auth_code" binding:"required"` // 验证码
}

type VerifyCodeRequest struct {
	Email string `json:"email"binding:"required,email"` // 邮箱地址
	Code  string `json:"code" binding:"required"`       // 验证码
}
