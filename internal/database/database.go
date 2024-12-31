package database

import (
	"library/internal/models"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var TestDB *gorm.DB

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
