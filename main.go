package go_homework4

import (
	"fmt"
	"go-homework4/config"
	"go-homework4/database"
	"go-homework4/routes"
	"log"
)

func main() {
	// 1、加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load configuration", err)
		return
	}
	fmt.Println("Config loaded successfully! ")

	// 2、连接数据库
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to connect database", err)
		return
	}
	fmt.Println("Database connected successfully! ")

	// 3、设置路由
	router := routes.SetupRouter()

	// 4、启动服务器
	router.Run(":" + config.AppConfig.DatabasePort)

}
