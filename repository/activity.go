package repository

import (
	"gorm.io/gorm"
	
)

// ActivityRepository 定义了活动相关的数据库操作接口
type activityRepository struct {
	db *gorm.DB
}

// NewActivityRepository 实例化一个新的活动仓库
func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{
		db: db,
	}
}

// GetUserActivity 获取当前用户的文章, 作为历史记录方便查看
func (a *activityRepository) GetUserActivity(userID uint, page, pageSize int) ([]model.Activity, error) {
	var activities []model.Activity
	if err := a.db.Where("user_id = ?", userID).Offset((page - 1) * pageSize).Limit(pageSize).Find(&activities).Error; err != nil {
		return nil, err
	}
	return activities, nil
}

func (a *activityRepository) Migrate() error {
	return a.db.AutoMigrate(&model.Activity{})
}
