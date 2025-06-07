package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"OnlineLibraryPortal/database"
	"OnlineLibraryPortal/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var jwtSecret = []byte("Teja")

type CustomClaims struct {
	jwt.StandardClaims
	Jti  string `json:"jti"`
	Role int    `json:"role"`
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Middleware start")

		authHeader := c.GetHeader("Authorization")
		fmt.Println("Authorization Header:", authHeader)

		if authHeader == "" {
			fmt.Println("Authorization header missing")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		fmt.Println("Auth header parts:", parts)

		if len(parts) != 2 || parts[0] != "Bearer" {
			fmt.Println("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		fmt.Println("Token string extracted:", tokenString)

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			fmt.Println("Inside key func")
			return jwtSecret, nil
		})

		if err != nil {
			fmt.Println("Error parsing token:", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			fmt.Println("Token invalid")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*CustomClaims)
		if !ok {
			fmt.Println("Failed to cast claims")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		fmt.Println("Claims extracted:", claims)
		fmt.Println("Subject (userID):", claims.Subject)
		fmt.Println("JTI:", claims.Jti)
		fmt.Println("Role:", claims.Role)

		var user models.User
		err = database.DB.Where("id = ? AND jti = ?", claims.Subject, claims.Jti).First(&user).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				fmt.Println("Token revoked or invalid: user not found with matching jti")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token revoked or invalid"})
				c.Abort()
				return
			}
			fmt.Println("Database error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}

		fmt.Println("User found:", user.ID)

		var fuser models.User
		if err := database.DB.First(&fuser, claims.Subject).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Set("userID", fuser.ID)
		c.Set("userRole", fuser.Role)

		fmt.Println(fuser.ID)
		fmt.Println(fuser.Role)
		fmt.Println("Middleware success, passing to next handler")

		c.Next()
	}
}
