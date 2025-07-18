package model

import "gorm.io/gorm"

// 点赞
//type Likes struct {
//	BaseModel
//	AuthorID uint `json:"author_id"` // 外键
//	Author   User `json:"author" gorm:"foreignKey:AuthorID"`
//	PostID   uint `json:"post_id"` // 外键
//	Post     Post `json:"post" gorm:"foreignKey:PostID"`
//}

type Likes struct {
	gorm.Model
	AuthorID uint `json:"author_id" gorm:"index;uniqueIndex:idx_author_post"`
	Author   User `json:"author" gorm:"foreignKey:AuthorID"`
	PostID   uint `json:"post_id" gorm:"index;uniqueIndex:idx_author_post"`
	Post     Post `json:"post" gorm:"foreignKey:PostID"`
}

/*
这个示例说明了如何在Post和Like结构体中使用gorm标签：

* ID字段使用gorm:"autoIncrement;primaryKey"标签，表示这是一个自动增长的主键。
* UserID和PostID字段使用gorm:"index"标签来创建索引，这将提高查询的速度。
* 使用foreignkey标签来表示外键关系。例如，在User结构体中，Likes字段表示一个User可以对多个Post进行点赞。

与此同时，在Post结构体中，Likes字段表示一个Post可以被多个User点赞。在这两种情况下，我们需要为这些关联关系指定正确的外键。
*/

func (*Likes) TableName() string {
	return "likes"
}
