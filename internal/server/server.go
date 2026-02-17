package server

import (
	"authapp/internal/handlers"
	"authapp/internal/service"
	"authapp/internal/storage/repository"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	Router *gin.Engine
}

func New() *Server {
	return &Server{
		Router: gin.Default(),
	}
}

func (s *Server) SetUp(db *sqlx.DB) {
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)
	
	apiGroup := s.Router.Group("/api/v1")
	
	usersGroup := apiGroup.Group("/users")
	usersGroup.POST("/", userHandler.CreateUser)
	usersGroup.GET("/secret", handlers.GetSecret)
	
	authGroup := apiGroup.Group("/auth")
	authGroup.POST("/login", userHandler.Authenticate)
}

func (s *Server) Run(port string) error {
	return s.Router.Run(port)
}
