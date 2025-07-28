package controller

import (
	"context"
	"develapar-server/middleware"
	"develapar-server/model/dto"
	"develapar-server/service"
	"develapar-server/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	service        service.ArticleService
	md             middleware.AuthMiddleware
	rg             *gin.RouterGroup
	errorHandler   middleware.ErrorHandler
	responseHelper *utils.ResponseHelper
}

// Helper function to extract user ID from context
func (c *ArticleController) getUserID(ctx *gin.Context) (int, error) {
	userIdRaw, exists := ctx.Get("userId")
	if !exists {
		return 0, fmt.Errorf("unauthorized")
	}

	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user ID type")
	}

	return int(userIdFloat), nil
}

// Helper function to parse article ID from URL parameter
func (c *ArticleController) parseArticleID(ctx *gin.Context) (int, error) {
	idStr := ctx.Param("article_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid article ID")
	}
	return id, nil
}

// @Summary Create a new article
// @Description Create a new blog article with tags
// @Tags Articles
// @Accept json
// @Produce json
// @Param payload body dto.CreateArticleRequest true "Article creation details"
// @Success 201 {object} dto.APIResponse{data=object{message=string,article=model.Article}} "Article successfully created"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /articles [post]
func (c *ArticleController) CreateArticleHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 30*time.Second)
	defer cancel()

	userId, err := c.getUserID(ginCtx)
	if err != nil {
		if err.Error() == "unauthorized" {
			appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrUnauthorized, "Authentication required")
			appErr.StatusCode = 401
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		} else {
			appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Invalid user ID type")
			appErr.StatusCode = 500
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		}
		return
	}

	var req dto.CreateArticleRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	data, err := c.service.CreateArticleWithTags(requestCtx, req, userId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "create article")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "create article")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to create article")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Article created successfully",
		"article": data,
	}
	c.responseHelper.SendCreated(ginCtx, responseData)
}

// @Summary Get all articles
// @Description Get a list of all blog articles
// @Tags Articles
// @Produce json
// @Success 200 {object} dto.APIResponse{data=object{message=string,articles=[]model.Article}} "List of articles"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /articles [get]
func (c *ArticleController) GetAllArticleHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	// Call service with context
	data, err := c.service.FindAll(requestCtx)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "get all articles")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "get all articles")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve articles")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message":  "Articles retrieved successfully",
		"articles": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Get all articles with pagination
// @Description Get a paginated list of all blog articles
// @Tags Articles
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} dto.APIResponse{data=object{message=string,articles=[]model.Article},pagination=dto.PaginationMetadata} "Paginated list of articles"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid pagination parameters"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /articles/paginated [get]
func (c *ArticleController) GetAllArticleWithPaginationHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 20*time.Second)
	defer cancel()

	// Get pagination parameters from query string
	page := 1
	limit := 10

	if pageStr := ginCtx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err != nil || p <= 0 {
			appErr := c.errorHandler.ValidationError(requestCtx, "page", "Page must be a positive integer")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		} else {
			page = p
		}
	}

	if limitStr := ginCtx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err != nil || l <= 0 || l > 100 {
			appErr := c.errorHandler.ValidationError(requestCtx, "limit", "Limit must be a positive integer between 1 and 100")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		} else {
			limit = l
		}
	}

	// Call service with pagination and context
	result, err := c.service.FindAllWithPagination(requestCtx, page, limit)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "get paginated articles")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "get paginated articles")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve paginated articles")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context and pagination
	responseData := gin.H{
		"message":  "Articles retrieved successfully",
		"articles": result.Data,
	}
	c.responseHelper.SendSuccessWithServicePagination(ginCtx, responseData, result.Metadata)
}

