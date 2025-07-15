package utils

import (
	"context"
	"fmt"
	"time"
)

// Error constants for different error types
const (
	ErrValidation    = "VALIDATION_ERROR"
	ErrNotFound      = "NOT_FOUND"
	ErrUnauthorized  = "UNAUTHORIZED"
	ErrForbidden     = "FORBIDDEN"
	ErrInternal      = "INTERNAL_ERROR"
	ErrDatabase      = "DATABASE_ERROR"
	ErrRateLimit     = "RATE_LIMIT_EXCEEDED"
	ErrTimeout       = "TIMEOUT_ERROR"
	ErrCancelled     = "REQUEST_CANCELLED"
	ErrConflict      = "CONFLICT_ERROR"
	ErrBadRequest    = "BAD_REQUEST"
)

// AppError represents a custom application error with context information
type AppError struct {
	Code       string            `json:"code"`
	Message    string            `json:"message"`
	Details    map[string]string `json:"details,omitempty"`
	StatusCode int               `json:"-"`
	Cause      error             `json:"-"`
	RequestID  string            `json:"request_id,omitempty"`
	UserID     string            `json:"user_id,omitempty"`
	Timestamp  time.Time         `json:"timestamp"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// ErrorWrapper interface for wrapping errors with context
type ErrorWrapper interface {
	WrapError(ctx context.Context, err error, code string, message string) *AppError
	ValidationError(ctx context.Context, field string, message string) *AppError
	NotFoundError(ctx context.Context, resource string) *AppError
	UnauthorizedError(ctx context.Context, message string) *AppError
	ForbiddenError(ctx context.Context, message string) *AppError
	InternalError(ctx context.Context, err error, message string) *AppError
	DatabaseError(ctx context.Context, err error, operation string) *AppError
	RateLimitError(ctx context.Context, limit int, window time.Duration) *AppError
	TimeoutError(ctx context.Context, operation string) *AppError
	CancellationError(ctx context.Context, operation string) *AppError
	ConflictError(ctx context.Context, resource string, message string) *AppError
	BadRequestError(ctx context.Context, message string) *AppError
}

// errorWrapper implements ErrorWrapper interface
type errorWrapper struct{}

// extractContextInfo extracts request ID and user ID from context
func (ew *errorWrapper) extractContextInfo(ctx context.Context) (string, string) {
	var requestID, userID string
	
	if ctx != nil {
		if rid, ok := ctx.Value("request_id").(string); ok {
			requestID = rid
		}
		if uid, ok := ctx.Value("user_id").(string); ok {
			userID = uid
		}
	}
	
	return requestID, userID
}

// WrapError wraps an existing error with context information
func (ew *errorWrapper) WrapError(ctx context.Context, err error, code string, message string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	return &AppError{
		Code:      code,
		Message:   message,
		Cause:     err,
		RequestID: requestID,
		UserID:    userID,
		Timestamp: time.Now(),
	}
}

// ValidationError creates a validation error with context
func (ew *errorWrapper) ValidationError(ctx context.Context, field string, message string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	details := make(map[string]string)
	if field != "" {
		details[field] = message
	}
	
	return &AppError{
		Code:       ErrValidation,
		Message:    "Validation failed",
		Details:    details,
		StatusCode: 400,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// NotFoundError creates a not found error with context
func (ew *errorWrapper) NotFoundError(ctx context.Context, resource string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	message := "Resource not found"
	if resource != "" {
		message = fmt.Sprintf("%s not found", resource)
	}
	
	return &AppError{
		Code:       ErrNotFound,
		Message:    message,
		StatusCode: 404,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// UnauthorizedError creates an unauthorized error with context
func (ew *errorWrapper) UnauthorizedError(ctx context.Context, message string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	if message == "" {
		message = "Authentication required"
	}
	
	return &AppError{
		Code:       ErrUnauthorized,
		Message:    message,
		StatusCode: 401,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// ForbiddenError creates a forbidden error with context
func (ew *errorWrapper) ForbiddenError(ctx context.Context, message string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	if message == "" {
		message = "Access forbidden"
	}
	
	return &AppError{
		Code:       ErrForbidden,
		Message:    message,
		StatusCode: 403,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// InternalError creates an internal server error with context
func (ew *errorWrapper) InternalError(ctx context.Context, err error, message string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	if message == "" {
		message = "Internal server error"
	}
	
	return &AppError{
		Code:       ErrInternal,
		Message:    message,
		Cause:      err,
		StatusCode: 500,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// DatabaseError creates a database error with context
func (ew *errorWrapper) DatabaseError(ctx context.Context, err error, operation string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	message := "Database operation failed"
	if operation != "" {
		message = fmt.Sprintf("Database %s operation failed", operation)
	}
	
	return &AppError{
		Code:       ErrDatabase,
		Message:    message,
		Cause:      err,
		StatusCode: 500,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// RateLimitError creates a rate limit error with context
func (ew *errorWrapper) RateLimitError(ctx context.Context, limit int, window time.Duration) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	message := fmt.Sprintf("Rate limit exceeded: %d requests per %v", limit, window)
	
	return &AppError{
		Code:       ErrRateLimit,
		Message:    message,
		StatusCode: 429,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// TimeoutError creates a timeout error with context
func (ew *errorWrapper) TimeoutError(ctx context.Context, operation string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	message := "Operation timed out"
	if operation != "" {
		message = fmt.Sprintf("%s operation timed out", operation)
	}
	
	return &AppError{
		Code:       ErrTimeout,
		Message:    message,
		StatusCode: 408,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// CancellationError creates a cancellation error with context
func (ew *errorWrapper) CancellationError(ctx context.Context, operation string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	message := "Operation was cancelled"
	if operation != "" {
		message = fmt.Sprintf("%s operation was cancelled", operation)
	}
	
	return &AppError{
		Code:       ErrCancelled,
		Message:    message,
		StatusCode: 499,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// ConflictError creates a conflict error with context
func (ew *errorWrapper) ConflictError(ctx context.Context, resource string, message string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	if message == "" {
		message = "Resource conflict"
		if resource != "" {
			message = fmt.Sprintf("%s already exists", resource)
		}
	}
	
	return &AppError{
		Code:       ErrConflict,
		Message:    message,
		StatusCode: 409,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// BadRequestError creates a bad request error with context
func (ew *errorWrapper) BadRequestError(ctx context.Context, message string) *AppError {
	requestID, userID := ew.extractContextInfo(ctx)
	
	if message == "" {
		message = "Bad request"
	}
	
	return &AppError{
		Code:       ErrBadRequest,
		Message:    message,
		StatusCode: 400,
		RequestID:  requestID,
		UserID:     userID,
		Timestamp:  time.Now(),
	}
}

// NewErrorWrapper creates a new error wrapper
func NewErrorWrapper() ErrorWrapper {
	return &errorWrapper{}
}

// Helper functions for common error scenarios

// IsTimeoutError checks if the error is a timeout error
func IsTimeoutError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrTimeout
	}
	return false
}

// IsCancellationError checks if the error is a cancellation error
func IsCancellationError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrCancelled
	}
	return false
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code == ErrValidation
	}
	return false
}

// GetStatusCode returns the HTTP status code for an error
func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok && appErr.StatusCode > 0 {
		return appErr.StatusCode
	}
	return 500 // Default to internal server error
}