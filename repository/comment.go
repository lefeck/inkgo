package repository

import (
	"gorm.io/gorm"
	"inkgo/model"
)

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (c *commentRepository) Add(comment *model.Comment) (*model.Comment, error) {
	err := c.db.Create(comment).Error
	return comment, err
}

func (c *commentRepository) Delete(id string) error {
	comment := &model.Comment{}
	if err := c.db.Delete(comment, id).Error; err != nil {
		return err
	}
	return nil
}

func (c *commentRepository) List(pid string) ([]model.Comment, error) {
	comments := make([]model.Comment, 0)
	err := c.db.Where("post_id = ?", pid).Find(&comments).Error
	return comments, err
}

// CountOfComments 统计评论数量
func (c *commentRepository) CountOfComments() (int64, error) {
	var count int64
	// select count(*) from comments where deleted_at IS NULL
	if err := c.db.Model(&model.Comment{}).Where("deleted_at IS NULL").Count(&count).Error; err != nil {
		return 0, nil
	}
	return count, nil
}

// 添加子评论
func (c *commentRepository) AddReply(comment *model.Comment) (*model.Comment, error) {
	err := c.db.Create(comment).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

// 自动创建表结构到db
func (a *commentRepository) Migrate() error {
	return a.db.AutoMigrate(&model.Comment{})
}