// @Summary Update an article
// @Description Update an existing article by ID
// @Tags Articles
// @Accept json
// @Produce json
// @Param article_id path int true "ID of the article to update"
// @Param payload body dto.UpdateArticleRequest true "Article update details"
// @Success 200 {object} dto.APIResponse{data=object{message=string,article=model.Article}} "Article updated successfully"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid article ID or payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 403 {object} dto.APIResponse{error=dto.ErrorResponse} "Forbidden (user does not own the article)"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Article not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /articles/{article_id} [put]
func (c *ArticleController) UpdateArticleHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 30*time.Second)
	defer cancel()

	userId, err := c.getUserID(ginCtx)
	if err != nil {
		if err.Error() == "unauthorized" {
			appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrUnauthorized, "Authentication required")
			appErr.StatusCode = 401
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		} else {
			appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Invalid user ID type")
			appErr.StatusCode = 500
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		}
		return
	}

	id, err := c.parseArticleID(ginCtx)
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "article_id", "Invalid article ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Check if user owns the article
	article, err := c.service.FindById(requestCtx, id)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "find article")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "find article")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrNotFound, "Article not found")
		appErr.StatusCode = 404
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	if article.User.Id != userId {
		appErr := c.errorHandler.WrapError(requestCtx, fmt.Errorf("user does not own article"), utils.ErrForbidden, "You do not own this article")
		appErr.StatusCode = 403
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Bind data from payload
	var req dto.UpdateArticleRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Update article with context
	updatedArticle, err := c.service.UpdateArticle(requestCtx, id, req)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "update article")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "update article")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to update article")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Article updated successfully",
		"article": updatedArticle,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}


// @Summary Get article by slug
// @Description Get article details by its slug
// @Tags Articles
// @Produce json
// @Param slug path string true "Slug of the article to retrieve"
// @Success 200 {object} dto.APIResponse{data=object{message=string,article=model.Article}} "Article details"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid slug"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Article not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /articles/{slug} [get]
func (c *ArticleController) GetBySlugHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	slug := ginCtx.Param("slug")
	if slug == "" {
		appErr := c.errorHandler.ValidationError(requestCtx, "slug", "Article slug is required")
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	article, err := c.service.FindBySlug(requestCtx, slug)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "get article by slug")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "get article by slug")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check for not found error
		if err.Error() == "not found" {
			appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrNotFound, "Article not found")
			appErr.StatusCode = 404
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve article")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Article retrieved successfully",
		"article": article,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Get articles by user ID
// @Description Get a list of articles by a specific user ID
// @Tags Articles
// @Produce json
// @Param user_id path int true "ID of the user whose articles to retrieve"
// @Success 200 {object} dto.APIResponse{data=object{message=string,articles=[]model.Article}} "List of articles by user"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid user ID"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /articles/author/{user_id} [get]
func (ac *ArticleController) GetByUserIdHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	userIdParam := ginCtx.Param("user_id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		appErr := ac.errorHandler.ValidationError(requestCtx, "user_id", "Invalid user ID: "+err.Error())
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	articles, err := ac.service.FindByUserId(requestCtx, userId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := ac.errorHandler.TimeoutError(requestCtx, "get articles by user ID")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := ac.errorHandler.CancellationError(requestCtx, "get articles by user ID")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := ac.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve articles by user ID")
		appErr.StatusCode = 500
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message":  "Articles retrieved successfully",
		"articles": articles,
	}
	ac.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Get articles by user ID with pagination
// @Description Get a paginated list of articles by a specific user ID
// @Tags Articles
// @Produce json
// @Param user_id path int true "ID of the user whose articles to retrieve"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} dto.APIResponse{data=object{message=string,articles=[]model.Article},pagination=dto.PaginationMetadata} "Paginated list of articles by user"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid user ID or pagination parameters"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /articles/author/{user_id}/paginated [get]
func (ac *ArticleController) GetByUserIdWithPaginationHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 20*time.Second)
	defer cancel()

	userIdParam := ginCtx.Param("user_id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		appErr := ac.errorHandler.ValidationError(requestCtx, "user_id", "Invalid user ID: "+err.Error())
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Get pagination parameters from query string
	page := 1
	limit := 10

	if pageStr := ginCtx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err != nil || p <= 0 {
			appErr := ac.errorHandler.ValidationError(requestCtx, "page", "Page must be a positive integer")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		} else {
			page = p
		}
	}

	if limitStr := ginCtx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err != nil || l <= 0 || l > 100 {
			appErr := ac.errorHandler.ValidationError(requestCtx, "limit", "Limit must be a positive integer between 1 and 100")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		} else {
			limit = l
		}
	}

	// Call service with pagination and context
	result, err := ac.service.FindByUserIdWithPagination(requestCtx, userId, page, limit)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := ac.errorHandler.TimeoutError(requestCtx, "get paginated articles by user ID")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := ac.errorHandler.CancellationError(requestCtx, "get paginated articles by user ID")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := ac.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve paginated articles by user ID")
		appErr.StatusCode = 500
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context and pagination
	responseData := gin.H{
		"message":  "Articles retrieved successfully",
		"articles": result.Data,
	}
	ac.responseHelper.SendSuccessWithServicePagination(ginCtx, responseData, result.Metadata)
}

