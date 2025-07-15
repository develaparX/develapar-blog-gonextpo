package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Context keys for type safety
type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	UserIDKey    contextKey = "user_id"
	StartTimeKey contextKey = "start_time"
)

// ContextManager interface for managing request context
type ContextManager interface {
	WithRequestID(ctx context.Context, requestID string) context.Context
	WithUserID(ctx context.Context, userID string) context.Context
	WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc)
	GetRequestID(ctx context.Context) string
	GetUserID(ctx context.Context) string
	GetStartTime(ctx context.Context) time.Time
}

// contextManager implements ContextManager interface
type contextManager struct{}

// WithRequestID adds request ID to context
func (cm *contextManager) WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID adds user ID to context
func (cm *contextManager) WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithTimeout creates a context with timeout
func (cm *contextManager) WithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

// GetRequestID retrieves request ID from context
func (cm *contextManager) GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}

// GetUserID retrieves user ID from context
func (cm *contextManager) GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetStartTime retrieves start time from context
func (cm *contextManager) GetStartTime(ctx context.Context) time.Time {
	if startTime, ok := ctx.Value(StartTimeKey).(time.Time); ok {
		return startTime
	}
	return time.Time{}
}

// NewContextManager creates a new context manager
func NewContextManager() ContextManager {
	return &contextManager{}
}

// ContextMiddleware interface for Gin middleware
type ContextMiddleware interface {
	InjectContext() gin.HandlerFunc
}

// contextMiddleware implements ContextMiddleware interface
type contextMiddleware struct {
	manager ContextManager
}

// InjectContext middleware injects request ID, user ID, and start time into context
func (cm *contextMiddleware) InjectContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID if not present
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Get user ID from existing context if available (set by auth middleware)
		userID := ""
		if userIDValue, exists := c.Get("userId"); exists {
			if uid, ok := userIDValue.(string); ok {
				userID = uid
			}
		}

		// Create enhanced context
		ctx := c.Request.Context()
		ctx = cm.manager.WithRequestID(ctx, requestID)
		ctx = cm.manager.WithUserID(ctx, userID)
		ctx = context.WithValue(ctx, StartTimeKey, time.Now())

		// Update request context
		c.Request = c.Request.WithContext(ctx)

		// Set values in Gin context for easy access
		c.Set("request_id", requestID)
		c.Set("user_id", userID)
		c.Set("start_time", time.Now())

		// Add request ID to response header
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// NewContextMiddleware creates a new context middleware
func NewContextMiddleware(manager ContextManager) ContextMiddleware {
	return &contextMiddleware{
		manager: manager,
	}
}