package handlers

import (
	"authapp/internal/service"
	"authapp/internal/storage/models"
	"github.com/gin-gonic/gin"
	"log"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	dto := &models.UserCreateDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := h.UserService.CreateUser(dto.Username, dto.Password)
	if err != nil {
		log.Print(err)
		c.JSON(500, gin.H{"error": "User creation failed"})
		return
	}

	c.JSON(201, gin.H{"id": user.ID, "username": user.Username})
}

func (h *UserHandler) Authenticate(c *gin.Context) {
	dto := &models.UserLoginDTO{}
	if err := c.ShouldBindJSON(dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jwtResp, err := h.UserService.Authenticate(dto.Username, dto.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, jwtResp)
}