// @Summary Get articles by category name
// @Description Get a list of articles by category name
// @Tags Articles
// @Produce json
// @Param category_name path string true "Name of the category to retrieve articles from"
// @Success 200 {object} dto.APIResponse{data=object{message=string,articles=[]model.Article}} "List of articles by category"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid category name"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /articles/category/{category_name} [get]
func (ac *ArticleController) GetByCategory(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	categoryName := ginCtx.Param("category_name")
	if categoryName == "" {
		appErr := ac.errorHandler.ValidationError(requestCtx, "category_name", "Category name is required")
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	articles, err := ac.service.FindByCategory(requestCtx, categoryName)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := ac.errorHandler.TimeoutError(requestCtx, "get articles by category")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := ac.errorHandler.CancellationError(requestCtx, "get articles by category")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := ac.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve articles by category")
		appErr.StatusCode = 500
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message":  "Articles retrieved successfully",
		"articles": articles,
	}
	ac.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Get articles by category name with pagination
// @Description Get a paginated list of articles by category name
// @Tags Articles
// @Produce json
// @Param category_name path string true "Name of the category to retrieve articles from"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} dto.APIResponse{data=object{message=string,articles=[]model.Article},pagination=dto.PaginationMetadata} "Paginated list of articles by category"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid pagination parameters"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /articles/category/{category_name}/paginated [get]
func (ac *ArticleController) GetByCategoryWithPaginationHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 20*time.Second)
	defer cancel()

	categoryName := ginCtx.Param("category_name")
	if categoryName == "" {
		appErr := ac.errorHandler.ValidationError(requestCtx, "category_name", "Category name is required")
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Get pagination parameters from query string
	page := 1
	limit := 10

	if pageStr := ginCtx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err != nil || p <= 0 {
			appErr := ac.errorHandler.ValidationError(requestCtx, "page", "Page must be a positive integer")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		} else {
			page = p
		}
	}

	if limitStr := ginCtx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err != nil || l <= 0 || l > 100 {
			appErr := ac.errorHandler.ValidationError(requestCtx, "limit", "Limit must be a positive integer between 1 and 100")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		} else {
			limit = l
		}
	}

	// Call service with pagination and context
	result, err := ac.service.FindByCategoryWithPagination(requestCtx, categoryName, page, limit)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := ac.errorHandler.TimeoutError(requestCtx, "get paginated articles by category")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := ac.errorHandler.CancellationError(requestCtx, "get paginated articles by category")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := ac.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve paginated articles by category")
		appErr.StatusCode = 500
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context and pagination
	responseData := gin.H{
		"message":  "Articles retrieved successfully",
		"articles": result.Data,
	}
	ac.responseHelper.SendSuccessWithServicePagination(ginCtx, responseData, result.Metadata)
}

