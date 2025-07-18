package controller

import (
	"github.com/gin-gonic/gin"
	"inkgo/model"
	"inkgo/service"
	"inkgo/utils"
	"net/http"
	"strconv"
)

type CategoryController struct {
	categoryService service.CategoryService
}

func NewCategoryController(categoryService service.CategoryService) Controller {
	return &CategoryController{
		categoryService: categoryService,
	}
}

func (ca *CategoryController) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "5"))
	categories, total, err := ca.categoryService.List(page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	utils.SuccessWithPage(c, categories, total, page, pageSize)
}

func (ca *CategoryController) Get(c *gin.Context) {
	aid := c.Param("id")
	category, err := ca.categoryService.Get(aid)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	utils.Success(c, category)
}
func (ca *CategoryController) Create(c *gin.Context) {
	category := &model.Category{}
	if err := c.ShouldBindJSON(category); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	category, err := ca.categoryService.Create(category)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	utils.Success(c, category)
}

func (ca *CategoryController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := ca.categoryService.Delete(id); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	utils.Success(c, nil)
}

func (ca *CategoryController) Update(c *gin.Context) {
	cid := c.Param("id")
	category := &model.Category{}
	if err := c.ShouldBindJSON(category); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	category, err := ca.categoryService.Update(cid, category)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	utils.Success(c, category)
}

func (ca *CategoryController) Name() string {
	return "categories"
}

func (ca *CategoryController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/categories", ca.List)
	api.GET("/category/:id", ca.Get)
	api.POST("/category", ca.Create)
	api.DELETE("/category/:id", ca.Delete)
	api.PUT("/category/:id", ca.Update)
}
