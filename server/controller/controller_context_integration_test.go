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

// MockUserServiceIntegration implements service.UserService for integration testing
type MockUserServiceIntegration struct {
	mock.Mock
}

func (m *MockUserServiceIntegration) Login(ctx context.Context, payload dto.LoginDto) (dto.LoginResponseDto, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

func (m *MockUserServiceIntegration) CreateNewUser(ctx context.Context, payload model.User) (model.User, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserServiceIntegration) FindUserById(ctx context.Context, userId string) (model.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserServiceIntegration) FindAllUser(ctx context.Context) ([]model.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserServiceIntegration) FindAllUserWithPagination(ctx context.Context, page, limit int) (service.PaginationResult, error) {
	args := m.Called(ctx, page, limit)
	return args.Get(0).(service.PaginationResult), args.Error(1)
}

func (m *MockUserServiceIntegration) RefreshToken(ctx context.Context, refreshToken string) (dto.LoginResponseDto, error) {
	args := m.Called(ctx, refreshToken)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

// MockCategoryServiceIntegration implements service.CategoryService for integration testing
type MockCategoryServiceIntegration struct {
	mock.Mock
}

func (m *MockCategoryServiceIntegration) CreateCategory(ctx context.Context, payload model.Category) (model.Category, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(model.Category), args.Error(1)
}

func (m *MockCategoryServiceIntegration) FindAll(ctx context.Context) ([]model.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Category), args.Error(1)
}

func (m *MockCategoryServiceIntegration) UpdateCategory(ctx context.Context, id int, req dto.UpdateCategoryRequest) (model.Category, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(model.Category), args.Error(1)
}

func (m *MockCategoryServiceIntegration) DeleteCategory(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockErrorHandlerIntegration implements middleware.ErrorHandler for integration testing
type MockErrorHandlerIntegration struct {
	mock.Mock
}

func (m *MockErrorHandlerIntegration) HandleError(ctx context.Context, c *gin.Context, err error) {
	m.Called(ctx, c, err)
	
	// Simulate error handling behavior
	var appErr *utils.AppError
	if ae, ok := err.(*utils.AppError); ok {
		appErr = ae
	} else {
		appErr = &utils.AppError{
			Code:       utils.ErrInternal,
			Message:    "Internal server error",
			StatusCode: 500,
			Cause:      err,
			Timestamp:  time.Now(),
		}
	}
	
	// Extract context info
	if ctx != nil {
		if rid, ok := ctx.Value("request_id").(string); ok {
			appErr.RequestID = rid
		}
		if uid, ok := ctx.Value("user_id").(string); ok {
			appErr.UserID = uid
		}
	}
	
	errorResponse := dto.ErrorResponseFromError(ctx, appErr.Code, appErr.Message, nil)
	c.JSON(utils.GetStatusCode(appErr), errorResponse)
	c.Abort()
}

func (m *MockErrorHandlerIntegration) WrapError(ctx context.Context, err error, code string, message string) *utils.AppError {
	args := m.Called(ctx, err, code, message)
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorHandlerIntegration) ValidationError(ctx context.Context, field string, message string) *utils.AppError {
	args := m.Called(ctx, field, message)
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorHandlerIntegration) TimeoutError(ctx context.Context, operation string) *utils.AppError {
	args := m.Called(ctx, operation)
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorHandlerIntegration) CancellationError(ctx context.Context, operation string) *utils.AppError {
	args := m.Called(ctx, operation)
	return args.Get(0).(*utils.AppError)
}

// MockAuthMiddlewareIntegration implements middleware.AuthMiddleware for integration testing
type MockAuthMiddlewareIntegration struct {
	mock.Mock
}

func (m *MockAuthMiddlewareIntegration) CheckToken(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Simulate successful authentication
		c.Set("userId", "test-user-123")
		c.Next()
	}
}

// setupIntegrationTestRouter creates a test router with context middleware
func setupIntegrationTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add context middleware
	contextManager := middleware.NewContextManager()
	contextMiddleware := middleware.NewContextMiddleware(contextManager)
	router.Use(contextMiddleware.InjectContext())
	
	return router
}

