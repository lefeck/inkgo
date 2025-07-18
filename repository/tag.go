package repository

import (
	"gorm.io/gorm"
	"inkgo/model"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{
		db: db,
	}
}

func (t *tagRepository) Get(tid uint) (model.Tag, error) {
	// 检查标签是否存在
	var tag model.Tag
	if err := t.db.Preload(model.PostAssociation).First(&tag, tid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return model.Tag{}, gorm.ErrRecordNotFound // 标签不存在
		}
		return model.Tag{}, err // 其他错误
	}
	return tag, nil
}

func (t *tagRepository) GetTagsByPost(Post *model.Post) ([]model.Tag, error) {
	tags := make([]model.Tag, 0)
	// 使用 GORM 的 Association 方法来获取 Post 关联的 Tags
	// 这里的 model.TagsAssociation 是 Post 模型中定义的关联标签的字段名
	err := t.db.Model(Post).Association(model.TagsAssociation).Find(&tags)
	return tags, err
}

func (t *tagRepository) List() ([]model.Tag, error) {
	var tags []model.Tag
	if err := t.db.Preload(model.PostAssociation).Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (t *tagRepository) Delete(tid uint) error {
	// 检查标签是否存在
	var tag model.Tag
	if err := t.db.First(&tag, tid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return gorm.ErrRecordNotFound // 标签不存在
		}
		return err // 其他错误
	}
	return t.db.Delete(&model.Tag{}, tid).Error
}

func (t *tagRepository) Create(tag *model.Tag) (*model.Tag, error) {
	// 检查标签名称是否已存在
	var existingTag model.Tag
	if err := t.db.Where("name = ?", tag.Name).First(&existingTag).Error; err == nil {
		return nil, gorm.ErrDuplicatedKey // 标签名称已存在
	}

	if err := t.db.Create(tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (t *tagRepository) Update(tag *model.Tag) (*model.Tag, error) {
	// 检查标签是否存在
	var existingTag model.Tag
	if err := t.db.First(&existingTag, tag.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound // 标签不存在
		}
		return nil, err // 其他错误
	}
	// 更新标签
	if err := t.db.Updates(tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

// 自动创建表结构到db
func (a *tagRepository) Migrate() error {
	return a.db.AutoMigrate(&model.Tag{})
}
