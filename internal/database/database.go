package database

import (
	"fmt"
	"library/internal/models"
	"library/logger"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var TestDB *gorm.DB

func ConnectWithRetry(maxRetries int, delay time.Duration) error {
	var database *gorm.DB
	var err error
	dsn := os.Getenv("DB_DSN")
	for i := 0; i < maxRetries; i++ {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			DB = database
			return nil
		}

		logger.InfoLog.Printf("Failed to connect to database (attempt %d/%d): %s", i+1, maxRetries, err)
		time.Sleep(delay)
	}
	return err
}

func ConnectDatabase() error {
	dsn := os.Getenv("DB_DSN")
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = database
	return nil
}

func CreateTrgmIndexes(db *gorm.DB) error {
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_books_title_trgm ON books USING GIN (title gin_trgm_ops);").Error; err != nil {
		return err
	}
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_books_description_trgm ON books USING GIN (description gin_trgm_ops);").Error; err != nil {
		return err
	}
	return nil
}

// Поиск книг
func SearchBooks(db *gorm.DB, searchString string, similarity float64, offset, limit int) ([]models.Book, int, error) {
	var books []models.Book
	query := db.Preload("Genres", func(db *gorm.DB) *gorm.DB {
		return db.Select("genres.id, genres.name")
	}).
		Where("similarity(lower(title), lower(?)) > ?", searchString, similarity).
		Or("similarity(lower(description), lower(?)) > ?", searchString, similarity).
		Or("lower(title) LIKE lower(?)", "%"+searchString+"%")
	var totalBooks int64
	if err := query.Model(&models.Book{}).Count(&totalBooks).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Find(&books).Error; err != nil {
		return nil, 0, err
	}

	return books, int(totalBooks), nil
}

func InitTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to open database: %v", err))
	}
	if err := db.AutoMigrate(&models.Book{}, &models.Genre{}, models.User{}); err != nil {
		panic(fmt.Sprintf("Failed to migrate database : %v", err))
	}

	TestDB = db
}

func CleanupTestDB() {
	if TestDB != nil {
		sqlDB, err := TestDB.DB()
		if err == nil {
			sqlDB.Close()
		}
		TestDB = nil
	}
}

func TestInitTestDB(t *testing.T) {
	assert := assert.New(t)

	InitTestDB()
	defer CleanupTestDB()
	assert.NotNil(TestDB, "TestDB should be initialized")
	err := TestDB.AutoMigrate(&models.Book{})
	assert.NoError(err, "AutoMigrate should succeed")
}

func TestCleanupTestDB(t *testing.T) {
	InitTestDB()
	CleanupTestDB()

	assert.Nil(t, TestDB, "TestDB should be nil after cleanup")
}
