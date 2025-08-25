package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/repository"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ProductService interface {
	// Product Categories
	CreateProductCategory(ctx context.Context, req dto.CreateProductCategoryRequest) (dto.ProductCategoryResponse, error)
	GetAllProductCategoriesWithPagination(ctx context.Context, page, limit int) (PaginationResult, error)
	GetProductCategoryById(ctx context.Context, id uuid.UUID) (dto.ProductCategoryResponse, error)
	GetProductCategoryBySlug(ctx context.Context, slug string) (dto.ProductCategoryResponse, error)
	UpdateProductCategory(ctx context.Context, id uuid.UUID, req dto.UpdateProductCategoryRequest) (dto.ProductCategoryResponse, error)
	DeleteProductCategory(ctx context.Context, id uuid.UUID) error

	// Products
	CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.ProductResponse, error)
	GetAllProductsWithPagination(ctx context.Context, page, limit int) (PaginationResult, error)
	GetProductById(ctx context.Context, id uuid.UUID) (dto.ProductResponse, error)
	GetProductsByCategoryWithPagination(ctx context.Context, categoryId uuid.UUID, page, limit int) (PaginationResult, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (dto.ProductResponse, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error

	// Product Affiliate Links
	CreateProductAffiliateLink(ctx context.Context, productId uuid.UUID, req dto.CreateProductAffiliateLinkRequest) (dto.ProductAffiliateLinkResponse, error)
	GetAffiliateLinksbyProductId(ctx context.Context, productId uuid.UUID) ([]dto.ProductAffiliateLinkResponse, error)
	UpdateProductAffiliateLink(ctx context.Context, productId, affiliateId uuid.UUID, req dto.UpdateProductAffiliateLinkRequest) (dto.ProductAffiliateLinkResponse, error)
	DeleteProductAffiliateLink(ctx context.Context, id uuid.UUID) error

	// Article Product Relations
	AddProductToArticle(ctx context.Context, articleId, productId uuid.UUID) error
	RemoveProductFromArticle(ctx context.Context, articleId, productId uuid.UUID) error
	GetProductsByArticleIdWithPagination(ctx context.Context, articleId uuid.UUID, page, limit int) (PaginationResult, error)
	GetArticlesByProductIdWithPagination(ctx context.Context, productId uuid.UUID, page, limit int) (PaginationResult, error)
}

type productService struct {
	productRepo repository.ProductRepository
	validation  ValidationService
	pagination  PaginationService
}

// Helper functions for conversion
func (s *productService) modelToProductCategoryResponse(category model.ProductCategory) dto.ProductCategoryResponse {
	return dto.ProductCategoryResponse{
		Id:          category.Id,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func (s *productService) modelToProductAffiliateLinkResponse(link model.ProductAffiliateLink) dto.ProductAffiliateLinkResponse {
	return dto.ProductAffiliateLinkResponse{
		Id:           link.Id,
		PlatformName: link.PlatformName,
		Url:          link.Url,
		CreatedAt:    link.CreatedAt,
		UpdatedAt:    link.UpdatedAt,
	}
}

func (s *productService) modelToProductResponse(product model.Product, affiliateLinks []model.ProductAffiliateLink) dto.ProductResponse {
	var categoryResponse *dto.ProductCategoryResponse
	if product.ProductCategory != nil {
		resp := s.modelToProductCategoryResponse(*product.ProductCategory)
		categoryResponse = &resp
	}

	var affiliateLinkResponses []dto.ProductAffiliateLinkResponse
	for _, link := range affiliateLinks {
		affiliateLinkResponses = append(affiliateLinkResponses, s.modelToProductAffiliateLinkResponse(link))
	}

	return dto.ProductResponse{
		Id:                product.Id,
		ProductCategoryId: product.ProductCategoryId,
		ProductCategory:   categoryResponse,
		Name:              product.Name,
		Description:       product.Description,
		ImageUrl:          product.ImageUrl,
		IsActive:          product.IsActive,
		AffiliateLinks:    affiliateLinkResponses,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
	}
}

func (s *productService) modelToProductListResponse(product model.Product) dto.ProductListResponse {
	var categoryResponse *dto.ProductCategoryResponse
	if product.ProductCategory != nil {
		resp := s.modelToProductCategoryResponse(*product.ProductCategory)
		categoryResponse = &resp
	}

	return dto.ProductListResponse{
		Id:                product.Id,
		ProductCategoryId: product.ProductCategoryId,
		ProductCategory:   categoryResponse,
		Name:              product.Name,
		Description:       product.Description,
		ImageUrl:          product.ImageUrl,
		IsActive:          product.IsActive,
		CreatedAt:         product.CreatedAt,
		UpdatedAt:         product.UpdatedAt,
	}
}

// Product Category implementations
func (s *productService) CreateProductCategory(ctx context.Context, req dto.CreateProductCategoryRequest) (dto.ProductCategoryResponse, error) {
	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		return dto.ProductCategoryResponse{}, errors.New("category name is required")
	}

	// Generate slug if not provided
	slug := req.Slug
	if strings.TrimSpace(slug) == "" {
		slug = s.generateSlug(req.Name)
	} else {
		slug = s.generateSlug(slug)
	}

	// Validate slug format
	if !s.isValidSlug(slug) {
		return dto.ProductCategoryResponse{}, errors.New("invalid slug format")
	}

	// Create model
	category := model.ProductCategory{
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	createdCategory, err := s.productRepo.CreateProductCategory(ctx, category)
	if err != nil {
		return dto.ProductCategoryResponse{}, err
	}

	return s.modelToProductCategoryResponse(createdCategory), nil
}

func (s *productService) GetAllProductCategoriesWithPagination(ctx context.Context, page, limit int) (PaginationResult, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return PaginationResult{}, ctx.Err()
	default:
	}

	// Parse and validate pagination query
	query, err := s.pagination.ParseQuery(ctx, page, limit, "created_at", "desc")
	if err != nil {
		return PaginationResult{}, err
	}

	// Get paginated categories from repository
	categories, total, repoErr := s.productRepo.GetAllProductCategoriesWithPagination(ctx, query.Offset, query.Limit)
	if repoErr != nil {
		if ctx.Err() != nil {
			return PaginationResult{}, ctx.Err()
		}
		return PaginationResult{}, repoErr
	}

	// Convert to response DTOs
	var responses []dto.ProductCategoryResponse
	for _, category := range categories {
		responses = append(responses, s.modelToProductCategoryResponse(category))
	}

	// Create pagination result
	result, paginationErr := s.pagination.Paginate(ctx, responses, total, query)
	if paginationErr != nil {
		return PaginationResult{}, paginationErr
	}

	return result, nil
}

func (s *productService) GetProductCategoryById(ctx context.Context, id uuid.UUID) (dto.ProductCategoryResponse, error) {
	if id == uuid.Nil {
		return dto.ProductCategoryResponse{}, errors.New("invalid category ID")
	}

	category, err := s.productRepo.GetProductCategoryById(ctx, id)
	if err != nil {
		return dto.ProductCategoryResponse{}, err
	}

	return s.modelToProductCategoryResponse(category), nil
}

func (s *productService) GetProductCategoryBySlug(ctx context.Context, slug string) (dto.ProductCategoryResponse, error) {
	if strings.TrimSpace(slug) == "" {
		return dto.ProductCategoryResponse{}, errors.New("slug is required")
	}

	category, err := s.productRepo.GetProductCategoryBySlug(ctx, slug)
	if err != nil {
		return dto.ProductCategoryResponse{}, err
	}

	return s.modelToProductCategoryResponse(category), nil
}

func (s *productService) UpdateProductCategory(ctx context.Context, id uuid.UUID, req dto.UpdateProductCategoryRequest) (dto.ProductCategoryResponse, error) {
	if id == uuid.Nil {
		return dto.ProductCategoryResponse{}, errors.New("invalid category ID")
	}

	// Get existing category
	existingCategory, err := s.productRepo.GetProductCategoryById(ctx, id)
	if err != nil {
		return dto.ProductCategoryResponse{}, err
	}

	// Update fields if provided
	if req.Name != nil {
		existingCategory.Name = *req.Name
	}
	if req.Slug != nil {
		existingCategory.Slug = s.generateSlug(*req.Slug)
	} else if req.Name != nil {
		// Auto-generate slug from name if name is updated but slug is not provided
		existingCategory.Slug = s.generateSlug(*req.Name)
	}
	if req.Description != nil {
		existingCategory.Description = req.Description
	}

	// Validate slug format
	if !s.isValidSlug(existingCategory.Slug) {
		return dto.ProductCategoryResponse{}, errors.New("invalid slug format")
	}

	existingCategory.UpdatedAt = time.Now()

	updatedCategory, err := s.productRepo.UpdateProductCategory(ctx, existingCategory)
	if err != nil {
		return dto.ProductCategoryResponse{}, err
	}

	return s.modelToProductCategoryResponse(updatedCategory), nil
}

func (s *productService) DeleteProductCategory(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid category ID")
	}
	return s.productRepo.DeleteProductCategory(ctx, id)
}

// Product implementations
func (s *productService) CreateProduct(ctx context.Context, req dto.CreateProductRequest) (dto.ProductResponse, error) {
	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		return dto.ProductResponse{}, errors.New("product name is required")
	}

	if req.ProductCategoryId == nil {
		return dto.ProductResponse{}, errors.New("product category ID is required")
	}

	// Validate image URL if provided
	if req.ImageUrl != nil && strings.TrimSpace(*req.ImageUrl) != "" {
		if !s.isValidURL(*req.ImageUrl) {
			return dto.ProductResponse{}, errors.New("invalid image URL format")
		}
	}

	// Set default value for IsActive
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	// Create model
	product := model.Product{
		ProductCategoryId: req.ProductCategoryId,
		Name:              req.Name,
		Description:       req.Description,
		ImageUrl:          req.ImageUrl,
		IsActive:          isActive,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	createdProduct, err := s.productRepo.CreateProduct(ctx, product)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Get product with affiliate links (empty for new product)
	return s.modelToProductResponse(createdProduct, []model.ProductAffiliateLink{}), nil
}

func (s *productService) GetAllProductsWithPagination(ctx context.Context, page, limit int) (PaginationResult, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return PaginationResult{}, ctx.Err()
	default:
	}

	// Parse and validate pagination query
	query, err := s.pagination.ParseQuery(ctx, page, limit, "created_at", "desc")
	if err != nil {
		return PaginationResult{}, err
	}

	// Get paginated products from repository
	products, total, repoErr := s.productRepo.GetAllProductsWithPagination(ctx, query.Offset, query.Limit)
	if repoErr != nil {
		if ctx.Err() != nil {
			return PaginationResult{}, ctx.Err()
		}
		return PaginationResult{}, repoErr
	}

	// Convert to response DTOs
	var responses []dto.ProductListResponse
	for _, product := range products {
		responses = append(responses, s.modelToProductListResponse(product))
	}

	// Create pagination result
	result, paginationErr := s.pagination.Paginate(ctx, responses, total, query)
	if paginationErr != nil {
		return PaginationResult{}, paginationErr
	}

	return result, nil
}