func TestIntegration_UserController_ContextPropagation(t *testing.T) {
	// Setup
	mockService := new(MockUserServiceIntegration)
	mockErrorHandler := new(MockErrorHandlerIntegration)
	
	router := setupIntegrationTestRouter()
	rg := router.Group("/api")
	controller := NewUserController(mockService, rg, mockErrorHandler)
	controller.Route()
	
	requestID := "integration-test-123"
	
	// Mock service with context verification
	mockService.On("FindAllUser", mock.AnythingOfType("*context.timerCtx")).Return([]model.User{
		{Id: 1, Name: "Test User", Email: "test@example.com"},
	}, nil)
	
	// Create request with context headers
	req, _ := http.NewRequest("GET", "/api/users/", nil)
	req.Header.Set("X-Request-ID", requestID)
	
	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Debug: Print response body and status
	t.Logf("Response Status: %d", w.Code)
	t.Logf("Response Body: %s", w.Body.String())
	t.Logf("Response Headers: %v", w.Header())
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify response contains context information
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, requestID, response.Meta.RequestID)
	assert.True(t, response.Meta.ProcessingTime > 0)
	
	// Verify response headers
	assert.Equal(t, requestID, w.Header().Get("X-Request-ID"))
	
	mockService.AssertExpectations(t)
}

func TestIntegration_UserController_LoginWithContextTimeout(t *testing.T) {
	// Setup
	mockService := new(MockUserServiceIntegration)
	mockErrorHandler := new(MockErrorHandlerIntegration)
	
	router := setupIntegrationTestRouter()
	rg := router.Group("/api")
	controller := NewUserController(mockService, rg, mockErrorHandler)
	controller.Route()
	
	requestID := "timeout-test-123"
	
	// Mock service to simulate timeout by returning context.DeadlineExceeded
	mockService.On("Login", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("dto.LoginDto")).Return(
		dto.LoginResponseDto{},
		context.DeadlineExceeded,
	)
	
	// Mock error handler - the controller will call WrapError for DeadlineExceeded
	expectedTimeoutError := &utils.AppError{
		Code:       utils.ErrUnauthorized,
		Message:    "Authentication failed",
		StatusCode: 401,
		RequestID:  requestID,
		Timestamp:  time.Now(),
	}
	
	mockErrorHandler.On("WrapError", mock.AnythingOfType("*context.timerCtx"), context.DeadlineExceeded, utils.ErrUnauthorized, "Authentication failed").Return(expectedTimeoutError)
	mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), expectedTimeoutError).Return()
	
	// Create request
	loginPayload := dto.LoginDto{
		Identifier: "test@example.com",
		Password:   "password123",
	}
	
	jsonPayload, _ := json.Marshal(loginPayload)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)
	
	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert timeout response (controller treats DeadlineExceeded as unauthorized)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, utils.ErrUnauthorized, response.Error.Code)
	assert.Equal(t, requestID, response.Error.RequestID)
	
	mockService.AssertExpectations(t)
	mockErrorHandler.AssertExpectations(t)
}

func TestIntegration_UserController_RegisterWithContextCancellation(t *testing.T) {
	// Setup
	mockService := new(MockUserServiceIntegration)
	mockErrorHandler := new(MockErrorHandlerIntegration)
	
	router := setupIntegrationTestRouter()
	rg := router.Group("/api")
	controller := NewUserController(mockService, rg, mockErrorHandler)
	controller.Route()
	
	requestID := "cancellation-test-123"
	
	// Mock service to detect cancellation
	mockService.On("CreateNewUser", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("model.User")).Return(
		model.User{},
		context.Canceled,
	)
	
	// Mock error handler for cancellation
	expectedCancelError := &utils.AppError{
		Code:       utils.ErrCancelled,
		Message:    "user registration operation was cancelled",
		StatusCode: 499,
		RequestID:  requestID,
		Timestamp:  time.Now(),
	}
	
	mockErrorHandler.On("CancellationError", mock.AnythingOfType("*context.timerCtx"), "user registration").Return(expectedCancelError)
	mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), expectedCancelError).Return()
	
	// Create request
	userPayload := model.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}
	
	jsonPayload, _ := json.Marshal(userPayload)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonPayload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)
	
	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert cancellation response
	assert.Equal(t, 499, w.Code) // Client Closed Request
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, utils.ErrCancelled, response.Error.Code)
	assert.Equal(t, requestID, response.Error.RequestID)
	
	mockService.AssertExpectations(t)
	mockErrorHandler.AssertExpectations(t)
}

