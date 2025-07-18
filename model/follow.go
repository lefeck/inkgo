package model

import "gorm.io/gorm"

// Follow关注
type Follow struct {
	gorm.Model
	UserID     uint `json:"user_id" gorm:"uniqueIndex:user_follow"`
	User       User `json:"-" gorm:"foreignKey:UserID"` // // 关注者
	FollowedID uint `json:"followed_id" gorm:"uniqueIndex:user_follow"`
	Followed   User `json:"-" gorm:"foreignKey:FollowedID"` // 被关注者
}

func (*Follow) TableName() string {
	return "follow"
}
