package controller

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
	rg      *gin.RouterGroup
}

func (u *UserController) loginHandler(ctx *gin.Context) {
	var payload dto.LoginDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := u.service.Login(payload)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Set HttpOnly cookie untuk refresh token
	ctx.SetCookie(
		"refreshToken",
		response.RefreshToken,
		60*60*24*7, // 7 hari
		"/",
		"localhost", // ganti ke domain di production
		false,       // secure: true jika HTTPS
		true,        // httpOnly
	)

	// Jangan kirim refresh token ke frontend
	ctx.JSON(http.StatusOK, gin.H{
		"message":     "Success Login",
		"accessToken": response.AccessToken,
	})
}

func (u *UserController) registerUser(ctx *gin.Context) {
	var payload model.User

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err},
		)
	}

	data, err := u.service.CreateNewUser(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err},
		)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Create New User",
		"data":    data,
	})
}

func (u *UserController) findUserByIdHandler(ctx *gin.Context) {
	userId := ctx.Param("user_id")

	user, err := u.service.FindUserById(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get User",
		"data":    user,
	})
}

func (u *UserController) findAllUserHandler(ctx *gin.Context) {
	user, err := u.service.FindAllUser()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get All User",
		"data":    user,
	})
}

func RefreshTokenHandler(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil dari cookie, bukan dari JSON body
		cookie, err := c.Cookie("refreshToken")
		if err != nil || cookie == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found in cookies"})
			return
		}

		// Proses refresh
		tokenResp, err := userService.RefreshToken(cookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		// Set ulang refreshToken ke cookie (opsional: rotate token)
		refreshExpiry := time.Now().Add(7 * 24 * time.Hour)
		c.SetCookie("refreshToken", tokenResp.RefreshToken, int(refreshExpiry.Sub(time.Now()).Seconds()), "/", "", true, true)

		// Return hanya access token ke frontend (refresh tetap di cookie)
		c.JSON(http.StatusOK, gin.H{
			"access_token": tokenResp.AccessToken,
		})
	}
}



func (u *UserController) Route() {
	router := u.rg.Group("/users")
	{
		router.GET("/", u.findAllUserHandler)
		router.GET("/:user_id", u.findUserByIdHandler)
	}

	r := u.rg.Group("/auth")
	{
		r.POST("/login", u.loginHandler)
		r.POST("/register", u.registerUser)
		r.POST("/refresh", RefreshTokenHandler(u.service))
	}

}

func NewUserController(uS service.UserService, rg *gin.RouterGroup) *UserController {
	return &UserController{

		service: uS,
		rg:      rg,
	}
}
