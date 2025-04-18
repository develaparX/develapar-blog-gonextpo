package controller

import (
	"develapar-server/model"
	"develapar-server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.UserService
	rg      *gin.RouterGroup
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

func (u *UserController) Route() {
	router := u.rg.Group("/users")

	router.POST("/register", u.registerUser)
}

func NewUserController(uS service.UserService, rg *gin.RouterGroup) *UserController {
	return &UserController{

		service: uS,
		rg:      rg,
	}
}
