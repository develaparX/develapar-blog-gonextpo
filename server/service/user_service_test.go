package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/utils"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/golang-jwt/jwt/v5"
)

// Mock UserRepository with context support
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateNewUser(ctx context.Context, user model.User) (model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(ctx context.Context, id int) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) GetAllUser(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) GetAllUserWithPagination(ctx context.Context, offset, limit int) ([]model.User, int, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]model.User), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (model.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockUserRepository) FindRefreshToken(ctx context.Context, token string) (model.RefreshToken, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(model.RefreshToken), args.Error(1)
}

func (m *MockUserRepository) UpdateRefreshToken(ctx context.Context, oldToken, newToken string, expiresAt time.Time) error {
	args := m.Called(ctx, oldToken, newToken, expiresAt)
	return args.Error(0)
}

func (m *MockUserRepository) ValidateRefreshToken(ctx context.Context, token string) (int, error) {
	args := m.Called(ctx, token)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteAllRefreshTOkensByUser(ctx context.Context, userId int) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, payload model.User) (model.User, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Mock JwtService
type MockJwtService struct {
	mock.Mock
}

func (m *MockJwtService) GenerateToken(user model.User) (dto.LoginResponseDto, error) {
	args := m.Called(user)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

func (m *MockJwtService) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockJwtService) ValidateToken(token string) (*dto.JwtTokenClaims, error) {
	args := m.Called(token)
	return args.Get(0).(*dto.JwtTokenClaims), args.Error(1)
}

func (m *MockJwtService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}

// Mock PasswordHasher
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) EncryptPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) ComparePasswordHash(passwordHash string, plainPassword string) error {
	args := m.Called(passwordHash, plainPassword)
	return args.Error(0)
}

// Mock PaginationService
type MockPaginationService struct {
	mock.Mock
}

func (m *MockPaginationService) Paginate(ctx context.Context, data interface{}, total int, query PaginationQuery) (PaginationResult, *utils.AppError) {
	args := m.Called(ctx, data, total, query)
	if args.Get(1) == nil {
		return args.Get(0).(PaginationResult), nil
	}
	return args.Get(0).(PaginationResult), args.Get(1).(*utils.AppError)
}

func (m *MockPaginationService) ValidatePagination(ctx context.Context, page, limit int) *utils.AppError {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockPaginationService) BuildMetadata(ctx context.Context, total, page, limit int) (PaginationMetadata, *utils.AppError) {
	args := m.Called(ctx, total, page, limit)
	if args.Get(1) == nil {
		return args.Get(0).(PaginationMetadata), nil
	}
	return args.Get(0).(PaginationMetadata), args.Get(1).(*utils.AppError)
}

func (m *MockPaginationService) ParseQuery(ctx context.Context, page, limit int, sortBy, sortDir string) (PaginationQuery, *utils.AppError) {
	args := m.Called(ctx, page, limit, sortBy, sortDir)
	if args.Get(1) == nil {
		return args.Get(0).(PaginationQuery), nil
	}
	return args.Get(0).(PaginationQuery), args.Get(1).(*utils.AppError)
}

func (m *MockPaginationService) CalculateOffset(page, limit int) int {
	args := m.Called(page, limit)
	return args.Int(0)
}

// Mock ValidationService
type MockValidationService struct {
	mock.Mock
}

func (m *MockValidationService) ValidateUser(ctx context.Context, user model.User) *utils.AppError {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockValidationService) ValidateArticle(ctx context.Context, article model.Article) *utils.AppError {
	args := m.Called(ctx, article)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockValidationService) ValidateComment(ctx context.Context, comment model.Comment) *utils.AppError {
	args := m.Called(ctx, comment)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockValidationService) ValidatePagination(ctx context.Context, page, limit int) *utils.AppError {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockValidationService) ValidateField(ctx context.Context, field string, value interface{}, rules string) *FieldError {
	args := m.Called(ctx, field, value, rules)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*FieldError)
}

func (m *MockValidationService) ValidateStruct(ctx context.Context, s interface{}) []FieldError {
	args := m.Called(ctx, s)
	return args.Get(0).([]FieldError)
}

