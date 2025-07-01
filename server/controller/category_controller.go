package controller

import (
	"develapar-server/middleware"
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
	md      middleware.AuthMiddleware
}

// @Summary Create a new category
// @Description Create a new category with a given name
// @Tags Categories
// @Accept json
// @Produce json
// @Param payload body model.Category true "Category creation details"
// @Success 200 {object} object{message=string,data=model.Category} "Category successfully created"
// @Failure 400 {object} object{message=string} "Invalid payload"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Security BearerAuth
// @Router /category [post]
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

// @Summary Get all categories
// @Description Get a list of all categories
// @Tags Categories
// @Produce json
// @Success 200 {object} object{message=string,data=[]model.Category} "List of categories"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Router /category [get]
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

// @Summary Update a category
// @Description Update an existing category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param cat_id path int true "ID of the category to update"
// @Param payload body dto.UpdateCategoryRequest true "Category update details"
// @Success 200 {object} object{message=string,data=model.Category} "Category updated successfully"
// @Failure 400 {object} object{error=string} "Invalid category ID or payload"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /category/{cat_id} [put]
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

// @Summary Delete a category
// @Description Delete a category by ID
// @Tags Categories
// @Produce json
// @Param cat_id path int true "ID of the category to delete"
// @Success 200 {object} object{message=string} "Category deleted successfully"
// @Failure 400 {object} object{error=string} "Invalid category ID"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /category/{cat_id} [delete]
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

	routerAuth := router.Group("/", c.md.CheckToken("admin"))
	routerAuth.POST("/", c.CreateCategoryHandler)
	routerAuth.PUT("/:cat_id", c.UpdateCategoryHandler)
	routerAuth.DELETE("/:cat_id", c.DeleteCategoryHandler)
}

func NewCategoryController(cS service.CategoryService, rg *gin.RouterGroup, md middleware.AuthMiddleware) *CategoryController {
	return &CategoryController{
		service: cS,
		rg:      rg,
		md:      md,
	}
}
