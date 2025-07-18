package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"inkgo/common"
	"inkgo/model"
	"inkgo/service"
	"inkgo/utils"
	"net/http"
	"strconv"
)

type FavoriteController struct {
	favoriteService service.FavoriteService
}

func (f *FavoriteController) Name() string {
	return "favorites"
}

func (f *FavoriteController) RegisterRoute(api *gin.RouterGroup) {
	api.GET("/favorites", f.ListFavorites)
	api.GET("/favorites/:id", f.GetFavorite)
	api.POST("/favorites", f.CreateFavorite)
	api.PUT("/favorites/:id", f.UpdateFavorite)
	api.DELETE("/favorites/:id", f.DeleteFavorite)
	api.POST("/favorite/:id/post", f.AddPostToFavorite)
	api.DELETE("/favorite/:id/post/:post_id", f.RemovePostFromFavorite)
	api.GET("/favorite/:id/posts", f.ListPostsInFavorite)

	api.GET("/favorite/:id/posts/count", f.CountPostsInFavorite)
}

func NewFavoriteController(favoriteService service.FavoriteService) Controller {
	return &FavoriteController{
		favoriteService: favoriteService,
	}
}

// ListFavorites 用于列出用户的收藏夹
func (f *FavoriteController) ListFavorites(c *gin.Context) {
	user := model.User{
		Model: gorm.Model{
			ID: 1, // This should be replaced with actual user retrieval logic
		},
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "5"))

	favorites, total, err := f.favoriteService.GetUserFavorites(strconv.Itoa(int(user.ID)), page, pageSize)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessWithPage(c, favorites, total, page, pageSize)
}

// GetFavorite 用于获取特定收藏夹的详细信息
func (f *FavoriteController) GetFavorite(c *gin.Context) {
	uid := common.GetUserID(c)
	fid := c.Param("id")
	favorite, err := f.favoriteService.GetFavoriteByID(strconv.Itoa(int(uid)), fid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, favorite)
}

func (f *FavoriteController) CreateFavorite(c *gin.Context) {
	favorite := &model.Favorite{}
	if err := c.ShouldBindJSON(favorite); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	user := model.User{
		Model: gorm.Model{
			ID: 1, // This should be replaced with actual user retrieval logic
		},
	}
	favorite.UserID = user.ID
	uid := common.GetUserID(c)

	newFavorite, err := f.favoriteService.CreateFavorite(strconv.Itoa(int(uid)), favorite)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, newFavorite)
}

func (f *FavoriteController) UpdateFavorite(c *gin.Context) {
	fid := c.Param("id")
	favorite := &model.Favorite{}
	if err := c.ShouldBindJSON(favorite); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}

	favoritID, _ := strconv.ParseUint(fid, 10, 64)
	favorite.ID = uint(favoritID)
	uid := common.GetUserID(c)
	updatedFavorite, err := f.favoriteService.UpdateFavorite(strconv.Itoa(int(uid)), favorite)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, updatedFavorite)

}

// DeleteFavorite 用于删除特定收藏夹
func (f *FavoriteController) DeleteFavorite(c *gin.Context) {
	fid := c.Param("id")
	uid := common.GetUserID(c)
	err := f.favoriteService.DeleteFavorite(strconv.Itoa(int(uid)), fid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, nil)
}

func (f *FavoriteController) ListPostsInFavorite(c *gin.Context) {
	fid := c.Param("id")
	uid := common.GetUserID(c)
	favorite, err := f.favoriteService.GetFavoriteByID(strconv.Itoa(int(uid)), fid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}

	if favorite == nil {
		utils.Error(c, http.StatusNotFound, err)
		return
	}

	utils.Success(c, favorite.Posts)
}

func (f *FavoriteController) AddPostToFavorite(c *gin.Context) {
	fid := c.Param("id")
	var request AddPostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	uid := common.GetUserID(c)
	pid := strconv.Itoa(int(request.PostID))
	err := f.favoriteService.AddPostToFavorite(strconv.Itoa(int(uid)), pid, fid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, nil)
}

type AddPostRequest struct {
	PostID uint `json:"post_id" binding:"required"`
}

func (f *FavoriteController) RemovePostFromFavorite(c *gin.Context) {
	fid := c.Param("id")
	var request AddPostRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.Error(c, http.StatusBadRequest, err)
		return
	}
	uid := common.GetUserID(c)
	err := f.favoriteService.RemovePostFromFavorite(strconv.Itoa(int(uid)), strconv.Itoa(int(request.PostID)), fid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, nil)
}

func (f *FavoriteController) CountPostsInFavorite(c *gin.Context) {
	fid := c.Param("id")
	uid := common.GetUserID(c)
	count, err := f.favoriteService.CountPostsInFavorite(strconv.Itoa(int(uid)), fid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, gin.H{"count": count})
}
