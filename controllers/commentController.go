package controllers

import (
	"github.com/gin-gonic/gin"
	"go-homework4/database"
	"go-homework4/models"
	"go-homework4/pkg/errno"
	"go-homework4/pkg/response"
	"strconv"
	"time"
)

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CommentResponse struct {
	ID        uint      `json:"id"`
	Content   string    `json:"content"`
	Username  string    `json:"username"`
	PostTitle string    `json:"post_title"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateComment(c *gin.Context) {
	uid, exists := c.Get("userId")
	if !exists {
		c.Error(errno.ErrUnauthorized)
		return
	}

	userId, ok := uid.(float64)
	if !ok {
		c.Error(errno.ErrUnauthorized)
		return
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Error(errno.InvalidParameter(err))
		return
	}

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(errno.InvalidParameter(err))
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postId).Error; err != nil {
		c.Error(errno.DB(err))
		return
	}

	comment := models.Comment{
		UserId:  uint(userId),
		PostId:  uint(postId),
		Content: req.Content,
	}

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
	if err := database.DB.Preload("User").Preload("Post").Where("post_id=?", postId).Order("created_at DESC").Offset(offset).Limit(size).Find(&comments).Error; err != nil {
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

func FindCommentByPostIdV2(c *gin.Context) {
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
	var comments []CommentResponse
	database.DB.Table("comments").
		Select("comments.id, comments.content, comments.created_at, users.username, posts.title as post_title").
		Joins("LEFT JOIN users ON users.id = comments.user_id").
		Joins("LEFT JOIN posts ON posts.id = comments.post_id").
		Where("comments.post_id = ?", postId).
		Order("comments.created_at DESC").
		Offset(offset).Limit(size).
		Scan(&comments)

	response.Success(c, gin.H{
		"list": comments,
		"pagination": gin.H{
			"page":  page,
			"size":  size,
			"total": total,
		},
	})

}
