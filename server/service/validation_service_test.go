package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/utils"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock ErrorWrapper
type MockErrorWrapper struct {
	mock.Mock
}

func (m *MockErrorWrapper) WrapError(ctx context.Context, err error, code string, message string) *utils.AppError {
	args := m.Called(ctx, err, code, message)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) ValidationError(ctx context.Context, field string, message string) *utils.AppError {
	args := m.Called(ctx, field, message)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) NotFoundError(ctx context.Context, resource string) *utils.AppError {
	args := m.Called(ctx, resource)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) UnauthorizedError(ctx context.Context, message string) *utils.AppError {
	args := m.Called(ctx, message)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) ForbiddenError(ctx context.Context, message string) *utils.AppError {
	args := m.Called(ctx, message)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) InternalError(ctx context.Context, err error, message string) *utils.AppError {
	args := m.Called(ctx, err, message)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) DatabaseError(ctx context.Context, err error, operation string) *utils.AppError {
	args := m.Called(ctx, err, operation)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) RateLimitError(ctx context.Context, limit int, window time.Duration) *utils.AppError {
	args := m.Called(ctx, limit, window)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) TimeoutError(ctx context.Context, operation string) *utils.AppError {
	args := m.Called(ctx, operation)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) CancellationError(ctx context.Context, operation string) *utils.AppError {
	args := m.Called(ctx, operation)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) ConflictError(ctx context.Context, resource string, message string) *utils.AppError {
	args := m.Called(ctx, resource, message)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

func (m *MockErrorWrapper) BadRequestError(ctx context.Context, message string) *utils.AppError {
	args := m.Called(ctx, message)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*utils.AppError)
}

// Context-aware test for ValidateUser
func TestValidateUser_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		user          model.User
		setupMocks    func(*MockErrorWrapper)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name: "successful user validation with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_123"),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "SecurePass123!",
				Role:     "user",
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for successful validation
			},
		},
		{
			name: "context cancellation during validation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "SecurePass123!",
				Role:     "user",
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				mockEW.On("CancellationError", mock.AnythingOfType("*context.cancelCtx"), "validation").Return(&utils.AppError{
					Code:    utils.ErrCancelled,
					Message: "Validation cancelled",
				}).Once()
			},
			expectCancel: true,
		},
		{
			name: "context timeout during validation",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "SecurePass123!",
				Role:     "user",
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "validation").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Validation timeout",
				}).Once()
			},
			expectTimeout: true,
		},
		{
			name: "validation error with context - empty name",
			ctx:  context.WithValue(context.Background(), "request_id", "req_empty_name"),
			user: model.User{
				Name:     "",
				Email:    "john.doe@example.com",
				Password: "SecurePass123!",
				Role:     "user",
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "User validation failed",
		},
		{
			name: "validation error with context - invalid email",
			ctx:  context.WithValue(context.Background(), "request_id", "req_invalid_email"),
			user: model.User{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "SecurePass123!",
				Role:     "user",
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "User validation failed",
		},
		{
			name: "validation error with context - weak password",
			ctx:  context.WithValue(context.Background(), "request_id", "req_weak_password"),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "weak",
				Role:     "user",
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "User validation failed",
		},
		{
			name: "validation error with context - invalid role",
			ctx:  context.WithValue(context.Background(), "request_id", "req_invalid_role"),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "SecurePass123!",
				Role:     "invalid_role",
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "User validation failed",
		},
		{
			name: "multiple validation errors with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_multiple_errors"),
			user: model.User{
				Name:     "",
				Email:    "invalid-email",
				Password: "weak",
				Role:     "invalid_role",
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "User validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockErrorWrapper := new(MockErrorWrapper)
			validationService := NewValidationService(mockErrorWrapper)

			// Setup mock expectations
			tt.setupMocks(mockErrorWrapper)

			// Execute test
			err := validationService.ValidateUser(tt.ctx, tt.user)

			// Verify results
			if tt.expectCancel {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrCancelled, err.Code)
			} else if tt.expectTimeout {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Message, tt.expectedError)
				assert.Equal(t, utils.ErrValidation, err.Code)
			} else {
				assert.Nil(t, err)
			}

			// Verify mock expectations
			mockErrorWrapper.AssertExpectations(t)
		})
	}
}

