package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService for testing
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateNewUser(ctx context.Context, user model.User) (model.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, loginDto dto.LoginDto) (dto.LoginResponseDto, error) {
	args := m.Called(ctx, loginDto)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

func (m *MockUserService) FindUserById(ctx context.Context, userId string) (model.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserService) FindAllUser(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserService) FindAllUserWithPagination(ctx context.Context, page, limit int) (service.PaginationResult, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(service.PaginationResult), args.Error(1)
}

func (m *MockUserService) UpdateUser(ctx context.Context, userId int, req dto.UpdateUserRequest) (model.User, error) {
	args := m.Called(ctx, userId, req)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserService) DeleteUser(ctx context.Context, userId int) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}

func (m *MockUserService) RefreshToken(ctx context.Context, refreshToken string) (dto.LoginResponseDto, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

// MockAuthMiddleware for testing
type MockAuthMiddleware struct {
	mock.Mock
}

func (m *MockAuthMiddleware) CheckToken(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// This will be overridden by the test setup middleware
		c.Next()
	}
}

// TestLogger implements middleware.Logger interface for testing
type TestLogger struct{}

func (l *TestLogger) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	// Simple test logger - just print to stdout
}

func (l *TestLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	// Simple test logger - just print to stdout
}

func (l *TestLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	// Simple test logger - just print to stdout
}

func TestUpdateUserHandler_Authorization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		contextUserID  int
		contextRole    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "User can update own profile",
			userID:         "1",
			contextUserID:  1,
			contextRole:    "user",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "User cannot update other user's profile",
			userID:         "2",
			contextUserID:  1,
			contextRole:    "user",
			expectedStatus: http.StatusForbidden,
			expectedError:  "Forbidden",
		},
		{
			name:           "Admin can update any user's profile",
			userID:         "2",
			contextUserID:  1,
			contextRole:    "admin",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Invalid user ID format",
			userID:         "invalid",
			contextUserID:  1,
			contextRole:    "user",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid user ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			mockAuth := new(MockAuthMiddleware)
			testLogger := &TestLogger{}
			errorHandler := middleware.NewErrorHandler(testLogger)

			// Setup router
			router := gin.New()
			rg := router.Group("/api")
			
			// Create controller
			controller := NewUserController(mockService, mockAuth, rg, errorHandler)

			// Setup routes first
			controller.Route()

			// Override the auth middleware to set context values
			router.Use(func(c *gin.Context) {
				c.Set("userId", tt.contextUserID)
				c.Set("role", tt.contextRole)
				c.Next()
			})

			// Mock service call for successful cases
			if tt.expectedStatus == http.StatusOK {
				expectedUser := model.User{
					Id:    1,
					Name:  "Updated User",
					Email: "updated@example.com",
				}
				mockService.On("UpdateUser", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("dto.UpdateUserRequest")).Return(expectedUser, nil)
			}

			// Create request
			updateReq := dto.UpdateUserRequest{
				Name: stringPtr("Updated Name"),
			}
			reqBody, _ := json.Marshal(updateReq)

			req, _ := http.NewRequest("PUT", "/api/users/"+tt.userID, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response["error"].(map[string]interface{})["message"].(string), tt.expectedError)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestUpdateUserHandler_MissingContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		setupContext   func(*gin.Context)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Missing user ID in context",
			setupContext: func(c *gin.Context) {
				c.Set("role", "user")
				// Don't set userId
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Failed to extract user ID from context",
		},
		{
			name: "Missing role in context",
			setupContext: func(c *gin.Context) {
				c.Set("userId", 1)
				// Don't set role
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Failed to extract user role from context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			mockAuth := new(MockAuthMiddleware)
			errorHandler := middleware.NewErrorHandler(nil) // Use default logger

			// Setup router
			router := gin.New()
			rg := router.Group("/api")
			
			// Create controller
			controller := NewUserController(mockService, mockAuth, rg, errorHandler)

			// Setup routes first
			controller.Route()

			// Mock auth middleware to set context values
			router.Use(func(c *gin.Context) {
				tt.setupContext(c)
				c.Next()
			})

			// Create request
			updateReq := dto.UpdateUserRequest{
				Name: stringPtr("Updated Name"),
			}
			reqBody, _ := json.Marshal(updateReq)

			req, _ := http.NewRequest("PUT", "/api/users/1", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Contains(t, response["error"].(map[string]interface{})["message"].(string), tt.expectedError)

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}