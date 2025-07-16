package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/utils"
	"strings"
	"testing"
	"time"
)

func TestValidationService_ValidateUser(t *testing.T) {
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)

	tests := []struct {
		name        string
		ctx         context.Context
		user        model.User
		wantErr     bool
		expectField string
	}{
		{
			name: "valid user",
			ctx:  context.WithValue(context.Background(), "request_id", "req_123"),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "SecurePass123!",
				Role:     "user",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			ctx:  context.WithValue(context.Background(), "request_id", "req_124"),
			user: model.User{
				Name:     "",
				Email:    "john.doe@example.com",
				Password: "SecurePass123!",
				Role:     "user",
			},
			wantErr:     true,
			expectField: "name",
		},
		{
			name: "short name",
			ctx:  context.WithValue(context.Background(), "request_id", "req_125"),
			user: model.User{
				Name:     "J",
				Email:    "john.doe@example.com",
				Password: "SecurePass123!",
				Role:     "user",
			},
			wantErr:     true,
			expectField: "name",
		},
		{
			name: "invalid email",
			ctx:  context.WithValue(context.Background(), "request_id", "req_126"),
			user: model.User{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "SecurePass123!",
				Role:     "user",
			},
			wantErr:     true,
			expectField: "email",
		},
		{
			name: "weak password",
			ctx:  context.WithValue(context.Background(), "request_id", "req_127"),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "weak",
				Role:     "user",
			},
			wantErr:     true,
			expectField: "password",
		},
		{
			name: "invalid role",
			ctx:  context.WithValue(context.Background(), "request_id", "req_128"),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "SecurePass123!",
				Role:     "invalid_role",
			},
			wantErr:     true,
			expectField: "role",
		},
		{
			name: "context cancellation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately
				return ctx
			}(),
			user: model.User{
				Name:     "John Doe",
				Email:    "john.doe@example.com",
				Password: "SecurePass123!",
				Role:     "user",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validationService.ValidateUser(tt.ctx, tt.user)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateUser() expected error, got nil")
					return
				}
				
				// Check if it's an AppError
				appErr := err
				if appErr.Code != utils.ErrValidation && appErr.Code != utils.ErrCancelled {
					t.Errorf("ValidateUser() expected validation or cancellation error, got %s", appErr.Code)
				}
				
				// Check if specific field error is present
				if tt.expectField != "" && appErr.Details != nil {
					if _, exists := appErr.Details[tt.expectField]; !exists {
						t.Errorf("ValidateUser() expected field error for %s, but not found in details", tt.expectField)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateUser() expected no error, got %v", err)
				}
			}
		})
	}
}

func TestValidationService_ValidatePagination(t *testing.T) {
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)

	tests := []struct {
		name    string
		ctx     context.Context
		page    int
		limit   int
		wantErr bool
	}{
		{
			name:    "valid pagination",
			ctx:     context.WithValue(context.Background(), "request_id", "req_200"),
			page:    1,
			limit:   10,
			wantErr: false,
		},
		{
			name:    "invalid page",
			ctx:     context.WithValue(context.Background(), "request_id", "req_201"),
			page:    0,
			limit:   10,
			wantErr: true,
		},
		{
			name:    "invalid limit - too small",
			ctx:     context.WithValue(context.Background(), "request_id", "req_202"),
			page:    1,
			limit:   0,
			wantErr: true,
		},
		{
			name:    "invalid limit - too large",
			ctx:     context.WithValue(context.Background(), "request_id", "req_203"),
			page:    1,
			limit:   101,
			wantErr: true,
		},
		{
			name: "context timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
				defer cancel()
				time.Sleep(2 * time.Nanosecond) // Ensure timeout
				return ctx
			}(),
			page:    1,
			limit:   10,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validationService.ValidatePagination(tt.ctx, tt.page, tt.limit)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidatePagination() expected error, got nil")
					return
				}
				
				// Check if it's an AppError
				appErr := err
				if appErr.Code != utils.ErrValidation && appErr.Code != utils.ErrTimeout {
					t.Errorf("ValidatePagination() expected validation or timeout error, got %s", appErr.Code)
				}
			} else {
				if err != nil {
					t.Errorf("ValidatePagination() expected no error, got %v", err)
				}
			}
		})
	}
}

