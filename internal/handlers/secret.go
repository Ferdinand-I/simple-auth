package handlers

import (
	"authapp/internal/service"
	"net/http"

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
	user, exists := c.Get("User")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authenticated"})
	}

	c.JSON(200, gin.H{"secret": "message", "user": user})
}
