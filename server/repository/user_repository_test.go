package repository

import (
	"context"
	"database/sql"
	"develapar-server/model"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Context-aware test for CreateNewUser
func TestCreateNewUser_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		user          model.User
		setupMock     func(sqlmock.Sqlmock)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name: "successful user creation with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_123"),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "hashedpassword",
				Role:     "user",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "role", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john.doe@example.com", "user", time.Now(), time.Now())
				mock.ExpectQuery(`INSERT INTO users \(name, email, password, role, updated_at\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING id, name, email, role, created_at, updated_at`).
					WithArgs("John Doe", "john.doe@example.com", "hashedpassword", "user", sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
		},
		{
			name: "context cancellation during user creation",
			ctx:  context.WithValue(context.Background(), "request_id", "req_cancel"),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "hashedpassword",
				Role:     "user",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users \(name, email, password, role, updated_at\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING id, name, email, role, created_at, updated_at`).
					WithArgs("John Doe", "john.doe@example.com", "hashedpassword", "user", sqlmock.AnyArg()).
					WillReturnError(context.Canceled)
			},
			expectCancel: true,
		},
		{
			name: "context timeout during user creation",
			ctx:  context.WithValue(context.Background(), "request_id", "req_timeout"),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "hashedpassword",
				Role:     "user",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users \(name, email, password, role, updated_at\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING id, name, email, role, created_at, updated_at`).
					WithArgs("John Doe", "john.doe@example.com", "hashedpassword", "user", sqlmock.AnyArg()).
					WillReturnError(context.DeadlineExceeded)
			},
			expectTimeout: true,
		},
		{
			name: "database error during user creation",
			ctx:  context.WithValue(context.Background(), "request_id", "req_db_error"),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "hashedpassword",
				Role:     "user",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users \(name, email, password, role, updated_at\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING id, name, email, role, created_at, updated_at`).
					WithArgs("John Doe", "john.doe@example.com", "hashedpassword", "user", sqlmock.AnyArg()).
					WillReturnError(errors.New("database connection failed"))
			},
			expectedError: "database connection failed",
		},
		{
			name: "duplicate email constraint violation",
			ctx:  context.WithValue(context.Background(), "request_id", "req_duplicate"),
			user: model.User{
				Name:     "John Doe",
				Email:    "existing@example.com",
				Password: "hashedpassword",
				Role:     "user",
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`INSERT INTO users \(name, email, password, role, updated_at\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING id, name, email, role, created_at, updated_at`).
					WithArgs("John Doe", "existing@example.com", "hashedpassword", "user", sqlmock.AnyArg()).
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
			},
			expectedError: "duplicate key value violates unique constraint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database mock
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Setup mock expectations
			tt.setupMock(mock)

			// Create repository
			repo := NewUserRepository(db)

			// Execute test
			result, err := repo.CreateNewUser(tt.ctx, tt.user)

			// Verify results
			if tt.expectCancel {
				assert.Error(t, err)
				assert.Equal(t, context.Canceled, err)
			} else if tt.expectTimeout {
				assert.Error(t, err)
				assert.Equal(t, context.DeadlineExceeded, err)
			} else if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "John Doe", result.Name)
				assert.Equal(t, "john.doe@example.com", result.Email)
				assert.Equal(t, "user", result.Role)
				assert.NotZero(t, result.Id)
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Context-aware test for GetUserById
func TestGetUserById_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		userID        int
		setupMock     func(sqlmock.Sqlmock)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name:   "successful user retrieval with context",
			ctx:    context.WithValue(context.Background(), "request_id", "req_123"),
			userID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john.doe@example.com", "hashedpassword", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE id=\$1`).
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
		{
			name:   "context cancellation during user retrieval",
			ctx:    context.WithValue(context.Background(), "request_id", "req_cancel"),
			userID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE id=\$1`).
					WithArgs(1).
					WillReturnError(context.Canceled)
			},
			expectCancel: true,
		},
		{
			name:   "context timeout during user retrieval",
			ctx:    context.WithValue(context.Background(), "request_id", "req_timeout"),
			userID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE id=\$1`).
					WithArgs(1).
					WillReturnError(context.DeadlineExceeded)
			},
			expectTimeout: true,
		},
		{
			name:   "user not found",
			ctx:    context.WithValue(context.Background(), "request_id", "req_not_found"),
			userID: 999,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE id=\$1`).
					WithArgs(999).
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: "no rows in result set",
		},
		{
			name:   "database error during user retrieval",
			ctx:    context.WithValue(context.Background(), "request_id", "req_db_error"),
			userID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE id=\$1`).
					WithArgs(1).
					WillReturnError(errors.New("database connection failed"))
			},
			expectedError: "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database mock
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Setup mock expectations
			tt.setupMock(mock)

			// Create repository
			repo := NewUserRepository(db)

			// Execute test
			result, err := repo.GetUserById(tt.ctx, tt.userID)

			// Verify results
			if tt.expectCancel {
				assert.Error(t, err)
				assert.Equal(t, context.Canceled, err)
			} else if tt.expectTimeout {
				assert.Error(t, err)
				assert.Equal(t, context.DeadlineExceeded, err)
			} else if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, result.Id)
				assert.Equal(t, "John Doe", result.Name)
				assert.Equal(t, "john.doe@example.com", result.Email)
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Context-aware test for GetByEmail
func TestGetByEmail_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		email         string
		setupMock     func(sqlmock.Sqlmock)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name:  "successful user retrieval by email with context",
			ctx:   context.WithValue(context.Background(), "request_id", "req_123"),
			email: "john.doe@example.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john.doe@example.com", "hashedpassword", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email=\$1`).
					WithArgs("john.doe@example.com").
					WillReturnRows(rows)
			},
		},
		{
			name:  "context cancellation during email lookup",
			ctx:   context.WithValue(context.Background(), "request_id", "req_cancel"),
			email: "john.doe@example.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email=\$1`).
					WithArgs("john.doe@example.com").
					WillReturnError(context.Canceled)
			},
			expectCancel: true,
		},
		{
			name: "context timeout during email lookup",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			email: "john.doe@example.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email=\$1`).
					WithArgs("john.doe@example.com").
					WillReturnError(context.DeadlineExceeded)
			},
			expectTimeout: true,
		},
		{
			name:  "user not found by email",
			ctx:   context.WithValue(context.Background(), "request_id", "req_not_found"),
			email: "nonexistent@example.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email=\$1`).
					WithArgs("nonexistent@example.com").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: "no rows in result set",
		},
		{
			name:  "database error during email lookup",
			ctx:   context.WithValue(context.Background(), "request_id", "req_db_error"),
			email: "john.doe@example.com",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, created_at, updated_at FROM users WHERE email=\$1`).
					WithArgs("john.doe@example.com").
					WillReturnError(errors.New("database connection failed"))
			},
			expectedError: "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database mock
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Setup mock expectations
			tt.setupMock(mock)

			// Create repository
			repo := NewUserRepository(db)

			// Execute test
			result, err := repo.GetByEmail(tt.ctx, tt.email)

			// Verify results
			if tt.expectCancel {
				assert.Error(t, err)
				assert.Equal(t, context.Canceled, err)
			} else if tt.expectTimeout {
				assert.Error(t, err)
				assert.Equal(t, context.DeadlineExceeded, err)
			} else if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, result.Id)
				assert.Equal(t, "John Doe", result.Name)
				assert.Equal(t, "john.doe@example.com", result.Email)
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Context-aware test for GetAllUser
func TestGetAllUser_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		setupMock     func(sqlmock.Sqlmock)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name: "successful retrieval of all users with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_123"),
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john.doe@example.com", "hashedpassword1", "user", time.Now(), time.Now()).
					AddRow(2, "Jane Smith", "jane.smith@example.com", "hashedpassword2", "admin", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT id, name, email, password, role, created_at, updated_at FROM users`).
					WillReturnRows(rows)
			},
		},
		{
			name: "context cancellation during user retrieval",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, role, created_at, updated_at FROM users`).
					WillReturnError(context.Canceled)
			},
			expectCancel: true,
		},
		{
			name: "context timeout during user retrieval",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, role, created_at, updated_at FROM users`).
					WillReturnError(context.DeadlineExceeded)
			},
			expectTimeout: true,
		},
		{
			name: "database error during user retrieval",
			ctx:  context.WithValue(context.Background(), "request_id", "req_db_error"),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, email, password, role, created_at, updated_at FROM users`).
					WillReturnError(errors.New("database connection failed"))
			},
			expectedError: "database connection failed",
		},
		{
			name: "empty result set",
			ctx:  context.WithValue(context.Background(), "request_id", "req_empty"),
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT id, name, email, password, role, created_at, updated_at FROM users`).
					WillReturnRows(rows)
			},
		},
		{
			name: "context cancellation during row iteration",
			ctx:  context.WithValue(context.Background(), "request_id", "req_cancel_iteration"),
			setupMock: func(mock sqlmock.Sqlmock) {
				// This test simulates context cancellation during row processing
				// We'll use a custom mock that returns rows but then cancels context
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john.doe@example.com", "hashedpassword1", "user", time.Now(), time.Now()).
					AddRow(2, "Jane Smith", "jane.smith@example.com", "hashedpassword2", "admin", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT id, name, email, password, role, created_at, updated_at FROM users`).
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database mock
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Setup mock expectations
			tt.setupMock(mock)

			// Create repository
			repo := NewUserRepository(db)

			// For the context cancellation during iteration test, we need to cancel the context
			// after the query starts but before iteration completes
			if tt.name == "context cancellation during row iteration" {
				ctx, cancel := context.WithCancel(tt.ctx)
				// Cancel the context immediately to simulate cancellation during iteration
				cancel()
				tt.ctx = ctx
				tt.expectCancel = true
			}

			// Execute test
			result, err := repo.GetAllUser(tt.ctx)

			// Verify results
			if tt.expectCancel {
				assert.Error(t, err)
				assert.Equal(t, context.Canceled, err)
			} else if tt.expectTimeout {
				assert.Error(t, err)
				assert.Equal(t, context.DeadlineExceeded, err)
			} else if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				if tt.name == "empty result set" {
					assert.Empty(t, result)
				} else {
					assert.Len(t, result, 2)
					assert.Equal(t, "John Doe", result[0].Name)
					assert.Equal(t, "Jane Smith", result[1].Name)
				}
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Context-aware test for GetAllUserWithPagination
func TestGetAllUserWithPagination_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		offset        int
		limit         int
		setupMock     func(sqlmock.Sqlmock)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name:   "successful paginated user retrieval with context",
			ctx:    context.WithValue(context.Background(), "request_id", "req_123"),
			offset: 0,
			limit:  10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock count query
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(25)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
					WillReturnRows(countRows)

				// Mock paginated query
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at", "updated_at"}).
					AddRow(1, "John Doe", "john.doe@example.com", "hashedpassword1", "user", time.Now(), time.Now()).
					AddRow(2, "Jane Smith", "jane.smith@example.com", "hashedpassword2", "admin", time.Now(), time.Now())
				mock.ExpectQuery(`SELECT id, name, email, password, role, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
					WithArgs(10, 0).
					WillReturnRows(rows)
			},
		},
		{
			name:   "context cancellation during count query",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			offset: 0,
			limit:  10,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
					WillReturnError(context.Canceled)
			},
			expectCancel: true,
		},
		{
			name:   "context timeout during count query",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			offset: 0,
			limit:  10,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
					WillReturnError(context.DeadlineExceeded)
			},
			expectTimeout: true,
		},
		{
			name:   "context cancellation during paginated query",
			ctx:    context.WithValue(context.Background(), "request_id", "req_cancel_paginated"),
			offset: 0,
			limit:  10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock successful count query
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(25)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
					WillReturnRows(countRows)

				// Mock cancelled paginated query
				mock.ExpectQuery(`SELECT id, name, email, password, role, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
					WithArgs(10, 0).
					WillReturnError(context.Canceled)
			},
			expectCancel: true,
		},
		{
			name:   "database error during count query",
			ctx:    context.WithValue(context.Background(), "request_id", "req_count_error"),
			offset: 0,
			limit:  10,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
					WillReturnError(errors.New("database connection failed"))
			},
			expectedError: "database connection failed",
		},
		{
			name:   "database error during paginated query",
			ctx:    context.WithValue(context.Background(), "request_id", "req_paginated_error"),
			offset: 0,
			limit:  10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock successful count query
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(25)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
					WillReturnRows(countRows)

				// Mock failed paginated query
				mock.ExpectQuery(`SELECT id, name, email, password, role, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
					WithArgs(10, 0).
					WillReturnError(errors.New("database connection failed"))
			},
			expectedError: "database connection failed",
		},
		{
			name:   "empty paginated result",
			ctx:    context.WithValue(context.Background(), "request_id", "req_empty_page"),
			offset: 100,
			limit:  10,
			setupMock: func(mock sqlmock.Sqlmock) {
				// Mock count query
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(25)
				mock.ExpectQuery(`SELECT COUNT\(\*\) FROM users`).
					WillReturnRows(countRows)

				// Mock empty paginated query
				rows := sqlmock.NewRows([]string{"id", "name", "email", "password", "role", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT id, name, email, password, role, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT \$1 OFFSET \$2`).
					WithArgs(10, 100).
					WillReturnRows(rows)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database mock
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Setup mock expectations
			tt.setupMock(mock)

			// Create repository
			repo := NewUserRepository(db)

			// Execute test
			result, total, err := repo.GetAllUserWithPagination(tt.ctx, tt.offset, tt.limit)

			// Verify results
			if tt.expectCancel {
				assert.Error(t, err)
				assert.Equal(t, context.Canceled, err)
				assert.Equal(t, 0, total)
				assert.Nil(t, result)
			} else if tt.expectTimeout {
				assert.Error(t, err)
				assert.Equal(t, context.DeadlineExceeded, err)
				assert.Equal(t, 0, total)
				assert.Nil(t, result)
			} else if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Equal(t, 0, total)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 25, total)
				if tt.name == "empty paginated result" {
					assert.Empty(t, result)
				} else {
					assert.Len(t, result, 2)
					assert.Equal(t, "John Doe", result[0].Name)
					assert.Equal(t, "Jane Smith", result[1].Name)
				}
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Context-aware test for SaveRefreshToken
func TestSaveRefreshToken_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		userID        int
		token         string
		expiresAt     time.Time
		setupMock     func(sqlmock.Sqlmock)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name:      "successful refresh token save with context",
			ctx:       context.WithValue(context.Background(), "request_id", "req_123"),
			userID:    1,
			token:     "refresh_token_123",
			expiresAt: time.Now().Add(7 * 24 * time.Hour),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO refresh_tokens \(user_id, token, expires_at\) VALUES \(\$1, \$2, \$3\)`).
					WithArgs(1, "refresh_token_123", sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
		},
		{
			name: "context cancellation during token save",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			userID:    1,
			token:     "refresh_token_123",
			expiresAt: time.Now().Add(7 * 24 * time.Hour),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO refresh_tokens \(user_id, token, expires_at\) VALUES \(\$1, \$2, \$3\)`).
					WithArgs(1, "refresh_token_123", sqlmock.AnyArg()).
					WillReturnError(context.Canceled)
			},
			expectCancel: true,
		},
		{
			name: "context timeout during token save",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			userID:    1,
			token:     "refresh_token_123",
			expiresAt: time.Now().Add(7 * 24 * time.Hour),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO refresh_tokens \(user_id, token, expires_at\) VALUES \(\$1, \$2, \$3\)`).
					WithArgs(1, "refresh_token_123", sqlmock.AnyArg()).
					WillReturnError(context.DeadlineExceeded)
			},
			expectTimeout: true,
		},
		{
			name:      "database error during token save",
			ctx:       context.WithValue(context.Background(), "request_id", "req_db_error"),
			userID:    1,
			token:     "refresh_token_123",
			expiresAt: time.Now().Add(7 * 24 * time.Hour),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO refresh_tokens \(user_id, token, expires_at\) VALUES \(\$1, \$2, \$3\)`).
					WithArgs(1, "refresh_token_123", sqlmock.AnyArg()).
					WillReturnError(errors.New("database connection failed"))
			},
			expectedError: "database connection failed",
		},
		{
			name:      "foreign key constraint violation",
			ctx:       context.WithValue(context.Background(), "request_id", "req_fk_error"),
			userID:    999, // Non-existent user
			token:     "refresh_token_123",
			expiresAt: time.Now().Add(7 * 24 * time.Hour),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO refresh_tokens \(user_id, token, expires_at\) VALUES \(\$1, \$2, \$3\)`).
					WithArgs(999, "refresh_token_123", sqlmock.AnyArg()).
					WillReturnError(errors.New("foreign key constraint violation"))
			},
			expectedError: "foreign key constraint violation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database mock
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Setup mock expectations
			tt.setupMock(mock)

			// Create repository
			repo := NewUserRepository(db)

			// Execute test
			err = repo.SaveRefreshToken(tt.ctx, tt.userID, tt.token, tt.expiresAt)

			// Verify results
			if tt.expectCancel {
				assert.Error(t, err)
				assert.Equal(t, context.Canceled, err)
			} else if tt.expectTimeout {
				assert.Error(t, err)
				assert.Equal(t, context.DeadlineExceeded, err)
			} else if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