func TestValidationService_ValidateField(t *testing.T) {
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)

	tests := []struct {
		name    string
		ctx     context.Context
		field   string
		value   interface{}
		rules   string
		wantErr bool
	}{
		{
			name:    "valid field",
			ctx:     context.WithValue(context.Background(), "request_id", "req_300"),
			field:   "name",
			value:   "John Doe",
			rules:   "required",
			wantErr: false,
		},
		{
			name:    "required field empty",
			ctx:     context.WithValue(context.Background(), "request_id", "req_301"),
			field:   "name",
			value:   "",
			rules:   "required",
			wantErr: true,
		},
		{
			name: "context cancellation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			field:   "name",
			value:   "John Doe",
			rules:   "required",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldErr := validationService.ValidateField(tt.ctx, tt.field, tt.value, tt.rules)
			
			if tt.wantErr {
				if fieldErr == nil {
					t.Errorf("ValidateField() expected error, got nil")
					return
				}
				
				if fieldErr.Field != tt.field {
					t.Errorf("ValidateField() expected field %s, got %s", tt.field, fieldErr.Field)
				}
			} else {
				if fieldErr != nil {
					t.Errorf("ValidateField() expected no error, got %v", fieldErr)
				}
			}
		})
	}
}

func TestValidationService_ValidateStruct(t *testing.T) {
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)

	tests := []struct {
		name    string
		ctx     context.Context
		s       interface{}
		wantErr bool
	}{
		{
			name:    "valid struct",
			ctx:     context.WithValue(context.Background(), "request_id", "req_400"),
			s:       &model.User{Name: "John"},
			wantErr: false,
		},
		{
			name:    "nil struct",
			ctx:     context.WithValue(context.Background(), "request_id", "req_401"),
			s:       nil,
			wantErr: true,
		},
		{
			name: "context cancellation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			s:       &model.User{Name: "John"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fieldErrors := validationService.ValidateStruct(tt.ctx, tt.s)
			
			if tt.wantErr {
				if len(fieldErrors) == 0 {
					t.Errorf("ValidateStruct() expected errors, got none")
				}
			} else {
				if len(fieldErrors) > 0 {
					t.Errorf("ValidateStruct() expected no errors, got %v", fieldErrors)
				}
			}
		})
	}
}

// Test helper functions
func TestValidationService_validateEmail(t *testing.T) {
	errorWrapper := utils.NewErrorWrapper()
	vs := &validationService{errorWrapper: errorWrapper}

	tests := []struct {
		name    string
		ctx     context.Context
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			ctx:     context.Background(),
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "invalid email format",
			ctx:     context.Background(),
			email:   "invalid-email",
			wantErr: true,
		},
		{
			name:    "email too long",
			ctx:     context.Background(),
			email:   "very-long-email-address-that-exceeds-the-maximum-allowed-length-for-email-addresses-which-is-254-characters-according-to-rfc-standards-and-this-email-is-definitely-longer-than-that-limit-so-it-should-fail-validation-and-this-part-makes-it-even-longer-to-ensure-it-exceeds-254-characters@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vs.validateEmail(tt.ctx, tt.email)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateEmail() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("validateEmail() expected no error, got %v", err)
				}
			}
		})
	}
}

func TestValidationService_validatePassword(t *testing.T) {
	errorWrapper := utils.NewErrorWrapper()
	vs := &validationService{errorWrapper: errorWrapper}

	tests := []struct {
		name     string
		ctx      context.Context
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			ctx:      context.Background(),
			password: "SecurePass123!",
			wantErr:  false,
		},
		{
			name:     "password too short",
			ctx:      context.Background(),
			password: "short",
			wantErr:  true,
		},
		{
			name:     "password missing uppercase",
			ctx:      context.Background(),
			password: "securepass123!",
			wantErr:  true,
		},
		{
			name:     "password missing lowercase",
			ctx:      context.Background(),
			password: "SECUREPASS123!",
			wantErr:  true,
		},
		{
			name:     "password missing number",
			ctx:      context.Background(),
			password: "SecurePass!",
			wantErr:  true,
		},
		{
			name:     "password missing special character",
			ctx:      context.Background(),
			password: "SecurePass123",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vs.validatePassword(tt.ctx, tt.password)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("validatePassword() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("validatePassword() expected no error, got %v", err)
				}
			}
		})
	}
}

