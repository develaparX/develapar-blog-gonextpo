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

// @Summary User login
// @Description Authenticate user and return access token
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body dto.LoginDto true "Login credentials"
// @Success 200 {object} object{message=string,accessToken=string} "Success Login"
// @Failure 401 {object} object{error=string} "Invalid credentials"
// @Router /auth/login [post]
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

// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags Users
// @Accept json
// @Produce json
// @Param payload body model.User true "User registration details"
// @Success 200 {object} object{message=string,data=model.User} "User successfully registered"
// @Failure 400 {object} object{message=string} "Invalid payload"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Router /auth/register [post]
func (u *UserController) registerUser(ctx *gin.Context) {
	var payload model.User

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error()},
		)
		return
	}

	data, err := u.service.CreateNewUser(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error()},
		)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Create New User",
		"data":    data,
	})
}

// @Summary Get user by ID
// @Description Get user details by their ID
// @Tags Users
// @Produce json
// @Param user_id path string true "ID of the user to retrieve"
// @Success 200 {object} object{message=string,data=model.User} "User details"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /users/{user_id} [get]
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

// @Summary Get all users
// @Description Get a list of all registered users
// @Tags Users
// @Produce json
// @Success 200 {object} object{message=string,data=[]model.User} "List of users"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /users [get]
func (u *UserController) findAllUserHandler(ctx *gin.Context) {
	user, err := u.service.FindAllUser()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get All User",
		"data":    user,
	})
}

// @Summary Refresh access token
// @Description Refresh access token using refresh token from cookie
// @Tags Users
// @Produce json
// @Success 200 {object} object{access_token=string} "Access token refreshed successfully"
// @Failure 401 {object} object{error=string} "Refresh token not found or invalid"
// @Router /auth/refresh [post]
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
