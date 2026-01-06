package controllers

import (
	"github.com/gin-gonic/gin"
	"go-homework4/database"
	"go-homework4/models"
	"net/http"
	"strconv"
)

//type CreateCommentRequest struct {
//	PostId uint `json:"postId"`
//	Content string `json:"content"`
//}

func CreateComment(c *gin.Context) {
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
			"code": 40001,
			"msg":  "Invalid post id",
		})
		return
	}

	var comment models.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 40001,
			"msg":  "Bad request parameters",
		})
		return
	}

	var post models.Post
	if err := database.DB.First(&post, postId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50001,
			"msg":  "Failed to found post",
		})
		return
	}

	comment.UserId = userId.(uint)
	comment.PostId = uint(postId)

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50001,
			"msg":  "Failed to create comment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": comment,
	})

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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 40001,
			"msg":  "Failed to count comments",
		})
		return
	}

	postId := c.Param("id")
	var comments []models.Comment
	if err := database.DB.Preload("USER.POST").Where("post_id=?", postId).Order("create_at DESC").Offset(offset).Limit(size).Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 50001,
			"msg":  "Failed to find comments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": gin.H{
			"list": comments,
			"pagination": gin.H{
				"page":         page,
				"size":         size,
				"total":        total,
			},
		},
	})

}