// Context-aware test for CreateNewUser
func TestCreateNewUser_WithContext(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		user           model.User
		setupMocks     func(*MockUserRepository, *MockJwtService, *MockPasswordHasher, *MockPaginationService, *MockValidationService)
		expectedError  string
		expectCancel   bool
		expectTimeout  bool
	}{
		{
			name: "successful user creation with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_123"),
			user: model.User{
				Name:     "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockVS.On("ValidateUser", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("model.User")).Return(nil).Once()
				mockPH.On("EncryptPassword", "password123").Return("hashedpassword", nil).Once()
				mockRepo.On("CreateNewUser", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("model.User")).Return(
					model.User{Id: 1, Name: "testuser", Email: "test@example.com", Password: "hashedpassword", Role: "user"}, nil).Once()
			},
		},
		{
			name: "context cancellation during user creation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately
				return ctx
			}(),
			user: model.User{
				Name:     "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed as context is cancelled immediately
			},
			expectCancel: true,
		},
		{
			name: "context timeout during user creation",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond) // Ensure timeout
				return ctx
			}(),
			user: model.User{
				Name:     "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed as context times out immediately
			},
			expectTimeout: true,
		},
		{
			name: "validation error with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_456"),
			user: model.User{
				Name:     "",
				Email:    "invalid-email",
				Password: "123",
				Role:     "invalid",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				validationErr := &utils.AppError{
					Code:      utils.ErrValidation,
					Message:   "User validation failed",
					RequestID: "req_456",
				}
				mockVS.On("ValidateUser", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("model.User")).Return(validationErr).Once()
			},
			expectedError: "User validation failed",
		},
		{
			name: "repository error with context propagation",
			ctx:  context.WithValue(context.Background(), "request_id", "req_789"),
			user: model.User{
				Name:     "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockVS.On("ValidateUser", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("model.User")).Return(nil).Once()
				mockPH.On("EncryptPassword", "password123").Return("hashedpassword", nil).Once()
				mockRepo.On("CreateNewUser", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("model.User")).Return(
					model.User{}, errors.New("database connection failed")).Once()
			},
			expectedError: "failed to create user",
		},
		{
			name: "context cancellation during repository operation",
			ctx:  context.WithValue(context.Background(), "request_id", "req_cancel"),
			user: model.User{
				Name:     "testuser",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockVS.On("ValidateUser", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("model.User")).Return(nil).Once()
				mockPH.On("EncryptPassword", "password123").Return("hashedpassword", nil).Once()
				mockRepo.On("CreateNewUser", mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("model.User")).Return(
					model.User{}, context.Canceled).Once()
			},
			expectCancel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockUserRepository)
			mockJwtService := new(MockJwtService)
			mockPasswordHasher := new(MockPasswordHasher)
			mockPaginationService := new(MockPaginationService)
			mockValidationService := new(MockValidationService)

			userService := NewUserservice(mockRepo, mockJwtService, mockPasswordHasher, mockPaginationService, mockValidationService)

			// Setup mock expectations
			tt.setupMocks(mockRepo, mockJwtService, mockPasswordHasher, mockPaginationService, mockValidationService)

			// Execute test
			result, err := userService.CreateNewUser(tt.ctx, tt.user)

			// Verify results
			if tt.expectCancel {
				assert.Error(t, err)
				// Check if the error is context.Canceled or contains it
				if err != context.Canceled {
					assert.Contains(t, err.Error(), "context canceled")
				}
			} else if tt.expectTimeout {
				assert.Error(t, err)
				// Check if the error is context.DeadlineExceeded or contains it
				if err != context.DeadlineExceeded {
					assert.Contains(t, err.Error(), "context deadline exceeded")
				}
			} else if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "testuser", result.Name)
				assert.Equal(t, "-", result.Password) // Password should be masked
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockJwtService.AssertExpectations(t)
			mockPasswordHasher.AssertExpectations(t)
			mockPaginationService.AssertExpectations(t)
			mockValidationService.AssertExpectations(t)
		})
	}
}

