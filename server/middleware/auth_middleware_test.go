package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"develapar-server/middleware"
	"develapar-server/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthMiddlewareTestSuite struct {
	suite.Suite
	jwtService *service.JwtServiceMock
	router     *gin.Engine
}

func (suite *AuthMiddlewareTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.jwtService = new(service.JwtServiceMock)
	authMiddleware := middleware.NewAuthMiddleware(suite.jwtService)
	suite.router = gin.New()
	suite.router.GET("/test", authMiddleware.CheckToken(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	suite.router.GET("/admin", authMiddleware.CheckToken("admin"), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	suite.jwtService.On("VerifyToken", "").Return(jwt.MapClaims{}, assert.AnError)
}

func (suite *AuthMiddlewareTestSuite) TestCheckToken_Success() {
	suite.jwtService.On("VerifyToken", "valid_token").Return(jwt.MapClaims{"userId": float64(1), "role": "user"}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer valid_token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *AuthMiddlewareTestSuite) TestCheckToken_NoToken() {
	suite.jwtService.On("VerifyToken", "").Return(nil, assert.AnError)

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *AuthMiddlewareTestSuite) TestCheckToken_InvalidToken() {
	suite.jwtService.On("VerifyToken", "invalid_token").Return(jwt.MapClaims{}, assert.AnError)

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid_token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *AuthMiddlewareTestSuite) TestCheckToken_Forbidden() {
	suite.jwtService.On("VerifyToken", "user_token").Return(jwt.MapClaims{"userId": float64(1), "role": "user"}, nil)

	req, _ := http.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer user_token")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusForbidden, w.Code)
}

func TestAuthMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}