// Context-aware test for ValidateArticle
func TestValidateArticle_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		article       model.Article
		setupMocks    func(*MockErrorWrapper)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name: "successful article validation with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_article_123"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length to pass validation.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for successful validation
			},
		},
		{
			name: "context cancellation during article validation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length to pass validation.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				mockEW.On("CancellationError", mock.AnythingOfType("*context.cancelCtx"), "validation").Return(&utils.AppError{
					Code:    utils.ErrCancelled,
					Message: "Article validation cancelled",
				}).Once()
			},
			expectCancel: true,
		},
		{
			name: "context timeout during article validation",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length to pass validation.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "validation").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Article validation timeout",
				}).Once()
			},
			expectTimeout: true,
		},
		{
			name: "validation error with context - empty title",
			ctx:  context.WithValue(context.Background(), "request_id", "req_empty_title"),
			article: model.Article{
				Title:   "",
				Content: "This is a valid article content with sufficient length to pass validation.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Article validation failed",
		},
		{
			name: "validation error with context - short title",
			ctx:  context.WithValue(context.Background(), "request_id", "req_short_title"),
			article: model.Article{
				Title:   "Hi",
				Content: "This is a valid article content with sufficient length to pass validation.",
				Slug:    "hi",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Article validation failed",
		},
		{
			name: "validation error with context - empty content",
			ctx:  context.WithValue(context.Background(), "request_id", "req_empty_content"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Article validation failed",
		},
		{
			name: "validation error with context - invalid slug",
			ctx:  context.WithValue(context.Background(), "request_id", "req_invalid_slug"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length to pass validation.",
				Slug:    "Invalid Slug With Spaces",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Article validation failed",
		},
		{
			name: "validation error with context - invalid user ID",
			ctx:  context.WithValue(context.Background(), "request_id", "req_invalid_user"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length to pass validation.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 0},
				Category: model.Category{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Article validation failed",
		},
		{
			name: "validation error with context - invalid category ID",
			ctx:  context.WithValue(context.Background(), "request_id", "req_invalid_category"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length to pass validation.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 0},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Article validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockErrorWrapper := new(MockErrorWrapper)
			validationService := NewValidationService(mockErrorWrapper)

			// Setup mock expectations
			tt.setupMocks(mockErrorWrapper)

			// Execute test
			err := validationService.ValidateArticle(tt.ctx, tt.article)

			// Verify results
			if tt.expectCancel {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrCancelled, err.Code)
			} else if tt.expectTimeout {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Message, tt.expectedError)
				assert.Equal(t, utils.ErrValidation, err.Code)
			} else {
				assert.Nil(t, err)
			}

			// Verify mock expectations
			mockErrorWrapper.AssertExpectations(t)
		})
	}
}

// Context-aware test for ValidateComment
func TestValidateComment_WithContext(t *testing.T) {
	tests := []struct {
		name          string
		ctx           context.Context
		comment       model.Comment
		setupMocks    func(*MockErrorWrapper)
		expectedError string
		expectCancel  bool
		expectTimeout bool
	}{
		{
			name: "successful comment validation with context",
			ctx:  context.WithValue(context.Background(), "request_id", "req_comment_123"),
			comment: model.Comment{
				Content: "This is a valid comment content.",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for successful validation
			},
		},
		{
			name: "context cancellation during comment validation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			comment: model.Comment{
				Content: "This is a valid comment content.",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				mockEW.On("CancellationError", mock.AnythingOfType("*context.cancelCtx"), "validation").Return(&utils.AppError{
					Code:    utils.ErrCancelled,
					Message: "Comment validation cancelled",
				}).Once()
			},
			expectCancel: true,
		},
		{
			name: "context timeout during comment validation",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				time.Sleep(2 * time.Nanosecond)
				return ctx
			}(),
			comment: model.Comment{
				Content: "This is a valid comment content.",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				mockEW.On("TimeoutError", mock.AnythingOfType("*context.timerCtx"), "validation").Return(&utils.AppError{
					Code:    utils.ErrTimeout,
					Message: "Comment validation timeout",
				}).Once()
			},
			expectTimeout: true,
		},
		{
			name: "validation error with context - empty content",
			ctx:  context.WithValue(context.Background(), "request_id", "req_empty_content"),
			comment: model.Comment{
				Content: "",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Comment validation failed",
		},
		{
			name: "validation error with context - content with HTML tags",
			ctx:  context.WithValue(context.Background(), "request_id", "req_html_content"),
			comment: model.Comment{
				Content: "This comment has <script>alert('xss')</script> HTML tags.",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Comment validation failed",
		},
		{
			name: "validation error with context - content with URLs",
			ctx:  context.WithValue(context.Background(), "request_id", "req_url_content"),
			comment: model.Comment{
				Content: "Check out this link: https://malicious-site.com",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Comment validation failed",
		},
		{
			name: "validation error with context - invalid user ID",
			ctx:  context.WithValue(context.Background(), "request_id", "req_invalid_user"),
			comment: model.Comment{
				Content: "This is a valid comment content.",
				User:    model.User{Id: 0},
				Article: model.Article{Id: 1},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Comment validation failed",
		},
		{
			name: "validation error with context - invalid article ID",
			ctx:  context.WithValue(context.Background(), "request_id", "req_invalid_article"),
			comment: model.Comment{
				Content: "This is a valid comment content.",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 0},
			},
			setupMocks: func(mockEW *MockErrorWrapper) {
				// No mocks needed for validation error
			},
			expectedError: "Comment validation failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockErrorWrapper := new(MockErrorWrapper)
			validationService := NewValidationService(mockErrorWrapper)

			// Setup mock expectations
			tt.setupMocks(mockErrorWrapper)

			// Execute test
			err := validationService.ValidateComment(tt.ctx, tt.comment)

			// Verify results
			if tt.expectCancel {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrCancelled, err.Code)
			} else if tt.expectTimeout {
				assert.NotNil(t, err)
				assert.Equal(t, utils.ErrTimeout, err.Code)
			} else if tt.expectedError != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Message, tt.expectedError)
				assert.Equal(t, utils.ErrValidation, err.Code)
			} else {
				assert.Nil(t, err)
			}

			// Verify mock expectations
			mockErrorWrapper.AssertExpectations(t)
		})
	}
}

