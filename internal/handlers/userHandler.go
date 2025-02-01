package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"library/internal/auth"
	"library/internal/models"
	"library/logger"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// RegisterUserRequest структура запроса для регистрации пользователя
// @Schema example={"name": "Vladislav", "email": "Laminano@mail.ru", "password":"123456", "mailing":true}
type RegisterUserRequest struct {
	Name     string `json:"name" binding:"required" example:"Vladislav"`
	Email    string `json:"email" binding:"required,email" example:"Laminano@mail.ru"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
	Mailing  bool   `json:"mailing" binding:"required" example:"true"`
}

// LoginRequest структура запроса для авторизации пользователя
// @Schema example={"email": "Laminano@mail.ru", "password":"123456"}
type LoginRequest struct {
	Email    string `json:"email" binding:"required" exmple:"Laminano@mail.ru"`
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
// @Param        user  body  RegisterUserRequest  true  "User Data" example({"name": "Vladislav", "email": "Laminano@mail.ru", "password":"123456", "mailing":true})
// @Success 201 {object} map[string]string
// @Router       /register [post]
func RegisterUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			Mailing:  request.Mailing,
			Role:     "reader",
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User registred successfully"})
	}
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
func LoginUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request LoginRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		var User models.User
		if err := db.Where("email = ?", request.Email).First(&User).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err, "message": "email is not registered"})
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

		timeSec, err := strconv.Atoi(os.Getenv("JWTCoo_expires_time_sec"))
		if err != nil {
			logger.ErrorLog.Println("Failed get `JWTCoo_expires_time_sec` in .env when login user\tError:", err)
		}
		c.SetCookie("jwt", token, timeSec, "/", os.Getenv("domain"), false, true)

		refreshToken := auth.GenerateRefreshToken()
		User.RefreshToken = refreshToken
		User.ExpiresAt = time.Now().Add(720 * time.Hour)
		if err = db.Save(&User).Error; err != nil {
			logger.ErrorLog.Println("Error save refresh token in db when logining\t Error:", err)
		}
		c.SetCookie("refreshToken", refreshToken, 2592000, "/", os.Getenv("domain"), false, true)

		c.JSON(http.StatusOK, gin.H{"message": "User authorization successfully"})
	}
}

// UnsubscribeMailing
// @Summary      Unsubscribe mailing
// @Description  Describes the user from the mailing list
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /unsubMailing [get]
// @name         UnsubscribeMailing
func UnsubscribeMailing(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		handlMailing(db, false, c)
	}
}

// SubscribeMailing
// @Summary      Subscribe mailing
// @Description  Subscribes a user to mailing lists
// @Tags         user
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /subMailing [get]
// @name         SubscribeMailing
func SubscribeMailing(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		handlMailing(db, true, c)
	}
}

// Это общая функция дл двух эндпоинтов, описанных выше(иначе swagger не разделяет 2 эндпоинта)
func handlMailing(db *gorm.DB, subscribe bool, c *gin.Context) {
	logger.InfoLog.Println("Getting jwt token from cookies")
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
		logger.InfoLog.Println("Error when getting jwt from cookies.\tJWT token == nil:", (tokenString == ""), "\nError:", err)
		c.Abort()
		return
	}

	logger.InfoLog.Println("Validating jwt token")
	token, err := jwt.ParseWithClaims(tokenString, &auth.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("jwtSecret")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		logger.InfoLog.Println("Error when validating jwt token.\ttoken.Valid:", token.Valid, "\tError: ", err)
		c.Abort()
		return
	}

	claims, ok := token.Claims.(*auth.MyClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Claims"})
		c.Abort()
		return
	}

	if !claims.Mailing && !subscribe {
		c.JSON(http.StatusOK, gin.H{"message": "You have already unsubscribed from the mailing list"})
		logger.InfoLog.Println("User already unsubscribed from the mailing list")
		c.Abort()
		return
	}

	if claims.Mailing && subscribe {
		c.JSON(http.StatusOK, gin.H{"message": "You have already subscribed to the mailing list"})
		logger.InfoLog.Println("User already subscribed to the mailing list")
		c.Abort()
		return
	}

	var user models.User

	if err := db.Where("id = ?", claims.Subject).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.SetCookie("jwt", "", -1, "/", os.Getenv("domain"), false, true)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
	}

	var message string

	if !subscribe {
		user.Mailing = false
		if err := db.Save(&user).Error; err != nil {
			logger.ErrorLog.Println("Failes unsubscribe from the mailing list\tError:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logger.InfoLog.Println("User unsubscribed from the mailing list. User ID:", claims.Subject)
		message = "You have unsubscribed from the mailing list"

	} else {
		user.Mailing = true
		if err := db.Save(&user).Error; err != nil {
			logger.ErrorLog.Println("Failes subscribe to the mailing list\tError:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		logger.InfoLog.Println("User subscribed to the mailing list. User ID:", claims.Subject)
		message = "You have subscribed to the mailing list"
	}
	logger.InfoLog.Println("Generating new jwt token with changed mailing =", user.Mailing)
	tokenString, err = auth.GenerateJWT(user)
	if err != nil {
		logger.ErrorLog.Println("Failing to generate new JWT token\tError:", err)
	}

	timeSec, err := strconv.Atoi(os.Getenv("JWTCoo_expires_time_sec"))
	if err != nil {
		logger.ErrorLog.Println("Failed get `JWTCoo_expires_time_sec` in .env when login user\tError:", err)
	}
	c.SetCookie("jwt", tokenString, timeSec, "/", os.Getenv("domain"), false, true)
	c.JSON(http.StatusOK, gin.H{"message": message})
}
