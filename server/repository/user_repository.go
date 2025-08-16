package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
	"errors"
	"time"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateNewUser(ctx context.Context, payload model.User) (model.User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetAllUser(ctx context.Context) ([]model.User, error)
	GetAllUserWithPagination(ctx context.Context, offset, limit int) ([]model.User, int, error)
	SaveRefreshToken(ctx context.Context, userId uuid.UUID, token string, expiresAt time.Time) error
	ValidateRefreshToken(ctx context.Context, token string) (uuid.UUID, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllRefreshTOkensByUser(ctx context.Context, userId uuid.UUID) error
	FindRefreshToken(ctx context.Context, token string) (model.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, oldToken string, newToken string, expiresAt time.Time) error
	UpdateUser(ctx context.Context, payload model.User) (model.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	db *sql.DB
}

func (r *userRepository) FindRefreshToken(ctx context.Context, token string) (model.RefreshToken, error) {
	var rt model.RefreshToken
	query := `SELECT id, user_id, token, expires_at, created_at, updated_at FROM refresh_tokens WHERE token = $1`
	row := r.db.QueryRowContext(ctx, query, token)

	err := row.Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt, &rt.UpdatedAt)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return model.RefreshToken{}, ctx.Err()
		}
		return model.RefreshToken{}, err
	}

	return rt, nil
}

func (r *userRepository) UpdateRefreshToken(ctx context.Context, oldToken, newToken string, expiresAt time.Time) error {
	query := `UPDATE refresh_tokens SET token = $1, expires_at = $2 WHERE token = $3`
	_, err := r.db.ExecContext(ctx, query, newToken, expiresAt, oldToken)
	return err
}

// DeleteAllRefreshTOkensByUser implements UserRepository.
func (u *userRepository) DeleteAllRefreshTOkensByUser(ctx context.Context, userId uuid.UUID) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := u.db.ExecContext(ctx, query, userId)
	return err
}

// DeleteRefreshToken implements UserRepository.
func (u *userRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := u.db.ExecContext(ctx, query, token)
	return err
}

// SaveRefreshToken implements UserRepository.
func (u *userRepository) SaveRefreshToken(ctx context.Context, userId uuid.UUID, token string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := u.db.ExecContext(ctx, query, userId, token, expiresAt)
	return err
}

// ValidateRefreshToken implements UserRepository.
func (u *userRepository) ValidateRefreshToken(ctx context.Context, token string) (uuid.UUID, error) {
	var userId uuid.UUID
	var expiresAt time.Time
	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1`
	err := u.db.QueryRowContext(ctx, query, token).Scan(&userId, &expiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.Nil, errors.New("refresh token not found")
		}
		return uuid.Nil, err
	}

	if time.Now().After(expiresAt) {
		_ = u.DeleteRefreshToken(ctx, token)
		return uuid.Nil, errors.New("refresh token expired")

	}

	return userId, nil
}

// GetByEmail implements UserRepository.
func (u *userRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User

	err := u.db.QueryRowContext(ctx, `SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email=$1`, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return model.User{}, ctx.Err()
		}
		return model.User{}, err
	}

	return user, nil
}

// GetAllUser implements UserRepository.
func (u *userRepository) GetAllUser(ctx context.Context) ([]model.User, error) {
	var listUser []model.User

	rows, err := u.db.QueryContext(ctx, `SELECT id, name, email, password, role, created_at, updated_at FROM users`)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// Check for context cancellation during iteration
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

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

// GetAllUserWithPagination implements UserRepository with pagination support
func (u *userRepository) GetAllUserWithPagination(ctx context.Context, offset, limit int) ([]model.User, int, error) {
	// First get the total count
	var totalCount int
	countQuery := `SELECT COUNT(*) FROM users`
	err := u.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}

	// Then get the paginated results
	var listUser []model.User
	query := `SELECT id, name, email, password, role, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := u.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		// Check for context cancellation during iteration
		select {
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		default:
		}

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
			return nil, 0, err
		}

		listUser = append(listUser, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return listUser, totalCount, nil
}

// GetUserById implements UserRepository.
func (u *userRepository) GetUserById(ctx context.Context, id uuid.UUID) (model.User, error) {
	var user model.User

	err := u.db.QueryRowContext(ctx, `SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id=$1`, id).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return model.User{}, ctx.Err()
		}
		return model.User{}, err
	}

	return user, nil
}

func (u *userRepository) CreateNewUser(ctx context.Context, payload model.User) (model.User, error) {
	newId := uuid.Must(uuid.NewV7())

	var user model.User
	err := u.db.QueryRowContext(ctx, `INSERT INTO users (id, name, email, password, role,created_at, updated_at) VALUES($1, $2, $3, $4, $5,$6, $7) RETURNING id, name, email, role, created_at, updated_at`, newId, payload.Name, payload.Email, payload.Password, payload.Role, time.Now(), time.Now()).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return model.User{}, ctx.Err()
		}
		return model.User{}, err
	}
	return user, nil
}

// UpdateUser implements UserRepository.
func (u *userRepository) UpdateUser(ctx context.Context, payload model.User) (model.User, error) {
	var user model.User
	err := u.db.QueryRowContext(ctx, `UPDATE users SET name = $1, email = $2, password = $3, updated_at = $4 WHERE id = $5 RETURNING id, name, email, role, created_at, updated_at`, payload.Name, payload.Email, payload.Password, time.Now(), payload.Id).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		// Check if context was cancelled or timed out
		if ctx.Err() != nil {
			return model.User{}, ctx.Err()
		}
		return model.User{}, err
	}
	return user, nil
}

// DeleteUser implements UserRepository.
func (u *userRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := u.db.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	return err
}

func NewUserRepository(database *sql.DB) UserRepository {
	return &userRepository{db: database}
}
