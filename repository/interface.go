package repository

import (
	"context"
	"inkgo/model"
)

// 工厂模式接口
type Repository interface {
	User() UserRepository
	Post() PostRepository
	Category() CategoryRepository
	Comment() CommentRepository
	Tag() TagRepository
	Like() LikeRepository
	Favorite() FavoriteRepository
	Follow() FollowRepository
	Activity() ActivityRepository
	Auth() AuthRepository
	Token() TokenRepository
	Close() error
	Ping(ctx context.Context) error
	Migrant
}

type Migrant interface {
	Migrate() error
}

// user实现的接口
type UserRepository interface {
	GetUserByID(uint) (*model.User, error)
	// 根据用户名、邮箱等查询（登录验证）
	FindByUserName(username string) (*model.User, error)
	FindByMobile(mobile string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	List(pageSize int, pageNum int) ([]model.User, int64, error)
	// 创建用户（注册）
	Create(user *model.User) (*model.User, error)
	// 注销用户
	Deactivate(userID uint) error
	// 登录（支持账号密码或 OAuth）
	FindByOauth(provider string, providerID string) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	// 更新用户密码
	UpdatePassword(userID uint, hashedPassword string) error
	//UpdateRole(userID uint, role string) error
	// 删除用户
	Delete(id uint) error
	// 查询是否存在用户名 / 邮箱（注册校验用）
	IsUsernameTaken(username string) (bool, error)
	IsEmailTaken(email string) (bool, error)

	Migrate() error
}

// Post 接口
type PostRepository interface {
	// GetPostByID 根据ID获取文章
	GetPostByID(uint) (*model.Post, error)
	// GetPostByName 根据文章名称获取文章
	GetPostByName(string) (*model.Post, error)
	// ListHasPublished 列出所有已发布的文章
	ListHasPublished(page int, pageSize int) ([]model.Post, int64, error)
	//ListDrafts 列出所有草稿文章
	ListDrafts(page int, pageSize int) ([]model.Post, int64, error)
	// Create 创建一篇文章
	Create(*model.User, *model.Post) (*model.Post, error)
	// Update 更新一篇文章
	Update(*model.Post) (*model.Post, error)
	// UpdateStatus 更新文章状态
	UpdateStatus(id uint, state model.PostState) (*model.Post, error)

	// Delete 删除一篇文章
	Delete(uint) error
	// IncView 增加文章的浏览量
	IncView(id uint) error
	// SortByViewCountDesc 按浏览量降序排序文章
	SortByViewCountDesc(page int, pageSize int) ([]model.Post, int64, error)

	// ListHotPosts  列出热门文章
	ListHotPosts(limit int) ([]model.Post, error)

	// ListRecentPosts 列出最近的文章
	ListRecentPosts(limit int) ([]model.Post, error)

	Migrate() error
}

type CategoryRepository interface {
	Delete(cid uint) error
	Create(category *model.Category) (*model.Category, error)
	Update(category *model.Category) (*model.Category, error)
	List(page, pageSize int) ([]model.Category, int64, error)
	Get(cid uint) (*model.Category, error)
	Migrate() error
}

type TagRepository interface {
	// GetTagsByPost 获取标签下的文章
	GetTagsByPost(post *model.Post) ([]model.Tag, error)
	Get(tid uint) (model.Tag, error)
	Create(tag *model.Tag) (*model.Tag, error)
	Delete(id uint) error
	List() ([]model.Tag, error)
	Update(tag *model.Tag) (*model.Tag, error)
	Migrate() error
}

type CommentRepository interface {
	Add(comment *model.Comment) (*model.Comment, error)
	Delete(id string) error
	List(aid string) ([]model.Comment, error)
	Migrate() error
}

type LikeRepository interface {
	LikePost(pid, uid uint) error
	UnLikePost(pid, uid uint) error
	CountLikes(pid uint) (int64, error)
	IsLiked(pid, uid uint) (bool, error)
	Migrate() error
}

// FavoriteRepository 是用户收藏相关的接口
type FavoriteRepository interface {
	// 收藏夹管理
	// CreateFavorite 创建一个新的收藏夹
	CreateFavorite(userID uint, favorite *model.Favorite) (*model.Favorite, error)
	// DeleteFavorite 删除收藏夹
	DeleteFavorite(userID uint, favoriteID uint) error
	// UpdateFavorite 更新收藏夹信息
	UpdateFavorite(userID uint, favorite *model.Favorite) (*model.Favorite, error)
	// GetFavoriteByID 根据收藏夹ID获取收藏夹信息
	GetFavoriteByID(userID uint, favoriteID uint) (*model.Favorite, error)
	// GetUserFavorites 获取某个用户创建的所有收藏夹列表
	GetUserFavorites(userID uint, page, pageSize int) ([]model.Favorite, int64, error)

	// 收藏操作（核心）
	// AddPostToFavorite 添加文章到收藏夹
	AddPostToFavorite(userID uint, postID uint, favoriteID uint) error
	// RemovePostFromFavorite 移除文章从收藏
	RemovePostFromFavorite(userID uint, postID uint, favoriteID uint) error
	// IsPostInFavorite 检查文章是否在收藏夹中
	IsPostInFavorite(userID uint, favoriteID, postID uint) (bool, error)

	// 收藏夹内容查询
	ListPostsInFavorite(userID uint, favoriteID uint) ([]model.Post, error) // 列出收藏夹中的文章
	CountPostsInFavorite(userID uint, favoriteID uint) (int64, error)       // 获取收藏夹中的文章数量

	// attach functions
	//	ListFavoritesOfPost(postID uint)	查询某篇文章被哪些收藏夹收藏（例如：后台统计）
	//BatchAddPostsToFavorite(favoriteID uint, postIDs []uint)	批量收藏多个文章（性能更优）
	//RemoveAllPostsFromFavorite(favoriteID uint)	清空收藏夹
	//FavoriteExists(favoriteID uint) (bool, error)	判断收藏夹是否存在（用于权限验证或幂等）
	Migrate() error
}

// Follow 是用户关注相关的接口
type FollowRepository interface {
	// 获取用户关注的用户数
	GetFollowedCount(userID uint) (int64, error)
	// 获取用户的粉丝数
	GetFollowerCount(userID uint) (int64, error)
	// 获取用户关注的用户列表
	GetFollowedList(userID uint, page, pageSize int) ([]model.User, error)
	// 获取用户的粉丝列表
	GetFollowerList(userID uint, page, pageSize int) ([]model.User, error)
	Migrate() error
}

type ActivityRepository interface {
	// 获取当前用户的文章, 作为历史记录方便查看
	GetUserActivity(userID uint, page, pageSize int) ([]model.Activity, error)
	Migrate() error
}
