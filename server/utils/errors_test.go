package utils

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAppError(t *testing.T) {
	t.Run("Error_WithCause", func(t *testing.T) {
		cause := errors.New("underlying error")
		appErr := &AppError{
			Code:    ErrInternal,
			Message: "Internal error",
			Cause:   cause,
		}

		expected := "INTERNAL_ERROR: Internal error (caused by: underlying error)"
		assert.Equal(t, expected, appErr.Error())
	})

	t.Run("Error_WithoutCause", func(t *testing.T) {
		appErr := &AppError{
			Code:    ErrValidation,
			Message: "Validation failed",
		}

		expected := "VALIDATION_ERROR: Validation failed"
		assert.Equal(t, expected, appErr.Error())
	})

	t.Run("Unwrap", func(t *testing.T) {
		cause := errors.New("underlying error")
		appErr := &AppError{
			Code:  ErrInternal,
			Cause: cause,
		}

		assert.Equal(t, cause, appErr.Unwrap())
	})
}

func TestErrorWrapper(t *testing.T) {
	wrapper := NewErrorWrapper()

	t.Run("WrapError_WithContext", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")
		ctx = context.WithValue(ctx, "user_id", "user-456")

		originalErr := errors.New("original error")
		appErr := wrapper.WrapError(ctx, originalErr, ErrInternal, "Wrapped error")

		assert.Equal(t, ErrInternal, appErr.Code)
		assert.Equal(t, "Wrapped error", appErr.Message)
		assert.Equal(t, originalErr, appErr.Cause)
		assert.Equal(t, "req-123", appErr.RequestID)
		assert.Equal(t, "user-456", appErr.UserID)
		assert.False(t, appErr.Timestamp.IsZero())
	})

	t.Run("ValidationError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")

		appErr := wrapper.ValidationError(ctx, "email", "Invalid email format")

		assert.Equal(t, ErrValidation, appErr.Code)
		assert.Equal(t, "Validation failed", appErr.Message)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "Invalid email format", appErr.Details["email"])
		assert.Equal(t, "req-123", appErr.RequestID)
	})

	t.Run("NotFoundError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")

		appErr := wrapper.NotFoundError(ctx, "User")

		assert.Equal(t, ErrNotFound, appErr.Code)
		assert.Equal(t, "User not found", appErr.Message)
		assert.Equal(t, 404, appErr.StatusCode)
		assert.Equal(t, "req-123", appErr.RequestID)
	})

	t.Run("UnauthorizedError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "user_id", "user-123")

		appErr := wrapper.UnauthorizedError(ctx, "Invalid token")

		assert.Equal(t, ErrUnauthorized, appErr.Code)
		assert.Equal(t, "Invalid token", appErr.Message)
		assert.Equal(t, 401, appErr.StatusCode)
		assert.Equal(t, "user-123", appErr.UserID)
	})

	t.Run("ForbiddenError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")

		appErr := wrapper.ForbiddenError(ctx, "")

		assert.Equal(t, ErrForbidden, appErr.Code)
		assert.Equal(t, "Access forbidden", appErr.Message)
		assert.Equal(t, 403, appErr.StatusCode)
		assert.Equal(t, "req-123", appErr.RequestID)
	})

	t.Run("InternalError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")

		originalErr := errors.New("database connection failed")
		appErr := wrapper.InternalError(ctx, originalErr, "")

		assert.Equal(t, ErrInternal, appErr.Code)
		assert.Equal(t, "Internal server error", appErr.Message)
		assert.Equal(t, 500, appErr.StatusCode)
		assert.Equal(t, originalErr, appErr.Cause)
		assert.Equal(t, "req-123", appErr.RequestID)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")

		originalErr := errors.New("connection timeout")
		appErr := wrapper.DatabaseError(ctx, originalErr, "SELECT")

		assert.Equal(t, ErrDatabase, appErr.Code)
		assert.Equal(t, "Database SELECT operation failed", appErr.Message)
		assert.Equal(t, 500, appErr.StatusCode)
		assert.Equal(t, originalErr, appErr.Cause)
		assert.Equal(t, "req-123", appErr.RequestID)
	})

	t.Run("RateLimitError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")

		appErr := wrapper.RateLimitError(ctx, 100, time.Minute)

		assert.Equal(t, ErrRateLimit, appErr.Code)
		assert.Equal(t, "Rate limit exceeded: 100 requests per 1m0s", appErr.Message)
		assert.Equal(t, 429, appErr.StatusCode)
		assert.Equal(t, "req-123", appErr.RequestID)
	})

	t.Run("TimeoutError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")

		appErr := wrapper.TimeoutError(ctx, "database query")

		assert.Equal(t, ErrTimeout, appErr.Code)
		assert.Equal(t, "database query operation timed out", appErr.Message)
		assert.Equal(t, 408, appErr.StatusCode)
		assert.Equal(t, "req-123", appErr.RequestID)
	})

	t.Run("CancellationError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")

		appErr := wrapper.CancellationError(ctx, "user creation")

		assert.Equal(t, ErrCancelled, appErr.Code)
		assert.Equal(t, "user creation operation was cancelled", appErr.Message)
		assert.Equal(t, 499, appErr.StatusCode)
		assert.Equal(t, "req-123", appErr.RequestID)
	})

	t.Run("ConflictError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")

		appErr := wrapper.ConflictError(ctx, "User", "")

		assert.Equal(t, ErrConflict, appErr.Code)
		assert.Equal(t, "User already exists", appErr.Message)
		assert.Equal(t, 409, appErr.StatusCode)
		assert.Equal(t, "req-123", appErr.RequestID)
	})

	t.Run("BadRequestError", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")

		appErr := wrapper.BadRequestError(ctx, "Invalid JSON format")

		assert.Equal(t, ErrBadRequest, appErr.Code)
		assert.Equal(t, "Invalid JSON format", appErr.Message)
		assert.Equal(t, 400, appErr.StatusCode)
		assert.Equal(t, "req-123", appErr.RequestID)
	})

	t.Run("ExtractContextInfo_EmptyContext", func(t *testing.T) {
		appErr := wrapper.ValidationError(context.Background(), "field", "message")

		assert.Empty(t, appErr.RequestID)
		assert.Empty(t, appErr.UserID)
	})

	t.Run("ExtractContextInfo_NilContext", func(t *testing.T) {
		appErr := wrapper.ValidationError(nil, "field", "message")

		assert.Empty(t, appErr.RequestID)
		assert.Empty(t, appErr.UserID)
	})
}

