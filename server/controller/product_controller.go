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
		routerProductCat := c.rg.Group("/product-categories")
		routerProductCat.GET("/", c.GetAllProductCategories)
		routerProductCat.GET("/:id", c.GetProductCategoryById)
		routerProductCat.GET("/s/:slug", c.GetProductCategoryBySlug)

		routerPCAuth := routerProductCat.Group("/", c.mD.CheckToken("admin"))
		routerPCAuth.POST("/", c.CreateProductCategory)
		routerPCAuth.PUT("/:id", c.UpdateProductCategory)
		routerPCAuth.DELETE("/:id", c.DeleteProductCategory)
	}

	{
		routerProduct := c.rg.Group("/products")
		// routerProduct.GET("/", c.GetAllProducts)
		routerProduct.GET("/:id", c.GetProductById)
		routerProduct.GET("/category/:id", c.GetProductsByCategory)
		routerProduct.GET("/article/:id", c.GetProductsByArticleId)

		routerPAuth := routerProduct.Group("/", c.mD.CheckToken("admin"))
		routerPAuth.POST("/", c.CreateProduct)
		routerPAuth.PUT("/:id", c.UpdateProduct)
		routerPAuth.DELETE("/:id", c.DeleteProduct)
		routerPAuth.POST("/:id/article/:articleId", c.AddProductToArticle)
		routerPAuth.DELETE("/:id/article/:articleId", c.RemoveProductFromArticle)

		routerAffiliate := routerProduct.Group("/:id/affiliate", c.mD.CheckToken("admin"))
		routerAffiliate.POST("/", c.CreateProductAffiliateLink)
		routerAffiliate.GET("/", c.GetAffiliateLinksbyProductId)
		routerAffiliate.PUT("/:affiliateId", c.UpdateProductAffiliateLink)
		routerAffiliate.DELETE("/:affiliateId", c.DeleteProductAffiliateLink)

	}

}

// func (c *ProductController) GetAllProducts(ginCtx *gin.Context) {
// 	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
// 	defer cancel()

// 	page := utils.StringToInt(ginCtx.Query("page"), 1)
// 	limit := utils.StringToInt(ginCtx.Query("limit"), 10)

// 	data, total, totalPages, err := c.s.GetAllProductsWithPagination(reqCtx, page, limit)
// 	if err != nil {
// 		if reqCtx.Err() == context.DeadlineExceeded {
// 			appErr := c.errorHandler.TimeoutError(reqCtx, "Get All Products")
// 			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
// 			return
// 		}
// 		if reqCtx.Err() == context.Canceled {
// 			appErr := c.errorHandler.CancellationError(reqCtx, "Get All Products")
// 			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
// 			return
// 		}
// 		if appErr, ok := err.(*utils.AppError); ok {
// 			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
// 			return
// 		}

// 		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to get products")
// 		appErr.StatusCode = 500
// 		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
// 		return
// 	}

// 	responseData := gin.H{
// 		"message":  "Products retrieved successfully",
// 		"products": data,
// 		"pagination": gin.H{
// 			"total_items": total,
// 			"total_pages": totalPages,
// 			"current_page": page,
// 			"limit": limit,
// 		},
// 	}
// 	c.responseHelper.SendSuccess(ginCtx, responseData)
// }

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

	data, err := c.s.GetProductsByCategory(reqCtx, uuId)
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

func (c *ProductController) UpdateProduct(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.Product
	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

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

	payload.Id = uuId

	data, err := c.s.UpdateProduct(reqCtx, payload)
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

func (c *ProductController) CreateProductAffiliateLink(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.ProductAffiliateLink
	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

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

	payload.ProductId = uuId

	data, err := c.s.CreateProductAffiliateLink(reqCtx, payload)
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

func (c *ProductController) UpdateProductAffiliateLink(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.ProductAffiliateLink
	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

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

	payload.ProductId = uuId
	payload.Id = affiliateUuId

	data, err := c.s.UpdateProductAffiliateLink(reqCtx, payload)
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

func (c *ProductController) DeleteProductAffiliateLink(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	id := ginCtx.Param("id")
	if id == "" {
		appErr := c.errorHandler.ValidationError(reqCtx, "id", "Product ID is required")
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return
	}

	_, err := uuid.Parse(id)
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

	data, err := c.s.GetProductsByArticleId(reqCtx, uuId)
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

	data, err := c.s.GetArticlesByProductId(reqCtx, uuId)
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

func (c *ProductController) CreateProduct(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.Product
	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)

		return
	}

	data, err := c.s.CreateProduct(reqCtx, payload)
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
