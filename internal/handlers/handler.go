package handlers

import (
	"errors"
	"library/internal/database"
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

// DeleteBookRequest структура запроса для удаления книги
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
func GetBooks(c *gin.Context) {
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
	if err := database.DB.Preload("Genres", func(db *gorm.DB) *gorm.DB {
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

// AddBook
// @Summary      Add a new book
// @Description  Add a new book to the library
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book  body  AddBookRequest  true  "Book Data"  example({"title": "Golang Basics", "author": "John Doe", "published_year": "2024", "genre": ["Учебная литература"], "description": "Эта книга — идеальный выбор для тех, кто хочет начать свое путешествие в программировании на языке Go."})
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /addBooks [post]
func AddBook(c *gin.Context) {
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
		if err := database.DB.Where("name = ?", genreName).First(&genre).Error; err != nil {
			// Если жанр не найден, создаем новый
			if errors.Is(err, gorm.ErrRecordNotFound) {
				genre = models.Genre{Name: genreName}
				if err := database.DB.Create(&genre).Error; err != nil {
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
	if err := database.DB.Create(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add book"})
		return
	}

	// Успешный ответ
	c.JSON(http.StatusCreated, gin.H{
		"message": "Book added successully!",
		"title":   request.Title,
	})
}

// DeleteBook
// @Summary      Delete the book
// @Description  deletes the book from the library
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book  body  DeleteBookRequest  true  "Book Data"  example({"id": "1"})
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /deleteBook [post]
func DeleteBook(c *gin.Context) {
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
	result := database.DB.First(&book, request.ID)
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
	if err := database.DB.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Book deleted successully!",
		"ID":      request.ID,
		"Title":   book.Title,
	})

}
