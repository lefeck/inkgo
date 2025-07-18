package service

import (
	"inkgo/model"
)

type UserService interface {
	// Login 用户登录
	LoginByPassword(identify, password string, user *model.User) (*model.User, error)
	//Register(user request.RegisterUser) (*model.User, error)
	Deactivate(user *model.User) error
	//Export(data *[]model.User, headerName []string, filename string, c *gin.Context) error

	GetUserByID(string) (*model.User, error)
	// 根据用户名、邮箱等查询（登录验证）
	FindByUserName(username string) (*model.User, error)
	FindByMobile(mobile string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)

	List(pageSize int, page int) ([]model.User, int64, error)
	// 创建用户（注册）
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	// 更新用户密码
	UpdatePassword(userID string, hashedPassword string) error
	//UpdateRole(userID string, role string) error
	// 删除用户
	Delete(id string) error
	// 登录（支持账号密码或 OAuth）
	FindByOauth(provider string, providerID string) (*model.User, error)
	CreateOAuthUser(user *model.User) (*model.User, error)

	// 查询是否存在用户名 / 邮箱（注册校验用）
	IsUsernameTaken(username string) (bool, error)
	IsEmailTaken(email string) (bool, error)
}

type PostService interface {
	GetPostByID(id string) (*model.Post, error)
	GetPostByName(name string) (*model.Post, error)
	HasPublished(page, pageSize int) ([]model.Post, int64, error)
	ListDrafts(page, pageSize int) ([]model.Post, int64, error)
	Create(*model.User, *model.Post) (*model.Post, error)
	Get(user *model.User, id string) (*model.Post, error)
	Update(id string, post *model.Post) (*model.Post, error)
	UpdateStatus(id string, state model.PostState) (*model.Post, error)
	Delete(id string) error
	SortByViewCountDesc(page int, pageSize int) ([]model.Post, int64, error)
	ListHotPosts(limit int) ([]model.Post, error)
	ListRecentPosts(limit int) ([]model.Post, error)
}

type LikeService interface {
	LikePost(user *model.User, pid string) error
	UnLikePost(user *model.User, pid string) error
	CountLikes(id string) (int64, error)
	IsLiked(user *model.User, id string) (bool, error)
}

type TagService interface {
	Get(id string) (*model.Tag, error)
	GetTagsByPost(id string) ([]model.Tag, error)
	List() ([]model.Tag, error)
	Create(tag string) (*model.Tag, error)
	Delete(id string) error
	Update(tag *model.Tag, id string) (*model.Tag, error)
}

type CategoryService interface {
	List(page, pageSize int) ([]model.Category, int64, error)
	Get(id string) (*model.Category, error)
	Create(category *model.Category) (*model.Category, error)
	Delete(id string) error
	Update(id string, category *model.Category) (*model.Category, error)
}

type CommentService interface {
	Add(comment *model.Comment, id string, user *model.User) (*model.Comment, error)
	Delete(id string) error
	List(aid string) ([]model.Comment, error)
}
