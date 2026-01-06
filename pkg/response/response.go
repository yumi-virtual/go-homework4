package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code int         `gorm:"code"`
	Msg  string      `gorm:"msg"`
	Data interface{} `gorm:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	// 存储到context 供日志使用
	c.Set("response_data", data)

	c.JSON(http.StatusOK, Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}
