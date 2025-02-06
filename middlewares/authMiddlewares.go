package middlewares

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				log.Println("Cookie not found:", err)
				c.AbortWithStatusJSON(401, gin.H{"error":"Unauthorized or missing token"})
				return
			}
			log.Println("error in getting token:", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Authentication error"})
			return
		}

		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			log.Println("JWT_SECRET is not set in .env")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
			return
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			log.Println("Invalid or expired token:", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		c.Set("userID", claims.Subject)

		c.Next()
	}
}