package controllers

import (
	"github.com/gin-gonic/gin"
	"go-homework4/database"
	"go-homework4/models"
	"net/http"
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
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 10001,
			"msg":  "Unauthorized",
		})
		return
	}

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 40001,
			"msg":  "Bad request parameters",
		})
	}

	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserId:  userId.(uint),
	}
	if err := database.DB.Model(&post).Where("title=?", post.Title).First(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 40002,
			"msg":  "Post already exists ",
		})
	}

	if err := database.DB.Create(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50001,
			"msg":  "Failed to create post",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "create post successfully!",
		"data": gin.H{
			"title":   post.Title,
			"content": post.Content,
		},
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50001,
			"msg":  "Failed to count posts",
		})
		return
	}

	var posts []models.Post
	if err := database.DB.Preload("User").Order("create_at DESC ").Offset(offser).Limit(size).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50002,
			"msg":  "Failed to get post list",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg": gin.H{
			"list": posts,
			"pagination": gin.H{
				"page":         page,
				"size":         size,
				"total":        total,
			},
		},
	})

}

func FirstPost(c *gin.Context) {
	postId := c.Param("id")
	var post models.Post
	if err := database.DB.Model(&post).Preload("User").First(&post, postId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50001,
			"msg":  "Failed to get post",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": post,
	})

}

func UpdatePost(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 10001,
			"msg":  "Unauthorized",
		})
		return
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 50001,
			"msg":  "Invalid post id",
		})
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code": 40001,
			"msg":  "Post not found ",
		})
		return
	}

	if userId != post.UserId {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 50005,
			"msg":  "You are not the author of the post",
		})
		return
	}

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 50001,
			"msg":  "Invalid request parameters",
		})
		return
	}

	post.Title = req.Title
	post.Content = req.Content

	if err := database.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50003,
			"msg":  "Failed to update post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": post,
	})

}

func DeletePost(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 10001,
			"msg":  "Unauthorized",
		})
		return
	}

	postId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 50001,
			"msg":  "Invalid post id",
		})
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50005,
			"msg":  "Failed to get post",
		})
		return
	}

	if userId != post.UserId {
		c.JSON(http.StatusForbidden, gin.H{
			"code": 50005,
			"msg":  "You are not the author of the post",
		})
		return
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50005,
			"msg":  "Failed to delete post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": post,
	})

}
