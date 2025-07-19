package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/service"
	"develapar-server/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserServiceWithContext is a mock for service.UserService with context support
type MockUserServiceWithContext struct {
	mock.Mock
}

func (m *MockUserServiceWithContext) CreateNewUser(ctx context.Context, payload model.User) (model.User, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserServiceWithContext) FindUserById(ctx context.Context, id string) (model.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserServiceWithContext) FindAllUser(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserServiceWithContext) FindAllUserWithPagination(ctx context.Context, page, limit int) (service.PaginationResult, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(service.PaginationResult), args.Error(1)
}

func (m *MockUserServiceWithContext) Login(ctx context.Context, payload dto.LoginDto) (dto.LoginResponseDto, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

func (m *MockUserServiceWithContext) RefreshToken(ctx context.Context, refreshToken string) (dto.LoginResponseDto, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

// MockErrorHandler is a mock for middleware.ErrorHandler
type MockErrorHandler struct {
	mock.Mock
}

func (m *MockErrorHandler) HandleError(ctx context.Context, c *gin.Context, err *utils.AppError) {
	m.Called(ctx, c, err)
	// Set appropriate HTTP status and response
	c.JSON(err.StatusCode, gin.H{"error": err.Message})
}

func (m *MockErrorHandler) ValidationError(ctx context.Context, field string, message string) *utils.AppError {
	args := m.Called(ctx, field, message)
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorHandler) TimeoutError(ctx context.Context, operation string) *utils.AppError {
	args := m.Called(ctx, operation)
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorHandler) CancellationError(ctx context.Context, operation string) *utils.AppError {
	args := m.Called(ctx, operation)
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorHandler) WrapError(ctx context.Context, err error, code string, message string) *utils.AppError {
	args := m.Called(ctx, err, code, message)
	return args.Get(0).(*utils.AppError)
}

// Context-aware test for Login Handler
func TestLoginHandler_WithContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        dto.LoginDto
		setupMocks     func(*MockUserServiceWithContext, *MockErrorHandler)
		expectedStatus int
		expectedError  string
		expectTimeout  bool
		expectCancel   bool
	}{
		{
			name: "successful login with context",
			payload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				loginResponse := dto.LoginResponseDto{
					AccessToken:  "mockAccessToken",
					RefreshToken: "mockRefreshToken",
				}
				mockService.On("Login", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("dto.LoginDto")).Return(loginResponse, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "context timeout during login",
			payload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("Login", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("dto.LoginDto")).Return(dto.LoginResponseDto{}, context.DeadlineExceeded).Once()
				timeoutErr := &utils.AppError{
					Code:       utils.ErrTimeout,
					Message:    "Login operation timed out",
					StatusCode: 408,
				}
				mockErrorHandler.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "login").Return(timeoutErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), timeoutErr).Once()
			},
			expectedStatus: http.StatusRequestTimeout,
			expectTimeout:  true,
		},
		{
			name: "context cancellation during login",
			payload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "password123",
			},
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("Login", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("dto.LoginDto")).Return(dto.LoginResponseDto{}, context.Canceled).Once()
				cancelErr := &utils.AppError{
					Code:       utils.ErrCancelled,
					Message:    "Login operation was cancelled",
					StatusCode: 499,
				}
				mockErrorHandler.On("CancellationError", mock.AnythingOfType("*context.timerCtx"), "login").Return(cancelErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), cancelErr).Once()
			},
			expectedStatus: 499,
			expectCancel:   true,
		},
		{
			name: "invalid credentials with context",
			payload: dto.LoginDto{
				Identifier: "test@example.com",
				Password:   "wrongpassword",
			},
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				serviceErr := &utils.AppError{
					Code:       utils.ErrUnauthorized,
					Message:    "Invalid credentials",
					StatusCode: 401,
				}
				mockService.On("Login", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("dto.LoginDto")).Return(dto.LoginResponseDto{}, serviceErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), serviceErr).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid credentials",
		},
		{
			name:    "invalid JSON payload with context",
			payload: dto.LoginDto{}, // Will be ignored, we'll send invalid JSON
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				validationErr := &utils.AppError{
					Code:       utils.ErrValidation,
					Message:    "Invalid request payload",
					StatusCode: 400,
				}
				mockErrorHandler.On("ValidationError", mock.AnythingOfType("*context.timerCtx"), "payload", mock.AnythingOfType("string")).Return(validationErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), validationErr).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request payload",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockService := new(MockUserServiceWithContext)
			mockErrorHandler := new(MockErrorHandler)
			
			// Setup mock expectations
			tt.setupMocks(mockService, mockErrorHandler)

			// Create router and controller
			router := gin.New()
			userController := &UserController{
				service:        mockService,
				rg:             router.Group("/api/v1"),
				errorHandler:   mockErrorHandler,
				responseHelper: utils.NewResponseHelper(),
			}
			userController.Route()

			// Prepare request
			var body []byte
			if tt.name == "invalid JSON payload with context" {
				body = []byte("invalid json")
			} else {
				body, _ = json.Marshal(tt.payload)
			}

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(rr, req)

			// Verify results
			if tt.expectTimeout {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectCancel {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectedError != "" {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				// For successful login, check response structure
				if tt.expectedStatus == http.StatusOK {
					var responseBody map[string]interface{}
					json.Unmarshal(rr.Body.Bytes(), &responseBody)
					assert.Contains(t, responseBody, "message")
					assert.Contains(t, responseBody, "access_token")
				}
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
			mockErrorHandler.AssertExpectations(t)
		})
	}
}

