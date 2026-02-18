package service

import (
	"authapp/internal/config"
	"authapp/internal/storage/models"
	"authapp/internal/storage/repository"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTResponseDTO struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type UserService struct {
	UserRepo *repository.UserRepo
	Cfg      *config.Config
}

func NewUserService(userRepo *repository.UserRepo, cfg *config.Config) *UserService {
	return &UserService{UserRepo: userRepo, Cfg: cfg}
}

func (s *UserService) CreateUser(username, password string) (*models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(passwordHash),
	}

	err = s.UserRepo.Create(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Authenticate(username, password string) (*JWTResponseDTO, error) {
	user, err := s.UserRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("User not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("Wrong password")
	}

	accessToken, refreshToken, err := s.generateTokens(user.ID)
	if err != nil {
		return nil, errors.New("Token generation failed")
	}

	jwtResp := &JWTResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return jwtResp, nil
}

func (s *UserService) generateTokens(userID int64) (string, string, error) {
	createToken := func(exp time.Duration) (string, error) {
		claims := &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
			Subject:   strconv.FormatInt(userID, 10),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		return token.SignedString(s.Cfg.Auth.SecretKey)
	}

	accessToken, err := createToken(s.Cfg.Auth.AccessTokenExpiresMinutes)
	if err != nil {
		log.Printf("Failed to generate access token: %v", err)
		return "", "", errors.New("Token generation failed")
	}

	refreshToken, err := createToken(s.Cfg.Auth.RefreshTokenExpiresDays)
	if err != nil {
		log.Printf("Failed to generate refresh token: %v", err)
		return "", "", errors.New("Token generation failed")
	}

	return accessToken, refreshToken, nil
}