// Context-aware test for FindRefreshToken
func TestFindRefreshToken_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		token         string
		setupMock     func(sqlmock.Sqlmock)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name:  "successful refresh token retrieval with context",
			ctx:   context.WithValue(context.Background(), "request_id", "req_123"),
			token: "refresh_token_123",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "user_id", "token", "expires_at", "created_at"}).
					AddRow(1, 1, "refresh_token_123", time.Now().Add(7*24*time.Hour), time.Now())
				mock.ExpectQuery(`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = \$1`).
					WithArgs("refresh_token_123").
					WillReturnRows(rows)
			},
		},
		{
			name: "context cancellation during token retrieval",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			token: "refresh_token_123",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = \$1`).
					WithArgs("refresh_token_123").
					WillReturnError(context.Canceled)
			},
			expectCancel: true,
		},
		{
			name: "context timeout during token retrieval",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			token: "refresh_token_123",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = \$1`).
					WithArgs("refresh_token_123").
					WillReturnError(context.DeadlineExceeded)
			},
			expectTimeout: true,
		},
		{
			name:  "refresh token not found",
			ctx:   context.WithValue(context.Background(), "request_id", "req_not_found"),
			token: "nonexistent_token",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = \$1`).
					WithArgs("nonexistent_token").
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: "no rows in result set",
		},
		{
			name:  "database error during token retrieval",
			ctx:   context.WithValue(context.Background(), "request_id", "req_db_error"),
			token: "refresh_token_123",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = \$1`).
					WithArgs("refresh_token_123").
					WillReturnError(errors.New("database connection failed"))
			},
			expectedError: "database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup database mock
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Setup mock expectations
			tt.setupMock(mock)

			// Create repository
			repo := NewUserRepository(db)

			// Execute test
			result, err := repo.FindRefreshToken(tt.ctx, tt.token)

			// Verify results
			if tt.expectCancel {
				assert.Error(t, err)
				assert.Equal(t, context.Canceled, err)
			} else if tt.expectTimeout {
				assert.Error(t, err)
				assert.Equal(t, context.DeadlineExceeded, err)
			} else if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 1, result.ID)
				assert.Equal(t, 1, result.UserID)
				assert.Equal(t, "refresh_token_123", result.Token)
				assert.True(t, result.ExpiresAt.After(time.Now()))
			}

			// Verify all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}