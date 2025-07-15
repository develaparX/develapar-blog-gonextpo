package middleware

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"develapar-server/utils"
)

// ErrorHandler interface for handling errors with context support
type ErrorHandler interface {
	HandleError(ctx context.Context, c *gin.Context, err error)
	WrapError(ctx context.Context, err error, code string, message string) *utils.AppError
	ValidationError(ctx context.Context, field string, message string) *utils.AppError
	TimeoutError(ctx context.Context, operation string) *utils.AppError
	CancellationError(ctx context.Context, operation string) *utils.AppError
}

// errorHandler implements ErrorHandler interface
type errorHandler struct {
	wrapper utils.ErrorWrapper
	logger  Logger
}

// Logger interface for error logging
type Logger interface {
	Error(ctx context.Context, msg string, err error, fields map[string]interface{})
	Warn(ctx context.Context, msg string, fields map[string]interface{})
	Info(ctx context.Context, msg string, fields map[string]interface{})
}

// defaultLogger implements Logger interface using standard log package
type defaultLogger struct{}

func (l *defaultLogger) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	requestID := ""
	userID := ""
	
	if ctx != nil {
		if rid, ok := ctx.Value("request_id").(string); ok {
			requestID = rid
		}
		if uid, ok := ctx.Value("user_id").(string); ok {
			userID = uid
		}
	}
	
	log.Printf("[ERROR] %s | RequestID: %s | UserID: %s | Error: %v | Fields: %+v", 
		msg, requestID, userID, err, fields)
}

func (l *defaultLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	requestID := ""
	if ctx != nil {
		if rid, ok := ctx.Value("request_id").(string); ok {
			requestID = rid
		}
	}
	
	log.Printf("[WARN] %s | RequestID: %s | Fields: %+v", msg, requestID, fields)
}

func (l *defaultLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	requestID := ""
	if ctx != nil {
		if rid, ok := ctx.Value("request_id").(string); ok {
			requestID = rid
		}
	}
	
	log.Printf("[INFO] %s | RequestID: %s | Fields: %+v", msg, requestID, fields)
}

// HandleError handles errors with context information
func (eh *errorHandler) HandleError(ctx context.Context, c *gin.Context, err error) {
	var appErr *utils.AppError
	
	// Check if it's already an AppError
	if ae, ok := err.(*utils.AppError); ok {
		appErr = ae
	} else {
		// Check for context-specific errors first
		if ctx != nil && ctx.Err() == context.DeadlineExceeded {
			appErr = eh.TimeoutError(ctx, "request")
		} else if ctx != nil && ctx.Err() == context.Canceled {
			appErr = eh.CancellationError(ctx, "request")
		} else {
			// Wrap as internal error
			appErr = eh.wrapper.InternalError(ctx, err, "An unexpected error occurred")
		}
	}
	
	// Log the error with context
	eh.logError(ctx, appErr)
	
	// Set response headers
	if requestID := appErr.RequestID; requestID != "" {
		c.Header("X-Request-ID", requestID)
	}
	
	// Create error response
	errorResponse := ErrorResponse{
		Success: false,
		Error: &ErrorDetail{
			Code:      appErr.Code,
			Message:   appErr.Message,
			Details:   appErr.Details,
			RequestID: appErr.RequestID,
			Timestamp: appErr.Timestamp,
		},
		Meta: &ResponseMetadata{
			RequestID:      appErr.RequestID,
			ProcessingTime: eh.calculateProcessingTime(ctx),
			Version:        "1.0.0",
			Timestamp:      time.Now(),
		},
	}
	
	// Determine status code
	statusCode := utils.GetStatusCode(appErr)
	
	c.JSON(statusCode, errorResponse)
	c.Abort()
}

// WrapError wraps an error with context information
func (eh *errorHandler) WrapError(ctx context.Context, err error, code string, message string) *utils.AppError {
	return eh.wrapper.WrapError(ctx, err, code, message)
}

