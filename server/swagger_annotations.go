package main

// This file contains all Swagger annotations for Product API endpoints
// Copy these annotations to the respective controller functions

/*
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
*/