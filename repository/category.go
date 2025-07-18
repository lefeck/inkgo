package repository

import (
	"errors"
	"gorm.io/gorm"
	"inkgo/model"
)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

// Get 获取指定分类
func (c *categoryRepository) Get(cid uint) (*model.Category, error) {
	var category model.Category
	if err := c.db.Preload(model.PostAssociation).First(&category, cid).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// Create 创建一个新的分类
func (c *categoryRepository) Create(category *model.Category) (*model.Category, error) {
	// 检查分类名称是否已存在

	if err := c.db.Where("name = ?", category.Name).First(&category).Error; err == nil {
		return nil, gorm.ErrDuplicatedKey // 分类名称已存在
	}

	if err := c.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

// Delete 删除一个分类
func (c *categoryRepository) Delete(cid uint) error {
	catgory := &model.Category{Model: gorm.Model{ID: cid}}
	// 检查分类的id 是否存在
	if err := c.db.First(&catgory).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound // 分类不存在
		}
		return err // 其他错误
	}

	if err := c.db.Delete(catgory).Error; err != nil {
		return err
	}
	return nil
}

// Update 更新一个分类
func (c *categoryRepository) Update(category *model.Category) (*model.Category, error) {
	var existingCategory model.Category
	if result := c.db.First(&existingCategory, category.ID); result.RowsAffected == 0 {
		return nil, errors.New("category name is not exist")
	}
	if err := c.db.Model(&existingCategory).Updates(&category).Error; err != nil {
		return nil, errors.New("category name is already exist")
	}
	return &existingCategory, nil
}

// List 列出所有分类, 分页设计
func (c *categoryRepository) List(page, pageSize int) ([]model.Category, int64, error) {
	categories := make([]model.Category, 0)
	var total int64

	db := c.db.Model(&model.Category{})
	// 计算总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	// 分页查询
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

// 自动创建表结构到db
func (a *categoryRepository) Migrate() error {
	return a.db.AutoMigrate(&model.Category{})
}
