package controller

import (
	"context"
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/service"
	"develapar-server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	service        service.CategoryService
	rg             *gin.RouterGroup
	md             middleware.AuthMiddleware
	errorHandler   middleware.ErrorHandler
	responseHelper *utils.ResponseHelper
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
// @Router /categories [post]
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
	responseData := gin.H{
		"message":  "Category created successfully",
		"category": data,
	}
	c.responseHelper.SendCreated(ginCtx, responseData)
}

// @Summary Get all categories
// @Description Get a list of all categories
// @Tags Categories
// @Produce json
// @Success 200 {object} middleware.SuccessResponse "List of categories"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /categories [get]
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
	responseData := gin.H{
		"message":    "Categories retrieved successfully",
		"categories": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Get category by ID
// @Description Get category details by its ID
// @Tags Categories
// @Produce json
// @Param category_id path int true "ID of the category to retrieve"
// @Success 200 {object} middleware.SuccessResponse "Category details"
// @Failure 400 {object} middleware.ErrorResponse "Invalid category ID"
// @Failure 404 {object} middleware.ErrorResponse "Category not found"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Router /categories/{category_id} [get]
func (c *CategoryController) GetCategoryByIdHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 10*time.Second)
	defer cancel()

	categoryIdStr := ginCtx.Param("category_id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "category_id", "Invalid category ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	category, err := c.service.FindById(requestCtx, categoryId)
	if err != nil {
		// Check for context-specific errors
		if requestCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(requestCtx, "get category")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}
		if requestCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(requestCtx, "get category")
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Check if it's already an AppError
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
			return
		}

		// Wrap as internal error
		appErr := c.errorHandler.WrapError(requestCtx, err, utils.ErrInternal, "Failed to retrieve category")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Create success response with context
	responseData := gin.H{
		"message":  "Category retrieved successfully",
		"category": category,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Update a category
// @Description Update an existing category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param category_id path int true "ID of the category to update"
// @Param payload body dto.UpdateCategoryRequest true "Category update details"
// @Success 200 {object} middleware.SuccessResponse "Category updated successfully"
// @Failure 400 {object} middleware.ErrorResponse "Invalid category ID or payload"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /categories/{category_id} [put]
func (c *CategoryController) UpdateCategoryHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	categoryIdStr := ginCtx.Param("category_id")
	id, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "category_id", "Invalid category ID: "+err.Error())
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
	responseData := gin.H{
		"message":  "Category updated successfully",
		"category": cat,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// @Summary Delete a category
// @Description Delete a category by ID
// @Tags Categories
// @Produce json
// @Param category_id path int true "ID of the category to delete"
// @Success 200 {object} middleware.SuccessResponse "Category deleted successfully"
// @Failure 400 {object} middleware.ErrorResponse "Invalid category ID"
// @Failure 408 {object} middleware.ErrorResponse "Request timeout"
// @Failure 500 {object} middleware.ErrorResponse "Internal server error"
// @Security BearerAuth
// @Router /categories/{category_id} [delete]
func (c *CategoryController) DeleteCategoryHandler(ginCtx *gin.Context) {
	// Get request context with timeout
	requestCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	categoryIdStr := ginCtx.Param("category_id")
	categoryId, err := strconv.Atoi(categoryIdStr)
	if err != nil {
		appErr := c.errorHandler.ValidationError(requestCtx, "category_id", "Invalid category ID: "+err.Error())
		c.errorHandler.HandleError(requestCtx, ginCtx, appErr)
		return
	}

	// Call service with context
	err = c.service.DeleteCategory(requestCtx, categoryId)
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
	responseData := gin.H{
		"message": "Category deleted successfully",
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

func (c *CategoryController) Route() {
	router := c.rg.Group("/categories")  // Changed from singular to plural
	router.GET("/", c.GetAllCategoryHandler)
	router.GET("/:category_id", c.GetCategoryByIdHandler)  // Added missing endpoint

	routerAuth := router.Group("/", c.md.CheckToken())
	routerAuth.POST("/", c.CreateCategoryHandler)
	routerAuth.PUT("/:category_id", c.UpdateCategoryHandler)    // Changed from cat_id to category_id
	routerAuth.DELETE("/:category_id", c.DeleteCategoryHandler) // Changed from cat_id to category_id
}

func NewCategoryController(cS service.CategoryService, rg *gin.RouterGroup, md middleware.AuthMiddleware, errorHandler middleware.ErrorHandler) *CategoryController {
	return &CategoryController{
		service:        cS,
		rg:             rg,
		md:             md,
		errorHandler:   errorHandler,
		responseHelper: utils.NewResponseHelper(),
	}
}
