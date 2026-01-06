package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go-homework4/config"
	"net/http"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1、从请求头拿Authorization
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 10001,
				"msg":  "缺少认证令牌",
			})
			return
		}

		// 2、验证并解析JWT
		token, err := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrHashUnavailable
			}
			return []byte(config.AppConfig.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 10003,
				"msg":  "无效的认证令牌",
			})
			return
		}

		// 3、提取用户id并把userId放进上下文
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userId, ok := claims["userId"].(float64); ok {
				c.Set("userId", userId)
				c.Next()
				return
			}
		}

		// 4、令牌有效但是无法解析claims
		c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{
			"code":10004,
			"msg":"令牌格式错误",
		})

	}

}
