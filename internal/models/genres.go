package models

import "gorm.io/gorm"

type Genre struct {
	gorm.Model  `swaggerignore:"true"`
	Name        string
	Description string
	Books       []Book `gorm:"many2many:book_genres"`
}
