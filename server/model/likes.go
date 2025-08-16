package model

import (
	"time"

	"github.com/google/uuid"
)

type Likes struct {
	Id        uuid.UUID `json:"id"`
	ArticleId uuid.UUID `json:"article_id"`
	UserId    uuid.UUID `json:"user_id"`
	Article   *Article  `json:"article,omitempty"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
