package utils

import (
	"fmt"
	"testing"
)

// TestCanModifyUser tests the CanModifyUser function with various user ID combinations
func TestCanModifyUser(t *testing.T) {
	tests := []struct {
		name               string
		requestingUserID   int
		requestingUserRole string
		targetUserID       int
		expected           bool
		description        string
	}{
		// Valid ownership scenarios
		{
			name:               "User can modify own account",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       1,
			expected:           true,
			description:        "Regular user should be able to modify their own account",
		},
		{
			name:               "User cannot modify other user's account",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       2,
			expected:           false,
			description:        "Regular user should not be able to modify another user's account",
		},
		
		// Admin role scenarios
		{
			name:               "Admin can modify any user account",
			requestingUserID:   1,
			requestingUserRole: "admin",
			targetUserID:       2,
			expected:           true,
			description:        "Admin should be able to modify any user's account",
		},
		{
			name:               "Admin can modify own account",
			requestingUserID:   1,
			requestingUserRole: "admin",
			targetUserID:       1,
			expected:           true,
			description:        "Admin should be able to modify their own account",
		},
		{
			name:               "Admin with uppercase role can modify any user",
			requestingUserID:   1,
			requestingUserRole: "ADMIN",
			targetUserID:       2,
			expected:           true,
			description:        "Admin role should be case-insensitive",
		},
		{
			name:               "Admin with mixed case role can modify any user",
			requestingUserID:   1,
			requestingUserRole: "Admin",
			targetUserID:       2,
			expected:           true,
			description:        "Admin role should be case-insensitive",
		},
		
		// Invalid user ID scenarios
		{
			name:               "Invalid requesting user ID (zero)",
			requestingUserID:   0,
			requestingUserRole: "user",
			targetUserID:       1,
			expected:           false,
			description:        "Zero requesting user ID should be invalid",
		},
		{
			name:               "Invalid requesting user ID (negative)",
			requestingUserID:   -1,
			requestingUserRole: "user",
			targetUserID:       1,
			expected:           false,
			description:        "Negative requesting user ID should be invalid",
		},
		{
			name:               "Invalid target user ID (zero)",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       0,
			expected:           false,
			description:        "Zero target user ID should be invalid",
		},
		{
			name:               "Invalid target user ID (negative)",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       -1,
			expected:           false,
			description:        "Negative target user ID should be invalid",
		},
		{
			name:               "Both user IDs invalid",
			requestingUserID:   0,
			requestingUserRole: "user",
			targetUserID:       0,
			expected:           false,
			description:        "Both invalid user IDs should return false",
		},
		
		// Edge cases with roles
		{
			name:               "Empty role cannot modify other user",
			requestingUserID:   1,
			requestingUserRole: "",
			targetUserID:       2,
			expected:           false,
			description:        "Empty role should not grant admin privileges",
		},
		{
			name:               "Empty role can modify own account",
			requestingUserID:   1,
			requestingUserRole: "",
			targetUserID:       1,
			expected:           true,
			description:        "Empty role should still allow self-modification",
		},
		{
			name:               "Invalid role cannot modify other user",
			requestingUserID:   1,
			requestingUserRole: "invalid_role",
			targetUserID:       2,
			expected:           false,
			description:        "Invalid role should not grant admin privileges",
		},
		{
			name:               "Moderator role cannot modify other user",
			requestingUserID:   1,
			requestingUserRole: "moderator",
			targetUserID:       2,
			expected:           false,
			description:        "Non-admin roles should not grant cross-user privileges",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CanModifyUser(tt.requestingUserID, tt.requestingUserRole, tt.targetUserID)
			if result != tt.expected {
				t.Errorf("CanModifyUser(%d, %q, %d) = %v, expected %v. %s",
					tt.requestingUserID, tt.requestingUserRole, tt.targetUserID,
					result, tt.expected, tt.description)
			}
		})
	}
}

