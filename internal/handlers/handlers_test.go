package handlers_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"library/internal/auth"
	"library/internal/database"
	"library/internal/handlers"
	"library/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestWelcomeHandler(t *testing.T) {
	router := gin.Default()
	router.GET("/", handlers.Welcome)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	expectedBody := `{"message":"Welcome to the Library API!"}`
	assert.JSONEq(t, expectedBody, recorder.Body.String())
}

func TestRegisterUser(t *testing.T) {
	database.InitTestDB()
	defer database.CleanupTestDB()
	db := database.TestDB

	router := gin.Default()
	router.POST("/register", handlers.RegisterUser(db))

	requestBody := map[string]string{
		"name":     "",
		"email":    "",
		"password": "",
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	requestBody = map[string]string{
		"name":     "Test Name",
		"email":    "Test@example.com",
		"password": "password123",
	}
	body, _ = json.Marshal(requestBody)

	req, err = http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	expectedResponse := `{"message":"User registred successfully"}`
	assert.JSONEq(t, expectedResponse, recorder.Body.String())

	var user models.User
	err = db.First(&user, "email = ?", "Test@example.com").Error
	assert.NoError(t, err)
	assert.Equal(t, "Test Name", user.Name)
	assert.Equal(t, "Test@example.com", user.Email)

	hasher := sha256.New()
	hasher.Write([]byte("password123"))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	assert.Equal(t, hashedPassword, user.Password)
}

func TestLoginUser(t *testing.T) {
	database.InitTestDB()
	defer database.CleanupTestDB()
	db := database.TestDB

	hasher := sha256.New()
	hasher.Write([]byte("password123"))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	user := models.User{
		Name:     "Test Name",
		Email:    "Test@example.com",
		Password: hashedPassword,
		Role:     "reader",
	}
	err := db.Create(&user).Error
	assert.NoError(t, err)

	router := gin.Default()
	router.POST("/login", handlers.LoginUser(db))

	requestBody := map[string]string{
		"email":    "Test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var token string
	err = json.Unmarshal(recorder.Body.Bytes(), &token)
	assert.NoError(t, err, "Failser to unmarshal response body")

	assert.NotEmpty(t, token, "Token Is Empty")

	claims := &auth.MyClaims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("your_secret_key"), nil
	})
	assert.NoError(t, err, "Failed to parse token")
	assert.True(t, parsedToken.Valid, "Token is invalid")

	assert.Equal(t, user.Role, claims.Role)

	requestBody = map[string]string{
		"email":    "Invalid@example.com",
		"password": "password123",
	}
	body, _ = json.Marshal(requestBody)
	req, err = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, "Logged in using a non-existent email")

	requestBody = map[string]string{
		"email":    "Test@example.com",
		"password": "invalid123",
	}
	body, _ = json.Marshal(requestBody)
	req, err = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code, "Logged in with an incorrect password")
}

func TestGetBook(t *testing.T) {
	database.InitTestDB()
	defer database.CleanupTestDB()
	db := database.TestDB

	book := models.Book{
		Title:         "Test title",
		Author:        "Test author",
		PublishedYear: "2025",
		Description:   "Test description",
	}

	err := db.Create(&book).Error
	assert.NoError(t, err)

	router := gin.Default()
	router.GET("/getBook", handlers.GetBook(db))

	req, err := http.NewRequest(http.MethodGet, "/getBook?bookId=1", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBook models.Book
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBook)
	assert.NoError(t, err)
	assert.Equal(t, book.Title, responseBook.Title)
	assert.Equal(t, book.Author, responseBook.Author)
	assert.Equal(t, book.PublishedYear, responseBook.PublishedYear)
	assert.Equal(t, book.Description, responseBook.Description)
	assert.Equal(t, uint(1), responseBook.ID)

	req, err = http.NewRequest(http.MethodGet, "/getBook?bookId=2", nil)
	assert.NoError(t, err)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.Equal(t, `{"error":"Book not found"}`, recorder.Body.String())

	req, err = http.NewRequest(http.MethodGet, "/getBook", nil)
	assert.NoError(t, err)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Equal(t, `{"error":"Missing bookId parameter"}`, recorder.Body.String())
}
