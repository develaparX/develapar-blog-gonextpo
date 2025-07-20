package utils

import (
	"fmt"
)

// AuthorizationError represents authorization-related errors
type AuthorizationError struct {
	Code    string
	Message string
	Details string
}

func (e *AuthorizationError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Authorization error codes (using existing constants from errors.go)
const (
	ErrCodeInvalidUserID = "INVALID_USER_ID"
	ErrCodeInvalidRole   = "INVALID_ROLE"
)

// CanModifyUser checks if a requesting user can modify a target user
// Returns true if the user can modify (either owns the account or is admin)
func CanModifyUser(requestingUserID int, requestingUserRole string, targetUserID int) bool {
	// Validate input parameters
	if requestingUserID <= 0 || targetUserID <= 0 {
		return false
	}

	// Admin users can modify any user
	if ValidateAdminRole(requestingUserRole) {
		return true
	}

	// Users can modify their own accounts
	return requestingUserID == targetUserID
}

// ValidateUserOwnership validates that a user owns the target account
// Returns nil if valid, error if not
func ValidateUserOwnership(requestingUserID int, targetUserID int) error {
	// Validate user IDs are positive integers
	if requestingUserID <= 0 {
		return &AuthorizationError{
			Code:    ErrCodeInvalidUserID,
			Message: "Invalid requesting user ID",
			Details: fmt.Sprintf("User ID must be positive, got: %d", requestingUserID),
		}
	}

	if targetUserID <= 0 {
		return &AuthorizationError{
			Code:    ErrCodeInvalidUserID,
			Message: "Invalid target user ID",
			Details: fmt.Sprintf("User ID must be positive, got: %d", targetUserID),
		}
	}

	// Check ownership
	if requestingUserID != targetUserID {
		return &AuthorizationError{
			Code:    ErrForbidden,
			Message: "Forbidden: You can only modify your own account",
			Details: fmt.Sprintf("Requesting user %d cannot modify user %d", requestingUserID, targetUserID),
		}
	}

	return nil
}

// ValidateAdminRole checks if the provided role is admin
// Returns true if admin, false otherwise
func ValidateAdminRole(userRole string) bool {
	// Handle empty or invalid role
	if userRole == "" {
		return false
	}

	// Check for admin role (case-insensitive)
	return userRole == "admin" || userRole == "Admin" || userRole == "ADMIN"
}

// ValidateUserPermissions performs comprehensive authorization validation
// Combines ownership and role validation with detailed error reporting
func ValidateUserPermissions(requestingUserID int, requestingUserRole string, targetUserID int) error {
	// Validate input parameters
	if requestingUserID <= 0 {
		return &AuthorizationError{
			Code:    ErrCodeInvalidUserID,
			Message: "Invalid requesting user ID",
			Details: fmt.Sprintf("User ID must be positive, got: %d", requestingUserID),
		}
	}

	if targetUserID <= 0 {
		return &AuthorizationError{
			Code:    ErrCodeInvalidUserID,
			Message: "Invalid target user ID",
			Details: fmt.Sprintf("User ID must be positive, got: %d", targetUserID),
		}
	}

	// Admin users have full permissions
	if ValidateAdminRole(requestingUserRole) {
		return nil
	}

	// Regular users can only modify their own accounts
	if requestingUserID != targetUserID {
		return &AuthorizationError{
			Code:    ErrForbidden,
			Message: "Forbidden: You can only modify your own account",
			Details: fmt.Sprintf("User %d attempted to modify user %d", requestingUserID, targetUserID),
		}
	}

	return nil
}

// IsAuthorizationError checks if an error is an authorization error
func IsAuthorizationError(err error) bool {
	_, ok := err.(*AuthorizationError)
	return ok
}

// GetAuthorizationErrorCode extracts the error code from an authorization error
func GetAuthorizationErrorCode(err error) string {
	if authErr, ok := err.(*AuthorizationError); ok {
		return authErr.Code
	}
	return ""
}