package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" json:"content"`

	UserId uint `gorm:"not null"`
	User   User `gorm:"foreignKey:UserId"`
	PostId uint `json:"postId"`
	Post   Post `gorm:"foreignKey:PostId"`
}
