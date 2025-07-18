package dto

import (
	"context"
	"time"
)

// APIResponse represents the standardized API response structure with context
type APIResponse struct {
	Success    bool                   `json:"success"`
	Data       interface{}            `json:"data,omitempty"`
	Error      *ErrorResponse         `json:"error,omitempty"`
	Pagination *PaginationMetadata    `json:"pagination,omitempty"`
	Meta       *ResponseMetadata      `json:"meta,omitempty"`
}

// ErrorResponse represents the standard error response format with context
type ErrorResponse struct {
	Code      string                 `json:"code"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// ResponseMetadata contains response metadata with context information
type ResponseMetadata struct {
	RequestID      string        `json:"request_id"`
	ProcessingTime time.Duration `json:"processing_time_ms"`
	Version        string        `json:"version"`
	Timestamp      time.Time     `json:"timestamp"`
}

// PaginationMetadata contains pagination information with context
type PaginationMetadata struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int    `json:"total"`
	TotalPages int    `json:"total_pages"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
	RequestID  string `json:"request_id,omitempty"`
}

// SuccessResponse creates a standardized success response with context
func SuccessResponse(ctx context.Context, data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Meta:    buildResponseMetadata(ctx),
	}
}

// SuccessResponseWithPagination creates a success response with pagination metadata
func SuccessResponseWithPagination(ctx context.Context, data interface{}, pagination *PaginationMetadata) APIResponse {
	// Add request ID to pagination metadata if available
	if pagination != nil && ctx != nil {
		if requestID, ok := ctx.Value("request_id").(string); ok {
			pagination.RequestID = requestID
		}
	}

	return APIResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination,
		Meta:       buildResponseMetadata(ctx),
	}
}

// ErrorResponseFromError creates a standardized error response from an error with context
func ErrorResponseFromError(ctx context.Context, code, message string, details map[string]interface{}) APIResponse {
	requestID := ""
	if ctx != nil {
		if rid, ok := ctx.Value("request_id").(string); ok {
			requestID = rid
		}
	}

	return APIResponse{
		Success: false,
		Error: &ErrorResponse{
			Code:      code,
			Message:   message,
			Details:   details,
			RequestID: requestID,
			Timestamp: time.Now(),
		},
		Meta: buildResponseMetadata(ctx),
	}
}

// ValidationErrorResponse creates a validation error response with context
func ValidationErrorResponse(ctx context.Context, fieldErrors map[string]interface{}) APIResponse {
	return ErrorResponseFromError(ctx, "VALIDATION_ERROR", "Invalid input data", fieldErrors)
}

// NotFoundErrorResponse creates a not found error response with context
func NotFoundErrorResponse(ctx context.Context, resource string) APIResponse {
	return ErrorResponseFromError(ctx, "NOT_FOUND", resource+" not found", nil)
}

// UnauthorizedErrorResponse creates an unauthorized error response with context
func UnauthorizedErrorResponse(ctx context.Context) APIResponse {
	return ErrorResponseFromError(ctx, "UNAUTHORIZED", "Authentication required", nil)
}

// ForbiddenErrorResponse creates a forbidden error response with context
func ForbiddenErrorResponse(ctx context.Context) APIResponse {
	return ErrorResponseFromError(ctx, "FORBIDDEN", "Access denied", nil)
}

// InternalErrorResponse creates an internal server error response with context
func InternalErrorResponse(ctx context.Context) APIResponse {
	return ErrorResponseFromError(ctx, "INTERNAL_ERROR", "An unexpected error occurred", nil)
}

// TimeoutErrorResponse creates a timeout error response with context
func TimeoutErrorResponse(ctx context.Context, operation string) APIResponse {
	details := map[string]interface{}{
		"operation": operation,
	}
	return ErrorResponseFromError(ctx, "TIMEOUT_ERROR", "Request timeout", details)
}

// CancellationErrorResponse creates a cancellation error response with context
func CancellationErrorResponse(ctx context.Context, operation string) APIResponse {
	details := map[string]interface{}{
		"operation": operation,
	}
	return ErrorResponseFromError(ctx, "REQUEST_CANCELLED", "Request was cancelled", details)
}

// RateLimitErrorResponse creates a rate limit error response with context
func RateLimitErrorResponse(ctx context.Context, retryAfter int) APIResponse {
	details := map[string]interface{}{
		"retry_after_seconds": retryAfter,
	}
	return ErrorResponseFromError(ctx, "RATE_LIMIT_EXCEEDED", "Rate limit exceeded", details)
}

// buildResponseMetadata builds response metadata from context
func buildResponseMetadata(ctx context.Context) *ResponseMetadata {
	requestID := ""
	processingTime := time.Duration(0)

	if ctx != nil {
		if rid, ok := ctx.Value("request_id").(string); ok {
			requestID = rid
		}
		if startTime, ok := ctx.Value("start_time").(time.Time); ok {
			processingTime = time.Since(startTime)
		}
	}

	return &ResponseMetadata{
		RequestID:      requestID,
		ProcessingTime: processingTime,
		Version:        "1.0.0",
		Timestamp:      time.Now(),
	}
}

// CreatePaginationMetadata creates pagination metadata with context
func CreatePaginationMetadata(ctx context.Context, page, limit, total int) *PaginationMetadata {
	totalPages := (total + limit - 1) / limit
	if totalPages == 0 {
		totalPages = 1
	}

	requestID := ""
	if ctx != nil {
		if rid, ok := ctx.Value("request_id").(string); ok {
			requestID = rid
		}
	}

	return &PaginationMetadata{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    page < totalPages,
		HasPrev:    page > 1,
		RequestID:  requestID,
	}
}