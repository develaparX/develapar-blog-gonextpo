package dto

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtTokenClaims struct {
	jwt.RegisteredClaims
	UserId uuid.UUID `json:"userId"`
	Role   string    `json:"role"`
}
