package repository

import (
	"gorm.io/gorm"
	"inkgo/model"
)

type followRepository struct {
	db *gorm.DB
}

func NewFollowRepository(db *gorm.DB) FollowRepository {
	return &followRepository{db: db}
}

// GetFollowedCount 获取用户关注的用户数
func (r *followRepository) GetFollowedCount(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Follow{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetFollowerCount 获取用户的粉丝数
func (r *followRepository) GetFollowerCount(userID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Follow{}).Where("followed_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetFollowedList获取用户的关注列表
func (r *followRepository) GetFollowedList(userID uint, page, pageSize int) ([]model.User, error) {
	var follows []model.Follow
	if err := r.db.Where("user_id = ?", userID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&follows).Error; err != nil {
		return nil, err
	}

	var followedUsers []model.User
	for _, follow := range follows {
		var user model.User
		if err := r.db.First(&user, follow.FollowedID).Error; err != nil {
			return nil, err
		}
		followedUsers = append(followedUsers, user)
	}

	return followedUsers, nil
}

// GetFollowerList 获取用户的粉丝列表
func (r *followRepository) GetFollowerList(userID uint, page, pageSize int) ([]model.User, error) {
	var follows []model.Follow
	if err := r.db.Where("followed_id = ?", userID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&follows).Error; err != nil {
		return nil, err
	}

	var followerUsers []model.User
	for _, follow := range follows {
		var user model.User
		if err := r.db.First(&user, follow.UserID).Error; err != nil {
			return nil, err
		}
		followerUsers = append(followerUsers, user)
	}

	return followerUsers, nil
}

func (r *followRepository) Migrate() error {
	return r.db.AutoMigrate(&model.Follow{})
}
