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

// MockArticleService adalah mock untuk service.ArticleService
type MockArticleService struct {
	mock.Mock
}

func (m *MockArticleService) CreateArticle(payload model.Article) (model.Article, error) {
	args := m.Called(payload)
	return args.Get(0).(model.Article), args.Error(1)
}

func (m *MockArticleService) CreateArticleWithTags(req dto.CreateArticleRequest, userID int) (model.Article, error) {
	args := m.Called(req, userID)
	return args.Get(0).(model.Article), args.Error(1)
}

func (m *MockArticleService) FindAll() ([]model.Article, error) {
	args := m.Called()
	return args.Get(0).([]model.Article), args.Error(1)
}

func (m *MockArticleService) FindById(id int) (model.Article, error) {
	args := m.Called(id)
	return args.Get(0).(model.Article), args.Error(1)
}

func (m *MockArticleService) UpdateArticle(id int, payload dto.UpdateArticleRequest) (model.Article, error) {
	args := m.Called(id, payload)
	return args.Get(0).(model.Article), args.Error(1)
}

func (m *MockArticleService) FindBySlug(slug string) (model.Article, error) {
	args := m.Called(slug)
	return args.Get(0).(model.Article), args.Error(1)
}

func (m *MockArticleService) FindByUserId(userId int) ([]model.Article, error) {
	args := m.Called(userId)
	return args.Get(0).([]model.Article), args.Error(1)
}

func (m *MockArticleService) FindByCategory(categoryName string) ([]model.Article, error) {
	args := m.Called(categoryName)
	return args.Get(0).([]model.Article), args.Error(1)
}

func (m *MockArticleService) DeleteArticle(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// MockAuthMiddleware adalah mock untuk middleware.AuthMiddleware
type MockAuthMiddleware struct {
	mock.Mock
}

func (m *MockAuthMiddleware) CheckToken(roles ...string) gin.HandlerFunc {
	args := m.Called(roles)
	return args.Get(0).(gin.HandlerFunc)
}

// DummyAuthMiddleware for public routes
type DummyAuthMiddleware struct{}

func (m *DummyAuthMiddleware) CheckToken(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

func TestCreateArticleHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Create Article", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.Set("userId", float64(1))
			c.Next()
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articlePayload := dto.CreateArticleRequest{
			Title:      "Test Article",
			Content:    "This is a test article content.",
			CategoryID: 1,
			Tags:       []string{"golang", "test"},
		}
		createdArticle := model.Article{
			Id:      1,
			Title:   "Test Article",
			Content: "This is a test article content.",
			User:    model.User{Id: 1},
		}

		mockArticleService.On("CreateArticleWithTags", mock.AnythingOfType("dto.CreateArticleRequest"), 1).Return(createdArticle, nil).Once()

		body, _ := json.Marshal(articlePayload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/article/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Success create article with tags", responseBody["message"])
		data := responseBody["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, "Test Article", data["title"])
		assert.Equal(t, "This is a test article content.", data["content"])
		assert.Equal(t, float64(1), data["user"].(map[string]interface{})["id"])

		mockArticleService.AssertExpectations(t)
		mockAuthMiddleware.AssertExpectations(t)
	})

	t.Run("Unauthorized - No userId in context", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articlePayload := dto.CreateArticleRequest{
			Title:      "Test Article",
			Content:    "This is a test article content.",
			CategoryID: 1,
		}
		body, _ := json.Marshal(articlePayload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/article/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Unauthorized", responseBody["message"])

		mockArticleService.AssertNotCalled(t, "CreateArticleWithTags")
		mockAuthMiddleware.AssertExpectations(t)
	})

	t.Run("Bad Request - Invalid JSON", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.Set("userId", float64(1))
			c.Next()
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		req, _ := http.NewRequest(http.MethodPost, "/api/v1/article/", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["message"], "Invalid payload")

		mockArticleService.AssertNotCalled(t, "CreateArticleWithTags")
		mockAuthMiddleware.AssertExpectations(t)
	})

	t.Run("Internal Server Error - CreateArticle fails", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.Set("userId", float64(1))
			c.Next()
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articlePayload := dto.CreateArticleRequest{
			Title:      "Test Article",
			Content:    "This is a test article content.",
			CategoryID: 1,
		}
		mockArticleService.On("CreateArticleWithTags", mock.AnythingOfType("dto.CreateArticleRequest"), 1).Return(model.Article{}, errors.New("database error")).Once()

		body, _ := json.Marshal(articlePayload)
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/article/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["message"], "Failed to create article")

		mockArticleService.AssertExpectations(t)
		mockAuthMiddleware.AssertExpectations(t)
	})
}

func TestGetAllArticleHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Get All Articles", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		dummyAuthMiddleware := new(DummyAuthMiddleware) // Use DummyAuthMiddleware

		articleController := NewArticleController(mockArticleService, dummyAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()
		expectedArticles := []model.Article{
			{Id: 1, Title: "Article 1", Content: "Content 1"},
			{Id: 2, Title: "Article 2", Content: "Content 2"},
		}

		mockArticleService.On("FindAll").Return(expectedArticles, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/article/", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Success Get All Articles", responseBody["message"])
		data := responseBody["data"].([]interface{})
		assert.Len(t, data, 2)

		mockArticleService.AssertExpectations(t)
	})

	t.Run("Internal Server Error - FindAll fails", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		dummyAuthMiddleware := new(DummyAuthMiddleware) // Use DummyAuthMiddleware

		articleController := NewArticleController(mockArticleService, dummyAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		mockArticleService.On("FindAll").Return([]model.Article{}, errors.New("database error")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/article/", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Failed to get articles: database error", responseBody["message"])

		mockArticleService.AssertExpectations(t)
	})
}

func TestGetArticleBySlugHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Get Article By Slug", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		dummyAuthMiddleware := new(DummyAuthMiddleware) // Use DummyAuthMiddleware

		articleController := NewArticleController(mockArticleService, dummyAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		expectedArticle := model.Article{
			Id:      1,
			Title:   "Test Article",
			Content: "This is a test article content.",
			Slug:    "test-article",
		}

		mockArticleService.On("FindBySlug", "test-article").Return(expectedArticle, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/article/test-article", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Success get article by slug", responseBody["message"])
		data := responseBody["data"].(map[string]interface{})
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, "Test Article", data["title"])
		assert.Equal(t, "This is a test article content.", data["content"])
		assert.Equal(t, "test-article", data["slug"])

		mockArticleService.AssertExpectations(t)
	})

	t.Run("Article Not Found", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		dummyAuthMiddleware := new(DummyAuthMiddleware) // Use DummyAuthMiddleware

		articleController := NewArticleController(mockArticleService, dummyAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		mockArticleService.On("FindBySlug", "non-existent-slug").Return(model.Article{}, errors.New("not found")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/article/non-existent-slug", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Article not found", responseBody["error"])

		mockArticleService.AssertExpectations(t)
	})

	t.Run("Internal Server Error - FindBySlug fails", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		dummyAuthMiddleware := new(DummyAuthMiddleware) // Use DummyAuthMiddleware

		articleController := NewArticleController(mockArticleService, dummyAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		mockArticleService.On("FindBySlug", "test-slug").Return(model.Article{}, errors.New("database error")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/article/test-slug", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Failed to get article by slug: database error", responseBody["error"])

		mockArticleService.AssertExpectations(t)
	})
}

func TestGetArticlesByUserIdHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Get Articles By User ID", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		dummyAuthMiddleware := new(DummyAuthMiddleware) // Use DummyAuthMiddleware

		articleController := NewArticleController(mockArticleService, dummyAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		expectedArticles := []model.Article{
			{Id: 1, Title: "Article 1", Content: "Content 1", User: model.User{Id: 1}},
			{Id: 2, Title: "Article 2", Content: "Content 2", User: model.User{Id: 1}},
		}

		mockArticleService.On("FindByUserId", 1).Return(expectedArticles, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/article/u/1", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Success get articles by user ID", responseBody["message"])
		data := responseBody["data"].([]interface{})
		assert.Len(t, data, 2)

		mockArticleService.AssertExpectations(t)
	})

	t.Run("Invalid User ID", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		dummyAuthMiddleware := new(DummyAuthMiddleware) // Use DummyAuthMiddleware

		articleController := NewArticleController(mockArticleService, dummyAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/article/u/invalid", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Invalid user ID", responseBody["error"])

		mockArticleService.AssertNotCalled(t, "FindByUserId")
	})

	t.Run("Internal Server Error - FindByUserId fails", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		dummyAuthMiddleware := new(DummyAuthMiddleware) // Use DummyAuthMiddleware

		articleController := NewArticleController(mockArticleService, dummyAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		mockArticleService.On("FindByUserId", 1).Return([]model.Article{}, errors.New("database error")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/article/u/1", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["error"], "Failed to get articles")

		mockArticleService.AssertExpectations(t)
	})
}

func TestGetArticlesByCategoryHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Get Articles By Category", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		dummyAuthMiddleware := new(DummyAuthMiddleware) // Use DummyAuthMiddleware

		articleController := NewArticleController(mockArticleService, dummyAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		expectedArticles := []model.Article{
			{Id: 1, Title: "Article 1", Content: "Content 1", Category: model.Category{Name: "Technology"}},
			{Id: 2, Title: "Article 2", Content: "Content 2", Category: model.Category{Name: "Technology"}},
		}

		mockArticleService.On("FindByCategory", "Technology").Return(expectedArticles, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/article/c/Technology", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]interface{}
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Success get articles by Category", responseBody["message"])
		data := responseBody["data"].([]interface{})
		assert.Len(t, data, 2)

		mockArticleService.AssertExpectations(t)
	})

	t.Run("Internal Server Error - FindByCategory fails", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		dummyAuthMiddleware := new(DummyAuthMiddleware) // Use DummyAuthMiddleware

		articleController := NewArticleController(mockArticleService, dummyAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		mockArticleService.On("FindByCategory", "Technology").Return([]model.Article{}, errors.New("database error")).Once()

		req, _ := http.NewRequest(http.MethodGet, "/api/v1/article/c/Technology", nil)
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["error"], "Failed to get articles")

		mockArticleService.AssertExpectations(t)
	})
}

func TestDeleteArticleHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success Delete Article", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.Set("userId", float64(1))
			c.Next()
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articleID := "1"
		userID := 1
		existingArticle := model.Article{
			Id:      1,
			Title:   "Original Title",
			Content: "Original Content",
			User:    model.User{Id: userID},
		}

		mockArticleService.On("FindById", 1).Return(existingArticle, nil).Once()
		mockArticleService.On("DeleteArticle", 1).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/article/"+articleID, nil)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Success delete article", responseBody["message"])

		mockArticleService.AssertExpectations(t)
		mockAuthMiddleware.AssertExpectations(t)
	})

	t.Run("Unauthorized - No userId in context", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articleID := "1"

		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/article/"+articleID, nil)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Unauthorized", responseBody["message"])

		mockArticleService.AssertNotCalled(t, "FindById")
		mockArticleService.AssertNotCalled(t, "DeleteArticle")
		mockAuthMiddleware.AssertExpectations(t)
	})

	t.Run("Invalid User ID Type in Context", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.Set("userId", "invalid_type") // Set invalid type
			c.Next()
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articleID := "1"

		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/article/"+articleID, nil)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Invalid user ID type", responseBody["message"])

		mockArticleService.AssertNotCalled(t, "FindById")
		mockArticleService.AssertNotCalled(t, "DeleteArticle")
		mockAuthMiddleware.AssertExpectations(t)
	})

	t.Run("Invalid Article ID", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.Set("userId", float64(1))
			c.Next()
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articleID := "invalid_id"

		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/article/"+articleID, nil)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Invalid article ID", responseBody["error"])

		mockArticleService.AssertNotCalled(t, "FindById")
		mockArticleService.AssertNotCalled(t, "DeleteArticle")
		mockAuthMiddleware.AssertExpectations(t)
	})

	t.Run("Article Not Found", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.Set("userId", float64(1))
			c.Next()
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articleID := "1"
		
		mockArticleService.On("FindById", 1).Return(model.Article{}, errors.New("not found")).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/api/v1/article/"+articleID, nil)
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "Article not found", responseBody["message"])

		mockArticleService.AssertExpectations(t)
		mockAuthMiddleware.AssertExpectations(t)
	})

	t.Run("Forbidden - User does not own article", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.Set("userId", float64(2)) // Different user ID
			c.Next()
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articleID := "1"
		userID := 1
		title := "Updated Title"
		content := "Updated Content"
		updatePayload := dto.UpdateArticleRequest{
			Title:   &title,
			Content: &content,
		}
		existingArticle := model.Article{
			Id:      1,
			Title:   "Original Title",
			Content: "Original Content",
			User:    model.User{Id: userID},
		}

		mockArticleService.On("FindById", 1).Return(existingArticle, nil).Once()

		body, _ := json.Marshal(updatePayload)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/article/"+articleID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Equal(t, "You do not own this article", responseBody["message"])

		mockArticleService.AssertExpectations(t)
		mockAuthMiddleware.AssertExpectations(t)
	})

	t.Run("Bad Request - Invalid JSON Payload", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.Set("userId", float64(1))
			c.Next()
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articleID := "1"
		userID := 1
		existingArticle := model.Article{
			Id:      1,
			Title:   "Original Title",
			Content: "Original Content",
			User:    model.User{Id: userID},
		}

		mockArticleService.On("FindById", 1).Return(existingArticle, nil).Once()

		req, _ := http.NewRequest(http.MethodPut, "/api/v1/article/"+articleID, bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["error"], "invalid character")

		mockArticleService.AssertExpectations(t)
		mockAuthMiddleware.AssertExpectations(t)
	})

	t.Run("Internal Server Error - UpdateArticle fails", func(t *testing.T) {
		router := gin.Default()
		mockArticleService := new(MockArticleService)
		mockAuthMiddleware := new(MockAuthMiddleware)

		mockAuthMiddleware.ExpectedCalls = nil // Clear previous mock
		mockAuthMiddleware.On("CheckToken", mock.Anything).Return(gin.HandlerFunc(func(c *gin.Context) {
			c.Set("userId", float64(1))
			c.Next()
		})).Once()

		articleController := NewArticleController(mockArticleService, mockAuthMiddleware, router.Group("/api/v1"))
		articleController.Route()

		articleID := "1"
		userID := 1
		title := "Updated Title"
		content := "Updated Content"
		updatePayload := dto.UpdateArticleRequest{
			Title:   &title,
			Content: &content,
		}
		existingArticle := model.Article{
			Id:      1,
			Title:   "Original Title",
			Content: "Original Content",
			User:    model.User{Id: userID},
		}

		mockArticleService.On("FindById", 1).Return(existingArticle, nil).Once()
		mockArticleService.On("UpdateArticle", 1, updatePayload).Return(model.Article{}, errors.New("database error")).Once()

		body, _ := json.Marshal(updatePayload)
		req, _ := http.NewRequest(http.MethodPut, "/api/v1/article/"+articleID, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var responseBody map[string]string
		json.Unmarshal(rr.Body.Bytes(), &responseBody)
		assert.Contains(t, responseBody["error"], "database error")

		mockArticleService.AssertExpectations(t)
		mockAuthMiddleware.AssertExpectations(t)
	})
}

