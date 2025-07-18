package service

import (
	"gorm.io/gorm"
	"inkgo/model"
	"inkgo/repository"
	"strconv"
)

type tagService struct {
	tagRepository repository.TagRepository
}

func NewTagService(tagRepository repository.TagRepository) TagService {
	return &tagService{
		tagRepository: tagRepository,
	}
}

func (t *tagService) GetTagsByPost(id string) ([]model.Tag, error) {
	tid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	Post := &model.Post{Model: gorm.Model{ID: uint(tid)}}
	return t.tagRepository.GetTagsByPost(Post)
}

// Get 获取指定标签
func (t *tagService) Get(id string) (*model.Tag, error) {
	tid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	tag, err := t.tagRepository.Get(uint(tid))
	return &tag, nil
}

func (t *tagService) Delete(id string) error {
	tid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return t.tagRepository.Delete(uint(tid))
}

func (t *tagService) Create(name string) (*model.Tag, error) {
	return t.tagRepository.Create(&model.Tag{Name: name})
}

func (t *tagService) List() ([]model.Tag, error) {
	return t.tagRepository.List()
}

func (t *tagService) Update(tag *model.Tag, id string) (*model.Tag, error) {
	tid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	tag.ID = uint(tid)
	return t.tagRepository.Update(tag)
}
