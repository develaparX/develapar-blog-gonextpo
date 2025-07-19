package controller

import (
	"context"
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/service"
	"develapar-server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type BookmarkController struct {
	service        service.BookmarkService
	rg             *gin.RouterGroup
	md             middleware.AuthMiddleware
	errorHandler   middleware.ErrorHandler
	responseHelper *utils.ResponseHelper
}

// @Summary Create a new bookmark
// @Description Create a new bookmark for an article
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param payload body model.Bookmark true "Bookmark creation details"
// @Success 201 {object} middleware.SuccessResponse "Bookmark successfully created"
// @Failure 400 {object} middleware.ErrorResponse "Invalid payload"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /bookmarks [post]
func (b *BookmarkController) CreateBookmarkHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.Bookmark
	userIdRaw, exists := ginCtx.Get("userId")
	if !exists {
		appErr := b.errorHandler.WrapError(requestCtx, nil, utils.ErrUnauthorized, "Authentication required")
		appErr.StatusCode = 401
		b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}
	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		appErr := b.errorHandler.WrapError(requestCtx, nil, utils.ErrInternal, "Invalid user ID type")
		appErr.StatusCode = 500
		b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}
	userId := int(userIdFloat)

	payload.User.Id = userId

	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := b.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	data, err := b.service.CreateBookmark(requestCtx, payload)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := b.errorHandler.TimeoutError(requestCtx, "create bookmark")
			b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := b.errorHandler.CancellationError(requestCtx, "create bookmark")
			b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := b.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to create bookmark")
		appErr.StatusCode = 500
		b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message":  "Bookmark created successfully",
		"bookmark": data,
	}
	b.responseHelper.SendCreated(ginCtx, responseData)
}

// @Summary Get bookmarks by user ID
// @Description Get a list of bookmarks for a specific user ID
// @Tags Bookmarks
// @Produce json
// @Param user_id path int true "ID of the user whose bookmarks to retrieve"
// @Success 200 {object} middleware.SuccessResponse "List of bookmarks for the user"
// @Failure 400 {object} middleware.ErrorResponse "Invalid user ID"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /bookmarks/{user_id} [get]
func (b *BookmarkController) GetBookmarkByUserId(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	userID := ginCtx.Param("user_id")
	if userID == "" {
		appErr := b.errorHandler.ValidationError(requestCtx, "user_id", "User ID is required")
		b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	bookmarks, err := b.service.FindByUserId(requestCtx, userID)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := b.errorHandler.TimeoutError(requestCtx, "get bookmarks by user ID")
			b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := b.errorHandler.CancellationError(requestCtx, "get bookmarks by user ID")
			b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := b.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve bookmarks")
		appErr.StatusCode = 500
		b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message":   "Bookmarks retrieved successfully",
		"bookmarks": bookmarks,
	}
	b.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Delete a bookmark
// @Description Delete a bookmark for an article by article ID
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param article_id path int true "Article ID to unbookmark"
// @Success 200 {object} middleware.SuccessResponse "Bookmark deleted successfully"
// @Failure 400 {object} middleware.ErrorResponse "Invalid article ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /bookmarks [delete]
func (b *BookmarkController) DeleteBookmarkHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	userIdRaw, exists := ginCtx.Get("userId")
	if !exists {
		appErr := b.errorHandler.WrapError(requestCtx, nil, utils.ErrUnauthorized, "Authentication required")
		appErr.StatusCode = 401
		b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}
	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		appErr := b.errorHandler.WrapError(requestCtx, nil, utils.ErrInternal, "Invalid user ID type")
		appErr.StatusCode = 500
		b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}
	userId := int(userIdFloat)

	articleIdParam := ginCtx.Param("article_id")
	articleId, err := strconv.Atoi(articleIdParam)
	if err != nil {
		appErr := b.errorHandler.ValidationError(requestCtx, "article_id", "Invalid article ID: "+err.Error())
		b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	err = b.service.DeleteBookmark(requestCtx, userId, articleId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := b.errorHandler.TimeoutError(requestCtx, "delete bookmark")
			b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := b.errorHandler.CancellationError(requestCtx, "delete bookmark")
			b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := b.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to delete bookmark")
		appErr.StatusCode = 500
		b.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Bookmark deleted successfully",
	}
	b.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Check if an article is bookmarked by the current user
// @Description Check if a specific article is bookmarked by the authenticated user
// @Tags Bookmarks
// @Produce json
// @Param article_id query int true "ID of the article to check"
// @Success 200 {object} middleware.SuccessResponse "Bookmark status"
// @Failure 400 {object} middleware.ErrorResponse "Invalid article ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /bookmarks/check [get]
func (c *BookmarkController) CheckBookmarkHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 10*time.Second)
	defer cancel()

	userIdRaw, exists := ginCtx.Get("userId")
	if !exists {
		appErr := c.errorHandler.WrapError(requestCtx, nil, utils.ErrUnauthorized, "Authentication required")
		appErr.StatusCode = 401
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}
	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		appErr := c.errorHandler.WrapError(requestCtx, nil, utils.ErrInternal, "Invalid user ID type")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}
	userId := int(userIdFloat)

	articleId, err := strconv.Atoi(ginCtx.Query("article_id"))
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "article_id", "Invalid article ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	bookmarked, err := c.service.IsBookmarked(requestCtx, userId, articleId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "check bookmark")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "check bookmark")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to check bookmark status")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"bookmarked": bookmarked,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

func (c *BookmarkController) Route() {
	router := c.rg.Group("/bookmarks")  // Changed from singular to plural
	router.GET("/:user_id", c.GetBookmarkByUserId)

	routerAuth := router.Group("/")
	routerAuth.Use(c.md.CheckToken())
	routerAuth.POST("/", c.CreateBookmarkHandler)
	routerAuth.DELETE("/", c.DeleteBookmarkHandler)
	routerAuth.GET("/check", c.CheckBookmarkHandler)
}

func NewBookmarkController(bS service.BookmarkService, rg *gin.RouterGroup, md middleware.AuthMiddleware, errorHandler middleware.ErrorHandler) *BookmarkController {
	return &BookmarkController{
		service:        bS,
		rg:             rg,
		md:             md,
		errorHandler:   errorHandler,
		responseHelper: utils.NewResponseHelper(),
	}
}
