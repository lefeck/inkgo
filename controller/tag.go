package controller

import (
	"github.com/gin-gonic/gin"
	"inkgo/model"
	"inkgo/service"
	"inkgo/utils"
	"net/http"
)

type TagController struct {
	tagService service.TagService
}

func NewTagController(tagService service.TagService) Controller {
	return &TagController{
		tagService: tagService,
	}
}

func (t *TagController) Get(c *gin.Context) {
	tid := c.Param("id")
	tags, err := t.tagService.Get(tid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, tags)
}

func (t *TagController) Create(c *gin.Context) {
	tag := &model.Tag{}
	if err := c.ShouldBindJSON(tag); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	tag, err := t.tagService.Create(tag.Name)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, tag)
}

func (t *TagController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := t.tagService.Delete(id); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, nil)
}

func (t *TagController) Update(c *gin.Context) {
	tag := &model.Tag{}
	tid := c.Param("id")
	if err := c.ShouldBind(tag); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	tag, err := t.tagService.Update(tag, tid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, tag)
}

func (t *TagController) List(c *gin.Context) {
	tag, err := t.tagService.List()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, tag)
}

func (t *TagController) Name() string {
	return "tags"
}

func (t *TagController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/tags", t.List)
	api.GET("/tag/:id", t.Get)
	api.POST("/tag", t.Create)
	api.DELETE("/tag/:id", t.Delete)
	api.PUT("/tag/:id", t.Update)
}
