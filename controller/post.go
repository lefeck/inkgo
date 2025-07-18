package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"inkgo/model"
	"inkgo/service"
	"inkgo/utils"
	"net/http"
	"strconv"
)

type PostController struct {
	postService service.PostService
}

func NewPostController(PostService service.PostService) Controller {
	return &PostController{
		postService: PostService,
	}
}

// GetPostByID 获取文章详情
func (p *PostController) GetPostByID(c *gin.Context) {
	pid := c.Param("id")
	post, err := p.postService.GetPostByID(pid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, post)
}

// Get 获取单个文章
//func (p *PostController) GetPostByID(c *gin.Context) {
//	pid := c.Param("id")
//	user := &model.User{}
//
//	post, err := p.postService.Get(user, pid)
//	if err != nil {
//		utils.Error(c, http.StatusInternalServerError, err)
//		return
//	}
//	utils.Success(c, post)
//}

// GetPostByName 获取文章详情
func (p *PostController) GetPostByName(c *gin.Context) {
	name := c.Param("name")
	post, err := p.postService.GetPostByName(name)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, post)
}

// List 已发布的文章列表
func (p *PostController) ListHasPublished(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "2"))
	posts, total, err := p.postService.HasPublished(page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessWithPage(c, posts, total, page, pageSize)
}

func (a *PostController) ListDrafts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "2"))
	posts, total, err := a.postService.ListDrafts(page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessWithPage(c, posts, total, page, pageSize)
}

// Create 创建文章
func (p *PostController) Create(c *gin.Context) {
	user := new(model.User)
	post := new(model.Post)
	if err := c.ShouldBindJSON(&post); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	post, err := p.postService.Create(user, post)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, post)
}

func (p *PostController) Delete(c *gin.Context) {
	pid := c.Param("id")

	if err := p.postService.Delete(pid); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, nil)
}

func (p *PostController) Update(c *gin.Context) {
	id := c.Param("id")
	post := &model.Post{}
	if err := c.ShouldBindJSON(&post); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	post, err := p.postService.Update(id, post)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, post)
}

func (p *PostController) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	state := c.Query("state")
	if state == "" {
		utils.Error(c, http.StatusBadRequest, errors.New("state is required"))
		return
	}

	post, err := p.postService.UpdateStatus(id, model.PostState(state))
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, post)
}

// SortByViewCountDesc 按照阅读数降序排列文章
func (a *PostController) SortByViewCountDesc(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "2"))
	posts, total, err := a.postService.SortByViewCountDesc(page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessWithPage(c, posts, total, page, pageSize)
}

// ListHotPosts 获取热门文章
func (a *PostController) ListHotPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	posts, err := a.postService.ListHotPosts(limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, posts)
}

// ListRecentPosts 获取最近的文章
func (a *PostController) ListRecentPosts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	posts, err := a.postService.ListRecentPosts(limit)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, posts)
}

func (a *PostController) Name() string {
	return "posts"
}

func (a *PostController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/posts/published", a.ListHasPublished)
	api.GET("/posts/drafts", a.ListDrafts)
	api.GET("/posts/hot", a.ListHotPosts)
	api.GET("/posts/recent", a.ListRecentPosts)
	api.GET("/posts/sort", a.SortByViewCountDesc)
	api.GET("/post/:id", a.GetPostByID)
	api.GET("/post/name/:name", a.GetPostByName)
	api.POST("/post", a.Create)
	api.DELETE("/post/:id", a.Delete)
	api.PUT("/post/:id", a.Update)
	api.PUT("/post/:id/status", a.UpdateStatus)
}
