package controller

import (
	"develapar-server/model"
	"develapar-server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	service service.CategoryService
	rg      *gin.RouterGroup
}

func (c *CategoryController) CreateCategoryHandler(ctx *gin.Context) {
	var payload model.Category
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid payload: " + err.Error(),
		})
		return
	}

	data, err := c.service.CreateCategory(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create category: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Create New Category",
		"data":    data,
	})
}

func (c *CategoryController) GetAllCategoryHandler(ctx *gin.Context) {
	data, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err},
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Create New User",
		"data":    data,
	})
}

func (c *CategoryController) Route() {
	router := c.rg.Group("/category")
	router.GET("/", c.GetAllCategoryHandler)
	router.POST("/a", c.CreateCategoryHandler)
}

func NewCategoryController(cS service.CategoryService, rg *gin.RouterGroup) *CategoryController {
	return &CategoryController{
		service: cS,
		rg:      rg,
	}
}
