package forms

import "app/model"

type TagForm struct {
	Name string `json:"name" form:"name" binding:"required,min=1,max=10"`
}

func (c *TagForm) GetTag() *model.Tag {
	return &model.Tag{
		Name: c.Name,
	}
}
