package controller

import (
	"github.com/gin-gonic/gin"
	"inkgo/common"
	"inkgo/service"
	"net/http"
)

type UploadController struct {
	uploadService service.IUploadService
}

func NewUploadController(uploadSerivce service.IUploadService) *UploadController {
	return &UploadController{uploadService: uploadSerivce}
}

// POST
// http://127.0.0.1:8080/api/v1/upload
// key: file, value: filename
func (uc *UploadController) Upload(c *gin.Context) {
	file, fileHeader, _ := c.Request.FormFile("file")
	url, err := uc.uploadService.UploadFile(file, fileHeader.Size)
	if err != nil {
		common.ResponseFailed(c, http.StatusInternalServerError, err)
		return
	}
	common.ResponseSuccess(c, url)
}