// Context-aware test for Register User Handler
func TestRegisterUser_WithContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		payload        model.User
		setupMocks     func(*MockUserServiceWithContext, *MockErrorHandler)
		expectedStatus int
		expectedError  string
		expectTimeout  bool
		expectCancel   bool
	}{
		{
			name: "successful user registration with context",
			payload: model.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				createdUser := model.User{
					Id:       1,
					Name:     "Test User",
					Email:    "test@example.com",
					Password: "-", // Masked password
					Role:     "user",
				}
				mockService.On("CreateNewUser", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("model.User")).Return(createdUser, nil).Once()
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "context timeout during user registration",
			payload: model.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("CreateNewUser", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("model.User")).Return(model.User{}, context.DeadlineExceeded).Once()
				timeoutErr := &utils.AppError{
					Code:       utils.ErrTimeout,
					Message:    "User registration operation timed out",
					StatusCode: 408,
				}
				mockErrorHandler.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "user registration").Return(timeoutErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), timeoutErr).Once()
			},
			expectedStatus: http.StatusRequestTimeout,
			expectTimeout:  true,
		},
		{
			name: "context cancellation during user registration",
			payload: model.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("CreateNewUser", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("model.User")).Return(model.User{}, context.Canceled).Once()
				cancelErr := &utils.AppError{
					Code:       utils.ErrCancelled,
					Message:    "User registration operation was cancelled",
					StatusCode: 499,
				}
				mockErrorHandler.On("CancellationError", mock.AnythingOfType("*context.timerCtx"), "user registration").Return(cancelErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), cancelErr).Once()
			},
			expectedStatus: 499,
			expectCancel:   true,
		},
		{
			name: "validation error with context",
			payload: model.User{
				Name:     "Test User",
				Email:    "test@example.com",
				Password: "password123",
				Role:     "user",
			},
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				validationErr := &utils.AppError{
					Code:       utils.ErrValidation,
					Message:    "User validation failed",
					StatusCode: 400,
				}
				mockService.On("CreateNewUser", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("model.User")).Return(model.User{}, validationErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), validationErr).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "User validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockService := new(MockUserServiceWithContext)
			mockErrorHandler := new(MockErrorHandler)
			
			// Setup mock expectations
			tt.setupMocks(mockService, mockErrorHandler)

			// Create router and controller
			router := gin.New()
			userController := &UserController{
				service:        mockService,
				rg:             router.Group("/api/v1"),
				errorHandler:   mockErrorHandler,
				responseHelper: utils.NewResponseHelper(),
			}
			userController.Route()

			// Prepare request
			body, _ := json.Marshal(tt.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(rr, req)

			// Verify results
			if tt.expectTimeout {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectCancel {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectedError != "" {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				// For successful registration, check response structure
				if tt.expectedStatus == http.StatusCreated {
					var responseBody map[string]interface{}
					json.Unmarshal(rr.Body.Bytes(), &responseBody)
					assert.Contains(t, responseBody, "message")
					assert.Contains(t, responseBody, "user")
				}
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
			mockErrorHandler.AssertExpectations(t)
		})
	}
}

