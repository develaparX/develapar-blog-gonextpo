package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"develapar-server/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteUserHandler_Authorization(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		contextUserID  int
		contextRole    string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "User can delete own account",
			userID:         "1",
			contextUserID:  1,
			contextRole:    "user",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "User cannot delete other user's account",
			userID:         "2",
			contextUserID:  1,
			contextRole:    "user",
			expectedStatus: http.StatusForbidden,
			expectedError:  "Forbidden",
		},
		{
			name:           "Admin can delete any user's account",
			userID:         "2",
			contextUserID:  1,
			contextRole:    "admin",
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name:           "Invalid user ID format",
			userID:         "invalid",
			contextUserID:  1,
			contextRole:    "user",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid user ID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			mockAuth := new(MockAuthMiddleware)
			testLogger := &TestLogger{}
			errorHandler := middleware.NewErrorHandler(testLogger)

			// Setup router with middleware first
			router := gin.New()
			
			// Apply context middleware before routes
			router.Use(func(c *gin.Context) {
				c.Set("userId", tt.contextUserID)
				c.Set("role", tt.contextRole)
				c.Next()
			})
			
			rg := router.Group("/api")
			
			// Create controller and setup routes
			controller := NewUserController(mockService, mockAuth, rg, errorHandler)
			controller.Route()

			// Mock service call for successful cases
			if tt.expectedStatus == http.StatusOK {
				mockService.On("DeleteUser", mock.Anything, mock.AnythingOfType("int")).Return(nil)
			}

			// Create request
			req, _ := http.NewRequest("DELETE", "/api/users/"+tt.userID, nil)

			// Execute request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				responseBody := w.Body.String()
				assert.Contains(t, responseBody, tt.expectedError)
			}

			// Verify mock expectations for successful cases
			if tt.expectedStatus == http.StatusOK {
				mockService.AssertExpectations(t)
			}
		})
	}
}