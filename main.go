package go_homework4

import (
	"fmt"
	"go-homework4/config"
	"go-homework4/database"
	"go-homework4/logger"
	"go-homework4/routes"
)

func main() {

	// 1、初始化日志
	logger.Init()

	// 2、加载配置
	if err := config.LoadConfig(); err != nil {
		logger.Error.Fatalf("Failed to load configuration", err)
		return
	}
	fmt.Println("Config loaded successfully! ")

	// 3、连接数据库
	if err := database.InitDB(); err != nil {
		logger.Error.Fatalf("Failed to connect database", err)
		return
	}
	fmt.Println("Database connected successfully! ")

	// 4、设置路由
	router := routes.SetupRouter()

	// 5、启动服务器
	router.Run(":" + config.AppConfig.DatabasePort)

}