func (s *productService) GetProductById(ctx context.Context, id uuid.UUID) (dto.ProductResponse, error) {
	if id == uuid.Nil {
		return dto.ProductResponse{}, errors.New("invalid product ID")
	}

	// Get product with affiliate links
	product, affiliateLinks, err := s.productRepo.GetProductByIdWithAffiliateLinks(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return s.modelToProductResponse(product, affiliateLinks), nil
}

func (s *productService) GetProductsByCategoryWithPagination(ctx context.Context, categoryId uuid.UUID, page, limit int) (PaginationResult, error) {
	if categoryId == uuid.Nil {
		return PaginationResult{}, errors.New("invalid category ID")
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return PaginationResult{}, ctx.Err()
	default:
	}

	// Parse and validate pagination query
	query, err := s.pagination.ParseQuery(ctx, page, limit, "created_at", "desc")
	if err != nil {
		return PaginationResult{}, err
	}

	// Get products with affiliate links and pagination
	products, affiliateLinksMap, total, repoErr := s.productRepo.GetProductsByCategoryWithAffiliateLinksAndPagination(ctx, categoryId, query.Offset, query.Limit)
	if repoErr != nil {
		if ctx.Err() != nil {
			return PaginationResult{}, ctx.Err()
		}
		return PaginationResult{}, repoErr
	}

	// Convert to response DTOs
	var responses []dto.ProductResponse
	for _, product := range products {
		affiliateLinks := affiliateLinksMap[product.Id]
		responses = append(responses, s.modelToProductResponse(product, affiliateLinks))
	}

	// Create pagination result
	result, paginationErr := s.pagination.Paginate(ctx, responses, total, query)
	if paginationErr != nil {
		return PaginationResult{}, paginationErr
	}

	return result, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uuid.UUID, req dto.UpdateProductRequest) (dto.ProductResponse, error) {
	if id == uuid.Nil {
		return dto.ProductResponse{}, errors.New("invalid product ID")
	}

	// Get existing product
	existingProduct, err := s.productRepo.GetProductById(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Update fields if provided
	if req.ProductCategoryId != nil {
		existingProduct.ProductCategoryId = req.ProductCategoryId
	}
	if req.Name != nil {
		existingProduct.Name = *req.Name
	}
	if req.Description != nil {
		existingProduct.Description = req.Description
	}
	if req.ImageUrl != nil {
		existingProduct.ImageUrl = req.ImageUrl
	}
	if req.IsActive != nil {
		existingProduct.IsActive = *req.IsActive
	}

	// Validate image URL if provided
	if existingProduct.ImageUrl != nil && strings.TrimSpace(*existingProduct.ImageUrl) != "" {
		if !s.isValidURL(*existingProduct.ImageUrl) {
			return dto.ProductResponse{}, errors.New("invalid image URL format")
		}
	}

	existingProduct.UpdatedAt = time.Now()

	updatedProduct, err := s.productRepo.UpdateProduct(ctx, existingProduct)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	// Get affiliate links
	affiliateLinks, err := s.productRepo.GetAffiliateLinksbyProductId(ctx, id)
	if err != nil {
		return dto.ProductResponse{}, err
	}

	return s.modelToProductResponse(updatedProduct, affiliateLinks), nil
}

func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid product ID")
	}
	return s.productRepo.DeleteProduct(ctx, id)
}

