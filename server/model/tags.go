package model

import (
	"time"

	"github.com/google/uuid"
)

type ArticleTag struct {
	ArticleId uuid.UUID `json:"article_id"`
	TagId     uuid.UUID `json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Tags struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
