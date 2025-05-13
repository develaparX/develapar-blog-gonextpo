package controller

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/service"
	"net/http"
	"strconv"

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
		"message": "Success Get All Category",
		"data":    data,
	})
}

func (c *CategoryController) UpdateCategoryHandler(ctx *gin.Context) {
	idCat := ctx.Param("cat_id")
	id, err := strconv.Atoi(idCat)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Category ID"})
		return
	}

	var req dto.UpdateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := c.service.UpdateCategory(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Category updated successfully",
		"data":    cat,
	})
}

func (c *CategoryController) DeleteCategoryHandler(ctx *gin.Context) {
	Id := ctx.Param("cat_id")

	catId, err := strconv.Atoi(Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Category ID"})
		return
	}

	err2 := c.service.DeleteCategory(catId)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Delete Category"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Delete Category",
	})
}

func (c *CategoryController) Route() {
	router := c.rg.Group("/category")
	router.GET("/", c.GetAllCategoryHandler)
	router.POST("/", c.CreateCategoryHandler)
	router.PUT("/:cat_id", c.UpdateCategoryHandler)
	router.DELETE("/:cat_id", c.DeleteCategoryHandler)
}

func NewCategoryController(cS service.CategoryService, rg *gin.RouterGroup) *CategoryController {
	return &CategoryController{
		service: cS,
		rg:      rg,
	}
}
