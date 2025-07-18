package utils

import (
	"context"
	"net/http"

	"develapar-server/model/dto"
	"github.com/gin-gonic/gin"
)

// ResponseHelper provides helper functions for sending standardized responses
type ResponseHelper struct{}

// NewResponseHelper creates a new response helper instance
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// SendSuccess sends a standardized success response with context
func (rh *ResponseHelper) SendSuccess(c *gin.Context, data interface{}) {
	ctx := c.Request.Context()
	response := dto.SuccessResponse(ctx, data)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusOK, response)
}

// SendSuccessWithPagination sends a success response with pagination metadata
func (rh *ResponseHelper) SendSuccessWithPagination(c *gin.Context, data interface{}, pagination *dto.PaginationMetadata) {
	ctx := c.Request.Context()
	response := dto.SuccessResponseWithPagination(ctx, data, pagination)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusOK, response)
}

// convertServicePaginationToDTO converts service pagination metadata to dto pagination metadata
func (rh *ResponseHelper) convertServicePaginationToDTO(ctx context.Context, servicePagination interface{}) *dto.PaginationMetadata {
	// Use reflection or type assertion to convert
	if sp, ok := servicePagination.(struct {
		Page       int    `json:"page"`
		Limit      int    `json:"limit"`
		Total      int    `json:"total"`
		TotalPages int    `json:"total_pages"`
		HasNext    bool   `json:"has_next"`
		HasPrev    bool   `json:"has_prev"`
		RequestID  string `json:"request_id,omitempty"`
	}); ok {
		return dto.CreatePaginationMetadata(ctx, sp.Page, sp.Limit, sp.Total)
	}
	return nil
}

// SendSuccessWithServicePagination sends a success response with service pagination metadata
func (rh *ResponseHelper) SendSuccessWithServicePagination(c *gin.Context, data interface{}, servicePagination interface{}) {
	ctx := c.Request.Context()
	
	// Convert service pagination to dto pagination
	dtoPagination := rh.convertServicePaginationToDTO(ctx, servicePagination)
	
	response := dto.SuccessResponseWithPagination(ctx, data, dtoPagination)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusOK, response)
}

// SendCreated sends a standardized created response with context
func (rh *ResponseHelper) SendCreated(c *gin.Context, data interface{}) {
	ctx := c.Request.Context()
	response := dto.SuccessResponse(ctx, data)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusCreated, response)
}

// SendNoContent sends a standardized no content response with context
func (rh *ResponseHelper) SendNoContent(c *gin.Context) {
	ctx := c.Request.Context()
	response := dto.SuccessResponse(ctx, nil)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusNoContent, response)
}

// SendValidationError sends a validation error response with context
func (rh *ResponseHelper) SendValidationError(c *gin.Context, fieldErrors map[string]interface{}) {
	ctx := c.Request.Context()
	response := dto.ValidationErrorResponse(ctx, fieldErrors)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusBadRequest, response)
}

// SendNotFound sends a not found error response with context
func (rh *ResponseHelper) SendNotFound(c *gin.Context, resource string) {
	ctx := c.Request.Context()
	response := dto.NotFoundErrorResponse(ctx, resource)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusNotFound, response)
}

// SendUnauthorized sends an unauthorized error response with context
func (rh *ResponseHelper) SendUnauthorized(c *gin.Context) {
	ctx := c.Request.Context()
	response := dto.UnauthorizedErrorResponse(ctx)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusUnauthorized, response)
}

// SendForbidden sends a forbidden error response with context
func (rh *ResponseHelper) SendForbidden(c *gin.Context) {
	ctx := c.Request.Context()
	response := dto.ForbiddenErrorResponse(ctx)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusForbidden, response)
}

// SendInternalError sends an internal server error response with context
func (rh *ResponseHelper) SendInternalError(c *gin.Context) {
	ctx := c.Request.Context()
	response := dto.InternalErrorResponse(ctx)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusInternalServerError, response)
}

// SendTimeout sends a timeout error response with context
func (rh *ResponseHelper) SendTimeout(c *gin.Context, operation string) {
	ctx := c.Request.Context()
	response := dto.TimeoutErrorResponse(ctx, operation)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(http.StatusRequestTimeout, response)
}

// SendCancellation sends a cancellation error response with context
func (rh *ResponseHelper) SendCancellation(c *gin.Context, operation string) {
	ctx := c.Request.Context()
	response := dto.CancellationErrorResponse(ctx, operation)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(499, response) // 499 Client Closed Request
}

// SendRateLimit sends a rate limit error response with context
func (rh *ResponseHelper) SendRateLimit(c *gin.Context, retryAfter int) {
	ctx := c.Request.Context()
	response := dto.RateLimitErrorResponse(ctx, retryAfter)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	// Add rate limit headers
	c.Header("Retry-After", string(rune(retryAfter)))
	c.Header("X-RateLimit-Limit", "100") // This should be configurable
	c.Header("X-RateLimit-Remaining", "0")
	
	c.JSON(http.StatusTooManyRequests, response)
}

// SendCustomError sends a custom error response with context
func (rh *ResponseHelper) SendCustomError(c *gin.Context, statusCode int, code, message string, details map[string]interface{}) {
	ctx := c.Request.Context()
	response := dto.ErrorResponseFromError(ctx, code, message, details)
	
	// Add request ID to response headers if available
	if response.Meta != nil && response.Meta.RequestID != "" {
		c.Header("X-Request-ID", response.Meta.RequestID)
	}
	
	c.JSON(statusCode, response)
}