func TestIntegration_CategoryController_ContextPropagation(t *testing.T) {
	// Setup
	mockService := new(MockCategoryServiceIntegration)
	mockErrorHandler := new(MockErrorHandlerIntegration)
	mockAuthMiddleware := new(MockAuthMiddlewareIntegration)
	
	router := setupIntegrationTestRouter()
	rg := router.Group("/api")
	controller := NewCategoryController(mockService, rg, mockAuthMiddleware, mockErrorHandler)
	controller.Route()
	
	requestID := "category-test-123"
	
	// Mock service with context verification
	mockService.On("FindAll", mock.AnythingOfType("*context.timerCtx")).Return([]model.Category{
		{Id: 1, Name: "Technology"},
	}, nil)
	
	// Create request
	req, _ := http.NewRequest("GET", "/api/category/", nil)
	req.Header.Set("X-Request-ID", requestID)
	
	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Verify response contains context information
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, requestID, response.Meta.RequestID)
	
	// Verify response headers
	assert.Equal(t, requestID, w.Header().Get("X-Request-ID"))
	
	mockService.AssertExpectations(t)
}

func TestIntegration_UserController_PaginationWithContext(t *testing.T) {
	// Setup
	mockService := new(MockUserServiceIntegration)
	mockErrorHandler := new(MockErrorHandlerIntegration)
	
	router := setupIntegrationTestRouter()
	rg := router.Group("/api")
	controller := NewUserController(mockService, rg, mockErrorHandler)
	controller.Route()
	
	requestID := "pagination-test-123"
	
	// Mock service with pagination
	expectedResult := service.PaginationResult{
		Data: []model.User{
			{Id: 1, Name: "User 1", Email: "user1@example.com"},
		},
		Metadata: service.PaginationMetadata{
			Page:       1,
			Limit:      10,
			Total:      1,
			TotalPages: 1,
			HasNext:    false,
			HasPrev:    false,
		},
	}
	
	mockService.On("FindAllUserWithPagination", mock.MatchedBy(func(ctx context.Context) bool {
		ctxRequestID, _ := ctx.Value("request_id").(string)
		return ctxRequestID == requestID
	}), 1, 10).Return(expectedResult, nil)
	
	// Create request
	req, _ := http.NewRequest("GET", "/api/users/paginated?page=1&limit=10", nil)
	req.Header.Set("X-Request-ID", requestID)
	
	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response with pagination and context
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
	assert.NotNil(t, response.Data)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, requestID, response.Meta.RequestID)
	
	mockService.AssertExpectations(t)
}

func TestIntegration_UserController_ValidationErrorWithContext(t *testing.T) {
	// Setup
	mockService := new(MockUserServiceIntegration)
	mockErrorHandler := new(MockErrorHandlerIntegration)
	
	router := setupIntegrationTestRouter()
	rg := router.Group("/api")
	controller := NewUserController(mockService, rg, mockErrorHandler)
	controller.Route()
	
	requestID := "validation-test-123"
	
	// Mock expectations for validation error
	expectedError := &utils.AppError{
		Code:       utils.ErrValidation,
		Message:    "Validation failed",
		StatusCode: 400,
		RequestID:  requestID,
		Timestamp:  time.Now(),
	}
	
	mockErrorHandler.On("ValidationError", mock.AnythingOfType("*context.timerCtx"), "payload", mock.AnythingOfType("string")).Return(expectedError)
	mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), expectedError).Return()
	
	// Create request with invalid JSON
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)
	
	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert validation error response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, utils.ErrValidation, response.Error.Code)
	assert.Equal(t, requestID, response.Error.RequestID)
	
	mockErrorHandler.AssertExpectations(t)
}

