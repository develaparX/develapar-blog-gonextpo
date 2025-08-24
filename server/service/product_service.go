package service

import (
	"context"
	"develapar-server/model"
	"develapar-server/repository"
	"errors"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type ProductService interface {
	// Product Categories
	CreateProductCategory(ctx context.Context, payload model.ProductCategory) (model.ProductCategory, error)
	GetAllProductCategories(ctx context.Context) ([]model.ProductCategory, error)
	GetProductCategoryById(ctx context.Context, id uuid.UUID) (model.ProductCategory, error)
	GetProductCategoryBySlug(ctx context.Context, slug string) (model.ProductCategory, error)
	UpdateProductCategory(ctx context.Context, payload model.ProductCategory) (model.ProductCategory, error)
	DeleteProductCategory(ctx context.Context, id uuid.UUID) error

	// Products
	CreateProduct(ctx context.Context, payload model.Product) (model.Product, error)
	GetAllProducts(ctx context.Context) ([]model.Product, error)
	// GetAllProductsWithPagination(ctx context.Context, page, limit int) ([]model.Product, int, int, error)
	GetProductById(ctx context.Context, id uuid.UUID) (model.Product, error)
	GetProductsByCategory(ctx context.Context, categoryId uuid.UUID) ([]model.Product, error)
	UpdateProduct(ctx context.Context, payload model.Product) (model.Product, error)
	DeleteProduct(ctx context.Context, id uuid.UUID) error

	// Product Affiliate Links
	CreateProductAffiliateLink(ctx context.Context, payload model.ProductAffiliateLink) (model.ProductAffiliateLink, error)
	GetAffiliateLinksbyProductId(ctx context.Context, productId uuid.UUID) ([]model.ProductAffiliateLink, error)
	UpdateProductAffiliateLink(ctx context.Context, payload model.ProductAffiliateLink) (model.ProductAffiliateLink, error)
	DeleteProductAffiliateLink(ctx context.Context, id uuid.UUID) error

	// Article Product Relations
	AddProductToArticle(ctx context.Context, articleId, productId uuid.UUID) error
	RemoveProductFromArticle(ctx context.Context, articleId, productId uuid.UUID) error
	GetProductsByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Product, error)
	GetArticlesByProductId(ctx context.Context, productId uuid.UUID) ([]model.Article, error)
}

type productService struct {
	productRepo repository.ProductRepository
	validation  ValidationService
	pagination  PaginationService
}

// CreateProductCategory implements ProductService
func (s *productService) CreateProductCategory(ctx context.Context, payload model.ProductCategory) (model.ProductCategory, error) {
	// Validate required fields
	if strings.TrimSpace(payload.Name) == "" {
		return model.ProductCategory{}, errors.New("category name is required")
	}

	// Generate slug if not provided
	if strings.TrimSpace(payload.Slug) == "" {
		payload.Slug = s.generateSlug(payload.Name)
	} else {
		payload.Slug = s.generateSlug(payload.Slug)
	}

	// Validate slug format
	if !s.isValidSlug(payload.Slug) {
		return model.ProductCategory{}, errors.New("invalid slug format")
	}

	return s.productRepo.CreateProductCategory(ctx, payload)
}

// GetAllProductCategories implements ProductService
func (s *productService) GetAllProductCategories(ctx context.Context) ([]model.ProductCategory, error) {
	return s.productRepo.GetAllProductCategories(ctx)
}

// GetProductCategoryById implements ProductService
func (s *productService) GetProductCategoryById(ctx context.Context, id uuid.UUID) (model.ProductCategory, error) {
	if id == uuid.Nil {
		return model.ProductCategory{}, errors.New("invalid category ID")
	}
	return s.productRepo.GetProductCategoryById(ctx, id)
}