// Context-aware test for Find User By ID Handler
func TestFindUserByIdHandler_WithContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		setupMocks     func(*MockUserServiceWithContext, *MockErrorHandler)
		expectedStatus int
		expectedError  string
		expectTimeout  bool
		expectCancel   bool
	}{
		{
			name:   "successful user retrieval with context",
			userID: "1",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				expectedUser := model.User{
					Id:    1,
					Name:  "Test User",
					Email: "test@example.com",
					Role:  "user",
				}
				mockService.On("FindUserById", mock.AnythingOfType("*context.timerCtx"), "1").Return(expectedUser, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "context timeout during user retrieval",
			userID: "1",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("FindUserById", mock.AnythingOfType("*context.timerCtx"), "1").Return(model.User{}, context.DeadlineExceeded).Once()
				timeoutErr := &utils.AppError{
					Code:       utils.ErrTimeout,
					Message:    "Get user operation timed out",
					StatusCode: 408,
				}
				mockErrorHandler.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "get user").Return(timeoutErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), timeoutErr).Once()
			},
			expectedStatus: http.StatusRequestTimeout,
			expectTimeout:  true,
		},
		{
			name:   "context cancellation during user retrieval",
			userID: "1",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("FindUserById", mock.AnythingOfType("*context.timerCtx"), "1").Return(model.User{}, context.Canceled).Once()
				cancelErr := &utils.AppError{
					Code:       utils.ErrCancelled,
					Message:    "Get user operation was cancelled",
					StatusCode: 499,
				}
				mockErrorHandler.On("CancellationError", mock.AnythingOfType("*context.timerCtx"), "get user").Return(cancelErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), cancelErr).Once()
			},
			expectedStatus: 499,
			expectCancel:   true,
		},
		{
			name:   "user not found with context",
			userID: "999",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				notFoundErr := &utils.AppError{
					Code:       utils.ErrNotFound,
					Message:    "User not found",
					StatusCode: 404,
				}
				mockService.On("FindUserById", mock.AnythingOfType("*context.timerCtx"), "999").Return(model.User{}, notFoundErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), notFoundErr).Once()
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "User not found",
		},
		{
			name:   "empty user ID with context validation",
			userID: "",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				validationErr := &utils.AppError{
					Code:       utils.ErrValidation,
					Message:    "User ID is required",
					StatusCode: 400,
				}
				mockErrorHandler.On("ValidationError", mock.AnythingOfType("*context.timerCtx"), "user_id", "User ID is required").Return(validationErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), validationErr).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "User ID is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockService := new(MockUserServiceWithContext)
			mockErrorHandler := new(MockErrorHandler)
			
			// Setup mock expectations
			tt.setupMocks(mockService, mockErrorHandler)

			// Create router and controller
			router := gin.New()
			userController := &UserController{
				service:        mockService,
				rg:             router.Group("/api/v1"),
				errorHandler:   mockErrorHandler,
				responseHelper: utils.NewResponseHelper(),
			}
			userController.Route()

			// Prepare request
			url := "/api/v1/users/" + tt.userID
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			rr := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(rr, req)

			// Verify results
			if tt.expectTimeout {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectCancel {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectedError != "" {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				// For successful retrieval, check response structure
				if tt.expectedStatus == http.StatusOK {
					var responseBody map[string]interface{}
					json.Unmarshal(rr.Body.Bytes(), &responseBody)
					assert.Contains(t, responseBody, "message")
					assert.Contains(t, responseBody, "user")
				}
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
			mockErrorHandler.AssertExpectations(t)
		})
	}
}

