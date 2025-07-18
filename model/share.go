package model

import "gorm.io/gorm"

// Share 分享文章
type Share struct {
	gorm.Model
	AuthorID  uint  `json:"author_id"`                         // 作者ID
	Author    User  `json:"author" gorm:"foreignKey:AuthorID"` // 关联的用户实体
	PostID    uint  `json:"post_id"`                           // 文章ID
	Post      Post  `json:"post" gorm:"foreignKey:PostID"`     // 关联的文章实体
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`  // 分享时间
}

func (*Share) TableName() string {
	return "share"
}
