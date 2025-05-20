package controller

import (
	"develapar-server/model"
	"develapar-server/model/dto"
	"develapar-server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleController struct {
	service service.ArticleService
	rg      *gin.RouterGroup
}

func (c *ArticleController) CreateArticleHandler(ctx *gin.Context) {
	var payload model.Article
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid payload: " + err.Error(),
		})
		return
	}

	data, err := c.service.CreateArticle(payload)
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

func (c *ArticleController) GetAllArticleHandler(ctx *gin.Context) {
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

func (c *ArticleController) UpdateArticleHandler(ctx *gin.Context) {
	idStr := ctx.Param("article_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	var req dto.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := c.service.UpdateArticle(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Article updated successfully",
		"data":    article,
	})
}

func (c *ArticleController) GetBySlugHandler(ctx *gin.Context) {
	slug := ctx.Param("slug")

	article, err := c.service.FindBySlug(slug)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get article by slug",
		"data":    article,
	})
}

func (ac *ArticleController) GetByUserIdHandler(ctx *gin.Context) {
	userIdParam := ctx.Param("user_id")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	articles, err := ac.service.FindByUserId(userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get articles by user ID",
		"data":    articles,
	})
}

func (ac *ArticleController) GetByCategory(ctx *gin.Context) {
	categoryName := ctx.Param("cat_name")
	articles, err := ac.service.FindByCategory(categoryName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get articles by Category",
		"data":    articles,
	})
}

func (ac *ArticleController) DeleteArticleHandler(ctx *gin.Context) {
	Id := ctx.Param("article_id")
	arcId, err := strconv.Atoi(Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	erro := ac.service.DeleteArticle(arcId)
	if erro != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Delete Article",
	})

}

func (c *ArticleController) Route() {

	router := c.rg.Group("/article")

	router.GET("/", c.GetAllArticleHandler)
	router.POST("/", c.CreateArticleHandler)
	router.PUT("/:article_id", c.UpdateArticleHandler)
	router.GET("/:slug", c.GetBySlugHandler)
	router.GET("/u/:user_id", c.GetByUserIdHandler)
	router.GET("/c/:cat_name", c.GetByCategory)
	router.DELETE("/:article_id", c.DeleteArticleHandler)
}

func NewArticleController(aS service.ArticleService, rg *gin.RouterGroup) *ArticleController {
	return &ArticleController{
		service: aS,
		rg:      rg,
	}
}
