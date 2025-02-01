package middleware_test

import (
	"fmt"
	"library/internal/auth"
	"library/internal/database"
	"library/internal/middleware"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// func TestJWTMiddlew(t *testing.T) {
// 	expirationTime := time.Now().Add(time.Minute)

// 	claims := &auth.MyClaims{
// 		Role: "testRole",
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			Subject:   fmt.Sprintf("%d", 2),
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 		},
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	signedToken, err := token.SignedString([]byte(os.Getenv("jwtSecret")))
// 	assert.NoError(t, err)

// 	router := gin.Default()
// 	router.GET("/test", middleware.JWTMiddleware(), func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"message": "Success"})
// 	})

// 	req, err := http.NewRequest(http.MethodGet, "/test", nil)
// 	assert.NoError(t, err)
// 	req.Header.Set("Authorization", signedToken)

// 	recorder := httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)
// 	assert.Equal(t, http.StatusOK, recorder.Code)

// 	req, err = http.NewRequest(http.MethodGet, "/test", nil)
// 	assert.NoError(t, err)

// 	recorder = httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)
// 	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

// 	invalidToken := "invalid.toke.value"
// 	req, err = http.NewRequest(http.MethodGet, "/test", nil)
// 	assert.NoError(t, err)
// 	req.Header.Set("Authorization", invalidToken)

// 	recorder = httptest.NewRecorder()
// 	router.ServeHTTP(recorder, req)
// 	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
// }

func TestRoleMiddleware(t *testing.T) {
	expirationTime := time.Now().Add(time.Minute)

	claims := &auth.MyClaims{
		Role: "testRole",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", 2),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("jwtSecret")))
	assert.NoError(t, err)

	router := gin.Default()
	router.GET("/test", middleware.RoleMiddleware(database.TestDB, "testRole"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	req, err := http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", signedToken)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Code)

	claims = &auth.MyClaims{
		Role: "Role",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", 2),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(os.Getenv("jwtSecret")))
	assert.NoError(t, err)

	req, err = http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", signedToken)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusForbidden, recorder.Code)

	req, err = http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err)
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)

	invalidToken := "invalid.toke.value"
	req, err = http.NewRequest(http.MethodGet, "/test", nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", invalidToken)

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}
