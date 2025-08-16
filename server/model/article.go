package model

import (
	"time"

	"github.com/google/uuid"
)

type Article struct {
	Id         uuid.UUID  `json:"id"`
	Title      string     `json:"title"`
	Slug       string     `json:"slug"`
	Content    string     `json:"content"`
	UserId     uuid.UUID  `json:"user_id"`
	User       *User      `json:"user,omitempty"`
	CategoryId *uuid.UUID `json:"category_id"`
	Category   *Category  `json:"category,omitempty"`
	Views      int        `json:"views"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
