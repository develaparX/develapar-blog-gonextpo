package controller

import (
	"context"
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/service"
	"develapar-server/utils"
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	service        service.CommentService
	rg             *gin.RouterGroup
	md             middleware.AuthMiddleware
	errorHandler   middleware.ErrorHandler
	responseHelper *utils.ResponseHelper
}

// @Summary Create a new comment
// @Description Create a new comment on an article
// @Tags Comments
// @Accept json
// @Produce json
// @Param payload body model.Comment true "Comment creation details"
// @Success 201 {object} dto.APIResponse{data=object{message=string,comment=model.Comment}} "Comment successfully created"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /comments [post]
func (c *CommentController) CreateCommentHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.Comment

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

	payload.User.Id = userId

	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	data, err := c.service.CreateComment(requestCtx, payload)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "create comment")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "create comment")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to create comment")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Comment created successfully",
		"comment": data,
	}
	c.responseHelper.SendCreated(ginCtx, responseData)
}

// @Summary Get comments by article ID
// @Description Get a list of comments for a specific article ID
// @Tags Comments
// @Produce json
// @Param article_id path int true "ID of the article to retrieve comments for"
// @Success 200 {object} dto.APIResponse{data=object{message=string,comments=[]model.Comment}} "List of comments for the article"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid article ID"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /comments/article/{article_id} [get]
func (c *CommentController) FindCommentByArticleIdHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	articleId, err := strconv.Atoi(ginCtx.Param("article_id"))
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "article_id", "Invalid article ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	comments, err := c.service.FindCommentByArticleId(requestCtx, articleId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "get comments by article ID")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "get comments by article ID")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve comments")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message":  "Comments retrieved successfully",
		"comments": comments,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Get comments by user ID
// @Description Get a list of comments by a specific user ID
// @Tags Comments
// @Produce json
// @Param user_id path int true "ID of the user whose comments to retrieve"
// @Success 200 {object} dto.APIResponse{data=object{message=string,comments=[]model.Comment}} "List of comments by the user"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid user ID"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /comments/user/{user_id} [get]
func (c *CommentController) FindCommentByUserIdHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	user_id, err := strconv.Atoi(ginCtx.Param("user_id"))
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "user_id", "Invalid user ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	comments, err := c.service.FindCommentByUserId(requestCtx, user_id)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "get comments by user ID")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "get comments by user ID")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve comments")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message":  "Comments retrieved successfully",
		"comments": comments,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Update a comment
// @Description Update an existing comment by ID
// @Tags Comments
// @Accept json
// @Produce json
// @Param comment_id path int true "ID of the comment to update"
// @Param payload body object{content=string} true "Comment update details"
// @Success 200 {object} dto.APIResponse{data=object{message=string}} "Comment updated successfully"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 403 {object} dto.APIResponse{error=dto.ErrorResponse} "Forbidden (user does not own the comment)"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /comments/{comment_id} [put]
func (c *CommentController) UpdateCommentHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	userId := ginCtx.GetInt("userId")
	commentId, err := strconv.Atoi(ginCtx.Param("comment_id"))
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "id", "Invalid comment ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	err = c.service.EditComment(requestCtx, commentId, req.Content, userId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "update comment")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "update comment")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check for specific service errors
		if errors.Is(err, service.ErrUnauthorized) {
			appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrForbidden, "You do not own this comment")
			appErr.StatusCode = 403
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to update comment")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Comment updated successfully",
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Delete a comment
// @Description Delete a comment by ID
// @Tags Comments
// @Produce json
// @Param comment_id path int true "ID of the comment to delete"
// @Success 200 {object} dto.APIResponse{data=object{message=string}} "Comment deleted successfully"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid comment ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /comments/{comment_id} [delete]
func (c *CommentController) DeleteCommentHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	userId := ginCtx.GetInt("userId")

	commentId, err := strconv.Atoi(ginCtx.Param("comment_id"))
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "comment_id", "Invalid comment ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	if err := c.service.DeleteComment(requestCtx, commentId, userId); err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "delete comment")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "delete comment")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to delete comment")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Comment deleted successfully",
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

func (c *CommentController) Route() {
	router := c.rg.Group("/comments")  // Changed from singular to plural
	router.GET("/article/:article_id", c.FindCommentByArticleIdHandler)  // Fixed typo: c:article_id to :article_id
	router.GET("/user/:user_id", c.FindCommentByUserIdHandler)

	routerAuth := router.Group("/", c.md.CheckToken())

	routerAuth.POST("/", c.CreateCommentHandler)
	routerAuth.PUT("/:comment_id", c.UpdateCommentHandler)    // Changed from :id to :comment_id for consistency
	routerAuth.DELETE("/:comment_id", c.DeleteCommentHandler) // Changed from :id to :comment_id for consistency
}

func NewCommentController(cS service.CommentService, rg *gin.RouterGroup, md middleware.AuthMiddleware, errorHandler middleware.ErrorHandler) *CommentController {
	return &CommentController{
		service:        cS,
		rg:             rg,
		md:             md,
		errorHandler:   errorHandler,
		responseHelper: utils.NewResponseHelper(),
	}
}
