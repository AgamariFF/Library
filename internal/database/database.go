package database

import (
	"fmt"
	"library/internal/models"
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

		fmt.Printf("Failed to connect to database (attempt %d/%d): %s", i+1, maxRetries, err)
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

func InitTestDB() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&models.Book{}, &models.Genre{}, models.User{}); err != nil {
		panic(err)
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
