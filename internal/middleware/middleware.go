package middleware

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"icu/internal/domain/users/models"
	"icu/internal/domain/users/service"
	"net/http"
)

type JWTMiddleware struct {
	userService  *service.UserService
	jwtSecretKey string
}

func NewJWTMiddleware(userService *service.UserService, jwtSecretKey string) *JWTMiddleware {
	return &JWTMiddleware{
		userService:  userService,
		jwtSecretKey: jwtSecretKey,
	}
}

func (m *JWTMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}

		user, err := m.verifyToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Добавляем пользователя в контекст Gin
		c.Set("user", user)

		// Продолжаем выполнение следующего обработчика
		c.Next()
	}
}
func (m *JWTMiddleware) verifyToken(tokenString string) (*models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.jwtSecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	email := claims["sub"].(string)
	user, err := m.userService.GetUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	return user, nil
}
