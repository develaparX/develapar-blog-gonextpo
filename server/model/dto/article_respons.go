package dto

import "time"

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
	Id    int    `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}
