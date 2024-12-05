package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"library/internal/auth"
	"library/internal/database"
	"library/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterUserRequest структура запроса для регистрации пользователя
// @Schema example={"name": "Vladislav", "email": "Laminano@mail.ru", "password":"123456"}
type RegisterUserRequest struct {
	Name     string `json:"name" binding:"required" example:"Vladislav"`
	Email    string `json:"email" binding:"required,email" example:"Laminano@mail.ru"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
}

// LoginRequest структура запроса для авторизации пользователя
// @Schema example={"email": "Laminano@mail.ru", "password":"123456"}
type LoginRequest struct {
	Email    string `json:"email" binding:"required" exmple:"Laminano@mail.ru`
	Password string `json:"password" binding:"required" example:"123456"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

// RegisterUser
// @Summary      Add a new User
// @Description  Add a new library User
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body  RegisterUserRequest  true  "User Data" example({"name": "Vladislav", "email": "Laminano@mail.ru", "password":"123456"})
// @Success 201 {object} map[string]string
// @Router       /register [post]
func RegisterUser(c *gin.Context) {
	var request RegisterUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(request.Password))
	hasherPassword := hex.EncodeToString(hasher.Sum(nil))

	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: hasherPassword,
		Role:     "reader",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registred successfully"})
}

// LoginUser
// @Summary      Performs user login
// @Description  Logs in an existing user
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user  body  LoginRequest  true  "User Data" example({"email": "Laminano@mail.ru", "password":"123456"})
// @Success 201 {object} map[string]string
// @Router       /login [post]
func LoginUser(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	var User models.User
	if err := database.DB.Where("email = ?", request.Email).First(&User).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(request.Password))
	hasherPassword := hex.EncodeToString(hasher.Sum(nil))
	if hasherPassword != User.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	token, err := auth.GenerateJWT(User)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, token)
}
