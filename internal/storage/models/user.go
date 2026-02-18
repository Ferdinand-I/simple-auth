package models

type User struct {
	ID           int64  `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
}

type UserCreateDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginDTO struct {
	UserCreateDTO
}

type UserResponseDTO struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
