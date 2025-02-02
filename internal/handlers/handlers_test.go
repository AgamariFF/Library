package handlers_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"library/internal/database"
	"library/internal/handlers"
	"library/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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

	requestBody := map[string]interface{}{
		"name":     "",
		"email":    "",
		"password": "",
		"mailing":  "",
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

	requestBody = map[string]interface{}{
		"name":     "Test Name",
		"email":    "Test@example.com",
		"password": "password123",
		"mailing":  true,
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

// func TestLoginUser(t *testing.T) {
// 	database.InitTestDB()
// 	defer database.CleanupTestDB()
// 	db := database.TestDB

// 	hasher := sha256.New()
// 	hasher.Write([]byte("password123"))
// 	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

// 	user := models.User{
// 		Name:     "Test Name",
// 		Email:    "Test@example.com",
// 		Password: hashedPassword,
// 		Role:     "reader",
// 		Mailing:  false,
// 	}
// 	err := db.Create(&user).Error
// 	assert.NoError(t, err)

// 	router := gin.Default()
// 	router.POST("/login", handlers.LoginUser(db))

// 	requestBody := map[string]string{
// 		"email":    "Test@example.com",
// 		"password": "password123",
// 	}
// 	body, _ := json.Marshal(requestBody)

// 	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
// 	assert.NoError(t, err)
// 	req.Header.Set("Content-Type", "application/json")

// 	recorder := httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)

// 	assert.Equal(t, http.StatusOK, recorder.Code)

// 	var token string
// 	err = json.Unmarshal(recorder.Body.Bytes(), &token)
// 	assert.NoError(t, err, "Failser to unmarshal response body")

// 	assert.NotEmpty(t, token, "Token Is Empty")

// 	claims := &auth.MyClaims{}
// 	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
// 		return []byte(os.Getenv("jwtSecret")), nil
// 	})
// 	assert.NoError(t, err, "Failed to parse token")
// 	assert.True(t, parsedToken.Valid, "Token is invalid")

// 	assert.Equal(t, user.Role, claims.Role)

// 	requestBody = map[string]string{
// 		"email":    "Invalid@example.com",
// 		"password": "password123",
// 	}
// 	body, _ = json.Marshal(requestBody)
// 	req, err = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
// 	assert.NoError(t, err)
// 	recorder = httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)
// 	assert.Equal(t, http.StatusUnauthorized, recorder.Code, "Logged in using a non-existent email")

// 	requestBody = map[string]string{
// 		"email":    "Test@example.com",
// 		"password": "invalid123",
// 	}
// 	body, _ = json.Marshal(requestBody)
// 	req, err = http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
// 	assert.NoError(t, err)
// 	recorder = httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)
// 	assert.Equal(t, http.StatusUnauthorized, recorder.Code, "Logged in with an incorrect password")
// }

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

// func TestAddBook(t *testing.T) {
// 	logger.InitLog()
// 	database.InitTestDB()
// 	db := database.TestDB
// 	defer database.CleanupTestDB()
// 	router := gin.Default()
// 	producer, _ := kafka.NewKafkaProducer([]string{"localhost:9092"}, "library-events")
// 	defer producer.Close()
// 	router.POST("/addBook", handlers.AddBook(db, producer))

// 	book := `{
//   "author": "John Doe",
//   "description": "Эта книга — идеальный выбор для тех, кто хочет начать свое путешествие в программировании на языке Go.",
//   "genre": [
//     "Учебная литература",
// 	"Программирование"
//   ],
//   "published_year": "2024",
//   "title": "Golang Basics"
// }`

// 	req, err := http.NewRequest(http.MethodPost, "/addBook", bytes.NewBuffer([]byte(book)))
// 	assert.NoError(t, err)
// 	req.Header.Set("Content-Type", "application/json")

// 	recorder := httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)

// 	assert.Equal(t, http.StatusCreated, recorder.Code)

// 	var gettedBook models.Book
// 	err = db.Where("title = ?", "Golang Basics").Preload("Genres").First(&gettedBook).Error
// 	assert.NoError(t, err)

// 	assert.Equal(t, "John Doe", gettedBook.Author)
// 	assert.Equal(t, "Эта книга — идеальный выбор для тех, кто хочет начать свое путешествие в программировании на языке Go.", gettedBook.Description)
// 	assert.Equal(t, "2024", gettedBook.PublishedYear)
// 	assert.Equal(t, "Golang Basics", gettedBook.Title)
// 	assert.Len(t, gettedBook.Genres, 2)

// 	genreNames := []string{gettedBook.Genres[0].Name, gettedBook.Genres[1].Name}
// 	assert.Contains(t, genreNames, "Учебная литература")
// 	assert.Contains(t, genreNames, "Программирование")

// 	var genresCount int64
// 	err = db.Model(&models.Genre{}).Count(&genresCount).Error
// 	assert.NoError(t, err)
// 	assert.Equal(t, int64(2), genresCount)

// 	var genre1, genre2 models.Genre
// 	err = db.Where("name = ?", "Учебная литература").First(&genre1).Error
// 	assert.NoError(t, err)
// 	err = db.Where("name = ?", "Программирование").First(&genre2).Error
// 	assert.NoError(t, err)

// 	assert.Equal(t, "Учебная литература", genre1.Name)
// 	assert.Equal(t, "Программирование", genre2.Name)

// 	var responseBody map[string]interface{}
// 	err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Book added successully!", responseBody["message"])
// 	assert.Equal(t, "Golang Basics", responseBody["title"])

// 	book = `{
//   "author": "",
//   "description": "Эта книга — идеальный выбор для тех, кто хочет начать свое путешествие в программировании на языке Go.",
//   "genre": [
//     "Учебная литература",
// 	"Программирование"
//   ],
//   "published_year": "2024",
//   "title": "Golang Basics"
// }`
// 	req, err = http.NewRequest(http.MethodPost, "/addBook", bytes.NewBuffer([]byte(book)))
// 	assert.NoError(t, err)
// 	req.Header.Set("Content-Type", "application/json")

// 	recorder = httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)

// 	assert.Equal(t, http.StatusBadRequest, recorder.Code)
// 	var errorResponseBody map[string]interface{}
// 	err = json.Unmarshal(recorder.Body.Bytes(), &errorResponseBody)
// 	assert.NoError(t, err)
// 	assert.Equal(t, `Key: 'AddBookRequest.Author' Error:Field validation for 'Author' failed on the 'required' tag`, errorResponseBody["error"])

// 	book = `{
//   "author": "Test author",
//   "description": "Test description",
//   "genre": [],
//   "published_year": "2022",
//   "title": "Test title"
// }`
// 	req, err = http.NewRequest(http.MethodPost, "/addBook", bytes.NewBuffer([]byte(book)))
// 	assert.NoError(t, err)
// 	req.Header.Set("Content-Type", "application/json")

// 	recorder = httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)
// 	gettedBook = models.Book{}
// 	assert.Equal(t, http.StatusCreated, recorder.Code)
// 	err = db.Where("title = ?", "Test title").Preload("Genres").First(&gettedBook).Error
// 	assert.NoError(t, err)

// 	assert.Equal(t, "Test title", gettedBook.Title)
// 	assert.Len(t, gettedBook.Genres, 0)

// 	book = `{
//   "author": "Test author0",
//   "description": "Test description0",
//   "genre": ["Тест", "Тест"],
//   "published_year": "2021",
//   "title": "Test title0"
// }`
// 	req, err = http.NewRequest(http.MethodPost, "/addBook", bytes.NewBuffer([]byte(book)))
// 	assert.NoError(t, err)
// 	req.Header.Set("Content-Type", "application/json")

// 	recorder = httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)
// 	gettedBook = models.Book{}
// 	assert.Equal(t, http.StatusCreated, recorder.Code)
// 	err = db.Where("title = ?", "Test title0").Preload("Genres").First(&gettedBook).Error
// 	assert.NoError(t, err)

// 	assert.Equal(t, "Test title0", gettedBook.Title)
// 	assert.Len(t, gettedBook.Genres, 1)
// 	genre1 = models.Genre{}
// 	err = db.Where("name = ?", "Тест").First(&genre1).Error
// 	assert.NoError(t, err)
// 	assert.Equal(t, "Тест", genre1.Name)

// 	book = `{
//   "author": "Test author1",
//   "description": "Test description1",
//   "genre": ["Тест",
//   "published_year": "2025",
//   "title": "Test title1"
// }`
// 	req, err = http.NewRequest(http.MethodPost, "/addBook", bytes.NewBuffer([]byte(book)))
// 	assert.NoError(t, err)
// 	req.Header.Set("Content-Type", "application/json")

// 	recorder = httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)
// 	gettedBook = models.Book{}
// 	assert.Equal(t, http.StatusBadRequest, recorder.Code)
// }

func TestGetBooks(t *testing.T) {
	database.InitTestDB()
	db := database.TestDB
	defer database.CleanupTestDB()

	book1 := models.Book{
		Title:         "Beta title",
		Author:        "Alpha author",
		PublishedYear: "2025",
		Genres: []models.Genre{
			{
				Name: "Test name1",
			},
			{
				Name: "Test name2",
			},
			{
				Name: "Test name3",
			},
		},
		Description: "Test description1",
	}

	err := db.Create(&book1).Error
	assert.NoError(t, err)

	book2 := models.Book{
		Title:         "Alpha title",
		Author:        "Charlie author",
		PublishedYear: "2023",
		Genres: []models.Genre{
			{
				Name: "Test name1",
			},
		},
		Description: "Test description2",
	}

	err = db.Create(&book2).Error
	assert.NoError(t, err)

	book3 := models.Book{
		Title:         "Charlie title",
		Author:        "Beta author",
		PublishedYear: "2021",
		Genres: []models.Genre{
			{
				Name: "Test name1",
			},
			{
				Name: "Test name2",
			},
		},
		Description: "Test description3",
	}

	err = db.Create(&book3).Error
	assert.NoError(t, err)

	router := gin.Default()
	router.GET("/getBooks", handlers.GetBooks(db))

	req, err := http.NewRequest(http.MethodGet, "/getBooks", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var responseBooks []models.Book
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBooks)
	assert.NoError(t, err)

	assert.Len(t, responseBooks, 3)
	assert.Equal(t, book1.Author, responseBooks[0].Author)
	assert.Equal(t, book1.PublishedYear, responseBooks[0].PublishedYear)
	assert.Equal(t, book1.Title, responseBooks[0].Title)
	assert.Len(t, responseBooks[0].Genres, 3)

	assert.Equal(t, book2.Author, responseBooks[1].Author)
	assert.Equal(t, book2.PublishedYear, responseBooks[1].PublishedYear)
	assert.Equal(t, book2.Title, responseBooks[1].Title)
	assert.Len(t, responseBooks[1].Genres, 1)

	assert.Equal(t, book3.Author, responseBooks[2].Author)
	assert.Equal(t, book3.PublishedYear, responseBooks[2].PublishedYear)
	assert.Equal(t, book3.Title, responseBooks[2].Title)
	assert.Len(t, responseBooks[2].Genres, 2)

	req, err = http.NewRequest(http.MethodGet, "/getBooks?sort=author", nil)
	assert.NoError(t, err)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	responseBooks = []models.Book{}
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBooks)
	assert.NoError(t, err)

	assert.Len(t, responseBooks, 3)
	assert.Equal(t, book1.Author, responseBooks[0].Author)
	assert.Equal(t, book1.PublishedYear, responseBooks[0].PublishedYear)
	assert.Equal(t, book1.Title, responseBooks[0].Title)

	assert.Equal(t, book2.Author, responseBooks[2].Author)
	assert.Equal(t, book2.PublishedYear, responseBooks[2].PublishedYear)
	assert.Equal(t, book2.Title, responseBooks[2].Title)

	assert.Equal(t, book3.Author, responseBooks[1].Author)
	assert.Equal(t, book3.PublishedYear, responseBooks[1].PublishedYear)
	assert.Equal(t, book3.Title, responseBooks[1].Title)

	req, err = http.NewRequest(http.MethodGet, "/getBooks?sort=year", nil)
	assert.NoError(t, err)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	responseBooks = []models.Book{}
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBooks)
	assert.NoError(t, err)

	assert.Len(t, responseBooks, 3)
	assert.Equal(t, book1.Author, responseBooks[2].Author)
	assert.Equal(t, book1.PublishedYear, responseBooks[2].PublishedYear)
	assert.Equal(t, book1.Title, responseBooks[2].Title)

	assert.Equal(t, book2.Author, responseBooks[1].Author)
	assert.Equal(t, book2.PublishedYear, responseBooks[1].PublishedYear)
	assert.Equal(t, book2.Title, responseBooks[1].Title)

	assert.Equal(t, book3.Author, responseBooks[0].Author)
	assert.Equal(t, book3.PublishedYear, responseBooks[0].PublishedYear)
	assert.Equal(t, book3.Title, responseBooks[0].Title)

	req, err = http.NewRequest(http.MethodGet, "/getBooks?sort=title", nil)
	assert.NoError(t, err)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	responseBooks = []models.Book{}
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBooks)
	assert.NoError(t, err)

	assert.Len(t, responseBooks, 3)
	assert.Equal(t, book1.Author, responseBooks[1].Author)
	assert.Equal(t, book1.PublishedYear, responseBooks[1].PublishedYear)
	assert.Equal(t, book1.Title, responseBooks[1].Title)

	assert.Equal(t, book2.Author, responseBooks[0].Author)
	assert.Equal(t, book2.PublishedYear, responseBooks[0].PublishedYear)
	assert.Equal(t, book2.Title, responseBooks[0].Title)

	assert.Equal(t, book3.Author, responseBooks[2].Author)
	assert.Equal(t, book3.PublishedYear, responseBooks[2].PublishedYear)
	assert.Equal(t, book3.Title, responseBooks[2].Title)

	req, err = http.NewRequest(http.MethodGet, "/getBooks?sort=zxc", nil)
	assert.NoError(t, err)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	responseBooks = []models.Book{}
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBooks)
	assert.NoError(t, err)

	assert.Len(t, responseBooks, 3)
	assert.Equal(t, book1.Author, responseBooks[0].Author)
	assert.Equal(t, book1.PublishedYear, responseBooks[0].PublishedYear)
	assert.Equal(t, book1.Title, responseBooks[0].Title)

	assert.Equal(t, book2.Author, responseBooks[1].Author)
	assert.Equal(t, book2.PublishedYear, responseBooks[1].PublishedYear)
	assert.Equal(t, book2.Title, responseBooks[1].Title)

	assert.Equal(t, book3.Author, responseBooks[2].Author)
	assert.Equal(t, book3.PublishedYear, responseBooks[2].PublishedYear)
	assert.Equal(t, book3.Title, responseBooks[2].Title)
}