// Context-aware test for FindUserById
func TestFindUserById_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		userID        string
		setupMocks    func(*MockUserRepository, *MockJwtService, *MockPasswordHasher, *MockPaginationService, *MockValidationService)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name:   "successful user retrieval with context",
			ctx:    context.WithValue(context.Background(), "request_id", "req_123"),
			userID: "1",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				expectedUser := model.User{Id: 1, Name: "testuser", Email: "test@example.com", Password: "hashedpassword"}
				mockRepo.On("GetUserById", mock.AnythingOfType("*context.valueCtx"), 1).Return(expectedUser, nil).Once()
			},
		},
		{
			name: "context cancellation during user retrieval",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			userID: "1",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed as context is cancelled immediately
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
			userID: "1",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed as context times out immediately
			},
			expectTimeout: true,
		},
		{
			name:   "invalid user ID format",
			ctx:    context.WithValue(context.Background(), "request_id", "req_456"),
			userID: "abc",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed for validation error
			},
			expectedError: "invalid user ID format",
		},
		{
			name:   "user ID must be greater than 0",
			ctx:    context.WithValue(context.Background(), "request_id", "req_789"),
			userID: "0",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed for validation error
			},
			expectedError: "user ID must be greater than 0",
		},
		{
			name:   "user not found in repository with context",
			ctx:    context.WithValue(context.Background(), "request_id", "req_notfound"),
			userID: "999",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockRepo.On("GetUserById", mock.AnythingOfType("*context.valueCtx"), 999).Return(model.User{}, errors.New("user not found")).Once()
			},
			expectedError: "failed to fetch user",
		},
		{
			name:   "context cancellation during repository operation",
			ctx:    context.WithValue(context.Background(), "request_id", "req_cancel"),
			userID: "1",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockRepo.On("GetUserById", mock.AnythingOfType("*context.valueCtx"), 1).Return(model.User{}, context.Canceled).Once()
			},
			expectCancel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockUserRepository)
			mockJwtService := new(MockJwtService)
			mockPasswordHasher := new(MockPasswordHasher)
			mockPaginationService := new(MockPaginationService)
			mockValidationService := new(MockValidationService)

			userService := NewUserservice(mockRepo, mockJwtService, mockPasswordHasher, mockPaginationService, mockValidationService)

			// Setup mock expectations
			tt.setupMocks(mockRepo, mockJwtService, mockPasswordHasher, mockPaginationService, mockValidationService)

			// Execute test
			result, err := userService.FindUserById(tt.ctx, tt.userID)

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
				assert.Equal(t, "testuser", result.Name)
				assert.Equal(t, "-", result.Password) // Password should be masked
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockJwtService.AssertExpectations(t)
			mockPasswordHasher.AssertExpectations(t)
			mockPaginationService.AssertExpectations(t)
			mockValidationService.AssertExpectations(t)
		})
	}
}

// Context-aware test for FindAllUser
func TestFindAllUser_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		setupMocks    func(*MockUserRepository, *MockJwtService, *MockPasswordHasher, *MockPaginationService, *MockValidationService)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name: "successful retrieval of all users with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_123"),
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				expectedUsers := []model.User{
					{Id: 1, Name: "user1", Email: "user1@example.com", Password: "hashedpassword1"},
					{Id: 2, Name: "user2", Email: "user2@example.com", Password: "hashedpassword2"},
				}
				mockRepo.On("GetAllUser", mock.AnythingOfType("*context.valueCtx")).Return(expectedUsers, nil).Once()
			},
		},
		{
			name: "context cancellation during user retrieval",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed as context is cancelled immediately
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
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed as context times out immediately
			},
			expectTimeout: true,
		},
		{
			name: "repository error with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_error"),
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockRepo.On("GetAllUser", mock.AnythingOfType("*context.valueCtx")).Return([]model.User{}, errors.New("database connection failed")).Once()
			},
			expectedError: "failed to fetch users",
		},
		{
			name: "context cancellation during repository operation",
			ctx:  context.WithValue(context.Background(), "request_id", "req_cancel"),
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockRepo.On("GetAllUser", mock.AnythingOfType("*context.valueCtx")).Return([]model.User{}, context.Canceled).Once()
			},
			expectCancel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockUserRepository)
			mockJwtService := new(MockJwtService)
			mockPasswordHasher := new(MockPasswordHasher)
			mockPaginationService := new(MockPaginationService)
			mockValidationService := new(MockValidationService)

			userService := NewUserservice(mockRepo, mockJwtService, mockPasswordHasher, mockPaginationService, mockValidationService)

			// Setup mock expectations
			tt.setupMocks(mockRepo, mockJwtService, mockPasswordHasher, mockPaginationService, mockValidationService)

			// Execute test
			result, err := userService.FindAllUser(tt.ctx)

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
				assert.Len(t, result, 2)
				// Verify passwords are masked
				for _, user := range result {
					assert.Equal(t, "-", user.Password)
				}
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockJwtService.AssertExpectations(t)
			mockPasswordHasher.AssertExpectations(t)
			mockPaginationService.AssertExpectations(t)
			mockValidationService.AssertExpectations(t)
		})
	}
}