func TestHelperFunctions(t *testing.T) {
	t.Run("IsTimeoutError", func(t *testing.T) {
		timeoutErr := &AppError{Code: ErrTimeout}
		validationErr := &AppError{Code: ErrValidation}
		regularErr := errors.New("regular error")

		assert.True(t, IsTimeoutError(timeoutErr))
		assert.False(t, IsTimeoutError(validationErr))
		assert.False(t, IsTimeoutError(regularErr))
	})

	t.Run("IsCancellationError", func(t *testing.T) {
		cancelErr := &AppError{Code: ErrCancelled}
		validationErr := &AppError{Code: ErrValidation}
		regularErr := errors.New("regular error")

		assert.True(t, IsCancellationError(cancelErr))
		assert.False(t, IsCancellationError(validationErr))
		assert.False(t, IsCancellationError(regularErr))
	})

	t.Run("IsValidationError", func(t *testing.T) {
		validationErr := &AppError{Code: ErrValidation}
		timeoutErr := &AppError{Code: ErrTimeout}
		regularErr := errors.New("regular error")

		assert.True(t, IsValidationError(validationErr))
		assert.False(t, IsValidationError(timeoutErr))
		assert.False(t, IsValidationError(regularErr))
	})

	t.Run("GetStatusCode", func(t *testing.T) {
		appErrWithStatus := &AppError{StatusCode: 404}
		appErrWithoutStatus := &AppError{}
		regularErr := errors.New("regular error")

		assert.Equal(t, 404, GetStatusCode(appErrWithStatus))
		assert.Equal(t, 500, GetStatusCode(appErrWithoutStatus))
		assert.Equal(t, 500, GetStatusCode(regularErr))
	})
}