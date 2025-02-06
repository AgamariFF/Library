package handlers

import (
	"encoding/json"
	"errors"
	"library/internal/database"
	"library/internal/kafka"
	"library/internal/models"
	"library/logger"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddBookRequest структура запроса для добавления книги
// @Schema example={"title": "Golang Basics", "author": "John Doe", "published_year": "2024", "genre": ["Учебная литература"], "description": "Эта книга — идеальный выбор для тех, кто хочет начать свое путешествие в программировании на языке Go."}
type AddBookRequest struct {
	Title          string   `json:"title" binding:"required" example:"Golang Basics"`                                                                             // Название книги
	Author         string   `json:"author" binding:"required" example:"John Doe"`                                                                                 // Автор
	Genre          []string `json:"genre" binding:"required" example:"Учебная литература"`                                                                        // Жанра
	Published_year string   `json:"published_year" binding:"required" example:"2024"`                                                                             // Год публикации
	Description    string   `json:"description" example:"Эта книга — идеальный выбор для тех, кто хочет начать свое путешествие в программировании на языке Go."` // Описание книги
}

// ModifyingBookRequest структура запроса для добавления книги
// @Schema example={"id": 1, "published_year": "2021", "author": "Jeff Bezos", "title": "Why Does the World Exist?", "description": "Explore the ultimate question: Why is there something rather than nothing? This thought-provoking journey through philosophy, science, and metaphysics challenges readers to ponder existence itself, blending deep inquiry with accessible insight. A must-read for curious minds."}
type ModifyingBookRequest struct {
	Id             uint     `json:"id" binding:"required" example:"1"`
	Title          string   `json:"title" example:"Why Does the World Exist?"`
	Author         string   `json:"author" example:"Jeff Bezos"`
	Genre          []string `json:"genre" example:"Детектив"`
	Published_year string   `json:"published_year" example:"2021"`
	Description    string   `json:"description" example:"Explore the ultimate question: Why is there something rather than nothing? This thought-provoking journey through philosophy, science, and metaphysics challenges readers to ponder existence itself, blending deep inquiry with accessible insight. A must-read for curious minds."`
}

// DeleteBookRequest структура запроса для удаления книги
// @Schema example={"id": 1}
type DeleteBookRequest struct {
	ID uint `json:"id" example:"1" binding:"required"`
}

// Welcome отображает главную страницу
// @Summary Show start page
// @Description Show the start page of the API
// @Tags book
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the Library API!",
	})
}

// GetBooks возвращает отсортированный или неотсортированный список книг с пагинацией
// @Summary Get list of books
// @Description Retrieve all books, optionally sorted by a specific field
// @Tags book
// @Accept json
// @Produce json
// @Param sort query string false "Field to sort by (e.g., 'title', 'author', 'published_year')"
// @Param page query int false "Page number for pagination (default: 1)"
// @Param limit query int false "Number of books per page (default: 10)"
// @Success 200 {object} map[string]interface{} "Returns a paginated and sorted list of books"
// @Failure 500 {object} map[string]string
// @Router /getBooks [get]
func GetBooks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var books []models.Book
		var response struct {
			Page       int `json:"page"`
			Limit      int `json:"limit"`
			TotalBooks int `json:"total_books"`
			TotalPages int `json:"total_pages"`
			Books      []struct {
				ID            uint   `json:"id"`
				Title         string `json:"title"`
				Author        string `json:"author"`
				PublishedYear string `json:"published_year"`
				Genres        []struct {
					ID   uint   `json:"id"`
					Name string `json:"name"`
				} `json:"genres"`
			} `json:"books"`
		}
		// Извлечение query-параметров

		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil {
			logger.ErrorLog.Println("Failed Atoi page\tError:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid page"})
		}
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil {
			logger.ErrorLog.Println("Failed Atoi limit\tError:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid limit"})
		}
		sortParam := c.Query("sort")

		offset := (page - 1) * limit

		var totalBooks int64
		db.Model(&models.Book{}).Count(&totalBooks)

		totalPages := int(math.Ceil(float64(totalBooks) / float64(limit)))

		// Получаем все книги из БД
		query := db.Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("genres.id, genres.name") // Выбираем только нужные поля
		})

		// Сортируем книги
		switch sortParam {
		case "author":
			query = query.Order("author")
		case "title":
			query = query.Order("title")
		case "year":
			query = query.Order("published_year")
		default:
			query = query.Order("id")
		}

		query = query.Offset(offset).Limit(limit)

		if err = query.Find(&books).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get books"})
		}

		response.Page = page
		response.Limit = limit
		response.TotalBooks = int(totalBooks)
		response.TotalPages = totalPages

		// Формируем ответ
		for _, book := range books {
			bookResponse := struct {
				ID            uint   `json:"id"`
				Title         string `json:"title"`
				Author        string `json:"author"`
				PublishedYear string `json:"published_year"`
				Genres        []struct {
					ID   uint   `json:"id"`
					Name string `json:"name"`
				} `json:"genres"`
			}{
				ID:            book.ID,
				Title:         book.Title,
				Author:        book.Author,
				PublishedYear: book.PublishedYear,
			}

			// Формируем список жанров
			for _, genre := range book.Genres {
				bookResponse.Genres = append(bookResponse.Genres, struct {
					ID   uint   `json:"id"`
					Name string `json:"name"`
				}{
					ID:   genre.ID,
					Name: genre.Name,
				})
			}
			response.Books = append(response.Books, bookResponse)
		}
		// Возврат ответа
		c.JSON(http.StatusOK, response)
	}
}

