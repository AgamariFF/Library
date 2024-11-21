package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

// GetBooks возвращает список книг
// @Summary Get list of books
// @Description Retrieve all books
// @Tags book
// @Accept json
// @Produce json
// @Success 200 {array} string
// @Router /getBooks [get]
func GetBooks(c *gin.Context) {
	books := []string{"Book 1", "Book 2", "Book 3"}
	c.JSON(http.StatusOK, gin.H{
		"books": books,
	})
}

// AddBook
// @Summary      Add a new book
// @Description  Add a new book to the library
// @Tags         book
// @Accept       json
// @Produce      json
// @Param        book  body      map[string]string  true  "Book Title"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Router       /addBooks [post]
func AddBook(c *gin.Context) {
	var request struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Book added successully!",
		"title":   request.Title,
	})
}