// GetProductCategoryBySlug implements ProductService
func (s *productService) GetProductCategoryBySlug(ctx context.Context, slug string) (model.ProductCategory, error) {
	if strings.TrimSpace(slug) == "" {
		return model.ProductCategory{}, errors.New("slug is required")
	}
	return s.productRepo.GetProductCategoryBySlug(ctx, slug)
}

// UpdateProductCategory implements ProductService
func (s *productService) UpdateProductCategory(ctx context.Context, payload model.ProductCategory) (model.ProductCategory, error) {
	if payload.Id == uuid.Nil {
		return model.ProductCategory{}, errors.New("invalid category ID")
	}

	if strings.TrimSpace(payload.Name) == "" {
		return model.ProductCategory{}, errors.New("category name is required")
	}

	// Generate slug if not provided
	if strings.TrimSpace(payload.Slug) == "" {
		payload.Slug = s.generateSlug(payload.Name)
	} else {
		payload.Slug = s.generateSlug(payload.Slug)
	}

	// Validate slug format
	if !s.isValidSlug(payload.Slug) {
		return model.ProductCategory{}, errors.New("invalid slug format")
	}

	return s.productRepo.UpdateProductCategory(ctx, payload)
}

// DeleteProductCategory implements ProductService
func (s *productService) DeleteProductCategory(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid category ID")
	}
	return s.productRepo.DeleteProductCategory(ctx, id)
}

// CreateProduct implements ProductService
func (s *productService) CreateProduct(ctx context.Context, payload model.Product) (model.Product, error) {
	// Validate required fields
	if strings.TrimSpace(payload.Name) == "" {
		return model.Product{}, errors.New("product name is required")
	}

	// Validate image URL if provided
	if payload.ImageUrl != nil && strings.TrimSpace(*payload.ImageUrl) != "" {
		if !s.isValidURL(*payload.ImageUrl) {
			return model.Product{}, errors.New("invalid image URL format")
		}
	}

	return s.productRepo.CreateProduct(ctx, payload)
}

// GetAllProducts implements ProductService
func (s *productService) GetAllProducts(ctx context.Context) ([]model.Product, error) {
	return s.productRepo.GetAllProducts(ctx)
}

// func (s *productService) GetAllProductsWithPagination(ctx context.Context, page, limit int) ([]model.Product, int, int, error) {
// 	offset := s.pagination.CalculateOffset(page, limit)

// 	products, totalCount, err := s.productRepo.GetAllProductsWithPagination(ctx, offset, limit)
// 	if err != nil {
// 		return nil, 0, 0, err
// 	}

// 	totalPages := s.pagination.CalculateTotalPages(totalCount, limit)
// 	return products, totalCount, totalPages, nil
// }

// GetProductById implements ProductService
func (s *productService) GetProductById(ctx context.Context, id uuid.UUID) (model.Product, error) {
	if id == uuid.Nil {
		return model.Product{}, errors.New("invalid product ID")
	}
	return s.productRepo.GetProductById(ctx, id)
}

// GetProductsByCategory implements ProductService
func (s *productService) GetProductsByCategory(ctx context.Context, categoryId uuid.UUID) ([]model.Product, error) {
	if categoryId == uuid.Nil {
		return nil, errors.New("invalid category ID")
	}
	return s.productRepo.GetProductsByCategory(ctx, categoryId)
}

// UpdateProduct implements ProductService
func (s *productService) UpdateProduct(ctx context.Context, payload model.Product) (model.Product, error) {
	if payload.Id == uuid.Nil {
		return model.Product{}, errors.New("invalid product ID")
	}

	if strings.TrimSpace(payload.Name) == "" {
		return model.Product{}, errors.New("product name is required")
	}

	// Validate image URL if provided
	if payload.ImageUrl != nil && strings.TrimSpace(*payload.ImageUrl) != "" {
		if !s.isValidURL(*payload.ImageUrl) {
			return model.Product{}, errors.New("invalid image URL format")
		}
	}

	return s.productRepo.UpdateProduct(ctx, payload)
}

