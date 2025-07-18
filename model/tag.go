package model

import "gorm.io/gorm"

// 标签
type Tag struct {
	gorm.Model
	Name  string `json:"name" gorm:"size:256;not null;unique"` // 标签名称
	Posts []Post `json:"posts" gorm:"many2many:tag_posts"`     // 标签和帖子之间的多对多关系
}

func (*Tag) TableName() string {
	return "tag"
}