// ValidationError creates a validation error with context
func (eh *errorHandler) ValidationError(ctx context.Context, field string, message string) *utils.AppError {
	return eh.wrapper.ValidationError(ctx, field, message)
}

// TimeoutError creates a timeout error with context
func (eh *errorHandler) TimeoutError(ctx context.Context, operation string) *utils.AppError {
	return eh.wrapper.TimeoutError(ctx, operation)
}

// CancellationError creates a cancellation error with context
func (eh *errorHandler) CancellationError(ctx context.Context, operation string) *utils.AppError {
	return eh.wrapper.CancellationError(ctx, operation)
}

// logError logs error with context information
func (eh *errorHandler) logError(ctx context.Context, appErr *utils.AppError) {
	fields := map[string]interface{}{
		"code":       appErr.Code,
		"status":     appErr.StatusCode,
		"request_id": appErr.RequestID,
		"user_id":    appErr.UserID,
		"timestamp":  appErr.Timestamp,
	}
	
	if appErr.Details != nil {
		fields["details"] = appErr.Details
	}
	
	// Log based on error severity
	switch appErr.Code {
	case utils.ErrValidation, utils.ErrNotFound, utils.ErrUnauthorized, utils.ErrForbidden:
		eh.logger.Warn(ctx, appErr.Message, fields)
	default:
		eh.logger.Error(ctx, appErr.Message, appErr.Cause, fields)
	}
}

// calculateProcessingTime calculates request processing time from context
func (eh *errorHandler) calculateProcessingTime(ctx context.Context) time.Duration {
	if ctx != nil {
		if startTime, ok := ctx.Value("start_time").(time.Time); ok {
			return time.Since(startTime)
		}
	}
	return 0
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger Logger) ErrorHandler {
	if logger == nil {
		logger = &defaultLogger{}
	}
	
	return &errorHandler{
		wrapper: utils.NewErrorWrapper(),
		logger:  logger,
	}
}

// ErrorMiddleware creates a Gin middleware for error handling
func ErrorMiddleware(handler ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		
		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Handle the last error (most recent)
			err := c.Errors.Last().Err
			ctx := c.Request.Context()
			
			handler.HandleError(ctx, c, err)
		}
	}
}

// RecoveryMiddleware creates a Gin middleware for panic recovery with context
func RecoveryMiddleware(handler ErrorHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctx := c.Request.Context()
				
				// Create internal error from panic
				var panicErr error
				if e, ok := err.(error); ok {
					panicErr = e
				} else {
					panicErr = utils.NewErrorWrapper().InternalError(ctx, nil, "Panic occurred")
				}
				
				handler.HandleError(ctx, c, panicErr)
			}
		}()
		
		c.Next()
	}
}

// Response structures for consistent API responses

// ErrorResponse represents the standard error response format
type ErrorResponse struct {
	Success bool              `json:"success"`
	Error   *ErrorDetail      `json:"error"`
	Meta    *ResponseMetadata `json:"meta,omitempty"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	Details   map[string]string `json:"details,omitempty"`
	RequestID string            `json:"request_id,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// ResponseMetadata contains response metadata
type ResponseMetadata struct {
	RequestID      string        `json:"request_id"`
	ProcessingTime time.Duration `json:"processing_time_ms"`
	Version        string        `json:"version"`
	Timestamp      time.Time     `json:"timestamp"`
}

// SuccessResponse represents the standard success response format
type SuccessResponse struct {
	Success bool              `json:"success"`
	Data    interface{}       `json:"data,omitempty"`
	Meta    *ResponseMetadata `json:"meta,omitempty"`
}

// Helper function to create success response
func CreateSuccessResponse(ctx context.Context, data interface{}) SuccessResponse {
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
	
	return SuccessResponse{
		Success: true,
		Data:    data,
		Meta: &ResponseMetadata{
			RequestID:      requestID,
			ProcessingTime: processingTime,
			Version:        "1.0.0",
			Timestamp:      time.Now(),
		},
	}
}