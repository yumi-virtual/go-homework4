package routes

import (
	"github.com/gin-gonic/gin"
	"go-homework4/controllers"
	"go-homework4/middleware"
)

func SetupRouter() *gin.Engine {

	// 创建默认引擎
	route := gin.Default()

	// 全局中间键
	// 错误统一处理
	route.Use(middleware.ErrorHandler())
	// 成功日志记录
	route.Use(middleware.SuccessLogger())

	public := route.Group("/api")
	{
		// 公共接口 无需登录
		public.POST("/register", controllers.Register)
		public.POST("/login", controllers.Login)

		// 查询
		public.GET("/posts", controllers.FindPost)
		public.GET("/posts/:id", controllers.FirstPost)

	}

	// 受保护接口 需要登录jwt
	protected := route.Group("/api")
	protected.Use(middleware.JWTAuth())
	{
		// 文章
		protected.POST("/posts", controllers.CreatePost)
		protected.POST("/posts/:id", controllers.UpdatePost)
		protected.DELETE("/posts/:id", controllers.DeletePost)

		// 评论
		protected.POST("/posts/:id/comments", controllers.CreateComment)
		protected.GET("/posts/:id/comments", controllers.FindCommentByPostId)

	}

	return route
}
