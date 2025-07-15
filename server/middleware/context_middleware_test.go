package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestContextManager(t *testing.T) {
	manager := NewContextManager()

	t.Run("WithRequestID", func(t *testing.T) {
		ctx := context.Background()
		requestID := "test-request-id"
		
		newCtx := manager.WithRequestID(ctx, requestID)
		
		assert.Equal(t, requestID, manager.GetRequestID(newCtx))
	})

	t.Run("WithUserID", func(t *testing.T) {
		ctx := context.Background()
		userID := "test-user-id"
		
		newCtx := manager.WithUserID(ctx, userID)
		
		assert.Equal(t, userID, manager.GetUserID(newCtx))
	})

	t.Run("WithTimeout", func(t *testing.T) {
		ctx := context.Background()
		timeout := 5 * time.Second
		
		newCtx, cancel := manager.WithTimeout(ctx, timeout)
		defer cancel()
		
		deadline, ok := newCtx.Deadline()
		assert.True(t, ok)
		assert.True(t, deadline.After(time.Now()))
	})

	t.Run("GetRequestID_Empty", func(t *testing.T) {
		ctx := context.Background()
		
		requestID := manager.GetRequestID(ctx)
		
		assert.Empty(t, requestID)
	})

	t.Run("GetUserID_Empty", func(t *testing.T) {
		ctx := context.Background()
		
		userID := manager.GetUserID(ctx)
		
		assert.Empty(t, userID)
	})

	t.Run("GetStartTime_Empty", func(t *testing.T) {
		ctx := context.Background()
		
		startTime := manager.GetStartTime(ctx)
		
		assert.True(t, startTime.IsZero())
	})
}

func TestContextMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("InjectContext_GeneratesRequestID", func(t *testing.T) {
		manager := NewContextManager()
		middleware := NewContextMiddleware(manager)
		
		router := gin.New()
		router.Use(middleware.InjectContext())
		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("request_id")
			assert.True(t, exists)
			assert.NotEmpty(t, requestID)
			
			// Check context
			ctx := c.Request.Context()
			assert.NotEmpty(t, manager.GetRequestID(ctx))
			
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	})

	t.Run("InjectContext_UsesProvidedRequestID", func(t *testing.T) {
		manager := NewContextManager()
		middleware := NewContextMiddleware(manager)
		expectedRequestID := "custom-request-id"
		
		router := gin.New()
		router.Use(middleware.InjectContext())
		router.GET("/test", func(c *gin.Context) {
			requestID, exists := c.Get("request_id")
			assert.True(t, exists)
			assert.Equal(t, expectedRequestID, requestID)
			
			// Check context
			ctx := c.Request.Context()
			assert.Equal(t, expectedRequestID, manager.GetRequestID(ctx))
			
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("X-Request-ID", expectedRequestID)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, expectedRequestID, w.Header().Get("X-Request-ID"))
	})

	t.Run("InjectContext_WithUserID", func(t *testing.T) {
		manager := NewContextManager()
		middleware := NewContextMiddleware(manager)
		expectedUserID := "test-user-123"
		
		router := gin.New()
		router.Use(func(c *gin.Context) {
			// Simulate auth middleware setting userId
			c.Set("userId", expectedUserID)
			c.Next()
		})
		router.Use(middleware.InjectContext())
		router.GET("/test", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			assert.True(t, exists)
			assert.Equal(t, expectedUserID, userID)
			
			// Check context
			ctx := c.Request.Context()
			assert.Equal(t, expectedUserID, manager.GetUserID(ctx))
			
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("InjectContext_SetsStartTime", func(t *testing.T) {
		manager := NewContextManager()
		middleware := NewContextMiddleware(manager)
		
		router := gin.New()
		router.Use(middleware.InjectContext())
		router.GET("/test", func(c *gin.Context) {
			startTime, exists := c.Get("start_time")
			assert.True(t, exists)
			assert.IsType(t, time.Time{}, startTime)
			
			// Check context
			ctx := c.Request.Context()
			contextStartTime := manager.GetStartTime(ctx)
			assert.False(t, contextStartTime.IsZero())
			
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
	})
}