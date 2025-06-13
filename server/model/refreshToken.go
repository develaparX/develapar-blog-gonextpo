package model

import "time"

type RefreshToken struct {
	ID        uint `gorm:"primaryKey"`
	UserID    int
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}