func TestIntegration_MultipleRequestsContextIsolation(t *testing.T) {
	// Setup
	mockService := new(MockUserServiceIntegration)
	mockErrorHandler := new(MockErrorHandlerIntegration)
	
	router := setupIntegrationTestRouter()
	rg := router.Group("/api")
	controller := NewUserController(mockService, rg, mockErrorHandler)
	controller.Route()
	
	// Test that contexts are isolated between requests
	requestID1 := "request-1"
	requestID2 := "request-2"
	
	// Mock service for both requests
	mockService.On("FindAllUser", mock.MatchedBy(func(ctx context.Context) bool {
		ctxRequestID, _ := ctx.Value("request_id").(string)
		return ctxRequestID == requestID1
	})).Return([]model.User{{Id: 1, Name: "User 1"}}, nil).Once()
	
	mockService.On("FindAllUser", mock.MatchedBy(func(ctx context.Context) bool {
		ctxRequestID, _ := ctx.Value("request_id").(string)
		return ctxRequestID == requestID2
	})).Return([]model.User{{Id: 2, Name: "User 2"}}, nil).Once()
	
	// Create first request
	req1, _ := http.NewRequest("GET", "/api/users", nil)
	req1.Header.Set("X-Request-ID", requestID1)
	
	// Create second request
	req2, _ := http.NewRequest("GET", "/api/users", nil)
	req2.Header.Set("X-Request-ID", requestID2)
	
	// Execute both requests
	w1 := httptest.NewRecorder()
	w2 := httptest.NewRecorder()
	
	router.ServeHTTP(w1, req1)
	router.ServeHTTP(w2, req2)
	
	// Assert both responses have correct context
	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, http.StatusOK, w2.Code)
	
	var response1, response2 dto.APIResponse
	json.Unmarshal(w1.Body.Bytes(), &response1)
	json.Unmarshal(w2.Body.Bytes(), &response2)
	
	assert.Equal(t, requestID1, response1.Meta.RequestID)
	assert.Equal(t, requestID2, response2.Meta.RequestID)
	
	assert.Equal(t, requestID1, w1.Header().Get("X-Request-ID"))
	assert.Equal(t, requestID2, w2.Header().Get("X-Request-ID"))
	
	mockService.AssertExpectations(t)
}

func TestIntegration_ContextMetadataInResponses(t *testing.T) {
	// Setup
	mockService := new(MockUserServiceIntegration)
	mockErrorHandler := new(MockErrorHandlerIntegration)
	
	router := setupIntegrationTestRouter()
	rg := router.Group("/api")
	controller := NewUserController(mockService, rg, mockErrorHandler)
	controller.Route()
	
	requestID := "metadata-test-123"
	
	// Mock service
	mockService.On("FindAllUser", mock.AnythingOfType("*context.timerCtx")).Return([]model.User{}, nil)
	
	// Create request
	req, _ := http.NewRequest("GET", "/api/users", nil)
	req.Header.Set("X-Request-ID", requestID)
	
	// Record start time
	startTime := time.Now()
	
	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert response metadata
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	
	// Verify metadata
	assert.NotNil(t, response.Meta)
	assert.Equal(t, requestID, response.Meta.RequestID)
	assert.True(t, response.Meta.ProcessingTime > 0)
	assert.Equal(t, "1.0.0", response.Meta.Version)
	assert.True(t, response.Meta.Timestamp.After(startTime))
	
	mockService.AssertExpectations(t)
}

func TestIntegration_ErrorContextPropagation(t *testing.T) {
	// Setup
	mockService := new(MockUserServiceIntegration)
	mockErrorHandler := new(MockErrorHandlerIntegration)
	
	router := setupIntegrationTestRouter()
	rg := router.Group("/api")
	controller := NewUserController(mockService, rg, mockErrorHandler)
	controller.Route()
	
	requestID := "error-test-123"
	userID := "user-789"
	
	// Mock service to return an error
	serviceError := &utils.AppError{
		Code:       utils.ErrNotFound,
		Message:    "User not found",
		StatusCode: 404,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
	
	mockService.On("FindUserById", mock.AnythingOfType("*context.timerCtx"), "999").Return(model.User{}, serviceError)
	mockErrorHandler.On("HandleError", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*gin.Context"), serviceError).Return()
	
	// Create request
	req, _ := http.NewRequest("GET", "/api/users/999", nil)
	req.Header.Set("X-Request-ID", requestID)
	
	// Add user context
	ctx := context.WithValue(req.Context(), "user_id", userID)
	req = req.WithContext(ctx)
	
	// Execute
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// Assert error response with context
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response dto.APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.False(t, response.Success)
	assert.NotNil(t, response.Error)
	assert.Equal(t, utils.ErrNotFound, response.Error.Code)
	assert.Equal(t, requestID, response.Error.RequestID)
	assert.NotNil(t, response.Meta)
	assert.Equal(t, requestID, response.Meta.RequestID)
	
	mockService.AssertExpectations(t)
	mockErrorHandler.AssertExpectations(t)
}