// GetBook возвращает информацию об одной книге
// @Summary      Get one book
// @Description  Get detailed information about a single book by ID
// @Description  JWT authentication via cookie.
// @Description	 The JWT token should be stored in a cookie named "jwt".
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        bookId  query    integer  true  "Book ID"
// @Success      200     {object} models.Book
// @Failure      400     {object} map[string]string
// @Failure      404     {object} map[string]string
// @Failure      500     {object} map[string]string
// @Router       /getBook [get]
func GetBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var book models.Book

		// Получение ID книги из query-параметра
		bookId := c.Query("bookId")
		if bookId == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Missing bookId parameter",
			})
			return
		}

		// Поиск книги в базе данных
		if err := db.Preload("Genres").Where("id = ?", bookId).First(&book).Error; err != nil {
			// Если книга не найдена
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Book not found",
				})
				return
			}

			// Если произошла другая ошибка
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to get book",
			})
			return
		}

		// Успешный ответ
		c.JSON(http.StatusOK, book)
	}
}

// AddBook
// @Summary      Add a new book
// @Description  JWT authentication via cookie only for admin.
// @Description	 The JWT token should be stored in a cookie named "jwt".
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book  body  AddBookRequest  true  "Book Data"  example({"title": "Golang Basics", "author": "John Doe", "published_year": "2024", "genre": ["Учебная литература"], "description": "Эта книга — идеальный выбор для тех, кто хочет начать свое путешествие в программировании на языке Go."})
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /addBook [post]
func AddBook(db *gorm.DB, producer *kafka.KafkaProducer) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.InfoLog.Println("Starting add book")
		// Структура для запроса
		var request AddBookRequest

		// Проверка на корректность данных запроса
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			logger.InfoLog.Println("Bad request when creating a book: " + err.Error())
			return
		} else {
			logger.InfoLog.Println("Succesfulle to decoding request")
		}

		// Поиск или создание жанров
		var genres []models.Genre
		for _, genreName := range request.Genre {
			genreName = strings.ToLower(genreName)
			genreName = strings.ToUpper(string(genreName[0:2])) + genreName[2:]
			var genre models.Genre
			// Попытка найти жанр
			if err := db.Where("name = ?", genreName).First(&genre).Error; err != nil {
				// Если жанр не найден, создаем новый
				logger.InfoLog.Println(`Trying to find a genre "` + genreName + `"`)
				if errors.Is(err, gorm.ErrRecordNotFound) {
					logger.InfoLog.Println("No existing genre was found when creating the book")
					genre = models.Genre{Name: genreName}
					if err := db.Create(&genre).Error; err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create genre"})
						logger.ErrorLog.Println("Failed to create genre when creating the book: " + err.Error())
						return
					} else {
						logger.InfoLog.Println(`Genre "` + genreName + `" was created`)
					}
				} else {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve genre"})
					logger.ErrorLog.Println("Failed to retrieve genre: " + err.Error())
					return
				}
			}
			genres = append(genres, genre)
		}

		// Создание экземпляра книги на основе данных запроса
		book := models.Book{
			Title:         request.Title,
			Author:        request.Author,
			PublishedYear: request.Published_year,
			Genres:        genres,
			Description:   request.Description,
		}

		// Сохранение книги в базе данных
		if err := db.Create(&book).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
			logger.ErrorLog.Println("Failed to add book " + err.Error())
			return
		} else {
			logger.InfoLog.Println(`Book "` + book.Title + `" created in database`)
		}
		if producer == nil {
			logger.ErrorLog.Panicln("Kafka producer is nil!")
		}
		event := map[string]interface{}{
			"event": "BookAdded",
			"data":  book,
		}
		eventBytes, _ := json.Marshal(event)
		logger.InfoLog.Println("JSON sent to Kafka: ", string(eventBytes))
		if err := producer.SendMessage(string(eventBytes)); err != nil {
			logger.ErrorLog.Println("Failed to send event to Kafka: " + err.Error())
		} else {
			logger.InfoLog.Println("Sending the event to kafka was successful")
		}

		// Успешный ответ
		c.JSON(http.StatusCreated, gin.H{
			"message": "Book added successully!",
			"title":   request.Title,
		})
		logger.InfoLog.Println("Adding book was successful")
	}
}

