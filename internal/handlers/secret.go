package handlers

import (
	"authapp/internal/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type SecretHandler struct {
	UserService service.UserService
}

func NewSecretHandler(userService service.UserService) *SecretHandler {
	return &SecretHandler{
		UserService: userService,
	}
}

func GetSecret(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
    if tokenString == "" {
      c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
      c.Abort()
      return
    }

    if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
      tokenString = tokenString[7:]
    }

    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
      return jwtSecret, nil
    })

    if err != nil || !token.Valid {
      c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
      c.Abort()
      return
    }

    if claims, ok := token.Claims.(*Claims); ok {
      c.Set("username", claims.Username)
      c.Next()
    } else {
      c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
      c.Abort()
}
