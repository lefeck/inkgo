package repository

import (
	"gorm.io/gorm"
	"inkgo/model"
)

type ShareRepository interface {
	Create(share *model.Share) (*model.Share, error)
	GetShareByID(id string) (*model.Share, error)
}

type shareRepository struct {
	db *gorm.DB
}

func NewShareRepository(db *gorm.DB) ShareRepository {
	return &shareRepository{
		db: db,
	}
}

// 创建一个分享
func (s *shareRepository) Create(share *model.Share) (*model.Share, error) {
	shares := model.Share{
		AuthorID: share.AuthorID,
		PostID:   share.PostID,
	}
	if err := s.db.Create(shares).Error; err != nil {
		return nil, err
	}
	return share, nil
}

// 获取分享
func (s *shareRepository) GetShareByID(shareID string) (*model.Share, error) {
	share := new(model.Share)
	err := s.db.Preload(model.AuthorAssociation).Preload(model.PostAssociation).Preload("Post.Author").
		First(&share, shareID).Error
	return share, err
}

// automatically create the table if it does not exist
func (s *shareRepository) Migrate() error {
	return s.db.AutoMigrate(&model.Share{})
}