// Context-aware test for Login
func TestLogin_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		loginPayload  dto.LoginDto
		setupMocks    func(*MockUserRepository, *MockJwtService, *MockPasswordHasher, *MockPaginationService, *MockValidationService)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name: "successful login with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_login_123"),
			loginPayload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				user := model.User{Id: 1, Email: "test@example.com", Password: "hashedpassword"}
				loginResponse := dto.LoginResponseDto{AccessToken: "access", RefreshToken: "refresh"}
				
				mockRepo.On("GetByEmail", mock.AnythingOfType("*context.valueCtx"), "test@example.com").Return(user, nil).Once()
				mockPH.On("ComparePasswordHash", "hashedpassword", "password123").Return(nil).Once()
				mockJwt.On("GenerateToken", mock.AnythingOfType("model.User")).Return(loginResponse, nil).Once()
				mockRepo.On("SaveRefreshToken", mock.AnythingOfType("*context.valueCtx"), 1, "refresh", mock.AnythingOfType("time.Time")).Return(nil).Once()
			},
		},
		{
			name: "context cancellation during login",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			loginPayload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed as context is cancelled immediately
			},
			expectCancel: true,
		},
		{
			name: "context timeout during login",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			loginPayload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed as context times out immediately
			},
			expectTimeout: true,
		},
		{
			name: "empty identifier validation error",
			ctx:  context.WithValue(context.Background(), "request_id", "req_empty_id"),
			loginPayload: dto.LoginDto{
				Identifier: "",
				Password:   "password123",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed for validation error
			},
			expectedError: "email is required",
		},
		{
			name: "empty password validation error",
			ctx:  context.WithValue(context.Background(), "request_id", "req_empty_pass"),
			loginPayload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed for validation error
			},
			expectedError: "password is required",
		},
		{
			name: "user not found with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_not_found"),
			loginPayload: dto.LoginDto{
				Identifier: "nonexistent@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockRepo.On("GetByEmail", mock.AnythingOfType("*context.valueCtx"), "nonexistent@example.com").Return(model.User{}, errors.New("user not found")).Once()
			},
			expectedError: "invalid credentials",
		},
		{
			name: "context cancellation during repository operation",
			ctx:  context.WithValue(context.Background(), "request_id", "req_cancel_repo"),
			loginPayload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockRepo.On("GetByEmail", mock.AnythingOfType("*context.valueCtx"), "test@example.com").Return(model.User{}, context.Canceled).Once()
			},
			expectCancel: true,
		},
		{
			name: "invalid password with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_invalid_pass"),
			loginPayload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "wrongpassword",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				user := model.User{Id: 1, Email: "test@example.com", Password: "hashedpassword"}
				mockRepo.On("GetByEmail", mock.AnythingOfType("*context.valueCtx"), "test@example.com").Return(user, nil).Once()
				mockPH.On("ComparePasswordHash", "hashedpassword", "wrongpassword").Return(errors.New("password mismatch")).Once()
			},
			expectedError: "invalid credentials",
		},
		{
			name: "token generation error with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_token_error"),
			loginPayload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				user := model.User{Id: 1, Email: "test@example.com", Password: "hashedpassword"}
				mockRepo.On("GetByEmail", mock.AnythingOfType("*context.valueCtx"), "test@example.com").Return(user, nil).Once()
				mockPH.On("ComparePasswordHash", "hashedpassword", "password123").Return(nil).Once()
				mockJwt.On("GenerateToken", mock.AnythingOfType("model.User")).Return(dto.LoginResponseDto{}, errors.New("token generation failed")).Once()
			},
			expectedError: "failed to create token",
		},
		{
			name: "refresh token save error with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_save_error"),
			loginPayload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				user := model.User{Id: 1, Email: "test@example.com", Password: "hashedpassword"}
				loginResponse := dto.LoginResponseDto{AccessToken: "access", RefreshToken: "refresh"}
				
				mockRepo.On("GetByEmail", mock.AnythingOfType("*context.valueCtx"), "test@example.com").Return(user, nil).Once()
				mockPH.On("ComparePasswordHash", "hashedpassword", "password123").Return(nil).Once()
				mockJwt.On("GenerateToken", mock.AnythingOfType("model.User")).Return(loginResponse, nil).Once()
				mockRepo.On("SaveRefreshToken", mock.AnythingOfType("*context.valueCtx"), 1, "refresh", mock.AnythingOfType("time.Time")).Return(errors.New("database error")).Once()
			},
			expectedError: "internal server error",
		},
		{
			name: "context cancellation during refresh token save",
			ctx:  context.WithValue(context.Background(), "request_id", "req_cancel_save"),
			loginPayload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				user := model.User{Id: 1, Email: "test@example.com", Password: "hashedpassword"}
				loginResponse := dto.LoginResponseDto{AccessToken: "access", RefreshToken: "refresh"}
				
				mockRepo.On("GetByEmail", mock.AnythingOfType("*context.valueCtx"), "test@example.com").Return(user, nil).Once()
				mockPH.On("ComparePasswordHash", "hashedpassword", "password123").Return(nil).Once()
				mockJwt.On("GenerateToken", mock.AnythingOfType("model.User")).Return(loginResponse, nil).Once()
				mockRepo.On("SaveRefreshToken", mock.AnythingOfType("*context.valueCtx"), 1, "refresh", mock.AnythingOfType("time.Time")).Return(context.Canceled).Once()
			},
			expectCancel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockUserRepository)
			mockJwtService := new(MockJwtService)
			mockPasswordHasher := new(MockPasswordHasher)
			mockPaginationService := new(MockPaginationService)
			mockValidationService := new(MockValidationService)

			userService := NewUserservice(mockRepo, mockJwtService, mockPasswordHasher, mockPaginationService, mockValidationService)

			// Setup mock expectations
			tt.setupMocks(mockRepo, mockJwtService, mockPasswordHasher, mockPaginationService, mockValidationService)

			// Execute test
			result, err := userService.Login(tt.ctx, tt.loginPayload)

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
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockJwtService.AssertExpectations(t)
			mockPasswordHasher.AssertExpectations(t)
			mockPaginationService.AssertExpectations(t)
			mockValidationService.AssertExpectations(t)
		})
	}
}