// DeleteBook
// @Summary      Delete the book
// @Description  deletes the book from the library
// @Description  JWT authentication via cookie only for admin.
// @Description	 The JWT token should be stored in a cookie named "jwt".
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book  body  DeleteBookRequest  true  "Book Data"  example({"id": 1})
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /deleteBook [delete]
func DeleteBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		// DeleteBookRequest структура запроса для удаления книги
		// @Schema example={"id": 1}
		var request DeleteBookRequest

		// Проверяем входящий JSON
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Пытаемся найти книгу
		var book models.Book
		result := db.First(&book, request.ID)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// Если книга не найдена
				c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
				return
			}
			// Прочие ошибки
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book", "details": result.Error})
			return
		}

		// Удаляем книгу
		if err := db.Delete(&book).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book", "details": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Book deleted successully!",
			"ID":      request.ID,
			"Title":   book.Title,
		})

	}
}

// Modifying book
// @Summary      Modifying book
// @Description  JWT authentication via cookie only for admin.
// @Description	 The JWT token should be stored in a cookie named "jwt".
// @Description JWT Bearer authentcation only admin
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book  body  ModifyingBookRequest  true  "Book Data"  example({"id": 1, "published_year": "2021", "author": "Jeff Bezos", "title": "Why Does the World Exist?", "genre": ["Детектив"], "description": "Explore the ultimate question: Why is there something rather than nothing? This thought-provoking journey through philosophy, science, and metaphysics challenges readers to ponder existence itself, blending deep inquiry with accessible insight. A must-read for curious minds."})
// @Router       /modifyingBook [post]
func ModifyingBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Структура для запроса
		var request ModifyingBookRequest
		var book models.Book

		// Проверка на корректность данных запроса
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Загрузка книги с предзагрузкой жанров
		if err := db.Preload("Genres").Where("id = ?", request.Id).First(&book).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
			return
		}

		var genres []models.Genre
		if len(request.Genre) != 0 {
			// Поиск или создание жанров
			for _, genreName := range request.Genre {
				genreName = strings.ToLower(genreName)
				genreName = strings.ToUpper(string(genreName[0:2])) + genreName[2:] // Заглавная первая буква
				var genre models.Genre

				// Попытка найти жанр
				if err := db.Where("name = ?", genreName).First(&genre).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						genre = models.Genre{Name: genreName}
						if err := db.Create(&genre).Error; err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create genre"})
							return
						}
					} else {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve genre"})
						return
					}
				}
				genres = append(genres, genre)
			}
		}

		if request.Title != "" {
			book.Title = request.Title
		}
		if request.Published_year != "" {
			book.PublishedYear = request.Published_year
		}
		if request.Author != "" {
			book.Author = request.Author
		}
		if request.Description != "" {
			book.Description = request.Description
		}

		if len(request.Genre) > 0 {
			if err := db.Model(&book).Association("Genres").Clear(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear genres"})
				return
			}
			book.Genres = genres
		}

		// Сохранение измененной книги
		if err := db.Save(&book).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Успешный ответ
		c.JSON(http.StatusOK, gin.H{
			"message": "Book changed successfully!",
			"title":   book.Title,
		})
	}
}

// SearchBooks возвращает информацию о книгах со схожим названием или описанием
// @Summary      Outputs an array of books
// @Description  Returns an array of books that are similar in name or description to the request
// @Tags         book
// @Accept       json
// @Produce      json
// @Param 	search query string false "Looking for a similar book"
// @Success      200     {object} models.Book
// @Failure      400     {object} map[string]string
// @Failure      404     {object} map[string]string
// @Failure      500     {object} map[string]string
// @Router       /SearchBooks [get]
func SearchBooksHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		searchString := c.Query("search")
		similarity := 0.1 // Порог схожести

		books, err := database.SearchBooks(db, searchString, similarity)
		if err != nil {
			logger.ErrorLog.Println("Failed to search books\tError:", err)
		}

		var response []struct {
			ID            uint   `json:"id"`
			Title         string `json:"title"`
			Author        string `json:"author"`
			PublishedYear string `json:"published_year"`
			Genres        []struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			} `json:"genres"`
		}

		for _, book := range books {
			bookResponse := struct {
				ID            uint   `json:"id"`
				Title         string `json:"title"`
				Author        string `json:"author"`
				PublishedYear string `json:"published_year"`
				Genres        []struct {
					ID   uint   `json:"id"`
					Name string `json:"name"`
				} `json:"genres"`
			}{
				ID:            book.ID,
				Title:         book.Title,
				Author:        book.Author,
				PublishedYear: book.PublishedYear,
			}

			// Формируем список жанров
			for _, genre := range book.Genres {
				bookResponse.Genres = append(bookResponse.Genres, struct {
					ID   uint   `json:"id"`
					Name string `json:"name"`
				}{
					ID:   genre.ID,
					Name: genre.Name,
				})
			}
			response = append(response, bookResponse)
		}

		c.JSON(http.StatusOK, gin.H{"books": response})
	}
}
