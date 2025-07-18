package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code     int         `json:"code"`                //状态码
	Message  string      `json:"message"`             //提示信息
	Data     interface{} `json:"data"`                //返回数据
	Total    *int64      `json:"total,omitempty"`     // 总数（分页时返回）
	Page     *int        `json:"page,omitempty"`      // 当前页（分页时返回）
	PageSize *int        `json:"page_size,omitempty"` // 页大小（分页时返回）
}

func JSON(c *gin.Context, code int, message string, data interface{}, total *int64, page, pageSize *int) {
	c.JSON(code, Response{
		Code:     code,
		Message:  message,
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, "success", data, nil, nil, nil)
}

func SuccessWithPage(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	JSON(c, http.StatusOK, "success", data, &total, &page, &pageSize)
}

func Error(c *gin.Context, code int, err error) {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	JSON(c, code, err.Error(), nil, nil, nil, nil)
}
