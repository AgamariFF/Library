package auth

import (
	"fmt"
	"library/internal/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

type MyClaims struct {
	Role    string `json:"role"`
	Mailing bool   `json:"mailing"`
	jwt.RegisteredClaims
}

func GenerateJWT(user models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

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