// Context-aware test for RefreshToken
func TestRefreshToken_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		refreshToken  string
		setupMocks    func(*MockUserRepository, *MockJwtService, *MockPasswordHasher, *MockPaginationService, *MockValidationService)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name:         "successful refresh token with context",
			ctx:          context.WithValue(context.Background(), "request_id", "req_refresh_123"),
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				user := model.User{Id: 1, Name: "testuser", Email: "test@example.com"}
				refreshTokenModel := model.RefreshToken{UserID: 1, Token: "valid_refresh_token", ExpiresAt: time.Now().Add(time.Hour)}
				newLoginResponse := dto.LoginResponseDto{AccessToken: "new_access", RefreshToken: "new_refresh"}
				
				mockRepo.On("FindRefreshToken", mock.AnythingOfType("*context.valueCtx"), "valid_refresh_token").Return(refreshTokenModel, nil).Once()
				mockRepo.On("GetUserById", mock.AnythingOfType("*context.valueCtx"), 1).Return(user, nil).Once()
				mockJwt.On("GenerateToken", user).Return(newLoginResponse, nil).Once()
				mockRepo.On("UpdateRefreshToken", mock.AnythingOfType("*context.valueCtx"), "valid_refresh_token", "new_refresh", mock.AnythingOfType("time.Time")).Return(nil).Once()
			},
		},
		{
			name: "context cancellation during refresh token",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed as context is cancelled immediately
			},
			expectCancel: true,
		},
		{
			name: "context timeout during refresh token",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed as context times out immediately
			},
			expectTimeout: true,
		},
		{
			name:         "invalid refresh token format",
			ctx:          context.WithValue(context.Background(), "request_id", "req_invalid_format"),
			refreshToken: "invalid%",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				// No mocks needed for format validation error
			},
			expectedError: "invalid refresh token format",
		},
		{
			name:         "refresh token not found with context",
			ctx:          context.WithValue(context.Background(), "request_id", "req_not_found"),
			refreshToken: "nonexistent_token",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockRepo.On("FindRefreshToken", mock.AnythingOfType("*context.valueCtx"), "nonexistent_token").Return(model.RefreshToken{}, errors.New("refresh token not found")).Once()
			},
			expectedError: "invalid refresh token",
		},
		{
			name:         "refresh token expired with context",
			ctx:          context.WithValue(context.Background(), "request_id", "req_expired"),
			refreshToken: "expired_token",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				expiredToken := model.RefreshToken{UserID: 1, Token: "expired_token", ExpiresAt: time.Now().Add(-time.Hour)}
				mockRepo.On("FindRefreshToken", mock.AnythingOfType("*context.valueCtx"), "expired_token").Return(expiredToken, nil).Once()
			},
			expectedError: "refresh token expired",
		},
		{
			name:         "context cancellation during repository operation",
			ctx:          context.WithValue(context.Background(), "request_id", "req_cancel_repo"),
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				mockRepo.On("FindRefreshToken", mock.AnythingOfType("*context.valueCtx"), "valid_refresh_token").Return(model.RefreshToken{}, context.Canceled).Once()
			},
			expectCancel: true,
		},
		{
			name:         "user not found for refresh token with context",
			ctx:          context.WithValue(context.Background(), "request_id", "req_user_not_found"),
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				refreshTokenModel := model.RefreshToken{UserID: 999, Token: "valid_refresh_token", ExpiresAt: time.Now().Add(time.Hour)}
				mockRepo.On("FindRefreshToken", mock.AnythingOfType("*context.valueCtx"), "valid_refresh_token").Return(refreshTokenModel, nil).Once()
				mockRepo.On("GetUserById", mock.AnythingOfType("*context.valueCtx"), 999).Return(model.User{}, errors.New("user not found")).Once()
			},
			expectedError: "user not found",
		},
		{
			name:         "token generation error with context",
			ctx:          context.WithValue(context.Background(), "request_id", "req_token_gen_error"),
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				user := model.User{Id: 1, Name: "testuser", Email: "test@example.com"}
				refreshTokenModel := model.RefreshToken{UserID: 1, Token: "valid_refresh_token", ExpiresAt: time.Now().Add(time.Hour)}
				
				mockRepo.On("FindRefreshToken", mock.AnythingOfType("*context.valueCtx"), "valid_refresh_token").Return(refreshTokenModel, nil).Once()
				mockRepo.On("GetUserById", mock.AnythingOfType("*context.valueCtx"), 1).Return(user, nil).Once()
				mockJwt.On("GenerateToken", user).Return(dto.LoginResponseDto{}, errors.New("token generation failed")).Once()
			},
			expectedError: "failed to generate new token",
		},
		{
			name:         "refresh token update error with context",
			ctx:          context.WithValue(context.Background(), "request_id", "req_update_error"),
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				user := model.User{Id: 1, Name: "testuser", Email: "test@example.com"}
				refreshTokenModel := model.RefreshToken{UserID: 1, Token: "valid_refresh_token", ExpiresAt: time.Now().Add(time.Hour)}
				newLoginResponse := dto.LoginResponseDto{AccessToken: "new_access", RefreshToken: "new_refresh"}
				
				mockRepo.On("FindRefreshToken", mock.AnythingOfType("*context.valueCtx"), "valid_refresh_token").Return(refreshTokenModel, nil).Once()
				mockRepo.On("GetUserById", mock.AnythingOfType("*context.valueCtx"), 1).Return(user, nil).Once()
				mockJwt.On("GenerateToken", user).Return(newLoginResponse, nil).Once()
				mockRepo.On("UpdateRefreshToken", mock.AnythingOfType("*context.valueCtx"), "valid_refresh_token", "new_refresh", mock.AnythingOfType("time.Time")).Return(errors.New("database update failed")).Once()
			},
			expectedError: "failed to update refresh token",
		},
		{
			name:         "context cancellation during token update",
			ctx:          context.WithValue(context.Background(), "request_id", "req_cancel_update"),
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockRepo *MockUserRepository, mockJwt *MockJwtService, mockPH *MockPasswordHasher, mockPS *MockPaginationService, mockVS *MockValidationService) {
				user := model.User{Id: 1, Name: "testuser", Email: "test@example.com"}
				refreshTokenModel := model.RefreshToken{UserID: 1, Token: "valid_refresh_token", ExpiresAt: time.Now().Add(time.Hour)}
				newLoginResponse := dto.LoginResponseDto{AccessToken: "new_access", RefreshToken: "new_refresh"}
				
				mockRepo.On("FindRefreshToken", mock.AnythingOfType("*context.valueCtx"), "valid_refresh_token").Return(refreshTokenModel, nil).Once()
				mockRepo.On("GetUserById", mock.AnythingOfType("*context.valueCtx"), 1).Return(user, nil).Once()
				mockJwt.On("GenerateToken", user).Return(newLoginResponse, nil).Once()
				mockRepo.On("UpdateRefreshToken", mock.AnythingOfType("*context.valueCtx"), "valid_refresh_token", "new_refresh", mock.AnythingOfType("time.Time")).Return(context.Canceled).Once()
			},
			expectCancel: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockRepo := new(MockUserRepository)
			mockJwtService := new(MockJwtService)
			mockPasswordHasher := new(MockPasswordHasher)
			mockPaginationService := new(MockPaginationService)
			mockValidationService := new(MockValidationService)

			userService := NewUserservice(mockRepo, mockJwtService, mockPasswordHasher, mockPaginationService, mockValidationService)

			// Setup mock expectations
			tt.setupMocks(mockRepo, mockJwtService, mockPasswordHasher, mockPaginationService, mockValidationService)

			// Execute test
			result, err := userService.RefreshToken(tt.ctx, tt.refreshToken)

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
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}

			// Verify mock expectations
			mockRepo.AssertExpectations(t)
			mockJwtService.AssertExpectations(t)
			mockPasswordHasher.AssertExpectations(t)
			mockPaginationService.AssertExpectations(t)
			mockValidationService.AssertExpectations(t)
		})
	}
}

