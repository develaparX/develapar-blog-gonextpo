package controller

import (
	"develapar-server/middleware"
	"develapar-server/model/dto"
	"develapar-server/service"
	"develapar-server/utils"
	"strconv"
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
		s:              pS,
		rg:             rg,
		mD:             mD,
		errorHandler:   errorHandler,
		responseHelper: utils.NewResponseHelper(),
	}
}

func (c *ProductController) Route() {
	// Product Categories routes
	{
		routerProductCat := c.rg.Group("/product-categories")
		routerProductCat.GET("/", c.GetAllProductCategories)
		routerProductCat.GET("/:id", c.GetProductCategoryById)
		routerProductCat.GET("/s/:slug", c.GetProductCategoryBySlug)

		routerPCAuth := routerProductCat.Group("/", c.mD.CheckToken("admin"))
		routerPCAuth.POST("/", c.CreateProductCategory)
		routerPCAuth.PUT("/:id", c.UpdateProductCategory)
		routerPCAuth.DELETE("/:id", c.DeleteProductCategory)
	}

	// Products routes
	{
		routerProduct := c.rg.Group("/products")
		routerProduct.GET("/", c.GetAllProducts)
		routerProduct.GET("/:id", c.GetProductById)
		routerProduct.GET("/category/:id", c.GetProductsByCategory)
		routerProduct.GET("/article/:id", c.GetProductsByArticleId)

		routerPAuth := routerProduct.Group("/", c.mD.CheckToken("admin"))
		routerPAuth.POST("/", c.CreateProduct)
		routerPAuth.PUT("/:id", c.UpdateProduct)
		routerPAuth.DELETE("/:id", c.DeleteProduct)
		routerPAuth.POST("/:id/article/:articleId", c.AddProductToArticle)
		routerPAuth.DELETE("/:id/article/:articleId", c.RemoveProductFromArticle)

		// Affiliate links routes
		routerAffiliate := routerProduct.Group("/:id/affiliate", c.mD.CheckToken("admin"))
		routerAffiliate.POST("/", c.CreateProductAffiliateLink)
		routerAffiliate.GET("/", c.GetAffiliateLinksbyProductId)
		routerAffiliate.PUT("/:affiliateId", c.UpdateProductAffiliateLink)
		routerAffiliate.DELETE("/:affiliateId", c.DeleteProductAffiliateLink)
	}
}

// Product Category Controllers

// CreateProductCategory godoc
// @Summary Create a new product category
// @Description Create a new product category with name, slug, and description
// @Tags Product Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.CreateProductCategoryRequest true "Product category creation details"
// @Success 201 {object} dto.APIResponse{data=object{message=string,product_category=dto.ProductCategoryResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /product-categories [post]
func (c *ProductController) CreateProductCategory(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var req dto.CreateProductCategoryRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	data, err := c.s.CreateProductCategory(reqCtx, req)
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

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to create product category")
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

// GetAllProductCategories godoc
// @Summary Get all product categories with pagination
// @Description Get a paginated list of all product categories
// @Tags Product Categories
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} dto.APIResponse{data=object{product_categories=[]dto.ProductCategoryResponse},pagination=dto.PaginationMetadata}
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /product-categories [get]
func (c *ProductController) GetAllProductCategories(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	// Parse pagination parameters
	page := 1
	limit := 10
	
	if pageStr := ginCtx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := ginCtx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	data, err := c.s.GetAllProductCategoriesWithPagination(reqCtx, page, limit)
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

	c.responseHelper.SendSuccess(ginCtx, data)
}

// GetProductCategoryById godoc
// @Summary Get product category by ID
// @Description Get a specific product category by its UUID
// @Tags Product Categories
// @Produce json
// @Param id path string true "Product Category ID (UUID)"
// @Success 200 {object} dto.APIResponse{data=object{message=string,product_category=dto.ProductCategoryResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid product category ID"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product category not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /product-categories/{id} [get]
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

