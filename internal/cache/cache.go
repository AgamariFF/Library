package cache

import (
	"context"
	"fmt"
	"library/internal/handlers"
	"library/internal/models"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
}

func GetClient() *redis.Client {
	return rdb
}

func CheckCashGetBooks(page, limit, sort string, totalBooks, totalPages int, db *gorm.DB) (handlers.ResponseGetBooks, error) {
	var response handlers.ResponseGetBooks
	var books []models.Book
	cacheKey := "books:" + page + ":" + limit + ":" + sort

	response.Limit = strconv.Atoi(limit)
	response.Page = strconv.Atoi(page)
	response.TotalBooks = totalBooks
	response.TotalPages = totalPages
	
	cachedData, err := rdb.Get(ctx, cacheKey).Result()
	if err != nil {
		if err := db.Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("genres.id, genres.name")
		}).Find(&books).Error; err != nil {
			return response, err
		}
		response.Books = books
	}
	if err := json.Unmarshal([]byte(cachedData), &response)
}