// TestValidateAdminRole tests admin role validation with different role values
func TestValidateAdminRole(t *testing.T) {
	tests := []struct {
		name     string
		userRole string
		expected bool
	}{
		// Valid admin roles
		{
			name:     "Lowercase admin role",
			userRole: "admin",
			expected: true,
		},
		{
			name:     "Uppercase admin role",
			userRole: "ADMIN",
			expected: true,
		},
		{
			name:     "Mixed case admin role",
			userRole: "Admin",
			expected: true,
		},
		
		// Invalid roles
		{
			name:     "Empty role",
			userRole: "",
			expected: false,
		},
		{
			name:     "User role",
			userRole: "user",
			expected: false,
		},
		{
			name:     "Moderator role",
			userRole: "moderator",
			expected: false,
		},
		{
			name:     "Invalid role",
			userRole: "invalid",
			expected: false,
		},
		{
			name:     "Admin with spaces",
			userRole: " admin ",
			expected: false,
		},
		{
			name:     "Admin substring",
			userRole: "administrator",
			expected: false,
		},
		{
			name:     "Numeric role",
			userRole: "123",
			expected: false,
		},
		{
			name:     "Special characters",
			userRole: "admin!",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateAdminRole(tt.userRole)
			if result != tt.expected {
				t.Errorf("ValidateAdminRole(%q) = %v, expected %v",
					tt.userRole, result, tt.expected)
			}
		})
	}
}

// TestValidateUserOwnership tests ownership validation with error handling
func TestValidateUserOwnership(t *testing.T) {
	tests := []struct {
		name             string
		requestingUserID int
		targetUserID     int
		expectError      bool
		expectedCode     string
		description      string
	}{
		// Valid ownership
		{
			name:             "Valid ownership - same user",
			requestingUserID: 1,
			targetUserID:     1,
			expectError:      false,
			expectedCode:     "",
			description:      "User should be able to validate ownership of their own account",
		},
		{
			name:             "Valid ownership - different valid users",
			requestingUserID: 5,
			targetUserID:     5,
			expectError:      false,
			expectedCode:     "",
			description:      "User should be able to validate ownership with same IDs",
		},
		
		// Invalid ownership
		{
			name:             "Invalid ownership - different users",
			requestingUserID: 1,
			targetUserID:     2,
			expectError:      true,
			expectedCode:     ErrForbidden,
			description:      "User should not be able to claim ownership of another user's account",
		},
		
		// Invalid requesting user ID
		{
			name:             "Invalid requesting user ID - zero",
			requestingUserID: 0,
			targetUserID:     1,
			expectError:      true,
			expectedCode:     ErrCodeInvalidUserID,
			description:      "Zero requesting user ID should return validation error",
		},
		{
			name:             "Invalid requesting user ID - negative",
			requestingUserID: -1,
			targetUserID:     1,
			expectError:      true,
			expectedCode:     ErrCodeInvalidUserID,
			description:      "Negative requesting user ID should return validation error",
		},
		{
			name:             "Invalid requesting user ID - large negative",
			requestingUserID: -999,
			targetUserID:     1,
			expectError:      true,
			expectedCode:     ErrCodeInvalidUserID,
			description:      "Large negative requesting user ID should return validation error",
		},
		
		// Invalid target user ID
		{
			name:             "Invalid target user ID - zero",
			requestingUserID: 1,
			targetUserID:     0,
			expectError:      true,
			expectedCode:     ErrCodeInvalidUserID,
			description:      "Zero target user ID should return validation error",
		},
		{
			name:             "Invalid target user ID - negative",
			requestingUserID: 1,
			targetUserID:     -1,
			expectError:      true,
			expectedCode:     ErrCodeInvalidUserID,
			description:      "Negative target user ID should return validation error",
		},
		{
			name:             "Invalid target user ID - large negative",
			requestingUserID: 1,
			targetUserID:     -999,
			expectError:      true,
			expectedCode:     ErrCodeInvalidUserID,
			description:      "Large negative target user ID should return validation error",
		},
		
		// Both IDs invalid
		{
			name:             "Both user IDs invalid - zero",
			requestingUserID: 0,
			targetUserID:     0,
			expectError:      true,
			expectedCode:     ErrCodeInvalidUserID,
			description:      "Both zero user IDs should return validation error for requesting user first",
		},
		{
			name:             "Both user IDs invalid - negative",
			requestingUserID: -1,
			targetUserID:     -2,
			expectError:      true,
			expectedCode:     ErrCodeInvalidUserID,
			description:      "Both negative user IDs should return validation error for requesting user first",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserOwnership(tt.requestingUserID, tt.targetUserID)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateUserOwnership(%d, %d) expected error but got nil. %s",
						tt.requestingUserID, tt.targetUserID, tt.description)
					return
				}
				
				// Check if it's an AuthorizationError
				authErr, ok := err.(*AuthorizationError)
				if !ok {
					t.Errorf("ValidateUserOwnership(%d, %d) expected AuthorizationError but got %T: %v",
						tt.requestingUserID, tt.targetUserID, err, err)
					return
				}
				
				// Check error code
				if authErr.Code != tt.expectedCode {
					t.Errorf("ValidateUserOwnership(%d, %d) expected error code %q but got %q",
						tt.requestingUserID, tt.targetUserID, tt.expectedCode, authErr.Code)
				}
				
				// Verify error message is not empty
				if authErr.Message == "" {
					t.Errorf("ValidateUserOwnership(%d, %d) error message should not be empty",
						tt.requestingUserID, tt.targetUserID)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateUserOwnership(%d, %d) expected no error but got: %v. %s",
						tt.requestingUserID, tt.targetUserID, err, tt.description)
				}
			}
		})
	}
}