// GetProductCategoryBySlug godoc
// @Summary Get product category by slug
// @Description Get a specific product category by its slug
// @Tags Product Categories
// @Produce json
// @Param slug path string true "Product Category Slug"
// @Success 200 {object} dto.APIResponse{data=object{message=string,product_category=dto.ProductCategoryResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid slug"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product category not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /product-categories/s/{slug} [get]
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

// UpdateProductCategory godoc
// @Summary Update product category
// @Description Update an existing product category by ID
// @Tags Product Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product Category ID (UUID)"
// @Param payload body dto.UpdateProductCategoryRequest true "Product category update details"
// @Success 200 {object} dto.APIResponse{data=object{message=string,product_category=dto.ProductCategoryResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload or ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product category not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /product-categories/{id} [put]
func (c *ProductController) UpdateProductCategory(ginCtx *gin.Context) {
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

	var req dto.UpdateProductCategoryRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	data, err := c.s.UpdateProductCategory(reqCtx, uuId, req)
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

// DeleteProductCategory godoc
// @Summary Delete product category
// @Description Delete a product category by ID
// @Tags Product Categories
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product Category ID (UUID)"
// @Success 200 {object} dto.APIResponse{data=object{message=string}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid product category ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product category not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /product-categories/{id} [delete]
func (c *ProductController) DeleteProductCategory(ginCtx *gin.Context) {
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

// Product Controllers

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with category, name, description, and image
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body dto.CreateProductRequest true "Product creation details"
// @Success 201 {object} dto.APIResponse{data=object{message=string,product=dto.ProductResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products [post]
func (c *ProductController) CreateProduct(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var req dto.CreateProductRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	data, err := c.s.CreateProduct(reqCtx, req)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Create Product")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Create Product")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to create product")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message": "Product created successfully",
		"product": data,
	}
	c.responseHelper.SendCreated(ginCtx, responseData)
}

// GetAllProducts godoc
// @Summary Get all products with pagination
// @Description Get a paginated list of all products
// @Tags Products
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} dto.APIResponse{data=object{products=[]dto.ProductListResponse},pagination=dto.PaginationMetadata}
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products [get]
func (c *ProductController) GetAllProducts(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	// Parse pagination parameters
	page := 1
	limit := 10
	
	if pageStr := ginCtx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := ginCtx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	data, err := c.s.GetAllProductsWithPagination(reqCtx, page, limit)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Get All Products")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Get All Products")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to get products")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	c.responseHelper.SendSuccess(ginCtx, data)
}

// GetProductById godoc
// @Summary Get product by ID
// @Description Get a specific product by its UUID including affiliate links
// @Tags Products
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} dto.APIResponse{data=object{message=string,product=dto.ProductResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid product ID"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/{id} [get]
func (c *ProductController) GetProductById(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	data, err := c.s.GetProductById(reqCtx, uuId)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Get Product By ID")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Get Product By ID")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to get product by ID")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message": "Product retrieved successfully",
		"product": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// GetProductsByCategory godoc
// @Summary Get products by category with pagination
// @Description Get a paginated list of products filtered by category ID
// @Tags Products
// @Produce json
// @Param id path string true "Category ID (UUID)"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10, max: 100)"
// @Success 200 {object} dto.APIResponse{data=object{message=string,products=object{products=[]dto.ProductListResponse}},pagination=dto.PaginationMetadata}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid category ID or pagination parameters"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/category/{id} [get]
func (c *ProductController) GetProductsByCategory(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Category ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Category ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 10
	
	if pageStr := ginCtx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := ginCtx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	data, err := c.s.GetProductsByCategoryWithPagination(reqCtx, uuId, page, limit)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Get Products By Category")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Get Products By Category")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to get products by category")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message":  "Products retrieved successfully",
		"products": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update an existing product by ID
// @Tags Products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Param payload body dto.UpdateProductRequest true "Product update details"
// @Success 200 {object} dto.APIResponse{data=object{message=string,product=dto.ProductResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload or ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/{id} [put]
func (c *ProductController) UpdateProduct(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	var req dto.UpdateProductRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	data, err := c.s.UpdateProduct(reqCtx, uuId, req)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Update Product")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Update Product")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to update product")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message": "Product updated successfully",
		"product": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product by ID
// @Tags Products
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} dto.APIResponse{data=object{message=string}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid product ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/{id} [delete]
func (c *ProductController) DeleteProduct(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	err = c.s.DeleteProduct(reqCtx, uuId)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Delete Product")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Delete Product")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to delete product")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message": "Product deleted successfully",
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// Product Affiliate Link Controllers

// CreateProductAffiliateLink godoc
// @Summary Create product affiliate link
// @Description Create a new affiliate link for a specific product
// @Tags Product Affiliate Links
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Param payload body dto.CreateProductAffiliateLinkRequest true "Affiliate link creation details"
// @Success 201 {object} dto.APIResponse{data=object{message=string,product_affiliate_link=dto.ProductAffiliateLinkResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload or product ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/{id}/affiliate [post]
func (c *ProductController) CreateProductAffiliateLink(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	var req dto.CreateProductAffiliateLinkRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	data, err := c.s.CreateProductAffiliateLink(reqCtx, uuId, req)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Create Product Affiliate Link")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Create Product Affiliate Link")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to create product affiliate link")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message":                "Product affiliate link created successfully",
		"product_affiliate_link": data,
	}
	c.responseHelper.SendCreated(ginCtx, responseData)
}

