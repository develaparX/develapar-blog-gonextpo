package dto

import "github.com/golang-jwt/jwt/v5"

type JwtTokenClaims struct {
	jwt.RegisteredClaims
	UserId int    `json:"userId"`
	Role   string `json:"role"`
}