// TestValidateUserPermissions tests comprehensive authorization validation
func TestValidateUserPermissions(t *testing.T) {
	tests := []struct {
		name               string
		requestingUserID   int
		requestingUserRole string
		targetUserID       int
		expectError        bool
		expectedCode       string
		description        string
	}{
		// Valid scenarios
		{
			name:               "Valid self-modification",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       1,
			expectError:        false,
			expectedCode:       "",
			description:        "User should be able to modify their own account",
		},
		{
			name:               "Valid admin modification of other user",
			requestingUserID:   1,
			requestingUserRole: "admin",
			targetUserID:       2,
			expectError:        false,
			expectedCode:       "",
			description:        "Admin should be able to modify any user's account",
		},
		{
			name:               "Valid admin self-modification",
			requestingUserID:   1,
			requestingUserRole: "admin",
			targetUserID:       1,
			expectError:        false,
			expectedCode:       "",
			description:        "Admin should be able to modify their own account",
		},
		{
			name:               "Valid admin with uppercase role",
			requestingUserID:   1,
			requestingUserRole: "ADMIN",
			targetUserID:       2,
			expectError:        false,
			expectedCode:       "",
			description:        "Admin role should be case-insensitive",
		},
		{
			name:               "Valid admin with mixed case role",
			requestingUserID:   1,
			requestingUserRole: "Admin",
			targetUserID:       2,
			expectError:        false,
			expectedCode:       "",
			description:        "Admin role should be case-insensitive",
		},
		
		// Invalid user ID scenarios
		{
			name:               "Invalid requesting user ID - zero",
			requestingUserID:   0,
			requestingUserRole: "user",
			targetUserID:       1,
			expectError:        true,
			expectedCode:       ErrCodeInvalidUserID,
			description:        "Zero requesting user ID should return validation error",
		},
		{
			name:               "Invalid requesting user ID - negative",
			requestingUserID:   -1,
			requestingUserRole: "user",
			targetUserID:       1,
			expectError:        true,
			expectedCode:       ErrCodeInvalidUserID,
			description:        "Negative requesting user ID should return validation error",
		},
		{
			name:               "Invalid target user ID - zero",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       0,
			expectError:        true,
			expectedCode:       ErrCodeInvalidUserID,
			description:        "Zero target user ID should return validation error",
		},
		{
			name:               "Invalid target user ID - negative",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       -1,
			expectError:        true,
			expectedCode:       ErrCodeInvalidUserID,
			description:        "Negative target user ID should return validation error",
		},
		
		// Authorization failures
		{
			name:               "Regular user cannot modify other user",
			requestingUserID:   1,
			requestingUserRole: "user",
			targetUserID:       2,
			expectError:        true,
			expectedCode:       ErrForbidden,
			description:        "Regular user should not be able to modify another user's account",
		},
		{
			name:               "Empty role cannot modify other user",
			requestingUserID:   1,
			requestingUserRole: "",
			targetUserID:       2,
			expectError:        true,
			expectedCode:       ErrForbidden,
			description:        "Empty role should not grant admin privileges",
		},
		{
			name:               "Invalid role cannot modify other user",
			requestingUserID:   1,
			requestingUserRole: "invalid_role",
			targetUserID:       2,
			expectError:        true,
			expectedCode:       ErrForbidden,
			description:        "Invalid role should not grant admin privileges",
		},
		{
			name:               "Moderator role cannot modify other user",
			requestingUserID:   1,
			requestingUserRole: "moderator",
			targetUserID:       2,
			expectError:        true,
			expectedCode:       ErrForbidden,
			description:        "Non-admin roles should not grant cross-user privileges",
		},
		
		// Edge cases
		{
			name:               "Both IDs invalid - requesting user checked first",
			requestingUserID:   0,
			requestingUserRole: "user",
			targetUserID:       0,
			expectError:        true,
			expectedCode:       ErrCodeInvalidUserID,
			description:        "Both invalid IDs should return error for requesting user first",
		},
		{
			name:               "Admin with invalid requesting ID",
			requestingUserID:   0,
			requestingUserRole: "admin",
			targetUserID:       1,
			expectError:        true,
			expectedCode:       ErrCodeInvalidUserID,
			description:        "Even admin role should not bypass user ID validation",
		},
		{
			name:               "Admin with invalid target ID",
			requestingUserID:   1,
			requestingUserRole: "admin",
			targetUserID:       0,
			expectError:        true,
			expectedCode:       ErrCodeInvalidUserID,
			description:        "Even admin role should not bypass target ID validation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUserPermissions(tt.requestingUserID, tt.requestingUserRole, tt.targetUserID)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("ValidateUserPermissions(%d, %q, %d) expected error but got nil. %s",
						tt.requestingUserID, tt.requestingUserRole, tt.targetUserID, tt.description)
					return
				}
				
				// Check if it's an AuthorizationError
				authErr, ok := err.(*AuthorizationError)
				if !ok {
					t.Errorf("ValidateUserPermissions(%d, %q, %d) expected AuthorizationError but got %T: %v",
						tt.requestingUserID, tt.requestingUserRole, tt.targetUserID, err, err)
					return
				}
				
				// Check error code
				if authErr.Code != tt.expectedCode {
					t.Errorf("ValidateUserPermissions(%d, %q, %d) expected error code %q but got %q",
						tt.requestingUserID, tt.requestingUserRole, tt.targetUserID, tt.expectedCode, authErr.Code)
				}
				
				// Verify error message is not empty
				if authErr.Message == "" {
					t.Errorf("ValidateUserPermissions(%d, %q, %d) error message should not be empty",
						tt.requestingUserID, tt.requestingUserRole, tt.targetUserID)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateUserPermissions(%d, %q, %d) expected no error but got: %v. %s",
						tt.requestingUserID, tt.requestingUserRole, tt.targetUserID, err, tt.description)
				}
			}
		})
	}
}

