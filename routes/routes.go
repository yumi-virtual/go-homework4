package routes

import (
	"github.com/gin-gonic/gin"
	"go-homework4/controllers"
	"go-homework4/middleware"
)

func SetupRouter() *gin.Engine{

	route:=gin.Default()
	auth:=route.Group("/api")

	// 公共接口 无需登录
	auth.POST("/register",controllers.Register)
	auth.POST("/login",controllers.Login)

	// 查询
	auth.GET("/posts",controllers.FindPost)
	auth.GET("/posts/:id",controllers.FirstPost)



	// 受保护接口 需要登录jwt
	auth.Use(middleware.JWTAuth())
	{
		// 文章
		auth.POST("/posts",controllers.CreatePost)
		auth.POST("/posts/:id",controllers.UpdatePost)
		auth.DELETE("/posts/:id",controllers.DeletePost)

		// 评论
		auth.POST("/posts/:id/comments",controllers.CreateComment)
		auth.GET("/posts/:id/comments",controllers.FindCommentByPostId)

	}

	return route
}
