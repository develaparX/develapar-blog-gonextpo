package dto

import "github.com/google/uuid"

type AssignTagsByNameDTO struct {
	ArticleID uuid.UUID `json:"article_id" binding:"required"`
	Tags      []string  `json:"tags" binding:"required"`
}
