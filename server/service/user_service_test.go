package service

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/golang-jwt/jwt/v5"
)

// Mock UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateNewUser(user model.User) (model.User, error) {
	args := m.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(id int) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) GetAllUser() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (model.User, error) {
	args := m.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserRepository) SaveRefreshToken(userID int, token string, expiresAt time.Time) error {
	args := m.Called(userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockUserRepository) FindRefreshToken(token string) (model.RefreshToken, error) {
	args := m.Called(token)
	return args.Get(0).(model.RefreshToken), args.Error(1)
}

func (m *MockUserRepository) UpdateRefreshToken(oldToken, newToken string, expiresAt time.Time) error {
	args := m.Called(oldToken, newToken, expiresAt)
	return args.Error(0)
}

func (m *MockUserRepository) ValidateRefreshToken(token string) (int, error) {
	args := m.Called(token)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) DeleteRefreshToken(token string) error {
	args := m.Called(token)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteAllRefreshTOkensByUser(userId int) error {
	args := m.Called(userId)
	return args.Error(0)
}

// Mock JwtService
type MockJwtService struct {
	mock.Mock
}

func (m *MockJwtService) GenerateToken(user model.User) (dto.LoginResponseDto, error) {
	args := m.Called(user)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

func (m *MockJwtService) GenerateRefreshToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockJwtService) ValidateToken(token string) (*dto.JwtTokenClaims, error) {
	args := m.Called(token)
	return args.Get(0).(*dto.JwtTokenClaims), args.Error(1)
}

func (m *MockJwtService) VerifyToken(tokenString string) (jwt.MapClaims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(jwt.MapClaims), args.Error(1)
}

// Mock PasswordHasher
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

func TestCreateNewUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJwtService := new(MockJwtService)
	mockPasswordHasher := new(MockPasswordHasher)
	userService := NewUserservice(mockRepo, mockJwtService, mockPasswordHasher)

	user := model.User{
		Name: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Test case 1: Successful user creation
	mockPasswordHasher.On("EncryptPassword", user.Password).Return("hashedpassword", nil).Once()
	mockRepo.On("CreateNewUser", mock.AnythingOfType("model.User")).Return(model.User{Id: 1, Name: "testuser", Email: "test@example.com", Password: "hashedpassword"}, nil).Once()
	
	createdUser, err := userService.CreateNewUser(user)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", createdUser.Name)
	assert.Equal(t, "-", createdUser.Password) // Password should be masked
	mockRepo.AssertExpectations(t)
	mockPasswordHasher.AssertExpectations(t)

	// Test case 2: Error during password encryption
	mockPasswordHasher.On("EncryptPassword", user.Password).Return("", errors.New("encryption error")).Once()

	_, err = userService.CreateNewUser(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to encrypt password")
	mockPasswordHasher.AssertExpectations(t)

	// Test case 3: Error during user creation in repository
	mockPasswordHasher.On("EncryptPassword", user.Password).Return("hashedpassword", nil).Once()
	mockRepo.On("CreateNewUser", mock.AnythingOfType("model.User")).Return(model.User{}, errors.New("db error")).Once()
	
	_, err = userService.CreateNewUser(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	mockRepo.AssertExpectations(t)
	mockPasswordHasher.AssertExpectations(t)
}

func TestFindUserById(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJwtService := new(MockJwtService)
	mockPasswordHasher := new(MockPasswordHasher)
	userService := NewUserservice(mockRepo, mockJwtService, mockPasswordHasher)

	// Test case 1: Successful user retrieval
	expectedUser := model.User{Id: 1, Name: "testuser"}
	mockRepo.On("GetUserById", 1).Return(expectedUser, nil).Once()

	user, err := userService.FindUserById("1")
	assert.NoError(t, err)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)

	// Test case 2: Invalid ID format
	_, err = userService.FindUserById("abc")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "strconv.Atoi: parsing \"abc\": invalid syntax")

	// Test case 3: User not found in repository
	mockRepo.On("GetUserById", 2).Return(model.User{}, errors.New("user not found")).Once()

	_, err = userService.FindUserById("2")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)
}

func TestFindAllUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJwtService := new(MockJwtService)
	mockPasswordHasher := new(MockPasswordHasher)
	userService := NewUserservice(mockRepo, mockJwtService, mockPasswordHasher)

	// Test case 1: Successful retrieval of all users
	expectedUsers := []model.User{{Id: 1, Name: "user1"}, {Id: 2, Name: "user2"}}
	mockRepo.On("GetAllUser").Return(expectedUsers, nil).Once()

	users, err := userService.FindAllUser()
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
	mockRepo.AssertExpectations(t)

	// Test case 2: Error during retrieval
	mockRepo.On("GetAllUser").Return([]model.User{}, errors.New("db error")).Once()

	_, err = userService.FindAllUser()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")
	mockRepo.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJwtService := new(MockJwtService)
	mockPasswordHasher := new(MockPasswordHasher)
	userService := NewUserservice(mockRepo, mockJwtService, mockPasswordHasher)

	loginPayload := dto.LoginDto{Identifier: "test@example.com", Password: "password123"}
	user := model.User{Id: 1, Email: "test@example.com", Password: "hashedpassword"}
	loginResponse := dto.LoginResponseDto{AccessToken: "access", RefreshToken: "refresh"}

	// Test case 1: Successful login
	mockRepo.On("GetByEmail", loginPayload.Identifier).Return(user, nil).Once()
	mockPasswordHasher.On("ComparePasswordHash", user.Password, loginPayload.Password).Return(nil).Once()
	mockJwtService.On("GenerateToken", mock.AnythingOfType("model.User")).Return(loginResponse, nil).Once()
	mockRepo.On("SaveRefreshToken", user.Id, loginResponse.RefreshToken, mock.AnythingOfType("time.Time")).Return(nil).Once()

	resp, err := userService.Login(loginPayload)
	assert.NoError(t, err)
	assert.Equal(t, loginResponse, resp)
	mockRepo.AssertExpectations(t)
	mockJwtService.AssertExpectations(t)
	mockPasswordHasher.AssertExpectations(t)

	// Test case 2: User not found
	mockRepo.On("GetByEmail", loginPayload.Identifier).Return(model.User{}, errors.New("not found")).Once()

	_, err = userService.Login(loginPayload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
	mockRepo.AssertExpectations(t)

	// Test case 3: Invalid password
	mockRepo.On("GetByEmail", loginPayload.Identifier).Return(user, nil).Once()
	mockPasswordHasher.On("ComparePasswordHash", user.Password, loginPayload.Password).Return(errors.New("password mismatch")).Once()

	_, err = userService.Login(loginPayload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
	mockRepo.AssertExpectations(t)
	mockPasswordHasher.AssertExpectations(t)

	// Test case 4: Error generating token
	mockRepo.On("GetByEmail", loginPayload.Identifier).Return(user, nil).Once()
	mockPasswordHasher.On("ComparePasswordHash", user.Password, loginPayload.Password).Return(nil).Once()
	mockJwtService.On("GenerateToken", mock.AnythingOfType("model.User")).Return(dto.LoginResponseDto{}, errors.New("token error")).Once()

	_, err = userService.Login(loginPayload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create token")
	mockRepo.AssertExpectations(t)
	mockJwtService.AssertExpectations(t)
	mockPasswordHasher.AssertExpectations(t)

	// Test case 5: Error saving refresh token
	mockRepo.On("GetByEmail", loginPayload.Identifier).Return(user, nil).Once()
	mockPasswordHasher.On("ComparePasswordHash", user.Password, loginPayload.Password).Return(nil).Once()
	mockJwtService.On("GenerateToken", mock.AnythingOfType("model.User")).Return(loginResponse, nil).Once()
	mockRepo.On("SaveRefreshToken", user.Id, loginResponse.RefreshToken, mock.AnythingOfType("time.Time")).Return(errors.New("db error")).Once()

	_, err = userService.Login(loginPayload)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "internal server error")
	mockRepo.AssertExpectations(t)
	mockJwtService.AssertExpectations(t)
	mockPasswordHasher.AssertExpectations(t)
}

func TestRefreshToken(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJwtService := new(MockJwtService)
	mockPasswordHasher := new(MockPasswordHasher)
	userService := NewUserservice(mockRepo, mockJwtService, mockPasswordHasher)

	refreshToken := "valid_refresh_token"
	user := model.User{Id: 1, Name: "testuser"}
	newLoginResponse := dto.LoginResponseDto{AccessToken: "new_access", RefreshToken: "new_refresh"}

	// Test case 1: Successful refresh
	mockRepo.On("FindRefreshToken", refreshToken).Return(model.RefreshToken{UserID: 1, Token: refreshToken, ExpiresAt: time.Now().Add(time.Hour)}, nil).Once()
	mockRepo.On("GetUserById", 1).Return(user, nil).Once()
	mockJwtService.On("GenerateToken", user).Return(newLoginResponse, nil).Once()
	mockRepo.On("UpdateRefreshToken", refreshToken, newLoginResponse.RefreshToken, mock.AnythingOfType("time.Time")).Return(nil).Once()

	resp, err := userService.RefreshToken(refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, newLoginResponse, resp)
	mockRepo.AssertExpectations(t)
	mockJwtService.AssertExpectations(t)

	// Test case 2: Invalid refresh token format (handled by service, no need to mock url.QueryUnescape)
	// The service itself calls url.QueryUnescape, so we just pass an invalid token
	_, err = userService.RefreshToken("invalid%")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid refresh token format")

	// Test case 3: Refresh token not found or expired
	mockRepo.On("FindRefreshToken", refreshToken).Return(model.RefreshToken{}, errors.New("not found")).Once()
	_, err = userService.RefreshToken(refreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	mockRepo.AssertExpectations(t)

	// Test case 4: Refresh token expired
	mockRepo.On("FindRefreshToken", refreshToken).Return(model.RefreshToken{UserID: 1, Token: refreshToken, ExpiresAt: time.Now().Add(-time.Hour)}, errors.New("token is expired")).Once()
	_, err = userService.RefreshToken(refreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token is expired") // Assuming the error message from FindRefreshToken for expired token
	mockRepo.AssertExpectations(t)

	// Test case 5: User not found for refresh token
	mockRepo.On("FindRefreshToken", refreshToken).Return(model.RefreshToken{UserID: 1, Token: refreshToken, ExpiresAt: time.Now().Add(time.Hour)}, nil).Once()
	mockRepo.On("GetUserById", 1).Return(model.User{}, errors.New("user not found")).Once()
	_, err = userService.RefreshToken(refreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
	mockRepo.AssertExpectations(t)

	// Test case 6: Error generating new token
	mockRepo.On("FindRefreshToken", refreshToken).Return(model.RefreshToken{UserID: 1, Token: refreshToken, ExpiresAt: time.Now().Add(time.Hour)}, nil).Once()
	mockRepo.On("GetUserById", 1).Return(user, nil).Once()
	mockJwtService.On("GenerateToken", user).Return(dto.LoginResponseDto{}, errors.New("generate error")).Once()
	_, err = userService.RefreshToken(refreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to generate new token")
	mockRepo.AssertExpectations(t)
	mockJwtService.AssertExpectations(t)

	// Test case 7: Error updating refresh token
	mockRepo.On("FindRefreshToken", refreshToken).Return(model.RefreshToken{UserID: 1, Token: refreshToken, ExpiresAt: time.Now().Add(time.Hour)}, nil).Once()
	mockRepo.On("GetUserById", 1).Return(user, nil).Once()
	mockJwtService.On("GenerateToken", user).Return(newLoginResponse, nil).Once()
	mockRepo.On("UpdateRefreshToken", refreshToken, newLoginResponse.RefreshToken, mock.AnythingOfType("time.Time")).Return(errors.New("update error")).Once()
	_, err = userService.RefreshToken(refreshToken)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update refresh token")
	mockRepo.AssertExpectations(t)
	mockJwtService.AssertExpectations(t)
}
