package utils

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPasswordHasher for testing
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) EncryptPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) ComparePasswordHash(passwordHash string, plainPassword string) error {
	args := m.Called(passwordHash, plainPassword)
	return args.Error(0)
}

func TestEncryptPassword(t *testing.T) {
	ph := NewPasswordHasher() // Use the real implementation for testing the function itself

	hashedPassword, err := ph.EncryptPassword("password")
	assert.NoError(t, err)
	assert.NotEmpty(t, hashedPassword)

	// Test case for error during encryption (not directly testable with bcrypt, but for completeness)
	// This part is more for demonstrating how you would mock an error if EncryptPassword had external dependencies
	mockPh := new(MockPasswordHasher)
	mockPh.On("EncryptPassword", "password").Return("", errors.New("encryption failed")).Once()
	_, err = mockPh.EncryptPassword("password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "encryption failed")
	mockPh.AssertExpectations(t)
}

func TestComparePasswordHash(t *testing.T) {
	ph := NewPasswordHasher() // Use the real implementation for testing the function itself

	hashedPassword, _ := ph.EncryptPassword("password")

	err := ph.ComparePasswordHash(hashedPassword, "password")
	assert.NoError(t, err)

	err = ph.ComparePasswordHash(hashedPassword, "wrongpassword")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "crypto/bcrypt: hashedPassword is not the hash of the given password")

	// Test case for error during comparison (not directly testable with bcrypt, but for completeness)
	mockPh := new(MockPasswordHasher)
	mockPh.On("ComparePasswordHash", "invalidhash", "password").Return(errors.New("comparison failed")).Once()
	err = mockPh.ComparePasswordHash("invalidhash", "password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "comparison failed")
	mockPh.AssertExpectations(t)
}