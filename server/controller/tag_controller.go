package controller

import (
	"context"
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/service"
	"develapar-server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
// @Success 201 {object} dto.APIResponse{data=model.Tags} "Tag successfully created"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
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
// @Success 200 {object} dto.APIResponse{data=object{message=string,tags=[]model.Tags}} "List of tags"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
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
// @Param tag_id path int true "ID of the tag to retrieve"
// @Success 200 {object} dto.APIResponse{data=object{message=string,tag=model.Tags}} "Tag details"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid tag ID"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /tags/{tag_id} [get]
func (t *TagController) GetByTagIdHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 10*time.Second)
	defer cancel()

	tagId, err := uuid.Parse(ginCtx.Param("tag_id"))
	if err != nil {
		appErr := t.errorHandler.ValidationError(requestCtx, "tag_id", "Invalid tag ID: "+err.Error())
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

// @Summary Update a tag
// @Description Update an existing tag by ID
// @Tags Tags
// @Accept json
// @Produce json
// @Param tag_id path int true "ID of the tag to update"
// @Param payload body model.Tags true "Tag update details"
// @Success 200 {object} dto.APIResponse{data=object{message=string,tag=model.Tags}} "Tag updated successfully"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid tag ID or payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Tag not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /tags/{tag_id} [put]
func (t *TagController) UpdateTagHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	tagId, err := uuid.Parse(ginCtx.Param("tag_id"))
	if err != nil {
		appErr := t.errorHandler.ValidationError(requestCtx, "tag_id", "Invalid tag ID: "+err.Error())
		t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	var payload model.Tags
	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := t.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	updatedTag, err := t.service.UpdateTag(requestCtx, tagId, payload)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := t.errorHandler.TimeoutError(requestCtx, "update tag")
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := t.errorHandler.CancellationError(requestCtx, "update tag")
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := t.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to update tag")
		appErr.StatusCode = 500
		t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Tag updated successfully",
		"tag":     updatedTag,
	}
	t.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Delete a tag
// @Description Delete a tag by ID
// @Tags Tags
// @Produce json
// @Param tag_id path int true "ID of the tag to delete"
// @Success 200 {object} dto.APIResponse{data=object{message=string}} "Tag deleted successfully"
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid tag ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Tag not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Security BearerAuth
// @Router /tags/{tag_id} [delete]
func (t *TagController) DeleteTagHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	tagId, err := uuid.Parse(ginCtx.Param("tag_id"))
	if err != nil {
		appErr := t.errorHandler.ValidationError(requestCtx, "tag_id", "Invalid tag ID: "+err.Error())
		t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	err = t.service.DeleteTag(requestCtx, tagId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := t.errorHandler.TimeoutError(requestCtx, "delete tag")
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := t.errorHandler.CancellationError(requestCtx, "delete tag")
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := t.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to delete tag")
		appErr.StatusCode = 500
		t.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message": "Tag deleted successfully",
	}
	t.responseHelper.SendSuccess(ginCtx, responseData)
}

func (t *TagController) Route() {
	router := t.rg.Group("/tags")
	router.GET("/:tag_id", t.GetByTagIdHandler) // Changed from tags_id to tag_id
	router.GET("/", t.GetAllTagHandler)

	routerAuth := router.Group("/", t.md.CheckToken())
	routerAuth.POST("/", t.CreateTagHandler)
	routerAuth.PUT("/:tag_id", t.UpdateTagHandler)    // Added missing update endpoint
	routerAuth.DELETE("/:tag_id", t.DeleteTagHandler) // Added missing delete endpoint
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
