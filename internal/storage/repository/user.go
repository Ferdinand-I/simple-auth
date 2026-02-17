package repository

import (
	"authapp/internal/storage/models"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}


func (r *UserRepo) Create(u *models.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", u.Username, u.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepo) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	err := r.db.Get(user, "SELECT id, username, password_hash FROM users WHERE username = $1", username)
	if err != nil {
		return nil, err
	}

	return user, nil
}
