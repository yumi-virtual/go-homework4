package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title   string `gorm:"not null" json:"title" binding:"required"`
	Content string `gorm:"type:text;not null" json:"content" binding:"required"`

	UserId   uint  `gorm:"not null"`
	User     User  `gorm:"foreignKey:UserId"`
	Comments []Comment  `gorm:"foreignKey:PostId"`
}