// Product Affiliate Link implementations
func (s *productService) CreateProductAffiliateLink(ctx context.Context, productId uuid.UUID, req dto.CreateProductAffiliateLinkRequest) (dto.ProductAffiliateLinkResponse, error) {
	if productId == uuid.Nil {
		return dto.ProductAffiliateLinkResponse{}, errors.New("invalid product ID")
	}

	if strings.TrimSpace(req.PlatformName) == "" {
		return dto.ProductAffiliateLinkResponse{}, errors.New("platform name is required")
	}

	if strings.TrimSpace(req.Url) == "" {
		return dto.ProductAffiliateLinkResponse{}, errors.New("URL is required")
	}

	if !s.isValidURL(req.Url) {
		return dto.ProductAffiliateLinkResponse{}, errors.New("invalid URL format")
	}

	// Create model
	affiliateLink := model.ProductAffiliateLink{
		ProductId:    productId,
		PlatformName: req.PlatformName,
		Url:          req.Url,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdLink, err := s.productRepo.CreateProductAffiliateLink(ctx, affiliateLink)
	if err != nil {
		return dto.ProductAffiliateLinkResponse{}, err
	}

	return s.modelToProductAffiliateLinkResponse(createdLink), nil
}

func (s *productService) GetAffiliateLinksbyProductId(ctx context.Context, productId uuid.UUID) ([]dto.ProductAffiliateLinkResponse, error) {
	if productId == uuid.Nil {
		return nil, errors.New("invalid product ID")
	}

	links, err := s.productRepo.GetAffiliateLinksbyProductId(ctx, productId)
	if err != nil {
		return nil, err
	}

	var responses []dto.ProductAffiliateLinkResponse
	for _, link := range links {
		responses = append(responses, s.modelToProductAffiliateLinkResponse(link))
	}

	return responses, nil
}

func (s *productService) UpdateProductAffiliateLink(ctx context.Context, productId, affiliateId uuid.UUID, req dto.UpdateProductAffiliateLinkRequest) (dto.ProductAffiliateLinkResponse, error) {
	if productId == uuid.Nil {
		return dto.ProductAffiliateLinkResponse{}, errors.New("invalid product ID")
	}
	if affiliateId == uuid.Nil {
		return dto.ProductAffiliateLinkResponse{}, errors.New("invalid affiliate link ID")
	}

	// Get existing affiliate link
	links, err := s.productRepo.GetAffiliateLinksbyProductId(ctx, productId)
	if err != nil {
		return dto.ProductAffiliateLinkResponse{}, err
	}

	var existingLink *model.ProductAffiliateLink
	for _, link := range links {
		if link.Id == affiliateId {
			existingLink = &link
			break
		}
	}

	if existingLink == nil {
		return dto.ProductAffiliateLinkResponse{}, errors.New("affiliate link not found")
	}

	// Update fields if provided
	if req.PlatformName != nil {
		existingLink.PlatformName = *req.PlatformName
	}
	if req.Url != nil {
		existingLink.Url = *req.Url
	}

	// Validate URL if provided
	if strings.TrimSpace(existingLink.Url) != "" && !s.isValidURL(existingLink.Url) {
		return dto.ProductAffiliateLinkResponse{}, errors.New("invalid URL format")
	}

	existingLink.UpdatedAt = time.Now()

	updatedLink, err := s.productRepo.UpdateProductAffiliateLink(ctx, *existingLink)
	if err != nil {
		return dto.ProductAffiliateLinkResponse{}, err
	}

	return s.modelToProductAffiliateLinkResponse(updatedLink), nil
}

func (s *productService) DeleteProductAffiliateLink(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid affiliate link ID")
	}
	return s.productRepo.DeleteProductAffiliateLink(ctx, id)
}