func TestValidationService_ValidateArticle(t *testing.T) {
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)

	tests := []struct {
		name        string
		ctx         context.Context
		article     model.Article
		wantErr     bool
		expectField string
	}{
		{
			name: "valid article",
			ctx:  context.WithValue(context.Background(), "request_id", "req_500"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			wantErr: false,
		},
		{
			name: "empty title",
			ctx:  context.WithValue(context.Background(), "request_id", "req_501"),
			article: model.Article{
				Title:   "",
				Content: "This is a valid article content with sufficient length.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			wantErr:     true,
			expectField: "title",
		},
		{
			name: "short title",
			ctx:  context.WithValue(context.Background(), "request_id", "req_502"),
			article: model.Article{
				Title:   "Hi",
				Content: "This is a valid article content with sufficient length.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			wantErr:     true,
			expectField: "title",
		},
		{
			name: "empty content",
			ctx:  context.WithValue(context.Background(), "request_id", "req_503"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			wantErr:     true,
			expectField: "content",
		},
		{
			name: "short content",
			ctx:  context.WithValue(context.Background(), "request_id", "req_504"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "Short",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			wantErr:     true,
			expectField: "content",
		},
		{
			name: "invalid slug",
			ctx:  context.WithValue(context.Background(), "request_id", "req_505"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length.",
				Slug:    "Invalid Slug With Spaces",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			wantErr:     true,
			expectField: "slug",
		},
		{
			name: "invalid user ID",
			ctx:  context.WithValue(context.Background(), "request_id", "req_506"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 0},
				Category: model.Category{Id: 1},
			},
			wantErr:     true,
			expectField: "user_id",
		},
		{
			name: "invalid category ID",
			ctx:  context.WithValue(context.Background(), "request_id", "req_507"),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 0},
			},
			wantErr:     true,
			expectField: "category_id",
		},
		{
			name: "context cancellation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately
				return ctx
			}(),
			article: model.Article{
				Title:   "Valid Article Title",
				Content: "This is a valid article content with sufficient length.",
				Slug:    "valid-article-title",
				User:    model.User{Id: 1},
				Category: model.Category{Id: 1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validationService.ValidateArticle(tt.ctx, tt.article)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateArticle() expected error, got nil")
					return
				}
				
				// Check if it's an AppError
				appErr := err
				if appErr.Code != utils.ErrValidation && appErr.Code != utils.ErrCancelled {
					t.Errorf("ValidateArticle() expected validation or cancellation error, got %s", appErr.Code)
				}
				
				// Check if specific field error is present
				if tt.expectField != "" && appErr.Details != nil {
					if _, exists := appErr.Details[tt.expectField]; !exists {
						t.Errorf("ValidateArticle() expected field error for %s, but not found in details", tt.expectField)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateArticle() expected no error, got %v", err)
				}
			}
		})
	}
}

func TestValidationService_validateSlug(t *testing.T) {
	errorWrapper := utils.NewErrorWrapper()
	vs := &validationService{errorWrapper: errorWrapper}

	tests := []struct {
		name    string
		ctx     context.Context
		slug    string
		wantErr bool
	}{
		{
			name:    "valid slug",
			ctx:     context.Background(),
			slug:    "valid-article-slug",
			wantErr: false,
		},
		{
			name:    "valid slug with numbers",
			ctx:     context.Background(),
			slug:    "article-123-title",
			wantErr: false,
		},
		{
			name:    "slug too short",
			ctx:     context.Background(),
			slug:    "ab",
			wantErr: true,
		},
		{
			name:    "slug with uppercase",
			ctx:     context.Background(),
			slug:    "Article-Title",
			wantErr: true,
		},
		{
			name:    "slug with spaces",
			ctx:     context.Background(),
			slug:    "article title",
			wantErr: true,
		},
		{
			name:    "slug with consecutive hyphens",
			ctx:     context.Background(),
			slug:    "article--title",
			wantErr: true,
		},
		{
			name:    "slug starting with hyphen",
			ctx:     context.Background(),
			slug:    "-article-title",
			wantErr: true,
		},
		{
			name:    "slug ending with hyphen",
			ctx:     context.Background(),
			slug:    "article-title-",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vs.validateSlug(tt.ctx, tt.slug)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateSlug() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("validateSlug() expected no error, got %v", err)
				}
			}
		})
	}
}

