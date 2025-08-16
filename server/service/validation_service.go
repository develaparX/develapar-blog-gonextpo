package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/utils"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// FieldError represents a validation error for a specific field with context information
type FieldError struct {
	Field     string `json:"field"`
	Message   string `json:"message"`
	Value     string `json:"value,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// ValidationResult holds the result of validation with context
type ValidationResult struct {
	IsValid   bool         `json:"is_valid"`
	Errors    []FieldError `json:"errors,omitempty"`
	RequestID string       `json:"request_id,omitempty"`
}

// ValidationService interface defines validation methods with context support
type ValidationService interface {
	ValidateUser(ctx context.Context, user model.User) *utils.AppError
	ValidateArticle(ctx context.Context, article model.Article) *utils.AppError
	ValidateComment(ctx context.Context, comment model.Comment) *utils.AppError
	ValidatePagination(ctx context.Context, page, limit int) *utils.AppError
	ValidateField(ctx context.Context, field string, value interface{}, rules string) *FieldError
	ValidateStruct(ctx context.Context, s interface{}) []FieldError
}

// validationService implements ValidationService interface
type validationService struct {
	errorWrapper utils.ErrorWrapper
}

// NewValidationService creates a new validation service instance
func NewValidationService(errorWrapper utils.ErrorWrapper) ValidationService {
	return &validationService{
		errorWrapper: errorWrapper,
	}
}

// extractRequestID extracts request ID from context
func (vs *validationService) extractRequestID(ctx context.Context) string {
	if ctx != nil {
		if rid, ok := ctx.Value("request_id").(string); ok {
			return rid
		}
	}
	return ""
}

// checkContextTimeout checks if context has timed out or been cancelled
func (vs *validationService) checkContextTimeout(ctx context.Context) error {
	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			return vs.errorWrapper.TimeoutError(ctx, "validation")
		}
		if ctx.Err() == context.Canceled {
			return vs.errorWrapper.CancellationError(ctx, "validation")
		}
		return ctx.Err()
	default:
		return nil
	}
}

// ValidateUser validates user data with context support
func (vs *validationService) ValidateUser(ctx context.Context, user model.User) *utils.AppError {
	// Check for context timeout/cancellation
	if err := vs.checkContextTimeout(ctx); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return appErr
		}
		return vs.errorWrapper.InternalError(ctx, err, "Context error during user validation")
	}

	var fieldErrors []FieldError
	requestID := vs.extractRequestID(ctx)

	// Validate name
	if strings.TrimSpace(user.Name) == "" {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "name",
			Message:   "Name is required",
			Value:     user.Name,
			RequestID: requestID,
		})
	} else if len(strings.TrimSpace(user.Name)) < 2 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "name",
			Message:   "Name must be at least 2 characters long",
			Value:     user.Name,
			RequestID: requestID,
		})
	} else if len(strings.TrimSpace(user.Name)) > 100 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "name",
			Message:   "Name must not exceed 100 characters",
			Value:     user.Name,
			RequestID: requestID,
		})
	}

	// Check context timeout after each validation step
	if err := vs.checkContextTimeout(ctx); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return appErr
		}
		return vs.errorWrapper.InternalError(ctx, err, "Context error during user validation")
	}

	// Validate email
	if strings.TrimSpace(user.Email) == "" {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "email",
			Message:   "Email is required",
			Value:     user.Email,
			RequestID: requestID,
		})
	} else {
		if err := vs.validateEmail(ctx, user.Email); err != nil {
			fieldErrors = append(fieldErrors, FieldError{
				Field:     "email",
				Message:   err.Error(),
				Value:     user.Email,
				RequestID: requestID,
			})
		}
	}

	// Check context timeout
	if err := vs.checkContextTimeout(ctx); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return appErr
		}
		return vs.errorWrapper.InternalError(ctx, err, "Context error during user validation")
	}

	// Validate password
	if strings.TrimSpace(user.Password) == "" {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "password",
			Message:   "Password is required",
			RequestID: requestID,
		})
	} else {
		if err := vs.validatePassword(ctx, user.Password); err != nil {
			fieldErrors = append(fieldErrors, FieldError{
				Field:     "password",
				Message:   err.Error(),
				RequestID: requestID,
			})
		}
	}

	// Validate role
	if strings.TrimSpace(user.Role) == "" {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "role",
			Message:   "Role is required",
			Value:     user.Role,
			RequestID: requestID,
		})
	} else {
		validRoles := []string{"admin", "user", "moderator"}
		isValidRole := false
		for _, role := range validRoles {
			if strings.ToLower(user.Role) == role {
				isValidRole = true
				break
			}
		}
		if !isValidRole {
			fieldErrors = append(fieldErrors, FieldError{
				Field:     "role",
				Message:   "Role must be one of: admin, user, moderator",
				Value:     user.Role,
				RequestID: requestID,
			})
		}
	}

	// Return validation error if there are field errors
	if len(fieldErrors) > 0 {
		details := make(map[string]string)
		for _, fieldErr := range fieldErrors {
			details[fieldErr.Field] = fieldErr.Message
		}

		return &utils.AppError{
			Code:       utils.ErrValidation,
			Message:    "User validation failed",
			Details:    details,
			StatusCode: 400,
			RequestID:  requestID,
			Timestamp:  time.Now(),
		}
	}

	return nil
}

// validateEmail validates email format with context timeout handling
func (vs *validationService) validateEmail(ctx context.Context, email string) error {
	// Check context timeout
	if err := vs.checkContextTimeout(ctx); err != nil {
		return err
	}

	// Basic email format validation
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email format")
	}

	// Additional email validation rules
	if len(email) > 254 {
		return fmt.Errorf("email address too long")
	}

	// Check for common email patterns
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// validatePassword validates password strength with context support
func (vs *validationService) validatePassword(ctx context.Context, password string) error {
	// Check context timeout
	if err := vs.checkContextTimeout(ctx); err != nil {
		return err
	}

	// Password length validation
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	if len(password) > 128 {
		return fmt.Errorf("password must not exceed 128 characters")
	}

	// Password strength validation
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return fmt.Errorf("password must contain at least one number")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// ValidateField validates a single field with context support
func (vs *validationService) ValidateField(ctx context.Context, field string, value interface{}, rules string) *FieldError {
	// Check context timeout
	if err := vs.checkContextTimeout(ctx); err != nil {
		return &FieldError{
			Field:     field,
			Message:   "Validation timeout",
			RequestID: vs.extractRequestID(ctx),
		}
	}

	// Basic field validation logic
	// This is a simplified implementation - in a real application,
	// you might want to use a more sophisticated validation library

	requestID := vs.extractRequestID(ctx)

	// Convert value to string for basic validation
	strValue := fmt.Sprintf("%v", value)

	// Parse rules (simplified - could be more complex)
	if strings.Contains(rules, "required") && strings.TrimSpace(strValue) == "" {
		return &FieldError{
			Field:     field,
			Message:   fmt.Sprintf("%s is required", field),
			Value:     strValue,
			RequestID: requestID,
		}
	}

	return nil
}

// ValidateStruct validates a struct with context support
func (vs *validationService) ValidateStruct(ctx context.Context, s interface{}) []FieldError {
	// Check context timeout
	if err := vs.checkContextTimeout(ctx); err != nil {
		return []FieldError{{
			Field:     "struct",
			Message:   "Validation timeout",
			RequestID: vs.extractRequestID(ctx),
		}}
	}

	// This is a placeholder implementation
	// In a real application, you would use reflection to validate struct fields
	// based on struct tags or other validation rules

	var fieldErrors []FieldError
	requestID := vs.extractRequestID(ctx)

	// Basic struct validation
	if s == nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "struct",
			Message:   "Struct cannot be nil",
			RequestID: requestID,
		})
	}

	return fieldErrors
}

// ValidatePagination validates pagination parameters with context support
func (vs *validationService) ValidatePagination(ctx context.Context, page, limit int) *utils.AppError {
	// Check context timeout
	if err := vs.checkContextTimeout(ctx); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return appErr
		}
		return vs.errorWrapper.InternalError(ctx, err, "Context error during pagination validation")
	}

	var fieldErrors []FieldError
	requestID := vs.extractRequestID(ctx)

	// Validate page
	if page < 1 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "page",
			Message:   "Page must be greater than 0",
			Value:     fmt.Sprintf("%d", page),
			RequestID: requestID,
		})
	}

	// Validate limit
	if limit < 1 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "limit",
			Message:   "Limit must be greater than 0",
			Value:     fmt.Sprintf("%d", limit),
			RequestID: requestID,
		})
	} else if limit > 100 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "limit",
			Message:   "Limit must not exceed 100",
			Value:     fmt.Sprintf("%d", limit),
			RequestID: requestID,
		})
	}

	// Return validation error if there are field errors
	if len(fieldErrors) > 0 {
		details := make(map[string]string)
		for _, fieldErr := range fieldErrors {
			details[fieldErr.Field] = fieldErr.Message
		}

		return &utils.AppError{
			Code:       utils.ErrValidation,
			Message:    "Pagination validation failed",
			Details:    details,
			StatusCode: 400,
			RequestID:  requestID,
			Timestamp:  time.Now(),
		}
	}

	return nil
}

// Placeholder implementations for ValidateArticle and ValidateComment
// These will be implemented in the next sub-tasks

// ValidateArticle validates article data with context support
func (vs *validationService) ValidateArticle(ctx context.Context, article model.Article) *utils.AppError {
	// Check for context timeout/cancellation
	if err := vs.checkContextTimeout(ctx); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return appErr
		}
		return vs.errorWrapper.InternalError(ctx, err, "Context error during article validation")
	}

	var fieldErrors []FieldError
	requestID := vs.extractRequestID(ctx)

	// Validate title
	if strings.TrimSpace(article.Title) == "" {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "title",
			Message:   "Title is required",
			Value:     article.Title,
			RequestID: requestID,
		})
	} else if len(strings.TrimSpace(article.Title)) < 3 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "title",
			Message:   "Title must be at least 3 characters long",
			Value:     article.Title,
			RequestID: requestID,
		})
	} else if len(strings.TrimSpace(article.Title)) > 200 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "title",
			Message:   "Title must not exceed 200 characters",
			Value:     article.Title,
			RequestID: requestID,
		})
	}

	// Check context timeout after each validation step
	if err := vs.checkContextTimeout(ctx); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return appErr
		}
		return vs.errorWrapper.InternalError(ctx, err, "Context error during article validation")
	}

	// Validate content
	if strings.TrimSpace(article.Content) == "" {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "content",
			Message:   "Content is required",
			Value:     article.Content,
			RequestID: requestID,
		})
	} else if len(strings.TrimSpace(article.Content)) < 10 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "content",
			Message:   "Content must be at least 10 characters long",
			Value:     article.Content,
			RequestID: requestID,
		})
	} else if len(strings.TrimSpace(article.Content)) > 50000 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "content",
			Message:   "Content must not exceed 50,000 characters",
			Value:     article.Content,
			RequestID: requestID,
		})
	}

	// Check context timeout
	if err := vs.checkContextTimeout(ctx); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return appErr
		}
		return vs.errorWrapper.InternalError(ctx, err, "Context error during article validation")
	}

	// Validate slug
	if strings.TrimSpace(article.Slug) == "" {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "slug",
			Message:   "Slug is required",
			Value:     article.Slug,
			RequestID: requestID,
		})
	} else {
		if err := vs.validateSlug(ctx, article.Slug); err != nil {
			fieldErrors = append(fieldErrors, FieldError{
				Field:     "slug",
				Message:   err.Error(),
				Value:     article.Slug,
				RequestID: requestID,
			})
		}
	}

	// Validate user (must have valid user ID)
	if article.User.Id == uuid.Nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "user_id",
			Message:   "Valid user ID is required",
			Value:     fmt.Sprintf("%d", article.User.Id),
			RequestID: requestID,
		})
	}

	// Validate category (must have valid category ID)
	if article.Category.Id == uuid.Nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "category_id",
			Message:   "Valid category ID is required",
			Value:     fmt.Sprintf("%d", article.Category.Id),
			RequestID: requestID,
		})
	}

	// Return validation error if there are field errors
	if len(fieldErrors) > 0 {
		details := make(map[string]string)
		for _, fieldErr := range fieldErrors {
			details[fieldErr.Field] = fieldErr.Message
		}

		return &utils.AppError{
			Code:       utils.ErrValidation,
			Message:    "Article validation failed",
			Details:    details,
			StatusCode: 400,
			RequestID:  requestID,
			Timestamp:  time.Now(),
		}
	}

	return nil
}

// validateSlug validates article slug with context support
func (vs *validationService) validateSlug(ctx context.Context, slug string) error {
	// Check context timeout
	if err := vs.checkContextTimeout(ctx); err != nil {
		return err
	}

	// Slug length validation
	if len(slug) < 3 {
		return fmt.Errorf("slug must be at least 3 characters long")
	}

	if len(slug) > 100 {
		return fmt.Errorf("slug must not exceed 100 characters")
	}

	// Slug format validation (URL-friendly)
	slugRegex := regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
	if !slugRegex.MatchString(slug) {
		return fmt.Errorf("slug must contain only lowercase letters, numbers, and hyphens, and cannot start or end with a hyphen")
	}

	// Check for consecutive hyphens
	if strings.Contains(slug, "--") {
		return fmt.Errorf("slug cannot contain consecutive hyphens")
	}

	return nil
}

// ValidateComment validates comment data with context support
func (vs *validationService) ValidateComment(ctx context.Context, comment model.Comment) *utils.AppError {
	// Check for context timeout/cancellation
	if err := vs.checkContextTimeout(ctx); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return appErr
		}
		return vs.errorWrapper.InternalError(ctx, err, "Context error during comment validation")
	}

	var fieldErrors []FieldError
	requestID := vs.extractRequestID(ctx)

	// Validate content
	if strings.TrimSpace(comment.Content) == "" {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "content",
			Message:   "Content is required",
			Value:     comment.Content,
			RequestID: requestID,
		})
	} else if len(strings.TrimSpace(comment.Content)) < 1 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "content",
			Message:   "Content must be at least 1 character long",
			Value:     comment.Content,
			RequestID: requestID,
		})
	} else if len(strings.TrimSpace(comment.Content)) > 1000 {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "content",
			Message:   "Content must not exceed 1,000 characters",
			Value:     comment.Content,
			RequestID: requestID,
		})
	}

	// Check context timeout after content validation
	if err := vs.checkContextTimeout(ctx); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return appErr
		}
		return vs.errorWrapper.InternalError(ctx, err, "Context error during comment validation")
	}

	// Validate user reference (must have valid user ID)
	if comment.User.Id == uuid.Nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "user_id",
			Message:   "Valid user ID is required",
			Value:     fmt.Sprintf("%d", comment.User.Id),
			RequestID: requestID,
		})
	}

	// Validate article reference (must have valid article ID)
	if comment.Article.Id == uuid.Nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "article_id",
			Message:   "Valid article ID is required",
			Value:     fmt.Sprintf("%d", comment.Article.Id),
			RequestID: requestID,
		})
	}

	// Check context timeout after reference validation
	if err := vs.checkContextTimeout(ctx); err != nil {
		if appErr, ok := err.(*utils.AppError); ok {
			return appErr
		}
		return vs.errorWrapper.InternalError(ctx, err, "Context error during comment validation")
	}

	// Additional content validation for inappropriate content
	if err := vs.validateCommentContent(ctx, comment.Content); err != nil {
		fieldErrors = append(fieldErrors, FieldError{
			Field:     "content",
			Message:   err.Error(),
			Value:     comment.Content,
			RequestID: requestID,
		})
	}

	// Return validation error if there are field errors
	if len(fieldErrors) > 0 {
		details := make(map[string]string)
		for _, fieldErr := range fieldErrors {
			details[fieldErr.Field] = fieldErr.Message
		}

		return &utils.AppError{
			Code:       utils.ErrValidation,
			Message:    "Comment validation failed",
			Details:    details,
			StatusCode: 400,
			RequestID:  requestID,
			Timestamp:  time.Now(),
		}
	}

	return nil
}

// validateCommentContent validates comment content for inappropriate content with context support
func (vs *validationService) validateCommentContent(ctx context.Context, content string) error {
	// Check context timeout
	if err := vs.checkContextTimeout(ctx); err != nil {
		return err
	}

	// Basic content validation - check for potentially harmful content
	content = strings.ToLower(strings.TrimSpace(content))

	// Check for excessive whitespace or special characters
	spaceCount := strings.Count(content, " ")
	if spaceCount > len(content)/3 { // More than 1/3 of content is spaces
		return fmt.Errorf("content contains too much whitespace")
	}

	// Check for HTML/script tags (basic XSS prevention)
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	if htmlTagRegex.MatchString(content) {
		return fmt.Errorf("HTML tags are not allowed in comments")
	}

	// Check for excessive repetition of characters (simplified approach)
	for i := 0; i < len(content)-10; i++ {
		char := content[i]
		count := 1
		for j := i + 1; j < len(content) && content[j] == char; j++ {
			count++
			if count > 10 {
				return fmt.Errorf("excessive repetition of characters is not allowed")
			}
		}
	}

	// Check for URLs (optional - depending on business rules)
	urlRegex := regexp.MustCompile(`https?://[^\s]+`)
	if urlRegex.MatchString(content) {
		return fmt.Errorf("URLs are not allowed in comments")
	}

	return nil
}