// Context-aware test for Find All Users Handler
func TestFindAllUserHandler_WithContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupMocks     func(*MockUserServiceWithContext, *MockErrorHandler)
		expectedStatus int
		expectedError  string
		expectTimeout  bool
		expectCancel   bool
	}{
		{
			name: "successful users retrieval with context",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				expectedUsers := []model.User{
					{Id: 1, Name: "User 1", Email: "user1@example.com", Role: "user"},
					{Id: 2, Name: "User 2", Email: "user2@example.com", Role: "admin"},
				}
				mockService.On("FindAllUser", mock.AnythingOfType("*context.timerCtx")).Return(expectedUsers, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "context timeout during users retrieval",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("FindAllUser", mock.AnythingOfType("*context.timerCtx")).Return([]model.User{}, context.DeadlineExceeded).Once()
				timeoutErr := &utils.AppError{
					Code:       utils.ErrTimeout,
					Message:    "Get all users operation timed out",
					StatusCode: 408,
				}
				mockErrorHandler.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "get all users").Return(timeoutErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), timeoutErr).Once()
			},
			expectedStatus: http.StatusRequestTimeout,
			expectTimeout:  true,
		},
		{
			name: "context cancellation during users retrieval",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("FindAllUser", mock.AnythingOfType("*context.timerCtx")).Return([]model.User{}, context.Canceled).Once()
				cancelErr := &utils.AppError{
					Code:       utils.ErrCancelled,
					Message:    "Get all users operation was cancelled",
					StatusCode: 499,
				}
				mockErrorHandler.On("CancellationError", mock.AnythingOfType("*context.timerCtx"), "get all users").Return(cancelErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), cancelErr).Once()
			},
			expectedStatus: 499,
			expectCancel:   true,
		},
		{
			name: "database error with context",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				dbErr := &utils.AppError{
					Code:       utils.ErrDatabase,
					Message:    "Database connection failed",
					StatusCode: 500,
				}
				mockService.On("FindAllUser", mock.AnythingOfType("*context.timerCtx")).Return([]model.User{}, dbErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), dbErr).Once()
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockService := new(MockUserServiceWithContext)
			mockErrorHandler := new(MockErrorHandler)
			
			// Setup mock expectations
			tt.setupMocks(mockService, mockErrorHandler)

			// Create router and controller
			router := gin.New()
			userController := &UserController{
				service:        mockService,
				rg:             router.Group("/api/v1"),
				errorHandler:   mockErrorHandler,
				responseHelper: utils.NewResponseHelper(),
			}
			userController.Route()

			// Prepare request
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/", nil)
			rr := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(rr, req)

			// Verify results
			if tt.expectTimeout {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectCancel {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectedError != "" {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				// For successful retrieval, check response structure
				if tt.expectedStatus == http.StatusOK {
					var responseBody map[string]interface{}
					json.Unmarshal(rr.Body.Bytes(), &responseBody)
					assert.Contains(t, responseBody, "message")
					assert.Contains(t, responseBody, "users")
				}
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
			mockErrorHandler.AssertExpectations(t)
		})
	}
}

