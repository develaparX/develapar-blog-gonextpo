package dto

import (
	"context"
	"time"
)

// Context keys for type safety (must match middleware keys)
type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	StartTimeKey contextKey = "start_time"
)

// APIResponse represents the standardized API response structure with context
// @Description Standard API response format with metadata and context information
type APIResponse struct {
	Success    bool                   `json:"success" example:"true"`                    // Indicates if the request was successful
	Data       interface{}            `json:"data,omitempty"`                           // Response data (varies by endpoint)
	Error      *ErrorResponse         `json:"error,omitempty"`                          // Error details (only present when success=false)
	Pagination *PaginationMetadata    `json:"pagination,omitempty"`                     // Pagination metadata (for paginated responses)
	Meta       *ResponseMetadata      `json:"meta,omitempty"`                           // Response metadata with request tracking
}

// ErrorResponse represents the standard error response format with context
// @Description Error response structure with detailed error information
type ErrorResponse struct {
	Code      string                 `json:"code" example:"VALIDATION_ERROR"`          // Error code for programmatic handling
	Message   string                 `json:"message" example:"Invalid input data"`     // Human-readable error message
	Details   map[string]interface{} `json:"details,omitempty"`                        // Additional error details
	RequestID string                 `json:"request_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"` // Request ID for tracking
	Timestamp time.Time              `json:"timestamp" example:"2025-07-24T20:43:16.123456789+07:00"`             // Error timestamp
}

// ResponseMetadata contains response metadata with context information
// @Description Response metadata containing request tracking and performance information
type ResponseMetadata struct {
	RequestID      string    `json:"request_id" example:"550e8400-e29b-41d4-a716-446655440000"`      // Unique request identifier for tracking
	ProcessingTime int64     `json:"processing_time_ms" example:"15000000"`                          // Request processing time in nanoseconds
	Version        string    `json:"version" example:"1.0.0"`                                        // API version
	Timestamp      time.Time `json:"timestamp" example:"2025-07-24T20:43:16.123456789+07:00"`       // Response generation timestamp
}

// PaginationMetadata contains pagination information with context
// @Description Pagination metadata for paginated responses
type PaginationMetadata struct {
	Page       int    `json:"page" example:"1"`                                                         // Current page number (1-based)
	Limit      int    `json:"limit" example:"10"`                                                       // Number of items per page
	Total      int    `json:"total" example:"100"`                                                      // Total number of items
	TotalPages int    `json:"total_pages" example:"10"`                                                 // Total number of pages
	HasNext    bool   `json:"has_next" example:"true"`                                                  // Whether there is a next page
	HasPrev    bool   `json:"has_prev" example:"false"`                                                 // Whether there is a previous page
	RequestID  string `json:"request_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`     // Request ID for tracking
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
		if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
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
		if rid, ok := ctx.Value(RequestIDKey).(string); ok {
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
	var processingTime int64 = 0

	if ctx != nil {
		if rid, ok := ctx.Value(RequestIDKey).(string); ok {
			requestID = rid
		}
		if startTime, ok := ctx.Value(StartTimeKey).(time.Time); ok {
			processingTime = int64(time.Since(startTime))
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
		if rid, ok := ctx.Value(RequestIDKey).(string); ok {
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