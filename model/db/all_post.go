package main

import (
	"gorm.io/gorm"
	"time"
)

const (
	UserAssociation         = "Users"
	UserAuthInfoAssociation = "AuthInfos"
)

type User struct {
	gorm.Model
	UserName string `json:"user_name" gorm:"size:64;uniqueIndex;not null"` // 用户名（唯一）
	//其它字段按需添加 json 标签，敏感字段（密码）用 json:"-" 隐藏
	Password string `json:"-" gorm:"size:64"`                             // 密码（存储加密后的值）
	Mobile   string `json:"mobile,omitempty" gorm:"size:20;uniqueIndex"`  // 手机号（唯一）
	Email    string `json:"email,omitempty" gorm:"size:128;;uniqueIndex"` // 邮箱 (唯一)
	Avatar   string `json:"avatar" gorm:"size:512;"`                      // 头像
	Intro    string `json:"intro" gorm:"size:256"`                        // 简介/签名
	//为后续进一步实现基于角色的权限系统（比如 admin 可以看所有收藏夹、user 只能看自己的），
	Role        Role      `gorm:"type:varchar(20);not null" json:"role"`  // 角色：user/admin 等
	Status      string    `gorm:"size:16;default:'active'" json:"status"` // 状态：active/banned/deleted 等
	LastLoginAt time.Time `json:"last_login_at"`                          // 上次登录时间
}

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

type PostState string

const (
	PostDraft     PostState = "draft"     // 草稿
	PostPublished PostState = "published" // 已发布
	PostArchived  PostState = "archived"  // 已归档（可选）
)

type Repost struct {
	ID           uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	Url          string `json:"url" gorm:"type:varchar(200);not null"` // 转发的文章链接
	IsAuthorized bool   `json:"is_authorized" gorm:"default:false"`    // 是否授权转载，默认 false
	RepostID     uint   `json:"repost_id"`                             // 被转发的文章ID
	Repost       Post   `json:"-" gorm:"foreignKey:RepostID"`          // 关联被转发的 Post 实体
}

// Follow关注
type Follow struct {
	gorm.Model
	UserID     uint `json:"user_id" gorm:"uniqueIndex:user_follow"`
	User       User `json:"-" gorm:"foreignKey:UserID"` // // 关注者
	FollowedID uint `json:"followed_id" gorm:"uniqueIndex:user_follow"`
	Followed   User `json:"-" gorm:"foreignKey:FollowedID"` // 被关注者
}

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

// 点赞
type Likes struct {
	gorm.Model
	AuthorID uint `json:"author_id"` // 外键
	Author   User `json:"author" gorm:"foreignKey:AuthorID"`
	PostID   uint `json:"post_id"` // 外键
	Post     Post `json:"post" gorm:"foreignKey:PostID"`
}

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// Share 分享文章
type Share struct {
	gorm.Model
	AuthorID  uint  `json:"author_id"`                         // 作者ID
	Author    User  `json:"author" gorm:"foreignKey:AuthorID"` // 关联的用户实体
	PostID    uint  `json:"post_id"`                           // 文章ID
	Post      Post  `json:"post" gorm:"foreignKey:PostID"`     // 关联的文章实体
	CreatedAt int64 `json:"created_at" gorm:"autoCreateTime"`  // 分享时间
}

// 标签
type Tag struct {
	gorm.Model
	Name  string `json:"name" gorm:"size:256;not null;unique"` // 标签名称
	Posts []Post `json:"posts" gorm:"many2many:tag_posts"`     // 标签和帖子之间的多对多关系
}

type Category struct {
	gorm.Model
	Name        string `json:"name" gorm:"size:50;not null;unique"`
	Description string `json:"description" gorm:"type:varchar(256);"` // 分类描述
	Image       string `json:"image" gorm:"type:varchar(200)"`
	Posts       []Post `json:"posts" gorm:"many2many:category_posts"` // Category 和 Post 之间的多对多关系
}

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

// Activity足迹, 记录用户浏览过的文章
type Activity struct {
	ID     uint   `json:"id" gorm:"autoIncrement;primaryKey"`
	PostID uint   `json:"post_id" gorm:"index"`       // 外键，关联 Post
	Post   Post   `json:"-" gorm:"foreignKey:PostID"` // 外键，关联 Post
	UserID uint   `json:"user_id" gorm:"index"`
	User   User   `json:"-" gorm:"foreignKey:UserID"`  // 外键，关联 User
	IsOpen bool   `json:"is_Open" gorm:"default:true"` // 是否公开足迹
	Device string `json:"device" gorm:"size:64"`       // 可选：记录访问设备
	IP     string `json:"ip" gorm:"size:64"`           // 可选：IP 地址
	BaseModel
}