// Test for DeleteUser authorization
func TestDeleteUser_Authorization(t *testing.T) {
	tests := []struct {
		name               string
		requestingUserID   int
		requestingUserRole string
		targetUserID       int
		expectError        bool
		expectedErrorMsg   string
		setupMocks         func(mockRepo *MockUserRepository)
	}{
		{
			name:               "successful self-deletion",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       1,
			expectError:        false,
			setupMocks: func(mockRepo *MockUserRepository) {
				mockRepo.On("DeleteUser", mock.AnythingOfType("*context.valueCtx"), 1).Return(nil).Once()
			},
		},
		{
			name:               "successful admin deletion of other user",
			requestingUserID:   1,
			requestingUserRole: "admin",
			targetUserID:       2,
			expectError:        false,
			setupMocks: func(mockRepo *MockUserRepository) {
				mockRepo.On("DeleteUser", mock.AnythingOfType("*context.valueCtx"), 2).Return(nil).Once()
			},
		},
		{
			name:               "authorization failure - regular user cannot delete other user",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       2,
			expectError:        true,
			expectedErrorMsg:   "Forbidden: You can only modify your own account",
			setupMocks: func(mockRepo *MockUserRepository) {
				// No repository mock needed as authorization should fail first
			},
		},
		{
			name:               "authorization failure - invalid requesting user ID",
			requestingUserID:   0,
			requestingUserRole: "user",
			targetUserID:       1,
			expectError:        true,
			expectedErrorMsg:   "Invalid requesting user ID",
			setupMocks: func(mockRepo *MockUserRepository) {
				// No repository mock needed as authorization should fail first
			},
		},
		{
			name:               "authorization failure - invalid target user ID",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       0,
			expectError:        true,
			expectedErrorMsg:   "Invalid target user ID",
			setupMocks: func(mockRepo *MockUserRepository) {
				// No repository mock needed as authorization should fail first
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockRepo := &MockUserRepository{}
			mockJwt := &MockJwtService{}
			mockPH := &MockPasswordHasher{}
			mockPS := &MockPaginationService{}
			mockVS := &MockValidationService{}

			// Setup mocks
			tt.setupMocks(mockRepo)

			// Create service
			userService := NewUserservice(mockRepo, mockJwt, mockPH, mockPS, mockVS)

			// Create context
			ctx := context.WithValue(context.Background(), "request_id", "test_req")

			// Call DeleteUser
			err := userService.DeleteUser(ctx, tt.requestingUserID, tt.requestingUserRole, tt.targetUserID)

			// Verify results
			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedErrorMsg != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorMsg)
				}
			} else {
				assert.NoError(t, err)
			}

			// Assert mock expectations
			mockRepo.AssertExpectations(t)
			mockJwt.AssertExpectations(t)
			mockPH.AssertExpectations(t)
			mockPS.AssertExpectations(t)
			mockVS.AssertExpectations(t)
		})
	}
}