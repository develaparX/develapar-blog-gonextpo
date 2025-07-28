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

type LikeController struct {
	service        service.LikeService
	rg             *gin.RouterGroup
	md             middleware.AuthMiddleware
	errorHandler   middleware.ErrorHandler
	responseHelper *utils.ResponseHelper
}

// @Summary Add a like to an article
// @Description Add a like to a specific article by the authenticated user
// @Tags Likes
// @Accept json
// @Produce json
// @Param payload body model.Likes true "Like creation details"
// @Success 201 {object} dto.APIResponse{data=object{message=string,like=model.Likes}} "Like successfully added"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /likes [post]
func (l *LikeController) AddLikeHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.Likes

	userIdRaw, exists := ginCtx.Get("userId")
	if !exists {
		appErr := l.errorHandler.WrapError(requestCtx, nil, utils.ErrUnauthorized, "Authentication required")
		appErr.StatusCode = 401
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}
	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		appErr := l.errorHandler.WrapError(requestCtx, nil, utils.ErrInternal, "Invalid user ID type")
		appErr.StatusCode = 500
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}
	userId := int(userIdFloat)

	payload.User.Id = userId

	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := l.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	data, err := l.service.CreateLike(requestCtx, payload)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := l.errorHandler.TimeoutError(requestCtx, "add like")
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := l.errorHandler.CancellationError(requestCtx, "add like")
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := l.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to add like")
		appErr.StatusCode = 500
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Like added successfully",
		"like":    data,
	}
	l.responseHelper.SendCreated(ginCtx, responseData)
}

// @Summary Get likes by article ID
// @Description Get a list of likes for a specific article ID
// @Tags Likes
// @Produce json
// @Param article_id path int true "ID of the article to retrieve likes for"
// @Success 200 {object} dto.APIResponse{data=object{message=string,likes=[]model.Likes}} "List of likes for the article"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid article ID"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /likes/article/{article_id} [get]
func (l *LikeController) GetLikeByArticleIdHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	articleID, err := strconv.Atoi(ginCtx.Param("article_id"))
	if err != nil {
		appErr := l.errorHandler.ValidationError(requestCtx, "article_id", "Invalid article ID: "+err.Error())
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	likes, err := l.service.FindLikeByArticleId(requestCtx, articleID)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := l.errorHandler.TimeoutError(requestCtx, "get likes by article ID")
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := l.errorHandler.CancellationError(requestCtx, "get likes by article ID")
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := l.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve likes")
		appErr.StatusCode = 500
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Likes retrieved successfully",
		"likes":   likes,
	}
	l.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Get likes by user ID
// @Description Get a list of likes by a specific user ID
// @Tags Likes
// @Produce json
// @Param user_id path int true "ID of the user whose likes to retrieve"
// @Success 200 {object} dto.APIResponse{data=object{message=string,likes=[]model.Likes}} "List of likes by the user"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid user ID"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /likes/user/{user_id} [get]
func (l *LikeController) GetLikeByUserIdHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	userID, err := strconv.Atoi(ginCtx.Param("user_id"))
	if err != nil {
		appErr := l.errorHandler.ValidationError(requestCtx, "user_id", "Invalid user ID: "+err.Error())
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	likes, err := l.service.FindLikeByUserId(requestCtx, userID)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := l.errorHandler.TimeoutError(requestCtx, "get likes by user ID")
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := l.errorHandler.CancellationError(requestCtx, "get likes by user ID")
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := l.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve likes")
		appErr.StatusCode = 500
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Likes retrieved successfully",
		"likes":   likes,
	}
	l.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Remove a like from an article
// @Description Remove a like from a specific article by the authenticated user
// @Tags Likes
// @Accept json
// @Produce json
// @Param payload body object{article_id=int} true "Article ID to unlike"
// @Success 200 {object} dto.APIResponse{data=object{message=string}} "Like deleted successfully"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid article ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /likes [delete]
func (l *LikeController) DeleteLikeHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.Likes

	userIdRaw, exists := ginCtx.Get("userId")
	if !exists {
		appErr := l.errorHandler.WrapError(requestCtx, nil, utils.ErrUnauthorized, "Authentication required")
		appErr.StatusCode = 401
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}
	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		appErr := l.errorHandler.WrapError(requestCtx, nil, utils.ErrInternal, "Invalid user ID type")
		appErr.StatusCode = 500
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}
	userId := int(userIdFloat)

	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := l.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	err := l.service.DeleteLike(requestCtx, userId, payload.Article.Id)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := l.errorHandler.TimeoutError(requestCtx, "delete like")
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := l.errorHandler.CancellationError(requestCtx, "delete like")
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := l.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to delete like")
		appErr.StatusCode = 500
		l.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Like deleted successfully",
	}
	l.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Check if an article is liked by the current user
// @Description Check if a specific article is liked by the authenticated user
// @Tags Likes
// @Produce json
// @Param article_id query int true "ID of the article to check"
// @Success 200 {object} dto.APIResponse{data=object{liked=bool}} "Like status"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid article ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /likes/check [get]
func (c *LikeController) CheckLikeHandler(ginCtx *gin.Context) {
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
	liked, err := c.service.IsLiked(requestCtx, userId, articleId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "check like")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "check like")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to check like status")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"liked": liked,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

func (l *LikeController) Route() {
	router := l.rg.Group("/likes")
	router.GET("/article/:article_id", l.GetLikeByArticleIdHandler)
	router.GET("/user/:user_id", l.GetLikeByUserIdHandler)

	routerAuth := router.Group("/", l.md.CheckToken())
	routerAuth.POST("/", l.AddLikeHandler)
	routerAuth.DELETE("/", l.DeleteLikeHandler)
	routerAuth.GET("/check", l.CheckLikeHandler)
}

func NewLikeController(lS service.LikeService, rg *gin.RouterGroup, md middleware.AuthMiddleware, errorHandler middleware.ErrorHandler) *LikeController {
	return &LikeController{
		service:        lS,
		rg:             rg,
		md:             md,
		errorHandler:   errorHandler,
		responseHelper: utils.NewResponseHelper(),
	}
}
