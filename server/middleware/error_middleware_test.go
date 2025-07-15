package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"develapar-server/utils"
)

// MockLogger for testing
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Error(ctx context.Context, msg string, err error, fields map[string]interface{}) {
	m.Called(ctx, msg, err, fields)
}

func (m *MockLogger) Warn(ctx context.Context, msg string, fields map[string]interface{}) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Info(ctx context.Context, msg string, fields map[string]interface{}) {
	m.Called(ctx, msg, fields)
}

func TestErrorHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("HandleError_AppError", func(t *testing.T) {
		handler := NewErrorHandler(nil) // Use default logger
		
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")
		ctx = context.WithValue(ctx, "start_time", time.Now().Add(-100*time.Millisecond))
		
		appErr := &utils.AppError{
			Code:       utils.ErrValidation,
			Message:    "Validation failed",
			StatusCode: 400,
			RequestID:  "req-123",
			Timestamp:  time.Now(),
		}
		
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			c.Request = c.Request.WithContext(ctx)
			handler.HandleError(ctx, c, appErr)
		})
		
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "req-123", w.Header().Get("X-Request-ID"))
	})

	t.Run("HandleError_ContextTimeout", func(t *testing.T) {
		handler := NewErrorHandler(nil) // Use default logger
		
		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		ctx = context.WithValue(ctx, "request_id", "req-123")
		time.Sleep(time.Millisecond) // Ensure timeout
		cancel()
		
		originalErr := errors.New("some operation failed")
		
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			c.Request = c.Request.WithContext(ctx)
			handler.HandleError(ctx, c, originalErr)
		})
		
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusRequestTimeout, w.Code)
	})

	t.Run("HandleError_ContextCancellation", func(t *testing.T) {
		handler := NewErrorHandler(nil) // Use default logger
		
		ctx, cancel := context.WithCancel(context.Background())
		ctx = context.WithValue(ctx, "request_id", "req-123")
		cancel() // Cancel immediately
		
		originalErr := errors.New("some operation failed")
		
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			c.Request = c.Request.WithContext(ctx)
			handler.HandleError(ctx, c, originalErr)
		})
		
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, 499, w.Code) // Client closed request
	})

	t.Run("HandleError_RegularError", func(t *testing.T) {
		handler := NewErrorHandler(nil) // Use default logger
		
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")
		
		originalErr := errors.New("database connection failed")
		
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			c.Request = c.Request.WithContext(ctx)
			handler.HandleError(ctx, c, originalErr)
		})
		
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestErrorMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("ErrorMiddleware_HandlesGinErrors", func(t *testing.T) {
		handler := NewErrorHandler(nil) // Use default logger
		
		router := gin.New()
		router.Use(ErrorMiddleware(handler))
		router.GET("/test", func(c *gin.Context) {
			c.Error(errors.New("test error"))
		})
		
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("ErrorMiddleware_NoErrors", func(t *testing.T) {
		handler := NewErrorHandler(nil) // Use default logger
		
		router := gin.New()
		router.Use(ErrorMiddleware(handler))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRecoveryMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("RecoveryMiddleware_HandlesPanic", func(t *testing.T) {
		handler := NewErrorHandler(nil) // Use default logger
		
		router := gin.New()
		router.Use(RecoveryMiddleware(handler))
		router.GET("/test", func(c *gin.Context) {
			panic("test panic")
		})
		
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("RecoveryMiddleware_NoPanic", func(t *testing.T) {
		handler := NewErrorHandler(nil) // Use default logger
		
		router := gin.New()
		router.Use(RecoveryMiddleware(handler))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})
		
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestDefaultLogger(t *testing.T) {
	logger := &defaultLogger{}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "req-123")
	ctx = context.WithValue(ctx, "user_id", "user-456")

	t.Run("Error", func(t *testing.T) {
		err := errors.New("test error")
		fields := map[string]interface{}{"key": "value"}
		
		// This should not panic
		logger.Error(ctx, "Test error message", err, fields)
	})

	t.Run("Warn", func(t *testing.T) {
		fields := map[string]interface{}{"key": "value"}
		
		// This should not panic
		logger.Warn(ctx, "Test warning message", fields)
	})

	t.Run("Info", func(t *testing.T) {
		fields := map[string]interface{}{"key": "value"}
		
		// This should not panic
		logger.Info(ctx, "Test info message", fields)
	})
}

func TestCreateSuccessResponse(t *testing.T) {
	t.Run("WithContext", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, "request_id", "req-123")
		ctx = context.WithValue(ctx, "start_time", time.Now().Add(-100*time.Millisecond))
		
		data := map[string]string{"key": "value"}
		response := CreateSuccessResponse(ctx, data)
		
		assert.True(t, response.Success)
		assert.Equal(t, data, response.Data)
		assert.Equal(t, "req-123", response.Meta.RequestID)
		assert.True(t, response.Meta.ProcessingTime > 0)
		assert.Equal(t, "1.0.0", response.Meta.Version)
		assert.False(t, response.Meta.Timestamp.IsZero())
	})

	t.Run("WithoutContext", func(t *testing.T) {
		data := map[string]string{"key": "value"}
		response := CreateSuccessResponse(nil, data)
		
		assert.True(t, response.Success)
		assert.Equal(t, data, response.Data)
		assert.Empty(t, response.Meta.RequestID)
		assert.Equal(t, time.Duration(0), response.Meta.ProcessingTime)
		assert.Equal(t, "1.0.0", response.Meta.Version)
		assert.False(t, response.Meta.Timestamp.IsZero())
	})
}

func TestNewErrorHandler(t *testing.T) {
	t.Run("WithLogger", func(t *testing.T) {
		mockLogger := &MockLogger{}
		handler := NewErrorHandler(mockLogger)
		
		assert.NotNil(t, handler)
	})

	t.Run("WithoutLogger", func(t *testing.T) {
		handler := NewErrorHandler(nil)
		
		assert.NotNil(t, handler)
	})
}