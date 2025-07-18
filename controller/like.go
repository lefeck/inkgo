package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"inkgo/model"
	"inkgo/service"
	"inkgo/utils"
	"net/http"
)

type LikeController struct {
	likeService service.LikeService
}

func NewLikeController(likeService service.LikeService) Controller {
	return &LikeController{
		likeService: likeService,
	}
}

func (l *LikeController) Name() string {
	return "likes"
}

func (l *LikeController) RegisterRoute(api *gin.RouterGroup) {
	api.POST("/post/:id/like", l.LikePost)
	api.DELETE("/post/:id/like", l.UnLikePost)
	api.GET("/post/:id/like/count", l.CountLikes)
}

func NewlikeController(likeService service.LikeService) Controller {
	return &LikeController{
		likeService: likeService,
	}
}

func (l *LikeController) LikePost(c *gin.Context) {

	user := &model.User{Model: gorm.Model{ID: 2}}
	pid := c.Param("id")
	if err := l.likeService.LikePost(user, pid); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, nil)
}

func (p *LikeController) UnLikePost(c *gin.Context) {
	user := &model.User{Model: gorm.Model{ID: 2}}
	pid := c.Param("id")
	if err := p.likeService.UnLikePost(user, pid); err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, nil)
}

// CountLikes 获取某个帖子的点赞数
func (l *LikeController) CountLikes(c *gin.Context) {
	pid := c.Param("id")
	count, err := l.likeService.CountLikes(pid)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, err)
		return
	}
	utils.Success(c, gin.H{"count": count})
}
