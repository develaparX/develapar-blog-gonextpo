package utils

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetUserIDFromGinContext safely extracts the user ID from Gin context
// Returns the user ID as an integer and an error if extraction fails
func GetUserIDFromGinContext(c *gin.Context) (int, error) {
	userIDValue, exists := c.Get("userId")
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	// Handle different possible types that might be stored in context
	switch v := userIDValue.(type) {
	case int:
		return v, nil
	case float64:
		// JWT claims are often parsed as float64
		return int(v), nil
	case string:
		// In case it's stored as string, we would need to parse it
		return 0, errors.New("user ID stored as string, expected numeric type")
	default:
		return 0, fmt.Errorf("user ID has unexpected type: %T", v)
	}
}

// GetUserRoleFromContext safely extracts the user role from Gin context
// Returns the user role as a string and an error if extraction fails
func GetUserRoleFromContext(c *gin.Context) (string, error) {
	roleValue, exists := c.Get("role")
	if !exists {
		return "", errors.New("user role not found in context")
	}

	role, ok := roleValue.(string)
	if !ok {
		return "", fmt.Errorf("user role has unexpected type: %T", roleValue)
	}

	if role == "" {
		return "", errors.New("user role is empty")
	}

	return role, nil
}