// DeleteProduct implements ProductService
func (s *productService) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid product ID")
	}
	return s.productRepo.DeleteProduct(ctx, id)
}

// CreateProductAffiliateLink implements ProductService
func (s *productService) CreateProductAffiliateLink(ctx context.Context, payload model.ProductAffiliateLink) (model.ProductAffiliateLink, error) {
	if payload.ProductId == uuid.Nil {
		return model.ProductAffiliateLink{}, errors.New("invalid product ID")
	}

	if strings.TrimSpace(payload.PlatformName) == "" {
		return model.ProductAffiliateLink{}, errors.New("platform name is required")
	}

	if strings.TrimSpace(payload.Url) == "" {
		return model.ProductAffiliateLink{}, errors.New("URL is required")
	}

	if !s.isValidURL(payload.Url) {
		return model.ProductAffiliateLink{}, errors.New("invalid URL format")
	}

	return s.productRepo.CreateProductAffiliateLink(ctx, payload)
}

// GetAffiliateLinksbyProductId implements ProductService
func (s *productService) GetAffiliateLinksbyProductId(ctx context.Context, productId uuid.UUID) ([]model.ProductAffiliateLink, error) {
	if productId == uuid.Nil {
		return nil, errors.New("invalid product ID")
	}
	return s.productRepo.GetAffiliateLinksbyProductId(ctx, productId)
}

// UpdateProductAffiliateLink implements ProductService
func (s *productService) UpdateProductAffiliateLink(ctx context.Context, payload model.ProductAffiliateLink) (model.ProductAffiliateLink, error) {
	if payload.Id == uuid.Nil {
		return model.ProductAffiliateLink{}, errors.New("invalid affiliate link ID")
	}

	if strings.TrimSpace(payload.PlatformName) == "" {
		return model.ProductAffiliateLink{}, errors.New("platform name is required")
	}

	if strings.TrimSpace(payload.Url) == "" {
		return model.ProductAffiliateLink{}, errors.New("URL is required")
	}

	if !s.isValidURL(payload.Url) {
		return model.ProductAffiliateLink{}, errors.New("invalid URL format")
	}

	return s.productRepo.UpdateProductAffiliateLink(ctx, payload)
}

// DeleteProductAffiliateLink implements ProductService
func (s *productService) DeleteProductAffiliateLink(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return errors.New("invalid affiliate link ID")
	}
	return s.productRepo.DeleteProductAffiliateLink(ctx, id)
}

// AddProductToArticle implements ProductService
func (s *productService) AddProductToArticle(ctx context.Context, articleId, productId uuid.UUID) error {
	if articleId == uuid.Nil {
		return errors.New("invalid article ID")
	}
	if productId == uuid.Nil {
		return errors.New("invalid product ID")
	}
	return s.productRepo.AddProductToArticle(ctx, articleId, productId)
}

// RemoveProductFromArticle implements ProductService
func (s *productService) RemoveProductFromArticle(ctx context.Context, articleId, productId uuid.UUID) error {
	if articleId == uuid.Nil {
		return errors.New("invalid article ID")
	}
	if productId == uuid.Nil {
		return errors.New("invalid product ID")
	}
	return s.productRepo.RemoveProductFromArticle(ctx, articleId, productId)
}

// GetProductsByArticleId implements ProductService
func (s *productService) GetProductsByArticleId(ctx context.Context, articleId uuid.UUID) ([]model.Product, error) {
	if articleId == uuid.Nil {
		return nil, errors.New("invalid article ID")
	}
	return s.productRepo.GetProductsByArticleId(ctx, articleId)
}

// GetArticlesByProductId implements ProductService
func (s *productService) GetArticlesByProductId(ctx context.Context, productId uuid.UUID) ([]model.Article, error) {
	if productId == uuid.Nil {
		return nil, errors.New("invalid product ID")
	}
	return s.productRepo.GetArticlesByProductId(ctx, productId)
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
