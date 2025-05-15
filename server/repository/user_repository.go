package repository

import (
	"database/sql"
	"develapar-server/model"
	"time"
)

type UserRepository interface {
	CreateNewUser(payload model.User) (model.User, error)
	GetUserById(id int) (model.User, error)
	GetByEmail(email string) (model.User, error)
	GetAllUser() ([]model.User, error)
}

type userRepository struct {
	db *sql.DB
}

// GetByEmail implements UserRepository.
func (u *userRepository) GetByEmail(email string) (model.User, error) {
	var user model.User

	err := u.db.QueryRow(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email=$1`, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

// GetAllUser implements UserRepository.
func (u *userRepository) GetAllUser() ([]model.User, error) {
	var listUser []model.User

	rows, err := u.db.Query(`SELECT id, name, email, password, role, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User

		err := rows.Scan(
			&user.Id,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		listUser = append(listUser, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return listUser, nil
}

// GetUserById implements UserRepository.
func (u *userRepository) GetUserById(id int) (model.User, error) {
	var user model.User

	err := u.db.QueryRow(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE id=$1`, id).Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return model.User{}, nil
	}

	return user, nil
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
