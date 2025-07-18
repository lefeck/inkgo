package repository

import (
	"errors"
	"gorm.io/gorm"
	"inkgo/model"
)

type favoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) FavoriteRepository {
	return &favoriteRepository{
		db: db,
	}
}

// Create 创建一个新的收藏夹
func (f *favoriteRepository) CreateFavorite(userID uint, favorite *model.Favorite) (*model.Favorite, error) {
	// 检查收藏夹名称是否已存在
	var existingFavorite model.Favorite
	// 这里假设收藏夹名称是唯一的
	if err := f.db.Where("name = ? AND user_id = ?", favorite.Name, userID).First(&existingFavorite).Error; err == nil {
		// 收藏夹已存在，返回错误
		return nil, gorm.ErrRecordNotFound
	}
	if err := f.db.Create(favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

// Delete 删除一个收藏夹
func (f *favoriteRepository) DeleteFavorite(userID uint, favoriteID uint) error {
	favorite := &model.Favorite{}
	// 检查收藏是否存在
	if err := f.db.First(favorite, favoriteID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return gorm.ErrRecordNotFound // 收藏夹不存在
		}
		return err // 其他错误
	}
	if err := f.db.Where("id = ? AND user_id = ?", favoriteID, userID).Delete(favorite).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新一个收藏夹
func (f *favoriteRepository) UpdateFavorite(userID uint, favorite *model.Favorite) (*model.Favorite, error) {

	// 检查收藏夹名称是否已存在
	var existingFavorite model.Favorite
	if err := f.db.Where("name = ? AND user_id = ? AND id != ?", favorite.Name, userID, favorite.ID).First(&existingFavorite).Error; err == nil {
		// 收藏夹名称已存在，返回错误
		return nil, gorm.ErrRecordNotFound
	}

	if err := f.db.Updates(favorite).Error; err != nil {
		return nil, err
	}
	return favorite, nil
}

// GetFavoriteByID 根据收藏夹ID获取收藏夹信息
func (f *favoriteRepository) GetFavoriteByID(userID uint, favoriteID uint) (*model.Favorite, error) {
	var favorite model.Favorite
	err := f.db.Preload(model.PostAssociation).
		Where("id = ? AND user_id = ?", favoriteID, userID).
		First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

// GetFavoritesByUserID 获取用户的收藏列表
func (f *favoriteRepository) GetFavoritesByUserID(userID uint) ([]model.Favorite, error) {
	var favorites []model.Favorite
	err := f.db.Where("user_id = ?", userID).Preload(model.PostAssociation).Find(&favorites).Error
	if err != nil {
		return nil, err
	}
	return favorites, nil
}

// AddPostToFavorite 添加文章到收藏夹
func (r *favoriteRepository) AddPostToFavorite(userID uint, postID uint, favoriteID uint) error {
	var favorite model.Favorite
	if err := r.db.Where("id = ? AND user_id = ?", favoriteID, userID).First(&favorite).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("收藏夹不存在或不属于当前用户")
		}
		return err
	}

	// 第二步：判断是否已收藏（避免重复）
	var count int64
	err := r.db.Model(&favorite).
		Where("favorite_id = ? AND post_id = ?", favoriteID, postID).
		Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // 已经收藏过，幂等返回
	}

	return r.db.Create(&favorite).Error
}

// RemovePostFromFavorite 从收藏夹中移除文章
func (f *favoriteRepository) RemovePostFromFavorite(userID uint, postID uint, favoriteID uint) error {

	var fav model.Favorite
	if err := f.db.Where("id = ? AND user_id = ?", favoriteID, userID).First(&fav).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("收藏夹不存在或不属于当前用户")
		}
		return err
	}
	var post model.Post
	if err := f.db.Where("id = ?", postID).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("文章不存在")
		}
		return err
	}

	// 使用关联删除方法从收藏夹中移除文章

	return f.db.Delete(&post).Error
}

// GetUserFavorites 获取用户收藏夹的文章
func (f *favoriteRepository) GetUserFavorites(userID uint, page, pageSize int) ([]model.Favorite, int64, error) {
	var favorites []model.Favorite

	// 计算总数
	var total int64
	err := f.db.Model(&model.Favorite{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = f.db.Where("user_id = ?", userID).
		Preload(model.PostAssociation).
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&favorites).Error
	if err != nil {
		return nil, 0, err
	}
	return favorites, total, nil
}

// IsPostInFavorite 检查文章是否在收藏夹中
func (f *favoriteRepository) IsPostInFavorite(userID uint, favoriteID, postID uint) (bool, error) {
	var count int64

	// 检查收藏夹是否存在
	var existingFavorite model.Favorite
	// 这里假设收藏夹名称是唯一的
	if err := f.db.Where("id = ? AND user_id = ?", favoriteID, userID).First(&existingFavorite).Error; err == nil {
		// 收藏夹已存在，返回错误
		return false, gorm.ErrRecordNotFound
	}

	err := f.db.Model(&model.Favorite{}).
		Where("id = ? AND EXISTS (SELECT 1 FROM favorite_posts WHERE favorite_id = ? AND post_id = ?)", favoriteID, favoriteID, postID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 验证收藏夹是否属于当前用户
func (r *favoriteRepository) isFavoriteOwnedByUser(userID, favoriteID uint) (bool, error) {
	var exists bool
	// select count(*) > 0
	err := r.db.Model(&model.Favorite{}).
		Select("count(*) > 0").
		Where("id = ? AND user_id = ?", favoriteID, userID).
		Find(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}

// ListPostsInFavorite 列出收藏夹中的文章
func (f *favoriteRepository) ListPostsInFavorite(userID uint, favoriteID uint) ([]model.Post, error) {

	exists, err := f.isFavoriteOwnedByUser(userID, favoriteID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("收藏夹不存在或不属于当前用户")
	}

	var posts []model.Post
	err = f.db.Model(&model.Favorite{}).
		Where("id = ?", favoriteID).
		Association(model.PostAssociation).
		Find(&posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

// CountPostsInFavorite 计算收藏夹中的文章数量
func (f *favoriteRepository) CountPostsInFavorite(userID uint, favoriteID uint) (int64, error) {
	exists, err := f.isFavoriteOwnedByUser(userID, favoriteID)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, errors.New("收藏夹不存在或不属于当前用户")
	}

	// 第二步：获取收藏夹中的文章数量
	var count int64

	err = f.db.Model(&model.Favorite{}).
		Where("id = ?", favoriteID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Migrate 自动创建表结构到db
func (f *favoriteRepository) Migrate() error {
	return f.db.AutoMigrate(&model.Favorite{})
}
