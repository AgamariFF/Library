package auth

import (
	"fmt"
	"library/internal/models"
	"library/logger"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var jwtSecret []byte

type MyClaims struct {
	Role    string `json:"role"`
	Mailing bool   `json:"mailing"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User) (string, error) {
	timeSec, err := strconv.Atoi(os.Getenv("JWTCoo_expires_time_sec"))
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(time.Duration(timeSec) * time.Second)

	claims := &MyClaims{
		Role:    user.Role,
		Mailing: user.Mailing,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.ID),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("jwtSecret")))
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	jwtSecret = []byte(os.Getenv("jwtSecret"))
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func GenerateRefreshToken() string {
	return uuid.New().String()
}

func UpdateJWTToken(c *gin.Context, db *gorm.DB) (string, error) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil || refreshToken == "" {
		logger.InfoLog.Println("Error when trying to update JWT token\tError:", err)
		return "", fmt.Errorf("refresh token not found")
	}

	var user models.User

	if err = db.Where("refresh_token = ?", refreshToken).First(&user).Error; err != nil {
		logger.InfoLog.Println("Error when trying to find refresher token in db\tError:", err)
		return "", fmt.Errorf("invalid refresh token")
	}

	if user.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("refresh token expired")
	}

	refreshToken = GenerateRefreshToken()
	user.RefreshToken = refreshToken
	user.ExpiresAt = time.Now().Add(720 * time.Hour)
	c.SetCookie("refreshToken", refreshToken, 2592000, "/", os.Getenv("domain"), false, true)
	if err = db.Save(user).Error; err != nil {
		logger.ErrorLog.Println("Faile to update save user with update refresh token and expired time id DB\tError:", err)
		return "", fmt.Errorf("failed to update refrest token")
	}

	jwtToken, err := GenerateJWT(user)
	if err != nil {
		logger.ErrorLog.Println("Failed to generate new JWT token\tError:", err)
		return "", err
	}
	timeSec, err := strconv.Atoi(os.Getenv("JWTCoo_expires_time_sec"))
	if err != nil {
		logger.ErrorLog.Println("Failed get `JWTCoo_expires_time_sec` in .env when login user\tError:", err)
	}
	c.SetCookie("jwt", jwtToken, timeSec, "/", os.Getenv("domain"), false, true)
	return jwtToken, nil
}
