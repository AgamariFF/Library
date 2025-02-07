package cache

import (
	"context"
	"library/internal/models"
	"library/logger"
	"math"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

var rdb *redis.Client
var Ctx = context.Background()

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

func ClearCache() {
    if err := rdb.FlushDB(Ctx).Err(); err != nil {
        logger.ErrorLog.Println("Failed to flush Redis database:", err)
    } else {
        logger.InfoLog.Println("Redis cache cleared successfully")
    }
}

func CheckCacheGetBooks(page, limit, sort string, db *gorm.DB) (models.ResponseGetBooks, error) {
	var response models.ResponseGetBooks
	var books []models.Book
	cacheKey := "books:" + page + ":" + limit + ":" + sort
	var err error

	response.Limit, err = strconv.Atoi(limit)
	if err != nil {
		return response, err
	}

	response.Page, err = strconv.Atoi(page)
	if err != nil {
		return response, err
	}

	cachedData, err := rdb.Get(Ctx, cacheKey).Result()
	if err != nil {
		logger.InfoLog.Println("No cache found when /getBooks by the key =", cacheKey)
		if err := db.Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("genres.id, genres.name")
		}).Order(sort).Offset((response.Page - 1) * response.Limit).Limit(response.Limit).Find(&books).Error; err != nil {
			return response, err
		}
		var totalBooks int64
		db.Model(&models.Book{}).Count(&totalBooks)
		response.TotalBooks = int(totalBooks)
		response.TotalPages = int(math.Ceil(float64(totalBooks) / float64(response.Limit)))
		var genreResponse models.GenreFroGetBooks
		var bookResponse models.BookForGetBooks
		for _, book := range books {
			bookResponse.Genres = nil
			for _, genre := range book.Genres {
				genreResponse.ID = genre.ID
				genreResponse.Name = genre.Name
				bookResponse.Genres = append(bookResponse.Genres, genreResponse)
			}
			bookResponse.Author = book.Author
			bookResponse.ID = book.ID
			bookResponse.PublishedYear = book.PublishedYear
			bookResponse.Title = book.Title
			response.Books = append(response.Books, bookResponse)
		}
		booksJSON, err := json.Marshal(response)
		if err != nil {
			return response, err
		}
		if err := rdb.Set(Ctx, cacheKey, booksJSON, 5*time.Minute).Err(); err != nil {
			return response, err
		}
		return response, err
	}

	logger.InfoLog.Println("The cache was successfully found during the cacheKey =", cacheKey)
	if err := json.Unmarshal([]byte(cachedData), &response); err != nil {
		return response, err
	}

	return response, nil
}
