package controller

import (
	"context"
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/service"
	"develapar-server/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	service      service.CategoryService
	rg           *gin.RouterGroup
	md           middleware.AuthMiddleware
	errorHandler middleware.ErrorHandler
}

// @Summary Create a new category
// @Description Create a new category with a given name
// @Tags Categories
// @Accept json
// @Produce json
// @Param payload body model.Category true "Category creation details"
// @Success 201 {object} middleware.SuccessResponse "Category successfully created"
// @Failure 400 {object} middleware.ErrorResponse "Invalid payload"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /category [post]
func (c *CategoryController) CreateCategoryHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.Category
	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	data, err := c.service.CreateCategory(requestCtx, payload)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "create category")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "create category")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to create category")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	successResponse := middleware.CreateSuccessResponse(requestCtx, gin.H{
		"message":  "Category created successfully",
		"category": data,
	})

	ginCtx.JSON(http.StatusCreated, successResponse)
}

// @Summary Get all categories
// @Description Get a list of all categories
// @Tags Categories
// @Produce json
// @Success 200 {object} middleware.SuccessResponse "List of categories"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /category [get]
func (c *CategoryController) GetAllCategoryHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 10*time.Second)
	defer cancel()

	// Call service with context
	data, err := c.service.FindAll(requestCtx)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "get all categories")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "get all categories")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve categories")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	successResponse := middleware.CreateSuccessResponse(requestCtx, gin.H{
		"message":    "Categories retrieved successfully",
		"categories": data,
	})

	ginCtx.JSON(http.StatusOK, successResponse)
}

// @Summary Update a category
// @Description Update an existing category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param cat_id path int true "ID of the category to update"
// @Param payload body dto.UpdateCategoryRequest true "Category update details"
// @Success 200 {object} middleware.SuccessResponse "Category updated successfully"
// @Failure 400 {object} middleware.ErrorResponse "Invalid category ID or payload"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /category/{cat_id} [put]
func (c *CategoryController) UpdateCategoryHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	idCat := ginCtx.Param("cat_id")
	id, err := strconv.Atoi(idCat)
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "cat_id", "Invalid category ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	var req dto.UpdateCategoryRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	cat, err := c.service.UpdateCategory(requestCtx, id, req)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "update category")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "update category")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to update category")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	successResponse := middleware.CreateSuccessResponse(requestCtx, gin.H{
		"message":  "Category updated successfully",
		"category": cat,
	})

	ginCtx.JSON(http.StatusOK, successResponse)
}

// @Summary Delete a category
// @Description Delete a category by ID
// @Tags Categories
// @Produce json
// @Param cat_id path int true "ID of the category to delete"
// @Success 200 {object} middleware.SuccessResponse "Category deleted successfully"
// @Failure 400 {object} middleware.ErrorResponse "Invalid category ID"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /category/{cat_id} [delete]
func (c *CategoryController) DeleteCategoryHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	Id := ginCtx.Param("cat_id")
	catId, err := strconv.Atoi(Id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "cat_id", "Invalid category ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	err = c.service.DeleteCategory(requestCtx, catId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "delete category")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "delete category")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to delete category")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	successResponse := middleware.CreateSuccessResponse(requestCtx, gin.H{
		"message": "Category deleted successfully",
	})

	ginCtx.JSON(http.StatusOK, successResponse)
}

func (c *CategoryController) Route() {
	router := c.rg.Group("/category")
	router.GET("/", c.GetAllCategoryHandler)

	routerAuth := router.Group("/", c.md.CheckToken("admin"))
	routerAuth.POST("/", c.CreateCategoryHandler)
	routerAuth.PUT("/:cat_id", c.UpdateCategoryHandler)
	routerAuth.DELETE("/:cat_id", c.DeleteCategoryHandler)
}

func NewCategoryController(cS service.CategoryService, rg *gin.RouterGroup, md middleware.AuthMiddleware, errorHandler middleware.ErrorHandler) *CategoryController {
	return &CategoryController{
		service:      cS,
		rg:           rg,
		md:           md,
		errorHandler: errorHandler,
	}
}
