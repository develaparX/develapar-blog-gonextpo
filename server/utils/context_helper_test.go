package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetUserIDFromGinContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		setupFunc   func(*gin.Context)
		expectedID  int
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid int user ID",
			setupFunc: func(c *gin.Context) {
				c.Set("userId", 123)
			},
			expectedID:  123,
			expectError: false,
		},
		{
			name: "Valid float64 user ID (JWT claims format)",
			setupFunc: func(c *gin.Context) {
				c.Set("userId", float64(456))
			},
			expectedID:  456,
			expectError: false,
		},
		{
			name: "Missing user ID in context",
			setupFunc: func(c *gin.Context) {
				// Don't set userId
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "user ID not found in context",
		},
		{
			name: "String user ID (unsupported)",
			setupFunc: func(c *gin.Context) {
				c.Set("userId", "123")
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "user ID stored as string, expected numeric type",
		},
		{
			name: "Invalid type user ID",
			setupFunc: func(c *gin.Context) {
				c.Set("userId", []int{123})
			},
			expectedID:  0,
			expectError: true,
			errorMsg:    "user ID has unexpected type: []int",
		},
		{
			name: "Zero user ID",
			setupFunc: func(c *gin.Context) {
				c.Set("userId", 0)
			},
			expectedID:  0,
			expectError: false,
		},
		{
			name: "Negative user ID",
			setupFunc: func(c *gin.Context) {
				c.Set("userId", -1)
			},
			expectedID:  -1,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			// Setup the context with test data
			tt.setupFunc(c)

			// Call the function under test
			userID, err := GetUserIDFromGinContext(c)

			// Assert results
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Equal(t, 0, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}

func TestGetUserRoleFromContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name         string
		setupFunc    func(*gin.Context)
		expectedRole string
		expectError  bool
		errorMsg     string
	}{
		{
			name: "Valid user role",
			setupFunc: func(c *gin.Context) {
				c.Set("role", "user")
			},
			expectedRole: "user",
			expectError:  false,
		},
		{
			name: "Valid admin role",
			setupFunc: func(c *gin.Context) {
				c.Set("role", "admin")
			},
			expectedRole: "admin",
			expectError:  false,
		},
		{
			name: "Missing role in context",
			setupFunc: func(c *gin.Context) {
				// Don't set role
			},
			expectedRole: "",
			expectError:  true,
			errorMsg:     "user role not found in context",
		},
		{
			name: "Empty role string",
			setupFunc: func(c *gin.Context) {
				c.Set("role", "")
			},
			expectedRole: "",
			expectError:  true,
			errorMsg:     "user role is empty",
		},
		{
			name: "Invalid type role",
			setupFunc: func(c *gin.Context) {
				c.Set("role", 123)
			},
			expectedRole: "",
			expectError:  true,
			errorMsg:     "user role has unexpected type: int",
		},
		{
			name: "Role with whitespace",
			setupFunc: func(c *gin.Context) {
				c.Set("role", "  admin  ")
			},
			expectedRole: "  admin  ",
			expectError:  false,
		},
		{
			name: "Custom role",
			setupFunc: func(c *gin.Context) {
				c.Set("role", "moderator")
			},
			expectedRole: "moderator",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test Gin context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			// Setup the context with test data
			tt.setupFunc(c)

			// Call the function under test
			role, err := GetUserRoleFromContext(c)

			// Assert results
			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Equal(t, "", role)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRole, role)
			}
		})
	}
}

// TestContextHelperIntegration tests both functions working together
func TestContextHelperIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Both user ID and role present", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		// Simulate what auth middleware would set
		c.Set("userId", float64(123)) // JWT claims are typically float64
		c.Set("role", "admin")

		userID, err := GetUserIDFromGinContext(c)
		assert.NoError(t, err)
		assert.Equal(t, 123, userID)

		role, err := GetUserRoleFromContext(c)
		assert.NoError(t, err)
		assert.Equal(t, "admin", role)
	})

	t.Run("Missing both user ID and role", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

		// Don't set anything in context

		userID, err := GetUserIDFromGinContext(c)
		assert.Error(t, err)
		assert.Equal(t, 0, userID)

		role, err := GetUserRoleFromContext(c)
		assert.Error(t, err)
		assert.Equal(t, "", role)
	})
}