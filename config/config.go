package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ServerPort       string `mapstructure:"SERVER_PORT"`
	DatabaseHost     string `mapstructure:"DATABASE_HOST"`
	DatabasePort     string `mapstructure:"DATABASE_PORT"`
	DatabaseUser     string `mapstructure:"DATABASE_USER"`
	DatabasePassword string `mapstructure:"DATABASE_PASSWORD"`
	DatabaseName     string `mapstructure:"DATABASE_NAME"`

	JWTSecret string `mapstructure:"JWT_SECRET"`
	LogLevel  string `mapstruture:"LOG_LEVEL"`
}

var AppConfig Config

func LoadConfig() error {
	// 配置文件名
	viper.SetConfigName(".env")
	// 配置文件类型
	viper.SetConfigType("env")
	// 设置文件路径
	viper.AddConfigPath(".")
	// 允许读取环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("No config file found")
		return err
	}
	if err:=viper.Unmarshal(&AppConfig);err != nil {
		log.Fatal("Failed to parse configuration ",err)
		return err
	}
	log.Println("config load success!")
	return nil
}
