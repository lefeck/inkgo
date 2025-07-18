package service

import (
	"inkgo/model"
	"inkgo/repository"
	"strconv"
)

type categoryService struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}

// List 列出所有分类
func (c *categoryService) List(page, pageSize int) ([]model.Category, int64, error) {
	categories, total, err := c.categoryRepository.List(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return categories, total, nil
}

// Get 获取指定分类
func (c *categoryService) Get(id string) (*model.Category, error) {
	cid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	category, err := c.categoryRepository.Get(uint(cid))
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (c *categoryService) Create(category *model.Category) (*model.Category, error) {
	return c.categoryRepository.Create(category)
}

func (c *categoryService) Delete(id string) error {
	cid, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return c.categoryRepository.Delete(uint(cid))
}

func (c *categoryService) Update(id string, category *model.Category) (*model.Category, error) {
	cid, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	category.ID = uint(cid)
	return c.categoryRepository.Update(category)
}
