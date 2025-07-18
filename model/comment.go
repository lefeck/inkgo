package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Content  string    `json:"content" gorm:"size:1024"`
	AuthorID uint      `json:"author_id"`
	Author   User      `json:"author" gorm:"foreignKey:AuthorID"`
	PostID   uint      `json:"postId"`
	Post     Post      `json:"post" gorm:"foreignKey:PostID"`
	ParentID *uint     `json:"parent_id"`
	Parent   *Comment  `json:"parent" gorm:"foreignKey:ParentID"`
	Replies  []Comment `json:"replies" gorm:"foreignKey:ParentID"`
}

/*
在此示例中，我们对Comment结构体进行了以下修改：
添加了ParentID字段，用于存储父级评论的ID。
添加了一个名为Parent的与Comment自身关联的结构体，它表示父级评论。
添加了一个名为Replies的Comment结构体切片，用于存储子评论。
*/

func (*Comment) TableName() string {
	return "comment"
}
