package controller

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
	rg      *gin.RouterGroup
}

func (u *UserController) loginHandler(ctx *gin.Context) {
	var payload dto.LoginDto
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := u.service.Login(payload)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusOK, gin.H{
			"message": "Success Login",
			"data":    response,
		},
	)
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
	}

}

func NewUserController(uS service.UserService, rg *gin.RouterGroup) *UserController {
	return &UserController{

		service: uS,
		rg:      rg,
	}
}
