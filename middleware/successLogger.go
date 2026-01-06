package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-homework4/logger"
)

func SuccessLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() != 200 || len(c.Errors) > 0 {
			return
		}

		data, _ := c.Get("response_data")
		logMsg := "code=0 msg=success"

		if m, ok := data.(gin.H); ok {
			if p, ok := m["pagination"].(gin.H); ok {
				logMsg += fmt.Sprintf(" page=%v size=%v total=%v", p["page"], p["size"], p["total"])
			}
		}
		logger.Info.Printf("%s method=%s path=%s", logMsg, c.Request.Method, c.Request.URL.Path)
	}
}