// Context-aware test for Find All Users With Pagination Handler
func TestFindAllUserWithPaginationHandler_WithContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    string
		setupMocks     func(*MockUserServiceWithContext, *MockErrorHandler)
		expectedStatus int
		expectedError  string
		expectTimeout  bool
		expectCancel   bool
	}{
		{
			name:        "successful paginated users retrieval with context",
			queryParams: "?page=1&limit=10",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				users := []model.User{
					{Id: 1, Name: "User 1", Email: "user1@example.com", Role: "user"},
					{Id: 2, Name: "User 2", Email: "user2@example.com", Role: "admin"},
				}
				paginationResult := service.PaginationResult{
					Data: users,
					Metadata: service.PaginationMetadata{
						Page:       1,
						Limit:      10,
						Total:      2,
						TotalPages: 1,
						HasNext:    false,
						HasPrev:    false,
					},
				}
				mockService.On("FindAllUserWithPagination", mock.AnythingOfType("*context.timerCtx"), 1, 10).Return(paginationResult, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "context timeout during paginated users retrieval",
			queryParams: "?page=1&limit=10",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("FindAllUserWithPagination", mock.AnythingOfType("*context.timerCtx"), 1, 10).Return(service.PaginationResult{}, context.DeadlineExceeded).Once()
				timeoutErr := &utils.AppError{
					Code:       utils.ErrTimeout,
					Message:    "Get paginated users operation timed out",
					StatusCode: 408,
				}
				mockErrorHandler.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "get paginated users").Return(timeoutErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), timeoutErr).Once()
			},
			expectedStatus: http.StatusRequestTimeout,
			expectTimeout:  true,
		},
		{
			name:        "context cancellation during paginated users retrieval",
			queryParams: "?page=1&limit=10",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("FindAllUserWithPagination", mock.AnythingOfType("*context.timerCtx"), 1, 10).Return(service.PaginationResult{}, context.Canceled).Once()
				cancelErr := &utils.AppError{
					Code:       utils.ErrCancelled,
					Message:    "Get paginated users operation was cancelled",
					StatusCode: 499,
				}
				mockErrorHandler.On("CancellationError", mock.AnythingOfType("*context.timerCtx"), "get paginated users").Return(cancelErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), cancelErr).Once()
			},
			expectedStatus: 499,
			expectCancel:   true,
		},
		{
			name:        "invalid pagination parameters with context",
			queryParams: "?page=0&limit=-1",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				validationErr := &utils.AppError{
					Code:       utils.ErrValidation,
					Message:    "Page must be a positive integer",
					StatusCode: 400,
				}
				mockErrorHandler.On("ValidationError", mock.AnythingOfType("*context.timerCtx"), "page", "Page must be a positive integer").Return(validationErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), validationErr).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Page must be a positive integer",
		},
		{
			name:        "limit too high with context validation",
			queryParams: "?page=1&limit=200",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				validationErr := &utils.AppError{
					Code:       utils.ErrValidation,
					Message:    "Limit must be a positive integer between 1 and 100",
					StatusCode: 400,
				}
				mockErrorHandler.On("ValidationError", mock.AnythingOfType("*context.timerCtx"), "limit", "Limit must be a positive integer between 1 and 100").Return(validationErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), validationErr).Once()
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Limit must be a positive integer between 1 and 100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockService := new(MockUserServiceWithContext)
			mockErrorHandler := new(MockErrorHandler)
			
			// Setup mock expectations
			tt.setupMocks(mockService, mockErrorHandler)

			// Create router and controller
			router := gin.New()
			userController := &UserController{
				service:        mockService,
				rg:             router.Group("/api/v1"),
				errorHandler:   mockErrorHandler,
				responseHelper: utils.NewResponseHelper(),
			}
			userController.Route()

			// Prepare request
			url := "/api/v1/users/paginated" + tt.queryParams
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			rr := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(rr, req)

			// Verify results
			if tt.expectTimeout {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectCancel {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectedError != "" {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				// For successful retrieval, check response structure
				if tt.expectedStatus == http.StatusOK {
					var responseBody map[string]interface{}
					json.Unmarshal(rr.Body.Bytes(), &responseBody)
					assert.Contains(t, responseBody, "message")
					assert.Contains(t, responseBody, "users")
				}
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
			mockErrorHandler.AssertExpectations(t)
		})
	}
}

