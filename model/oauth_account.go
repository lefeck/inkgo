package model

import (
	"gorm.io/gorm"
	"time"
)

// OAuthAccount 代表一个用户的第三方OAuth账号信息
type OAuthAccount struct {
	gorm.Model
	UserID       uint      `json:"user_id" gorm:"uniqueIndex:idx_user_provider_id;not null"`              // 用户ID
	User         User      `json:"user" gorm:"foreignKey:UserID;references:ID"`                           // 用户关联
	Provider     string    `gorm:"size:64;not null;uniqueIndex:idx_user_provider_id" json:"provider"`     // github, wechat, etc
	ProviderID   string    `gorm:"size:128;not null;uniqueIndex:idx_user_provider_id" json:"provider_id"` // 第三方平台用户ID
	AccessToken  string    `json:"-" gorm:"size:256"`                                                     // 访问令牌
	RefreshToken string    `json:"-" gorm:"size:256"`                                                     // 刷新令牌
	Expiry       time.Time `json:"-"`                                                                     // 令牌过期时间
	URL          string    `json:"url" gorm:"size:256"`                                                   // 授权的URL
}

func (*OAuthAccount) TableName() string {
	return "oauth_account"
}
