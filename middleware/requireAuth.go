package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"stokq-backend/config"
	"stokq-backend/dto"
	"stokq-backend/models"
)

func RequireAuth(c *gin.Context) {
	// Get the authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Authorization header is required",
		})
		c.Abort()
		return
	}

	// Check if the header starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Authorization header must start with Bearer",
		})
		c.Abort()
		return
	}

	// Extract the token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse and validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Invalid token",
		})
		c.Abort()
		return
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Invalid token claims",
		})
		c.Abort()
		return
	}

	// Check token expiration
	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Token has expired",
			})
			c.Abort()
			return
		}
	}

	// Get user ID from claims
	userID, ok := claims["sub"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "Invalid user ID in token",
		})
		c.Abort()
		return
	}

	// Find the user
	var user models.User
	if err := config.DB.First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Error: "User not found",
		})
		c.Abort()
		return
	}

	// Attach user to context
	c.Set("user", user)
	c.Next()
}
