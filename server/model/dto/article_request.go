package dto

import "github.com/google/uuid"

type CreateArticleRequest struct {
	Title      string    `json:"title" binding:"required"`
	Content    string    `json:"content" binding:"required"`
	CategoryID uuid.UUID `json:"category_id" binding:"required"`
	Tags       []string  `json:"tags,omitempty"`
}

type UpdateArticleRequest struct {
	Title      *string    `json:"title"`
	Content    *string    `json:"content"`
	CategoryID *uuid.UUID `json:"category_id"`
	Tags       []string   `json:"tags,omitempty"`
}

type UpdateCategoryRequest struct {
	Name *string `json:"name"`
}