// Article Product Relations
func (s *productService) AddProductToArticle(ctx context.Context, articleId, productId uuid.UUID) error {
	if articleId == uuid.Nil {
		return errors.New("invalid article ID")
	}
	if productId == uuid.Nil {
		return errors.New("invalid product ID")
	}
	return s.productRepo.AddProductToArticle(ctx, articleId, productId)
}

func (s *productService) RemoveProductFromArticle(ctx context.Context, articleId, productId uuid.UUID) error {
	if articleId == uuid.Nil {
		return errors.New("invalid article ID")
	}
	if productId == uuid.Nil {
		return errors.New("invalid product ID")
	}
	return s.productRepo.RemoveProductFromArticle(ctx, articleId, productId)
}

func (s *productService) GetProductsByArticleIdWithPagination(ctx context.Context, articleId uuid.UUID, page, limit int) (PaginationResult, error) {
	if articleId == uuid.Nil {
		return PaginationResult{}, errors.New("invalid article ID")
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return PaginationResult{}, ctx.Err()
	default:
	}

	// Parse and validate pagination query
	query, err := s.pagination.ParseQuery(ctx, page, limit, "created_at", "desc")
	if err != nil {
		return PaginationResult{}, err
	}

	// Get products with affiliate links and pagination
	products, affiliateLinksMap, total, repoErr := s.productRepo.GetProductsByArticleIdWithAffiliateLinksAndPagination(ctx, articleId, query.Offset, query.Limit)
	if repoErr != nil {
		if ctx.Err() != nil {
			return PaginationResult{}, ctx.Err()
		}
		return PaginationResult{}, repoErr
	}

	// Convert to response DTOs
	var responses []dto.ProductResponse
	for _, product := range products {
		affiliateLinks := affiliateLinksMap[product.Id]
		responses = append(responses, s.modelToProductResponse(product, affiliateLinks))
	}

	// Create pagination result
	result, paginationErr := s.pagination.Paginate(ctx, responses, total, query)
	if paginationErr != nil {
		return PaginationResult{}, paginationErr
	}

	return result, nil
}

