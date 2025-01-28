package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `swaggerignore:"true"`
	Name       string `gorm:"size:100" json:"name" binding:"required"`
	Email      string `gorm:"unique; not null" json:"email" binding:"required,email"`
	Role       string `gorm:"not null" json:"role"`
	Mailing    bool   `gorm:"not null" json:"mailing" binding:"required"`
	Password   string `json:"-"`
}
