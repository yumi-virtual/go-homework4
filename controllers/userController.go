package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go-homework4/config"
	"go-homework4/database"
	"go-homework4/models"
	"go-homework4/pkg/errno"
	"go-homework4/pkg/response"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type LoginUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(errno.ErrInvalidParameter)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Error(errno.ErrBadPassword)
		return
	}
	user.Password = string(password)

	var exists int64
	database.DB.Model(&user).Where("username=?", user.Username).Count(&exists)
	if exists > 0 {
		c.Error(errno.ErrUserExists)
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	response.Success(c, gin.H{
		"userId":   user.ID,
		"username": user.Username,
		"email":    user.Email,
	})

}

func Login(c *gin.Context) {
	var user LoginUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(errno.InvalidParameter(err))
		return
	}

	var loginUser models.User
	if err := database.DB.Model(&models.User{}).Where("username=?", user.Username).First(&loginUser).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(user.Password)); err != nil {
		c.Error(errno.ErrInvalidNameOrPassword)
		return
	}

	config := config.AppConfig
	sign := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   loginUser.ID,
		"username": loginUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := sign.SignedString([]byte(config.JWTSecret))
	if err != nil {
		c.Error(errno.Internal(err))
		return
	}

	response.Success(c, gin.H{
		"token":  token,
		"userId": loginUser.ID,
		"expire": time.Now().Add(time.Hour * 24).Format(time.RFC3339),
	})

}
