package dto

import (
	"time"

	"github.com/google/uuid"
)

// Product Category Response
type ProductCategoryResponse struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Product Affiliate Link Response
type ProductAffiliateLinkResponse struct {
	Id           uuid.UUID `json:"id"`
	PlatformName string    `json:"platform_name"`
	Url          string    `json:"url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Product Response (includes affiliate links)
type ProductResponse struct {
	Id                uuid.UUID                      `json:"id"`
	ProductCategoryId *uuid.UUID                     `json:"product_category_id"`
	ProductCategory   *ProductCategoryResponse       `json:"product_category,omitempty"`
	Name              string                         `json:"name"`
	Description       *string                        `json:"description"`
	ImageUrl          *string                        `json:"image_url"`
	IsActive          bool                           `json:"is_active"`
	AffiliateLinks    []ProductAffiliateLinkResponse `json:"affiliate_links"`
	CreatedAt         time.Time                      `json:"created_at"`
	UpdatedAt         time.Time                      `json:"updated_at"`
}

// Product List Response (without affiliate links for performance)
type ProductListResponse struct {
	Id                uuid.UUID                `json:"id"`
	ProductCategoryId *uuid.UUID               `json:"product_category_id"`
	ProductCategory   *ProductCategoryResponse `json:"product_category,omitempty"`
	Name              string                   `json:"name"`
	Description       *string                  `json:"description"`
	ImageUrl          *string                  `json:"image_url"`
	IsActive          bool                     `json:"is_active"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
}