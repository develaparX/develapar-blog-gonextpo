package controller

import (
	"context"
	"develapar-server/middleware"
	"develapar-server/model/dto"
	"develapar-server/service"
	"develapar-server/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleTagController struct {
	service        service.ArticleTagService
	rg             *gin.RouterGroup
	md             middleware.AuthMiddleware
	errorHandler   middleware.ErrorHandler
	responseHelper *utils.ResponseHelper
}

type AssignTagRequest struct {
	ArticleID int   `json:"article_id"`
	TagIDs    []int `json:"tag_ids"`
}

// @Summary Assign tags to an article by tag names
// @Description Assigns a list of tags (by name) to a specific article
// @Tags Tags
// @Accept json
// @Produce json
// @Param payload body dto.AssignTagsByNameDTO true "Article ID and list of tag names"
// @Success 200 {object} middleware.SuccessResponse "Tags assigned successfully"
// @Failure 400 {object} middleware.ErrorResponse "Invalid payload"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /article-to-tag [post]
func (c *ArticleTagController) AssignTagToArticleByNameHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var req dto.AssignTagsByNameDTO
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	err := c.service.AsignTagsByName(requestCtx, req.ArticleID, req.Tags)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "assign tags by name")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "assign tags by name")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to assign tags")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	successResponse := middleware.CreateSuccessResponse(requestCtx, gin.H{
		"message": "Tags assigned successfully",
	})

	ginCtx.JSON(http.StatusOK, successResponse)
}

func (c *ArticleTagController) AssignTagToArticleByIdHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var req AssignTagRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	err := c.service.AssignTags(requestCtx, req.ArticleID, req.TagIDs)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "assign tags by ID")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "assign tags by ID")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to assign tags")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	successResponse := middleware.CreateSuccessResponse(requestCtx, gin.H{
		"message": "Tags assigned successfully",
	})

	ginCtx.JSON(http.StatusOK, successResponse)
}

// @Summary Get tags by article ID
// @Description Get a list of tags associated with a specific article ID
// @Tags Tags
// @Produce json
// @Param article_id path int true "ID of the article to retrieve tags for"
// @Success 200 {object} middleware.SuccessResponse "List of tags for the article"
// @Failure 400 {object} middleware.ErrorResponse "Invalid article ID"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /article-to-tag/tags/{article_id} [get]
func (c *ArticleTagController) GetTagsByArticleIDHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 10*time.Second)
	defer cancel()

	articleID, err := strconv.Atoi(ginCtx.Param("article_id"))
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "article_id", "Invalid article ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	tags, err := c.service.FindTagByArticleId(requestCtx, articleID)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "get tags by article ID")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "get tags by article ID")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve tags")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	successResponse := middleware.CreateSuccessResponse(requestCtx, gin.H{
		"tags": tags,
	})

	ginCtx.JSON(http.StatusOK, successResponse)
}

// @Summary Get articles by tag ID
// @Description Get a list of articles associated with a specific tag ID
// @Tags Articles
// @Produce json
// @Param tag_id path int true "ID of the tag to retrieve articles for"
// @Success 200 {object} middleware.SuccessResponse "List of articles with the specified tag"
// @Failure 400 {object} middleware.ErrorResponse "Invalid tag ID"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /article-to-tag/article/{tag_id} [get]
func (c *ArticleTagController) GetArticlesByTagIDHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	tagID, err := strconv.Atoi(ginCtx.Param("tag_id"))
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "tag_id", "Invalid tag ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	articles, err := c.service.FindArticleByTagId(requestCtx, tagID)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "get articles by tag ID")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "get articles by tag ID")
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
	successResponse := middleware.CreateSuccessResponse(requestCtx, gin.H{
		"articles": articles,
	})

	ginCtx.JSON(http.StatusOK, successResponse)
}

// @Summary Remove a tag from an article
// @Description Remove a specific tag from an article by their IDs
// @Tags Tags
// @Produce json
// @Param article_id path int true "ID of the article"
// @Param tag_id path int true "ID of the tag to remove"
// @Success 200 {object} middleware.SuccessResponse "Tag removed from article successfully"
// @Failure 400 {object} middleware.ErrorResponse "Invalid article ID or tag ID"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /article-to-tag/articles/{article_id}/tags/{tag_id} [delete]
func (c *ArticleTagController) RemoveTagFromArticleHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	articleId, err := strconv.Atoi(ginCtx.Param("article_id"))
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "article_id", "Invalid article ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	tagId, err := strconv.Atoi(ginCtx.Param("tag_id"))
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "tag_id", "Invalid tag ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	err = c.service.RemoveTagFromArticle(requestCtx, articleId, tagId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "remove tag from article")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "remove tag from article")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to remove tag from article")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	successResponse := middleware.CreateSuccessResponse(requestCtx, gin.H{
		"message": "Tag removed from article successfully",
	})

	ginCtx.JSON(http.StatusOK, successResponse)
}

func (at *ArticleTagController) Route() {
	router := at.rg.Group("/article-to-tag")
	router.GET("/tags/:article_id", at.GetTagsByArticleIDHandler)
	router.GET("/article/:tag_id", at.GetArticlesByTagIDHandler)

	routerAuth := router.Group("/")
	routerAuth.Use(at.md.CheckToken())
	routerAuth.POST("/", at.AssignTagToArticleByNameHandler)
	routerAuth.DELETE("/articles/:article_id/tags/:tag_id", at.RemoveTagFromArticleHandler)

}

func NewArticleTagController(s service.ArticleTagService, rg *gin.RouterGroup, md middleware.AuthMiddleware, errorHandler middleware.ErrorHandler) *ArticleTagController {
	return &ArticleTagController{
		service:        s,
		rg:             rg,
		md:             md,
		errorHandler:   errorHandler,
		responseHelper: utils.NewResponseHelper(),
	}
}
