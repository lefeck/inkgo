package model

import (
	"gorm.io/gorm"
)

const (
	AuthorAssociation   = "Author"
	LikesAssociation    = "Likes"
	TagsAssociation     = "Tags"
	CategoryAssociation = "Categories"
	PostAssociation     = "Posts"
	CommentsAssociation = "Comments"
)

//type PostState int
//
//const (
//	Published PostState = iota // 已发布
//	Draft                      // 草稿
//)
//
//var PostStates = map[PostState]string{
//	Published: "published", // 已发布
//	Draft:     "draft",     // 草稿
//}
//
//func (ps PostState) String() string {
//	if val, ok := PostStates[ps]; ok {
//		return val
//	}
//	return "unknown"
//}
//
//// jsonmarshon
//func (ps PostState) MarshalJSON() ([]byte, error) {
//	return json.Marshal(ps.String())
//}

type PostState string

const (
	PostDraft     PostState = "draft"     // 草稿
	PostPublished PostState = "published" // 已发布
	PostArchived  PostState = "archived"  // 已归档（可选）
)

// Post 文章模型
type Post struct {
	gorm.Model
	Title      string     `json:"title" gorm:"type:varchar(100);not null"`       // 文章标题
	Content    string     `json:"content" gorm:"type:text;not null"`             // 文章内容
	Cover      string     `json:"cover" gorm:"not null"`                         // 封面图片
	AuthorID   uint       `json:"author_id"`                                     // 外键
	Author     User       `json:"author" gorm:"foreignKey:AuthorID"`             // 关联 User 实体
	Tags       []Tag      `json:"tags" gorm:"many2many:tag_posts"`               // Post 和 Tag 之间的多对多关系
	Categories []Category `json:"categories" gorm:"many2many:category_posts"`    // Post 和 Category 之间的多对多关系
	Comments   []Comment  `json:"comments,omitempty"`                            // Post 和 Comment 之间的一对多关系
	ViewCount  uint       `json:"view_count" gorm:"type:uint"`                   // 文章阅读量
	LikeCount  uint       `json:"like_count" gorm:"type:uint"`                   // 点赞数
	UserLiked  bool       `json:"user_liked" gorm:"-"`                           // 用户是否点赞
	Original   bool       `json:"original"`                                      // 是否原创，true:原创，false:转载
	State      PostState  `json:"state" gorm:"type:varchar(20);default:'draft'"` // 文章的发布状态，已发布/草稿
}

//Repost      *Post      `json:"repost,omitempty" gorm:"foreignKey:RepostID"`    // 转载的文章

func (p *Post) TableName() string {
	return "post"
}
