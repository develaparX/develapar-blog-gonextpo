package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"develapar-server/model"
	"develapar-server/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService adalah mock untuk service.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateNewUser(payload model.User) (model.User, error) {
	args := m.Called(payload)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserService) FindUserById(id string) (model.User, error) {
	args := m.Called(id)
	return args.Get(0).(model.User), args.Error(1)
}

func (m *MockUserService) FindAllUser() ([]model.User, error) {
	args := m.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (m *MockUserService) Login(payload dto.LoginDto) (dto.LoginResponseDto, error) {
	args := m.Called(payload)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

func (m *MockUserService) RefreshToken(refreshToken string) (dto.LoginResponseDto, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(dto.LoginResponseDto), args.Error(1)
}

func TestLoginHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(MockUserService)
	router := gin.Default()
	userController := NewUserController(mockUserService, router.Group("/api/v1"))
	userController.Route()

	t.Run("Success Login", func(t *testing.T) {
		loginPayload := dto.LoginDto{
			Identifier: "test@example.com",
			Password:   "password123",
		}
		loginResponse := dto.LoginResponseDto{
			AccessToken:  "mockAccessToken",
			RefreshToken: "mockRefreshToken",
		}

		mockUserService.On("Login", loginPayload).Return(loginResponse, nil).Once()

		body, _ := json.Marshal(loginPayload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Success Login", responseBody["message"])
		assert.Equal(t, "mockAccessToken", responseBody["accessToken"])
		assert.NotEmpty(t, rr.Header().Get("Set-Cookie"))

		mockUserService.AssertExpectations(t)
	})

	t.Run("Invalid Credentials", func(t *testing.T) {
		loginPayload := dto.LoginDto{
			Identifier: "test@example.com",
			Password:   "wrongpassword",
		}

		mockUserService.On("Login", loginPayload).Return(dto.LoginResponseDto{}, errors.New("invalid credentials")).Once()

		body, _ := json.Marshal(loginPayload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "invalid credentials", responseBody["error"])

		mockUserService.AssertExpectations(t)
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["error"], "invalid character")

		mockUserService.AssertNotCalled(t, "Login")
	})
}

func TestRegisterUser(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(MockUserService)
	router := gin.Default()
	userController := NewUserController(mockUserService, router.Group("/api/v1"))
	userController.Route()

	t.Run("Success Register User", func(t *testing.T) {
		userPayload := model.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		createdUser := userPayload
		createdUser.Id = 1
		createdUser.Password = "-"

		mockUserService.On("CreateNewUser", mock.AnythingOfType("model.User")).Return(createdUser, nil).Once()

		body, _ := json.Marshal(userPayload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Success Create New User", responseBody["message"])
		data := responseBody["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, "Test User", data["name"])
		assert.Equal(t, "test@example.com", data["email"])
		assert.Equal(t, "-", data["password"])

		mockUserService.AssertExpectations(t)
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["message"], "invalid character")

		mockUserService.AssertNotCalled(t, "CreateNewUser")
	})

	t.Run("Internal Server Error - CreateNewUser fails", func(t *testing.T) {
		userPayload := model.User{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}

		mockUserService.On("CreateNewUser", mock.AnythingOfType("model.User")).Return(model.User{}, errors.New("database error")).Once()

		body, _ := json.Marshal(userPayload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "database error", responseBody["message"])

		mockUserService.AssertExpectations(t)
	})
}

func TestFindUserByIdHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(MockUserService)
	router := gin.Default()
	userController := NewUserController(mockUserService, router.Group("/api/v1"))
	userController.Route()

	t.Run("Success Find User By ID", func(t *testing.T) {
		userID := "1"
		expectedUser := model.User{
			Id:    1,
			Name:  "Test User",
			Email: "test@example.com",
		}

		mockUserService.On("FindUserById", userID).Return(expectedUser, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/"+userID, nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Success Get User", responseBody["message"])
		data := responseBody["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, "Test User", data["name"])
		assert.Equal(t, "test@example.com", data["email"])

		mockUserService.AssertExpectations(t)
	})

	t.Run("Internal Server Error - FindUserById fails", func(t *testing.T) {
		userID := "1"

		mockUserService.On("FindUserById", userID).Return(model.User{}, errors.New("user not found")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/"+userID, nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "user not found", responseBody["error"])

		mockUserService.AssertExpectations(t)
	})
}

func TestFindAllUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(MockUserService)
	router := gin.Default()
	userController := NewUserController(mockUserService, router.Group("/api/v1"))
	userController.Route()

	t.Run("Success Find All Users", func(t *testing.T) {
		expectedUsers := []model.User{
			{Id: 1, Name: "User 1", Email: "user1@example.com"},
			{Id: 2, Name: "User 2", Email: "user2@example.com"},
		}

		mockUserService.On("FindAllUser").Return(expectedUsers, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Success Get All User", responseBody["message"])
		data := responseBody["data"].([]interface{})
		assert.Len(t, data, 2)

		mockUserService.AssertExpectations(t)
	})

	t.Run("Internal Server Error - FindAllUser fails", func(t *testing.T) {
		mockUserService.On("FindAllUser").Return([]model.User{}, errors.New("database error")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "database error", responseBody["error"])

		mockUserService.AssertExpectations(t)
	})
}

func TestRefreshTokenHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUserService := new(MockUserService)
	router := gin.Default()
	router.POST("/api/v1/auth/refresh", RefreshTokenHandler(mockUserService))

	t.Run("Success Refresh Token", func(t *testing.T) {
		refreshToken := "validRefreshToken"
		loginResponse := dto.LoginResponseDto{
			AccessToken:  "newAccessToken",
			RefreshToken: "newRefreshToken",
		}

		mockUserService.On("RefreshToken", refreshToken).Return(loginResponse, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "newAccessToken", responseBody["access_token"])
		assert.NotEmpty(t, rr.Header().Get("Set-Cookie"))

		mockUserService.AssertExpectations(t)
	})

	t.Run("Refresh Token Not Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "refresh token not found in cookies", responseBody["error"])

		mockUserService.AssertNotCalled(t, "RefreshToken")
	})

	t.Run("Refresh Token Invalid", func(t *testing.T) {
		refreshToken := "invalidRefreshToken"
		mockUserService.On("RefreshToken", refreshToken).Return(dto.LoginResponseDto{}, errors.New("invalid refresh token")).Once()

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refreshToken", Value: refreshToken})
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "invalid refresh token", responseBody["error"])

		mockUserService.AssertExpectations(t)
	})
}
