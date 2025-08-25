package dto

import "github.com/google/uuid"

// Product Category DTOs
type CreateProductCategoryRequest struct {
	Name        string  `json:"name" binding:"required"`
	Slug        string  `json:"slug,omitempty"`
	Description *string `json:"description"`
}

type UpdateProductCategoryRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
}

// Product DTOs
type CreateProductRequest struct {
	ProductCategoryId *uuid.UUID `json:"product_category_id" binding:"required"`
	Name              string     `json:"name" binding:"required"`
	Description       *string    `json:"description"`
	ImageUrl          *string    `json:"image_url"`
	IsActive          *bool      `json:"is_active"`
}

type UpdateProductRequest struct {
	ProductCategoryId *uuid.UUID `json:"product_category_id"`
	Name              *string    `json:"name"`
	Description       *string    `json:"description"`
	ImageUrl          *string    `json:"image_url"`
	IsActive          *bool      `json:"is_active"`
}

// Product Affiliate Link DTOs
type CreateProductAffiliateLinkRequest struct {
	PlatformName string `json:"platform_name" binding:"required"`
	Url          string `json:"url" binding:"required"`
}

type UpdateProductAffiliateLinkRequest struct {
	PlatformName *string `json:"platform_name"`
	Url          *string `json:"url"`
}