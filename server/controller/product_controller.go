package controller

import (
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/service"
	"develapar-server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type ProductController struct {
	s              service.ProductService
	rg             *gin.RouterGroup
	mD             middleware.AuthMiddleware
	errorHandler   middleware.ErrorHandler
	responseHelper *utils.ResponseHelper
}

func NewProductController(
	pS service.ProductService,
	rg *gin.RouterGroup,
	mD middleware.AuthMiddleware,
	errorHandler middleware.ErrorHandler,
) *ProductController {
	return &ProductController{
		s:            pS,
		rg:           rg,
		mD:           mD,
		errorHandler: errorHandler,
	}
}

func (c *ProductController) Route() {

	{
		router := c.rg.Group("/product-categories")
		router.GET("/", c.GetAllProductCategories)
		router.GET("/:id", c.GetProductCategoryById)
		router.GET("/s/:slug", c.GetProductCategoryBySlug)

		routerAuth := router.Group("/", c.mD.CheckToken())
		routerAuth.POST("/", c.CreateProductCategory)
		routerAuth.PUT("/:id", c.UpdateProductCategory)
		routerAuth.DELETE("/:id", c.DeleteProductCategory)
	}
}

func (c *ProductController) CreateProductCategory(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.ProductCategory
	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)

		return
	}

	data, err := c.s.CreateProductCategory(reqCtx, payload)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Create Product Category")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Create Product Category")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to create category")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return

	}

	responseData := gin.H{
		"message":          "Product Category created successfully",
		"product_category": data,
	}
	c.responseHelper.SendCreated(ginCtx, responseData)
}

func (c *ProductController) GetAllProductCategories(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	data, err := c.s.GetAllProductCategories(reqCtx)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Get All Product Categories")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Get All Product Categories")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to get product categories")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message":            "Product Categories retrieved successfully",
		"product_categories": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

func (c *ProductController) GetProductCategoryById(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product Category ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product Category ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	data, err := c.s.GetProductCategoryById(reqCtx, uuId)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Get Product Category By ID")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Get Product Category By ID")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to get product category by ID")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message":          "Product Category retrieved successfully",
		"product_category": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

func (c *ProductController) GetProductCategoryBySlug(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	slug := ginCtx.Param("slug")
	if slug == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "slug", "Product Category slug is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	data, err := c.s.GetProductCategoryBySlug(reqCtx, slug)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Get Product Category By Slug")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Get Product Category By Slug")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to get product category by slug")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message":          "Product Category retrieved successfully",
		"product_category": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

func (c *ProductController) UpdateProductCategory(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.ProductCategory
	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product Category ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product Category ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	payload.Id = uuId

	data, err := c.s.UpdateProductCategory(reqCtx, payload)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Update Product Category")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Update Product Category")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to update product category")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message":          "Product Category updated successfully",
		"product_category": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

func (c *ProductController) DeleteProductCategory(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product Category ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	// Parse the UUID from the ID parameter
	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product Category ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	err = c.s.DeleteProductCategory(reqCtx, uuId)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Delete Product Category")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Delete Product Category")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to delete product category")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message": "Product Category deleted successfully",
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}
