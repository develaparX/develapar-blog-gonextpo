package model

import (
	"time"

	"github.com/google/uuid"
)

type ProductCategory struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
	Id                uuid.UUID        `json:"id"`
	ProductCategoryId *uuid.UUID       `json:"product_category_id" binding:"required"`
	ProductCategory   *ProductCategory `json:"product_category,omitempty" binding:required`
	Name              string           `json:"name"`
	Description       *string          `json:"description"`
	ImageUrl          *string          `json:"image_url"`
	IsActive          bool             `json:"is_active"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

type ProductAffiliateLink struct {
	Id           uuid.UUID `json:"id"`
	ProductId    uuid.UUID `json:"product_id"`
	Product      *Product  `json:"product,omitempty"`
	PlatformName string    `json:"platform_name"`
	Url          string    `json:"url"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ArticleProduct struct {
	ArticleId uuid.UUID `json:"article_id"`
	ProductId uuid.UUID `json:"product_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
