package repository

import (
	"database/sql"
	"develapar-server/model"
	"errors"
	"time"
)

type UserRepository interface {
	CreateNewUser(payload model.User) (model.User, error)
	GetUserById(id int) (model.User, error)
	GetByEmail(email string) (model.User, error)
	GetAllUser() ([]model.User, error)
	SaveRefreshToken(userId int, token string, expiresAt time.Time) error
	ValidateRefreshToken(token string) (int, error)
	DeleteRefreshToken(token string) error
	DeleteAllRefreshTOkensByUser(userId int) error
	FindRefreshToken(token string) (model.RefreshToken, error)
	UpdateRefreshToken(oldToken string, newToken string, expiresAt time.Time) error
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) FindRefreshToken(token string) (model.RefreshToken, error) {
	var rt model.RefreshToken
	query := `SELECT id, user_id, token, expires_at, created_at, updated_at FROM refresh_tokens WHERE token = $1`
	row := r.db.QueryRow(query, token)

	err := row.Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt, &rt.UpdatedAt)
	if err != nil {
		return model.RefreshToken{}, err
	}

	return rt, nil
}


func (r *userRepository) UpdateRefreshToken(oldToken, newToken string, expiresAt time.Time) error {
	query := `UPDATE refresh_tokens SET token = $1, expires_at = $2, updated_at = NOW() WHERE token = $3`
	_, err := r.db.Exec(query, newToken, expiresAt, oldToken)
	return err
}



// DeleteAllRefreshTOkensByUser implements UserRepository.
func (u *userRepository) DeleteAllRefreshTOkensByUser(userId int) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := u.db.Exec(query, userId)
	return err
}

// DeleteRefreshToken implements UserRepository.
func (u *userRepository) DeleteRefreshToken(token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := u.db.Exec(query, token)
	return err
}

// SaveRefreshToken implements UserRepository.
func (u *userRepository) SaveRefreshToken(userId int, token string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := u.db.Exec(query, userId, token, expiresAt)
	return err
}

// ValidateRefreshToken implements UserRepository.
func (u *userRepository) ValidateRefreshToken(token string) (int, error) {
	var userId int
	var expiresAt time.Time
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1`
	err := u.db.QueryRow(query, token).Scan(&userId, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("refresh token not found")
		}
		return 0, err
	}

	if time.Now().After(expiresAt) {
		_ = u.DeleteRefreshToken(token)
		return 0, errors.New("refresh token expired")

	}

	return userId, nil
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
