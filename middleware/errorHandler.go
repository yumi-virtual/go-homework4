package middleware

import (
	"github.com/gin-gonic/gin"
	"go-homework4/logger"
	"go-homework4/pkg/errno"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err
		if e, ok := err.(*errno.Error); ok {
			if e.Err != nil {
				logger.Error.Printf("code=%d msg=%s err=%v ", e.Code, e.Msg, e.Err)
			} else {
				logger.Warn.Printf("code=%d msg=%s ", e.Code, e.Msg)
			}

			c.JSON(e.HTTPCode, gin.H{
				"code": e.Code,
				"msg":  e.Msg,
			})
			return
		}

		logger.Error.Printf("panic=%v", err)
		c.JSON(500, gin.H{
			"code": errno.ErrInternal.Code,
			"msg":  errno.ErrInternal.Msg,
		})

	}
}