// TestIsAuthorizationError tests the error type checking utility
func TestIsAuthorizationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name: "Valid AuthorizationError",
			err: &AuthorizationError{
				Code:    ErrForbidden,
				Message: "Test error",
			},
			expected: true,
		},
		{
			name:     "Nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "Generic error",
			err:      fmt.Errorf("generic error"),
			expected: false,
		},
		{
			name: "AppError from errors.go",
			err: &AppError{
				Code:    ErrForbidden,
				Message: "App error",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAuthorizationError(tt.err)
			if result != tt.expected {
				t.Errorf("IsAuthorizationError(%v) = %v, expected %v",
					tt.err, result, tt.expected)
			}
		})
	}
}

// TestGetAuthorizationErrorCode tests error code extraction utility
func TestGetAuthorizationErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name: "Valid AuthorizationError with code",
			err: &AuthorizationError{
				Code:    ErrForbidden,
				Message: "Test error",
			},
			expected: ErrForbidden,
		},
		{
			name: "Valid AuthorizationError with different code",
			err: &AuthorizationError{
				Code:    ErrCodeInvalidUserID,
				Message: "Invalid user ID",
			},
			expected: ErrCodeInvalidUserID,
		},
		{
			name:     "Nil error",
			err:      nil,
			expected: "",
		},
		{
			name:     "Generic error",
			err:      fmt.Errorf("generic error"),
			expected: "",
		},
		{
			name: "AppError from errors.go",
			err: &AppError{
				Code:    ErrForbidden,
				Message: "App error",
			},
			expected: "",
		},
		{
			name: "AuthorizationError with empty code",
			err: &AuthorizationError{
				Code:    "",
				Message: "Test error",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetAuthorizationErrorCode(tt.err)
			if result != tt.expected {
				t.Errorf("GetAuthorizationErrorCode(%v) = %q, expected %q",
					tt.err, result, tt.expected)
			}
		})
	}
}