// @Summary Delete an article
// @Description Delete an article by ID
// @Tags Articles
// @Produce json
// @Param article_id path int true "ID of the article to delete"
// @Success 200 {object} dto.APIResponse{data=object{message=string}} "Article deleted successfully"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid article ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 403 {object} dto.APIResponse{error=dto.ErrorResponse} "Forbidden (user does not own the article)"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Article not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /articles/{article_id} [delete]
func (ac *ArticleController) DeleteArticleHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 30*time.Second)
	defer cancel()

	userId, err := ac.getUserID(ginCtx)
	if err != nil {
		if err.Error() == "unauthorized" {
			appErr := ac.errorHandler.WrapError(requestCtx, err, utils.ErrUnauthorized, "Authentication required")
			appErr.StatusCode = 401
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		} else {
			appErr := ac.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Invalid user ID type")
			appErr.StatusCode = 500
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		}
		return
	}

	articleId, err := ac.parseArticleID(ginCtx)
	if err != nil {
		appErr := ac.errorHandler.ValidationError(requestCtx, "article_id", "Invalid article ID: "+err.Error())
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Check ownership
	article, err := ac.service.FindById(requestCtx, articleId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := ac.errorHandler.TimeoutError(requestCtx, "find article")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := ac.errorHandler.CancellationError(requestCtx, "find article")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		appErr := ac.errorHandler.WrapError(requestCtx, err, utils.ErrNotFound, "Article not found")
		appErr.StatusCode = 404
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	if article.User.Id != userId {
		appErr := ac.errorHandler.WrapError(requestCtx, fmt.Errorf("user does not own article"), utils.ErrForbidden, "You do not own this article")
		appErr.StatusCode = 403
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Delete article with context
	err = ac.service.DeleteArticle(requestCtx, articleId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := ac.errorHandler.TimeoutError(requestCtx, "delete article")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := ac.errorHandler.CancellationError(requestCtx, "delete article")
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		appErr := ac.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to delete article")
		appErr.StatusCode = 500
		ac.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Article deleted successfully",
	}
	ac.responseHelper.SendSuccess(ginCtx, responseData)
}


func (c *ArticleController) Route() {
	// Public routes - Changed from singular to plural
	publicRoutes := c.rg.Group("/articles")
	publicRoutes.GET("/", c.GetAllArticleHandler)
	publicRoutes.GET("/paginated", c.GetAllArticleWithPaginationHandler)
	publicRoutes.GET("/:slug", c.GetBySlugHandler)
	publicRoutes.GET("/author/:user_id", c.GetByUserIdHandler)                        // Changed from /u/ to /author/
	publicRoutes.GET("/author/:user_id/paginated", c.GetByUserIdWithPaginationHandler) // Changed from /u/ to /author/
	publicRoutes.GET("/category/:category_name", c.GetByCategory)                      // Changed from /c/:cat_name to /category/:category_name
	publicRoutes.GET("/category/:category_name/paginated", c.GetByCategoryWithPaginationHandler) // Changed from /c/:cat_name to /category/:category_name

	// Protected routes - Changed from singular to plural
	protectedRoutes := c.rg.Group("/articles")
	protectedRoutes.Use(c.md.CheckToken("user","admin")) // hanya butuh login
	protectedRoutes.POST("/", c.CreateArticleHandler)
	protectedRoutes.PUT("/:article_id", c.UpdateArticleHandler)
	protectedRoutes.DELETE("/:article_id", c.DeleteArticleHandler)
}


func NewArticleController(aS service.ArticleService, md middleware.AuthMiddleware, rg *gin.RouterGroup, errorHandler middleware.ErrorHandler) *ArticleController {
	return &ArticleController{
		service:        aS,
		md:             md,
		rg:             rg,
		errorHandler:   errorHandler,
		responseHelper: utils.NewResponseHelper(),
	}
}