func (s *productService) GetArticlesByProductIdWithPagination(ctx context.Context, productId uuid.UUID, page, limit int) (PaginationResult, error) {
	if productId == uuid.Nil {
		return PaginationResult{}, errors.New("invalid product ID")
	}

	// Check context cancellation
	select {
	case <-ctx.Done():
		return PaginationResult{}, ctx.Err()
	default:
	}

	// Parse and validate pagination query
	query, err := s.pagination.ParseQuery(ctx, page, limit, "created_at", "desc")
	if err != nil {
		return PaginationResult{}, err
	}

	// Get articles with pagination
	articles, total, repoErr := s.productRepo.GetArticlesByProductIdWithPagination(ctx, productId, query.Offset, query.Limit)
	if repoErr != nil {
		if ctx.Err() != nil {
			return PaginationResult{}, ctx.Err()
		}
		return PaginationResult{}, repoErr
	}

	// Create pagination result
	result, paginationErr := s.pagination.Paginate(ctx, articles, total, query)
	if paginationErr != nil {
		return PaginationResult{}, paginationErr
	}

	return result, nil
}

// Helper functions
func (s *productService) generateSlug(input string) string {
	// Convert to lowercase
	slug := strings.ToLower(input)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}

func (s *productService) isValidSlug(slug string) bool {
	// Slug should only contain lowercase letters, numbers, and hyphens
	// Should not start or end with hyphen
	match, _ := regexp.MatchString(`^[a-z0-9]+(-[a-z0-9]+)*$`, slug)
	return match
}

func (s *productService) isValidURL(url string) bool {
	// Basic URL validation
	match, _ := regexp.MatchString(`^https?://[^\s/$.?#].[^\s]*$`, url)
	return match
}

func NewProductService(productRepo repository.ProductRepository, validation ValidationService, pagination PaginationService) ProductService {
	return &productService{
		productRepo: productRepo,
		validation:  validation,
		pagination:  pagination,
	}
}