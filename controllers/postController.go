package controllers

import (
	"github.com/gin-gonic/gin"
	"go-homework4/database"
	"go-homework4/models"
	"go-homework4/pkg/errno"
	"go-homework4/pkg/response"
	"strconv"
)

type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=2,max=200"`
	Content string `json:"content" binding:"required,min=10"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" binding:"required,min=2,max=200"`
	Content string `json:"content" binding:"required,min=10"`
}

func CreatePost(c *gin.Context) {
	// 从jwt中拿userId
	userId, exists := c.Get("userId")
	if !exists {
		c.Error(errno.ErrUnauthorized)
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errno.ErrInvalidParameter)
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserId:  userId.(uint),
	}
	if err := database.DB.Model(&post).Where("title=?", post.Title).First(&post).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	response.Success(c, gin.H{
		"title":   post.Title,
		"content": post.Content,
	})
}

func FindPost(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	offser := (page - 1) * size

	var total int64
	if err := database.DB.Model(&models.Post{}).Count(&total).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	var posts []models.Post
	if err := database.DB.Preload("User").Order("create_at DESC ").Offset(offser).Limit(size).Find(&posts).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	response.Success(c, gin.H{
		"list": posts,
		"pagination": gin.H{
			"page":  page,
			"size":  size,
			"total": total,
		},
	})

}

func FirstPost(c *gin.Context) {
	postId := c.Param("id")
	var post models.Post
	if err := database.DB.Model(&post).Preload("User").First(&post, postId).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	response.Success(c, post)

}

func UpdatePost(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Error(errno.ErrUnauthorized)
		return
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errno.ErrInvalidParameter)
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postId).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	if userId != post.UserId {
		c.Error(errno.ErrStatusForbidden)
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errno.ErrInvalidParameter)
		return
	}

	post.Title = req.Title
	post.Content = req.Content

	if err := database.DB.Save(&post).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	response.Success(c, post)
}

func DeletePost(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Error(errno.ErrUnauthorized)
		return
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errno.ErrInvalidParameter)
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postId).Error; err != nil {
		c.Error(errno.ErrPostNotFound)
		return
	}

	if userId != post.UserId {
		c.Error(errno.ErrStatusForbidden)
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	response.Success(c, post)

}
