package model

import "gorm.io/gorm"

// Activity足迹, 记录用户浏览过的文章
type Activity struct {
	gorm.Model
	PostID uint   `json:"post_id" gorm:"index"`       // 外键，关联 Post
	Post   Post   `json:"-" gorm:"foreignKey:PostID"` // 外键，关联 Post
	UserID uint   `json:"user_id" gorm:"index"`
	User   User   `json:"-" gorm:"foreignKey:UserID"`  // 外键，关联 User
	IsOpen bool   `json:"is_Open" gorm:"default:true"` // 是否公开足迹
	Device string `json:"device" gorm:"size:64"`       // 可选：记录访问设备
	IP     string `json:"ip" gorm:"size:64"`           // 可选：IP 地址
}

func (*Activity) TableName() string {
	return "activity"
}
