package controllers

import (
	"github.com/gin-gonic/gin"
	"go-homework4/database"
	"go-homework4/models"
	"go-homework4/pkg/errno"
	"go-homework4/pkg/response"
	"strconv"
)

//type CreateCommentRequest struct {
//	PostId uint `json:"postId"`
//	Content string `json:"content"`
//}

func CreateComment(c *gin.Context) {
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

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.Error(errno.ErrInvalidParameter)
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postId).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	comment.UserId = userId.(uint)
	comment.PostId = uint(postId)

	if err := database.DB.Create(&comment).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	response.Success(c, comment)

}

func FindCommentByPostId(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}
	offset := (page - 1) * size

	var total int64
	if err := database.DB.Model(&models.Comment{}).Count(&total).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	postId := c.Param("id")
	var comments []models.Comment
	if err := database.DB.Preload("USER.POST").Where("post_id=?", postId).Order("create_at DESC").Offset(offset).Limit(size).Find(&comments).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	response.Success(c, gin.H{
		"list": comments,
		"pagination": gin.H{
			"page":  page,
			"size":  size,
			"total": total,
		},
	})

}
