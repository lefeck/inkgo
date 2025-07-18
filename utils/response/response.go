package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Response struct {
	Code    int         `json:"code"`    //状态码
	Message string      `json:"message"` //提示信息
	Data    interface{} `json:"data"`    //返回数据
}

func JSON(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, "success", data)
}

func Error(c *gin.Context, code int, message string) {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	JSON(c, code, message, nil)
}

type PageData struct {
	List     interface{} `json:"list"`      // 数据列表
	Total    int64       `json:"total"`     // 总数量
	Page     int         `json:"page"`      // 当前页码
	PageSize int         `json:"page_size"` // 每页数量
}

// ParsePagination 从请求参数解析分页信息
func ParsePagination(c *gin.Context) (page int, pageSize int, offset int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	offset = (page - 1) * pageSize
	return
}
