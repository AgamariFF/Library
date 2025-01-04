package handlers

import (
	"errors"
	"library/internal/models"
	"net/http"
	"sort"
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
// @Schema example={"id": 1, "published_year": "2024"}
type ModifyingBookRequest struct {
	Id             uint     `json:"id" binding:"required" example:"1"`
	Title          string   `json:"title" swaggerignore:"true"`       // Название книги
	Author         string   `json:"author" swaggerignore:"true"`      // Автор
	Genre          []string `json:"genre" swaggerignore:"true"`       // Жанра
	Published_year string   `json:"published_year" example:"2021"`    // Год публикации
	Description    string   `json:"description" swaggerignore:"true"` // Описание книги
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

// GetBooks возвращает отсортированный или неотсортированный список книг
// @Summary Get list of books
// @Description Retrieve all books, optionally sorted by a specific field
// @Tags book
// @Accept json
// @Produce json
// @Param sort query string false "Field to sort by (e.g., 'id')"
// @Success 200 {array} models.Book
// @Failure 500 {object} map[string]string
// @Router /getBooks [get]
func GetBooks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var books []models.Book
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
		// Извлечение query-параметров
		sortParam := c.Query("sort")

		// Получаем все книги из БД
		if err := db.Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Select("genres.id, genres.name") // Выбираем только нужные поля
		}).Find(&books).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to get books",
			})
			return
		}

		// Сортируем книги
		switch sortParam {
		case "author":
			sort.Slice(books, func(i, j int) bool {
				return books[i].Author < books[j].Author
			})
		case "title":
			sort.Slice(books, func(i, j int) bool {
				return books[i].Title < books[j].Title
			})
		case "year":
			sort.Slice(books, func(i, j int) bool {
				return books[i].PublishedYear < books[j].PublishedYear
			})
		}

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
			response = append(response, bookResponse)
		}
		// Возврат ответа
		c.JSON(http.StatusOK, response)
	}
}

// GetBook возвращает информацию об одной книге
// @Summary      Get one book
// @Description JWT Bearer authentcation
// @Description  Get detailed information about a single book by ID
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        bookId  query    integer  true  "Book ID"
// @Success      200     {object} models.Book
// @Failure      400     {object} map[string]string
// @Failure      404     {object} map[string]string
// @Failure      500     {object} map[string]string
// @Security BearerAuth
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
		if err := db.Where("id = ?", bookId).First(&book).Error; err != nil {
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
// @Description  Add a new book to the library
// @Description JWT Bearer authentcation only admin
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book  body  AddBookRequest  true  "Book Data"  example({"title": "Golang Basics", "author": "John Doe", "published_year": "2024", "genre": ["Учебная литература"], "description": "Эта книга — идеальный выбор для тех, кто хочет начать свое путешествие в программировании на языке Go."})
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Security BearerAuth
// @Router       /addBook [post]
func AddBook(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Структура для запроса
		var request AddBookRequest

		// Проверка на корректность данных запроса
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
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
			return
		}

		// Успешный ответ
		c.JSON(http.StatusCreated, gin.H{
			"message": "Book added successully!",
			"title":   request.Title,
		})
	}
}

// DeleteBook
// @Summary      Delete the book
// @Description  deletes the book from the library
// @Description JWT Bearer authentcation only admin
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book  body  DeleteBookRequest  true  "Book Data"  example({"id": 1})
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security BearerAuth
// @Router       /deleteBook [post]
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
// @Description  Modifies the data of an existing workbook
// @Description JWT Bearer authentcation only admin
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book  body  ModifyingBookRequest  true  "Modifying book"  example({"id: "1", "published_year": "2021"})
// @Security BearerAuth
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