// TestAuthorizationErrorString tests the Error() method of AuthorizationError
func TestAuthorizationErrorString(t *testing.T) {
	tests := []struct {
		name     string
		err      *AuthorizationError
		expected string
	}{
		{
			name: "Error with details",
			err: &AuthorizationError{
				Code:    ErrForbidden,
				Message: "Access denied",
				Details: "User 1 cannot modify user 2",
			},
			expected: "FORBIDDEN: Access denied (User 1 cannot modify user 2)",
		},
		{
			name: "Error without details",
			err: &AuthorizationError{
				Code:    ErrCodeInvalidUserID,
				Message: "Invalid user ID",
				Details: "",
			},
			expected: "INVALID_USER_ID: Invalid user ID",
		},
		{
			name: "Error with empty message",
			err: &AuthorizationError{
				Code:    ErrForbidden,
				Message: "",
				Details: "Some details",
			},
			expected: "FORBIDDEN:  (Some details)",
		},
		{
			name: "Error with empty code",
			err: &AuthorizationError{
				Code:    "",
				Message: "Some message",
				Details: "",
			},
			expected: ": Some message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("AuthorizationError.Error() = %q, expected %q",
					result, tt.expected)
			}
		})
	}
}

// TestEdgeCasesAndBoundaryConditions tests edge cases and boundary conditions
func TestEdgeCasesAndBoundaryConditions(t *testing.T) {
	t.Run("Large user IDs", func(t *testing.T) {
		largeID := 2147483647 // Max int32
		result := CanModifyUser(largeID, "user", largeID)
		if !result {
			t.Errorf("CanModifyUser should handle large user IDs correctly")
		}
	})

	t.Run("Role with special characters", func(t *testing.T) {
		result := ValidateAdminRole("admin@#$")
		if result {
			t.Errorf("ValidateAdminRole should not accept roles with special characters")
		}
	})

	t.Run("Very long role string", func(t *testing.T) {
		longRole := string(make([]byte, 1000))
		for i := range longRole {
			longRole = longRole[:i] + "a" + longRole[i+1:]
		}
		result := ValidateAdminRole(longRole)
		if result {
			t.Errorf("ValidateAdminRole should not accept very long role strings")
		}
	})

	t.Run("Unicode characters in role", func(t *testing.T) {
		result := ValidateAdminRole("admin🔒")
		if result {
			t.Errorf("ValidateAdminRole should not accept roles with unicode characters")
		}
	})

	t.Run("Whitespace-only role", func(t *testing.T) {
		result := ValidateAdminRole("   ")
		if result {
			t.Errorf("ValidateAdminRole should not accept whitespace-only roles")
		}
	})

	t.Run("Tab and newline characters in role", func(t *testing.T) {
		result := ValidateAdminRole("admin\t\n")
		if result {
			t.Errorf("ValidateAdminRole should not accept roles with tab/newline characters")
		}
	})
}