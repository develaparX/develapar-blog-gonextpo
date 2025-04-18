package model

import "time"

type Bookmark struct {
	Id        int       `json:"id"`
	Article   Article   `json:"article"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}
