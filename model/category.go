package model

import "gorm.io/gorm"

// 分类
type Category struct {
	gorm.Model
	Name        string `json:"name" gorm:"size:50;not null;unique"`
	Description string `json:"description" gorm:"type:varchar(256);"` // 分类描述
	Image       string `json:"image" gorm:"type:varchar(200)"`
	Posts       []Post `json:"posts" gorm:"many2many:category_posts"` // Category 和 Post 之间的多对多关系
}

func (*Category) TableName() string {
	return "category"
}
