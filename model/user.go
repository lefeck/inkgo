package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

//type BaseModel struct {
//	CreatedAt time.Time      `json:"createdAt"`
//	UpdatedAt time.Time      `json:"updatedAt"`
//	DeletedAt gorm.DeletedAt `json:"-"` // soft delete
//}

const (
	UserAssociation         = "Users"
	UserAuthInfoAssociation = "AuthInfos"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type User struct {
	gorm.Model
	UserName string `json:"user_name" gorm:"size:64;uniqueIndex;not null"` // 用户名（唯一）
	//其它字段按需添加 json 标签，敏感字段（密码）用 json:"-" 隐藏
	Password string `json:"-" gorm:"size:64"`                             // 密码（存储加密后的值）
	Mobile   string `json:"mobile,omitempty" gorm:"size:20;uniqueIndex"`  // 手机号（唯一）
	Email    string `json:"email,omitempty" gorm:"size:128;;uniqueIndex"` // 邮箱 (唯一)
	Avatar   string `json:"avatar" gorm:"size:512;"`                      // 头像
	Intro    string `json:"intro" gorm:"size:256"`                        // 简介/签名
	//为后续进一步实现基于角色的权限系统（比如 admin 可以看所有收藏夹、user 只能看自己的），
	Role          Role           `gorm:"type:varchar(20);not null;index;default:user" json:"role"` // 角色：user/admin 等
	Status        string         `gorm:"size:16;default:'active'" json:"status"`                   // 状态：active/banned/deleted 等
	LastLoginAt   *time.Time     `json:"last_login_at" gorm:"default:null"`                        // 上次登录时间
	OAuthAccounts []OAuthAccount `json:"oauth_accounts" gorm:"foreignKey:UserID;references:ID"`    // 用户的认证信息
}

//	Favorites []Favorite `gorm:"foreignKey:UserID"`

func (*User) TableName() string {
	return "user"
}

func (u *User) CacheKey() string {
	return u.TableName() + ":id"
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
