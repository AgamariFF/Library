package models_test

import (
	"library/internal/database"
	"library/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateTestData(t *testing.T) {
	// Инициализация тестовой базы данных
	database.InitTestDB()
	defer database.CleanupTestDB()

	// Тестовые данные
	book := models.Book{
		Title:         "Test Book",
		Author:        "Test Author",
		PublishedYear: "2022",
		Description:   "A test book",
		Genres: []models.Genre{
			{Name: "Fiction"},
		},
	}

	// Создание книги
	err := database.DB.Create(&book).Error
	assert.NoError(t, err)
	var result models.Book
	err = database.TestDB.Preload("Genres").First(&result).Error
	assert.NoError(t, err)
	assert.Equal(t, "Test Book", result.Title)
	assert.Equal(t, "Test Author", result.Author)
	assert.Equal(t, "2022", result.PublishedYear)
	assert.Equal(t, "A test book", result.Description)
}

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Book{}, &models.Genre{}, &models.User{})
	return db
}

func clearDB(db *gorm.DB) {
	db.Exec("DELETE FROM Books")
	db.Exec("DELETE FROM GENRES")
}

func TestBookCreation(t *testing.T) {
	db := setupTestDB()
	defer clearDB(db)

	book := models.Book{
		Title:  "Test book",
		Author: "Test Author",
	}

	err := db.Create(&book).Error
	assert.NoError(t, err)
	var result models.Book
	db.First(&result, "title = ?", "Test Book")

	assert.Equal(t, "Test Book", result.Author)
}