// GetAffiliateLinksbyProductId godoc
// @Summary Get affiliate links by product ID
// @Description Get all affiliate links for a specific product
// @Tags Product Affiliate Links
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} dto.APIResponse{data=object{message=string,product_affiliate_links=[]dto.ProductAffiliateLinkResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid product ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/{id}/affiliate [get]
func (c *ProductController) GetAffiliateLinksbyProductId(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	data, err := c.s.GetAffiliateLinksbyProductId(reqCtx, uuId)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Get Affiliate Links by Product ID")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Get Affiliate Links by Product ID")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to get affiliate links by product ID")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message":                 "Product affiliate links retrieved successfully",
		"product_affiliate_links": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// UpdateProductAffiliateLink godoc
// @Summary Update product affiliate link
// @Description Update an existing affiliate link for a product
// @Tags Product Affiliate Links
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Param affiliateId path string true "Affiliate Link ID (UUID)"
// @Param payload body dto.UpdateProductAffiliateLinkRequest true "Affiliate link update details"
// @Success 200 {object} dto.APIResponse{data=object{message=string,product_affiliate_link=dto.ProductAffiliateLinkResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid payload, product ID, or affiliate ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product or affiliate link not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/{id}/affiliate/{affiliateId} [put]
func (c *ProductController) UpdateProductAffiliateLink(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	affiliateId := ginCtx.Param("affiliateId")
	if affiliateId == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "affiliateId", "Affiliate ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	affiliateUuId, err := uuid.Parse(affiliateId)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "affiliateId", "Invalid Affiliate ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	var req dto.UpdateProductAffiliateLinkRequest
	if err := ginCtx.ShouldBindJSON(&req); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	data, err := c.s.UpdateProductAffiliateLink(reqCtx, uuId, affiliateUuId, req)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Update Product Affiliate Link")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Update Product Affiliate Link")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to update product affiliate link")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message":                "Product affiliate link updated successfully",
		"product_affiliate_link": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// DeleteProductAffiliateLink godoc
// @Summary Delete product affiliate link
// @Description Delete an affiliate link by ID
// @Tags Product Affiliate Links
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Param affiliateId path string true "Affiliate Link ID (UUID)"
// @Success 200 {object} dto.APIResponse{data=object{message=string}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid affiliate ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Affiliate link not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/{id}/affiliate/{affiliateId} [delete]
func (c *ProductController) DeleteProductAffiliateLink(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	affiliateId := ginCtx.Param("affiliateId")
	if affiliateId == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "affiliateId", "Affiliate ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	affiliateUuId, err := uuid.Parse(affiliateId)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "affiliateId", "Invalid Affiliate ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	err = c.s.DeleteProductAffiliateLink(reqCtx, affiliateUuId)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Delete Product Affiliate Link")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Delete Product Affiliate Link")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to delete product affiliate link")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message": "Product affiliate link deleted successfully",
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// Article Product Relations Controllers

// AddProductToArticle godoc
// @Summary Add product to article
// @Description Associate a product with an article
// @Tags Product Article Relations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Param articleId path string true "Article ID (UUID)"
// @Success 200 {object} dto.APIResponse{data=object{message=string}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid product ID or article ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product or article not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/{id}/article/{articleId} [post]
func (c *ProductController) AddProductToArticle(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	productUuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	articleId := ginCtx.Param("articleId")
	if articleId == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "articleId", "Article ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	articleUuId, err := uuid.Parse(articleId)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "articleId", "Invalid Article ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	err = c.s.AddProductToArticle(reqCtx, articleUuId, productUuId)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Add Product To Article")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Add Product To Article")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to add product to article")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message": "Product added to article successfully",
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// RemoveProductFromArticle godoc
// @Summary Remove product from article
// @Description Remove the association between a product and an article
// @Tags Product Article Relations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID (UUID)"
// @Param articleId path string true "Article ID (UUID)"
// @Success 200 {object} dto.APIResponse{data=object{message=string}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid product ID or article ID"
// @Failure 401 {object} dto.APIResponse{error=dto.ErrorResponse} "Unauthorized"
// @Failure 404 {object} dto.APIResponse{error=dto.ErrorResponse} "Product or article not found"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/{id}/article/{articleId} [delete]
func (c *ProductController) RemoveProductFromArticle(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	productUuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	articleId := ginCtx.Param("articleId")
	if articleId == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "articleId", "Article ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	articleUuId, err := uuid.Parse(articleId)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "articleId", "Invalid Article ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	err = c.s.RemoveProductFromArticle(reqCtx, articleUuId, productUuId)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Remove Product From Article")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Remove Product From Article")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to remove product from article")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message": "Product removed from article successfully",
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// GetProductsByArticleId godoc
// @Summary Get products by article ID
// @Description Get all products associated with a specific article
// @Tags Products
// @Produce json
// @Param id path string true "Article ID (UUID)"
// @Success 200 {object} dto.APIResponse{data=object{message=string,products=[]dto.ProductListResponse}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid article ID"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/article/{id} [get]
func (c *ProductController) GetProductsByArticleId(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Article ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Article ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 10
	
	if pageStr := ginCtx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := ginCtx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	data, err := c.s.GetProductsByArticleIdWithPagination(reqCtx, uuId, page, limit)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Get Products By Article ID")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Get Products By Article ID")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to get products by article ID")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message":  "Products retrieved successfully",
		"products": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}

// GetArticlesByProductId godoc
// @Summary Get articles by product ID
// @Description Get all articles associated with a specific product
// @Tags Product Article Relations
// @Produce json
// @Param id path string true "Product ID (UUID)"
// @Success 200 {object} dto.APIResponse{data=object{message=string,articles=[]model.Article}}
// @Failure 400 {object} dto.APIResponse{error=dto.ErrorResponse} "Invalid product ID"
// @Failure 408 {object} dto.APIResponse{error=dto.ErrorResponse} "Request timeout"
// @Failure 500 {object} dto.APIResponse{error=dto.ErrorResponse} "Internal server error"
// @Router /products/{id}/articles [get]
func (c *ProductController) GetArticlesByProductId(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	uuId, err := uuid.Parse(id)
	if err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Invalid Product ID format: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	// Parse pagination parameters
	page := 1
	limit := 10
	
	if pageStr := ginCtx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	
	if limitStr := ginCtx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	data, err := c.s.GetArticlesByProductIdWithPagination(reqCtx, uuId, page, limit)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Get Articles By Product ID")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Get Articles By Product ID")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to get articles by product ID")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	responseData := gin.H{
		"message":  "Articles retrieved successfully",
		"articles": data,
	}
	c.responseHelper.SendSuccess(ginCtx, responseData)
}