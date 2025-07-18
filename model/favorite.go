package model

import (
	"gorm.io/gorm"
)

/*
“收藏夹”功能的设计通常是指用户可以收藏自己喜欢的文章。这个功能本质上是一个多对多（many-to-many）关系：一个用户可以收藏多篇文章，一篇文章也可以被多个用户收藏。

收藏夹（Favorites）：不是一个单独的表，而是一个用户和文章之间的中间表。

关联对象：

用户（User）

文章（Post 或 Article）
*/

// Favorite 收藏文章 用于记录当前用户收藏哪些文章
type Favorite struct {
	gorm.Model
	// Name 字段只用 unique —— 意味着全系统唯一，这可能会导致不同用户不能取相同的收藏夹名字。
	// 更合理的做法： 添加一个联合唯一约束 (user_id, name)：
	Name   string `json:"name" gorm:"size:256;not null;uniqueIndex:idx_user_name"` // 标签名称
	Desc   string `json:"desc" gorm:"type:varchar(256);"`                          // 收藏夹描述
	Public bool   `json:"public" gorm:"default:true;comment:是否公开"`                 // 是否公开收藏夹
	// 为了满足允许一个收藏夹重复收藏同一篇文章
	Posts []Post `json:"posts" gorm:"many2many:favorite_posts"`
	//Posts  []FavoritePost `json:"posts"`                                                   // 收藏和帖子之间的多对多关系
	UserID uint `json:"user_id" gorm:"not null;index;uniqueIndex:idx_user_name"` // 属主
	User   User `json:"user" gorm:"foreignKey:UserID"`
}

/*
默认 GORM 的 many2many 会去重。如果你想允许重复收藏（不常见），就得自己定义中间表结构，不用自动 many2many。为了满足允许一个收藏夹重复收藏同一篇文章
如果你还想记录“每篇文章被收藏的时间”，则要手动定义中间表 favorite_posts: 如下

Posts     []Post    `json:"posts" gorm:"many2many:favorite_posts"`                   // 收藏和帖子之间的多对多关系

*/

func (*Favorite) TableName() string {
	return "favorite"
}
