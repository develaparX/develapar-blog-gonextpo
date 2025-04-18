package repository

import (
	"database/sql"
	"develapar-server/model"
	"time"
)

type UserRepository interface {
	CreateNewUser(payload model.User) (model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) CreateNewUser(payload model.User) (model.User, error) {
	var user model.User
	err := u.db.QueryRow(`INSERT INTO users (name, email, password, role, updated_at) VALUES($1, $2, $3, $4, $5) RETURNING id, name, email, role, created_at, updated_at`, payload.Name, payload.Email, payload.Password, payload.Role, time.Now()).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func NewUserRepository(database *sql.DB) UserRepository {
	return &userRepository{db: database}
}
