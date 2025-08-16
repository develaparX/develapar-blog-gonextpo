package utils

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetUserIDFromGinContext safely extracts the user ID from Gin context
// Returns the user ID as an integer and an error if extraction fails
func GetUserIDFromGinContext(c *gin.Context) (uuid.UUID, error) {
	// 1. Ambil nilai dari context
	userIDValue, exists := c.Get("userId")
	if !exists {
		// Kembalikan UUID nil jika tidak ada
		return uuid.Nil, errors.New("user ID not found in context")
	}

	// 2. Handle tipe data yang mungkin ada di context
	switch v := userIDValue.(type) {
	case string:
		// Jika tipenya string, parse ke UUID
		parsedUUID, err := uuid.Parse(v)
		if err != nil {
			// Gagal parse, kembalikan error
			return uuid.Nil, fmt.Errorf("failed to parse user ID string '%s' as UUID: %w", v, err)
		}
		return parsedUUID, nil
	case uuid.UUID:
		// Jika tipenya sudah UUID (kasus ideal), langsung kembalikan
		return v, nil
	default:
		// Jika tipenya bukan string atau UUID, ini adalah error
		return uuid.Nil, fmt.Errorf("user ID has an unexpected type: %T", v)
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
