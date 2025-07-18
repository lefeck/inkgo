package repository

import (
	"gorm.io/gorm"
	"inkgo/model"
)

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{
		db: db,
	}
}

// Add 添加用户对某个帖子的点赞
func (l *likeRepository) LikePost(pid, uid uint) error {
	// 检查是否已经点赞
	exists, err := l.IsLiked(uid, pid)
	if err != nil {
		return err
	}
	if exists {
		// 减少帖子点赞数
		if err := l.DecLike(pid); err != nil {
			return err
		}
	}

	// 插入点赞记录
	like := &model.Likes{
		PostID:   pid,
		AuthorID: uid,
	}
	if err := l.db.Create(like).Error; err != nil {
		return err
	}
	// 增加帖子点赞数
	if err := l.IncLike(pid); err != nil {
		return err
	}

	return nil
}

// Delete 删除用户对某个帖子的点赞
func (l *likeRepository) UnLikePost(pid, uid uint) error {
	// 检查是否存在点赞记录
	exists, err := l.IsLiked(uid, pid)
	if err != nil {
		return err
	}
	if !exists {
		// 增加帖子点赞数
		if err := l.IncLike(pid); err != nil {
			return err
		}
	}
	// 删除点赞记录
	like := &model.Likes{}
	err = l.db.Where("post_id = ? and author_id = ?", pid, uid).Delete(like).Error
	if err != nil {
		return err
	}
	// 减少帖子点赞数
	if err := l.DecLike(pid); err != nil {
		return err
	}
	return nil
}

// IncLike 增加帖子的点赞数
func (r *likeRepository) IncLike(pid uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", pid).
		UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *likeRepository) DecLike(pid uint) error {
	return r.db.Model(&model.Post{}).Where("id = ?", pid).
		UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
}

// CountLikes 统计某个帖子的点赞数
func (l *likeRepository) CountLikes(pid uint) (int64, error) {
	var count int64
	like := &model.Likes{}

	if err := l.db.Model(like).Where("post_id = ?", pid).Count(&count).Error; err != nil {
		return count, err
	}
	return count, nil
}

// IsLiked 判断是否已经点赞
func (r *likeRepository) IsLiked(userID uint, postID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Likes{}).Where("author_id = ? AND post_id = ?", userID, postID).Count(&count).Error
	return count > 0, err
}

// 自动创建表结构到db
func (a *likeRepository) Migrate() error {
	return a.db.AutoMigrate(&model.Likes{})
}
