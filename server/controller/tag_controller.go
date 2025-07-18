package controller

import (
	"context"
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/service"
	"develapar-server/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	service        service.TagService
	rg             *gin.RouterGroup
	md             middleware.AuthMiddleware
	errorHandler   middleware.ErrorHandler
	responseHelper *utils.ResponseHelper
}

// @Summary Create a new tag
// @Description Create a new tag with a given name
// @Tags Tags
// @Accept json
// @Produce json
// @Param payload body model.Tags true "Tag creation details"
// @Success 201 {object} middleware.SuccessResponse "Tag successfully created"
// @Failure 400 {object} middleware.ErrorResponse "Invalid payload"
// @Failure 401 {object} middleware.ErrorResponse "Unauthorized"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /tags [post]
func (t *TagController) CreateTagHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.Tags
	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := t.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	data, err := t.service.CreateTag(requestCtx, payload)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := t.errorHandler.TimeoutError(requestCtx, "create tag")
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := t.errorHandler.CancellationError(requestCtx, "create tag")
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := t.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to create tag")
		appErr.StatusCode = 500
		t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Tag created successfully",
		"tag":     data,
	}
	t.responseHelper.SendCreated(ginCtx, responseData)
}

// @Summary Get all tags
// @Description Get a list of all tags
// @Tags Tags
// @Produce json
// @Success 200 {object} middleware.SuccessResponse "List of tags"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /tags [get]
func (t *TagController) GetAllTagHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 10*time.Second)
	defer cancel()

	// Call service with context
	data, err := t.service.FindAll(requestCtx)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := t.errorHandler.TimeoutError(requestCtx, "get all tags")
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := t.errorHandler.CancellationError(requestCtx, "get all tags")
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := t.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve tags")
		appErr.StatusCode = 500
		t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Tags retrieved successfully",
		"tags":    data,
	}
	t.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Get tag by ID
// @Description Get tag details by its ID
// @Tags Tags
// @Produce json
// @Param tags_id path int true "ID of the tag to retrieve"
// @Success 200 {object} middleware.SuccessResponse "Tag details"
// @Failure 400 {object} middleware.ErrorResponse "Invalid tag ID"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /tags/{tags_id} [get]
func (t *TagController) GetByTagIdHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 10*time.Second)
	defer cancel()

	tagId, err := strconv.Atoi(ginCtx.Param("tags_id"))
	if err != nil {
		appErr := t.errorHandler.ValidationError(requestCtx, "tags_id", "Invalid tag ID: "+err.Error())
		t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	tags, err := t.service.FindById(requestCtx, tagId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := t.errorHandler.TimeoutError(requestCtx, "get tag by ID")
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := t.errorHandler.CancellationError(requestCtx, "get tag by ID")
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := t.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve tag")
		appErr.StatusCode = 500
		t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Tag retrieved successfully",
		"tag":     tags,
	}
	t.responseHelper.SendSuccess(ginCtx, responseData)
}

func (t *TagController) Route() {
	router := t.rg.Group("/tags")
	router.GET("/:tags_id", t.GetByTagIdHandler)
	router.GET("/", t.GetAllTagHandler)

	routerAuth := router.Group("/", t.md.CheckToken())
	routerAuth.POST("/", t.CreateTagHandler)
}

func NewTagController(tS service.TagService, rg *gin.RouterGroup, md middleware.AuthMiddleware, errorHandler middleware.ErrorHandler) *TagController {
	return &TagController{
		service:        tS,
		rg:             rg,
		md:             md,
		errorHandler:   errorHandler,
		responseHelper: utils.NewResponseHelper(),
	}
}
