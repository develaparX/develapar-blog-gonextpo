package dto

import (
	"time"

	"develapar-server/model"

	"github.com/google/uuid"
)

type CommentResponse struct {
	Id        int             `json:"id"`
	Content   string          `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	User      UserResponse    `json:"user"`
	Article   ArticleResponse `json:"article"`
}

type UserResponse struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ArticleResponse struct {
	Id         uuid.UUID       `json:"id"`
	Title      string          `json:"title"`
	Slug       string          `json:"slug"`
	Content    string          `json:"content"`
	UserId     uuid.UUID       `json:"user_id"`
	User       *model.User     `json:"user,omitempty"`
	CategoryId uuid.UUID       `json:"category_id"`
	Category   *model.Category `json:"category,omitempty"`
	Views      int             `json:"views"`
	Status     string          `json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	Tags       []model.Tags    `json:"tags"`
}
