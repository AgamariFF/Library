package middleware

import (
	"fmt"
	"library/internal/auth"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// func JWTMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		logger.InfoLog.Println("Getting jwt token from cookies")
// 		token, err := c.Cookie("jwt")
// 		if err != nil {
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
// 			logger.InfoLog.Println("Error when getting jwt from cookies.\tJWT token == nil:", (token == ""), "\nError:", err)
// 			c.Abort()
// 			return
// 		}

// 		logger.InfoLog.Println("Validating jwt token")
// 		_, err = auth.ValidateJWT(token)
// 		if err != nil {
// 			logger.InfoLog.Println("Error when validating jwt token.\tError:", err)
// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 			c.Abort()
// 			return
// 		}
// 		// c.Set("userID", claims["id"])
// 		// c.Set("userRole", claims["role"])
// 		c.Next()
// 	}
// }

func RoleMiddleware(db *gorm.DB, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt")
		if err != nil {
			tokenString, err = auth.UpdateJWTToken(c, db)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
				c.Abort()
				return
			}
		}

		token, _ := jwt.ParseWithClaims(tokenString, &auth.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("jwtSecret")), nil
		})

		if err != nil || !token.Valid {
			tokenString, err = auth.UpdateJWTToken(c, db)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}
			token, _ = jwt.ParseWithClaims(tokenString, &auth.MyClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("jwtSecret")), nil
			})
		}

		claims, ok := token.Claims.(*auth.MyClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Claims"})
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if claims.Role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have access to this resource"})
		c.Abort()
	}
}