// Context-aware test for Refresh Token Handler
func TestRefreshTokenHandler_WithContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		refreshToken   string
		setupMocks     func(*MockUserServiceWithContext, *MockErrorHandler)
		expectedStatus int
		expectedError  string
		expectTimeout  bool
		expectCancel   bool
	}{
		{
			name:         "successful token refresh with context",
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				refreshResponse := dto.LoginResponseDto{
					AccessToken:  "new_access_token",
					RefreshToken: "new_refresh_token",
					ExpiresIn:    3600,
					TokenType:    "Bearer",
				}
				mockService.On("RefreshToken", mock.AnythingOfType("*context.timerCtx"), "valid_refresh_token").Return(refreshResponse, nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "context timeout during token refresh",
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("RefreshToken", mock.AnythingOfType("*context.timerCtx"), "valid_refresh_token").Return(dto.LoginResponseDto{}, context.DeadlineExceeded).Once()
				timeoutErr := &utils.AppError{
					Code:       utils.ErrTimeout,
					Message:    "Token refresh operation timed out",
					StatusCode: 408,
				}
				mockErrorHandler.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "token refresh").Return(timeoutErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), timeoutErr).Once()
			},
			expectedStatus: http.StatusRequestTimeout,
			expectTimeout:  true,
		},
		{
			name:         "context cancellation during token refresh",
			refreshToken: "valid_refresh_token",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				mockService.On("RefreshToken", mock.AnythingOfType("*context.timerCtx"), "valid_refresh_token").Return(dto.LoginResponseDto{}, context.Canceled).Once()
				cancelErr := &utils.AppError{
					Code:       utils.ErrCancelled,
					Message:    "Token refresh operation was cancelled",
					StatusCode: 499,
				}
				mockErrorHandler.On("CancellationError", mock.AnythingOfType("*context.timerCtx"), "token refresh").Return(cancelErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), cancelErr).Once()
			},
			expectedStatus: 499,
			expectCancel:   true,
		},
		{
			name:         "invalid refresh token with context",
			refreshToken: "invalid_refresh_token",
			setupMocks: func(mockService *MockUserServiceWithContext, mockErrorHandler *MockErrorHandler) {
				unauthorizedErr := &utils.AppError{
					Code:       utils.ErrUnauthorized,
					Message:    "Invalid refresh token",
					StatusCode: 401,
				}
				mockService.On("RefreshToken", mock.AnythingOfType("*context.timerCtx"), "invalid_refresh_token").Return(dto.LoginResponseDto{}, unauthorizedErr).Once()
				mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), unauthorizedErr).Once()
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid refresh token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockService := new(MockUserServiceWithContext)
			mockErrorHandler := new(MockErrorHandler)
			
			// Setup mock expectations
			tt.setupMocks(mockService, mockErrorHandler)

			// Create router and controller
			router := gin.New()
			userController := &UserController{
				service:        mockService,
				rg:             router.Group("/api/v1"),
				errorHandler:   mockErrorHandler,
				responseHelper: utils.NewResponseHelper(),
			}
			userController.Route()

			// Prepare request with refresh token in cookie
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
			req.AddCookie(&http.Cookie{
				Name:  "refresh_token",
				Value: tt.refreshToken,
			})
			rr := httptest.NewRecorder()

			// Execute request
			router.ServeHTTP(rr, req)

			// Verify results
			if tt.expectTimeout {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectCancel {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else if tt.expectedError != "" {
				assert.Equal(t, tt.expectedStatus, rr.Code)
			} else {
				assert.Equal(t, tt.expectedStatus, rr.Code)
				// For successful refresh, check response structure
				if tt.expectedStatus == http.StatusOK {
					var responseBody map[string]interface{}
					json.Unmarshal(rr.Body.Bytes(), &responseBody)
					assert.Contains(t, responseBody, "message")
					assert.Contains(t, responseBody, "access_token")
					// Check that new refresh token cookie is set
					assert.NotEmpty(t, rr.Header().Get("Set-Cookie"))
				}
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
			mockErrorHandler.AssertExpectations(t)
		})
	}
}