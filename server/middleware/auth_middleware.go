package middleware

import (
	"develapar-server/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	CheckToken(roles ...string) gin.HandlerFunc
}

type authMiddleware struct {
	jwtService service.JwtService
}

func (a *authMiddleware) CheckToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		token := strings.Replace(header, "Bearer ", "", -1)

		claims, err := a.jwtService.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
			return
		}

		ctx.Set("userId", claims["userId"])
		ctx.Set("role", claims["role"]) // taruh juga kalau perlu di handler

		if len(roles) > 0 {
			var validRole bool
			for _, r := range roles {
				if r == claims["role"] {
					validRole = true
					break
				}
			}
			if !validRole {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Forbidden Access"})
				return
			}
		}

		ctx.Next()
	}
}


func NewAuthMiddleware(jwtService service.JwtService) AuthMiddleware {
	return &authMiddleware{jwtService: jwtService}
}
