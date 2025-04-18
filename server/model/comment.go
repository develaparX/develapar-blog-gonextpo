package model

import "time"

type Comment struct {
	Id        int       `json:"id"`
	Article   Article   `json:"article"`
	User      User      `json:"user"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
