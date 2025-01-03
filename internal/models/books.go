package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model    `swaggerignore:"true"`
	Title         string
	Author        string
	PublishedYear string  `json:"published_year"`
	Genres        []Genre `gorm:"many2many:book_genres"`
	Description   string
}
