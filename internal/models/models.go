package models

import (
	"time"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model    `swaggerignore:"true"`
	Title         string
	Author        string
	PublishedYear string  `json:"published_year"`
	Genres        []Genre `gorm:"many2many:book_genres"`
	Description   string
}

type Genre struct {
	gorm.Model  `swaggerignore:"true"`
	Name        string
	Description string
	Books       []Book `gorm:"many2many:book_genres"`
}

type User struct {
	gorm.Model `swaggerignore:"true"`
	Name       string `gorm:"size:100" json:"name" binding:"required"`
	Email      string `gorm:"unique; not null" json:"email" binding:"required,email"`
	Role       string `gorm:"not null" json:"role"`
	Mailing    bool   `gorm:"not null" json:"mailing" binding:"required"`
	Password   string `json:"-"`

	// Поля сессии
	RefreshToken string    `gorm:"not null" json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type GenreFroGetBooks struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type BookForGetBooks struct {
	ID            uint               `json:"id"`
	Title         string             `json:"title"`
	Author        string             `json:"author"`
	PublishedYear string             `json:"published_year"`
	Genres        []GenreFroGetBooks `json:"genres"`
}

// ResponseGetBooks структура ответа при GET запросе /getBooks
type ResponseGetBooks struct {
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalBooks int               `json:"total_books"`
	TotalPages int               `json:"total_pages"`
	Books      []BookForGetBooks `json:"books"`
}