func TestDeleteBook(t *testing.T) {
	database.InitTestDB()
	db := database.TestDB
	defer database.CleanupTestDB()

	book := models.Book{
		Title:         "Test title",
		Author:        "Test author",
		PublishedYear: "2025",
		Description:   "Test description",
	}

	err := db.Create(&book).Error
	assert.NoError(t, err)

	router := gin.Default()
	router.DELETE("/DeleteBook", handlers.DeleteBook(db))

	requestBody := handlers.DeleteBookRequest{
		ID: 2,
	}
	body, _ := json.Marshal(requestBody)

	req, err := http.NewRequest(http.MethodDelete, "/DeleteBook", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)
	var responseBody map[string]interface{}
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "Book not found", responseBody["error"])

	requestBody = handlers.DeleteBookRequest{
		ID: 1,
	}
	body, _ = json.Marshal(requestBody)

	req, err = http.NewRequest(http.MethodDelete, "/DeleteBook", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	err = json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, "Book deleted successully!", responseBody["message"])
	err = db.Where("title = ?", "Test title").First(&book).Error
	assert.NotNil(t, err)
}

func TestModifyingBook(t *testing.T) {
	database.InitTestDB()
	db := database.TestDB
	defer database.CleanupTestDB()

	book := models.Book{
		Title:         "Test title",
		Author:        "Test author",
		PublishedYear: "2025",
		Genres: []models.Genre{
			{
				Name: "Тест1",
			},
			{
				Name: "Тест2",
			},
		},
		Description: "Test description",
	}

	err := db.Create(&book).Error
	assert.NoError(t, err)

	requestBody := handlers.ModifyingBookRequest{
		Id:             2,
		Title:          "Changed title",
		Author:         "Changed author",
		Genre:          []string{"Изменен1", "Изменен2", "Изменен3"},
		Published_year: "2000",
		Description:    "Changed description",
	}

	body, _ := json.Marshal(requestBody)

	router := gin.Default()
	router.POST("/ModifyingBook", handlers.ModifyingBook(db))
	req, err := http.NewRequest(http.MethodPost, "/ModifyingBook", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	var responseBody map[string]interface{}
	json.Unmarshal(recorder.Body.Bytes(), &responseBody)

	assert.Equal(t, "Book not found", responseBody["error"])

	requestBody = handlers.ModifyingBookRequest{
		Id:             1,
		Title:          "Changed title",
		Author:         "Changed author",
		Genre:          []string{"Изменен1", "Изменен2", "Изменен3"},
		Published_year: "2000",
		Description:    "Changed description",
	}

	body, _ = json.Marshal(requestBody)

	req, err = http.NewRequest(http.MethodPost, "/ModifyingBook", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	responseBody = nil
	json.Unmarshal(recorder.Body.Bytes(), &responseBody)
	assert.Equal(t, "Book changed successfully!", responseBody["message"])

	var responseBook models.Book
	err = db.Preload("Genres").Where("title = ?", requestBody.Title).First(&responseBook).Error
	assert.NoError(t, err)
	assert.Equal(t, requestBody.Title, requestBody.Title)
	assert.Equal(t, requestBody.Author, requestBody.Author)
	assert.Equal(t, requestBody.Published_year, requestBody.Published_year)
	assert.Equal(t, requestBody.Description, requestBody.Description)
	assert.Len(t, responseBook.Genres, 3)
	assert.Equal(t, "Изменен1", responseBook.Genres[0].Name)

	invalidRequestBody := `{
  "genre": ["Тест",
}`

	req, err = http.NewRequest(http.MethodPost, "/ModifyingBook", bytes.NewBuffer([]byte(invalidRequestBody)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)

}