func TestValidationService_ValidateComment(t *testing.T) {
	errorWrapper := utils.NewErrorWrapper()
	validationService := NewValidationService(errorWrapper)

	tests := []struct {
		name        string
		ctx         context.Context
		comment     model.Comment
		wantErr     bool
		expectField string
	}{
		{
			name: "valid comment",
			ctx:  context.WithValue(context.Background(), "request_id", "req_600"),
			comment: model.Comment{
				Content: "This is a valid comment.",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			wantErr: false,
		},
		{
			name: "empty content",
			ctx:  context.WithValue(context.Background(), "request_id", "req_601"),
			comment: model.Comment{
				Content: "",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			wantErr:     true,
			expectField: "content",
		},
		{
			name: "content too long",
			ctx:  context.WithValue(context.Background(), "request_id", "req_602"),
			comment: model.Comment{
				Content: strings.Repeat("a", 1001), // 1001 characters
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			wantErr:     true,
			expectField: "content",
		},
		{
			name: "invalid user ID",
			ctx:  context.WithValue(context.Background(), "request_id", "req_603"),
			comment: model.Comment{
				Content: "This is a valid comment.",
				User:    model.User{Id: 0},
				Article: model.Article{Id: 1},
			},
			wantErr:     true,
			expectField: "user_id",
		},
		{
			name: "invalid article ID",
			ctx:  context.WithValue(context.Background(), "request_id", "req_604"),
			comment: model.Comment{
				Content: "This is a valid comment.",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 0},
			},
			wantErr:     true,
			expectField: "article_id",
		},
		{
			name: "content with HTML tags",
			ctx:  context.WithValue(context.Background(), "request_id", "req_605"),
			comment: model.Comment{
				Content: "This is a comment with <script>alert('xss')</script> tags.",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			wantErr:     true,
			expectField: "content",
		},
		{
			name: "content with URLs",
			ctx:  context.WithValue(context.Background(), "request_id", "req_606"),
			comment: model.Comment{
				Content: "Check out this link: https://example.com",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			wantErr:     true,
			expectField: "content",
		},
		{
			name: "content with excessive repetition",
			ctx:  context.WithValue(context.Background(), "request_id", "req_607"),
			comment: model.Comment{
				Content: "This is aaaaaaaaaaaaaaa comment.",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			wantErr:     true,
			expectField: "content",
		},
		{
			name: "context cancellation",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately
				return ctx
			}(),
			comment: model.Comment{
				Content: "This is a valid comment.",
				User:    model.User{Id: 1},
				Article: model.Article{Id: 1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validationService.ValidateComment(tt.ctx, tt.comment)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateComment() expected error, got nil")
					return
				}
				
				// Check if it's an AppError
				appErr := err
				if appErr.Code != utils.ErrValidation && appErr.Code != utils.ErrCancelled {
					t.Errorf("ValidateComment() expected validation or cancellation error, got %s", appErr.Code)
				}
				
				// Check if specific field error is present
				if tt.expectField != "" && appErr.Details != nil {
					if _, exists := appErr.Details[tt.expectField]; !exists {
						t.Errorf("ValidateComment() expected field error for %s, but not found in details", tt.expectField)
					}
				}
			} else {
				if err != nil {
					t.Errorf("ValidateComment() expected no error, got %v", err)
				}
			}
		})
	}
}

func TestValidationService_validateCommentContent(t *testing.T) {
	errorWrapper := utils.NewErrorWrapper()
	vs := &validationService{errorWrapper: errorWrapper}

	tests := []struct {
		name    string
		ctx     context.Context
		content string
		wantErr bool
	}{
		{
			name:    "valid content",
			ctx:     context.Background(),
			content: "This is a valid comment.",
			wantErr: false,
		},
		{
			name:    "content with HTML tags",
			ctx:     context.Background(),
			content: "This has <b>bold</b> text.",
			wantErr: true,
		},
		{
			name:    "content with script tags",
			ctx:     context.Background(),
			content: "Malicious <script>alert('xss')</script> content.",
			wantErr: true,
		},
		{
			name:    "content with URLs",
			ctx:     context.Background(),
			content: "Check this out: https://example.com",
			wantErr: true,
		},
		{
			name:    "content with excessive repetition",
			ctx:     context.Background(),
			content: "This is sooooooooooooo repetitive.",
			wantErr: true,
		},
		{
			name:    "content with too much whitespace",
			ctx:     context.Background(),
			content: "This   has   too   much   whitespace   everywhere   in   the   text.",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := vs.validateCommentContent(tt.ctx, tt.content)
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateCommentContent() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("validateCommentContent() expected no error, got %v", err)
				}
			}
		})
	}
}