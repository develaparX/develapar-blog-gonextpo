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

type ArticleController struct {
	service service.ArticleService
	md middleware.AuthMiddleware
	rg      *gin.RouterGroup
}

func (c *ArticleController) CreateArticleHandler(ctx *gin.Context) {
	userIdRaw, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID type"})
		return
	}
	userId := int(userIdFloat)

	var payload model.Article
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid payload: " + err.Error(),
		})
		return
	}

	payload.User.Id = userId // assign author ID from token

	data, err := c.service.CreateArticle(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create article: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success create article",
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
	// Ambil user ID dari JWT context
	userIdRaw, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID type"})
		return
	}
	userId := int(userIdFloat)

	// Ambil ID artikel dari URL param
	idStr := ctx.Param("article_id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	// Cek apakah user adalah pemilik artikel
	article, err := c.service.FindById(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Article not found"})
		return
	}
	if article.User.Id != userId {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "You do not own this article"})
		return
	}

	// Bind data dari payload
	var req dto.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update artikel
	updatedArticle, err := c.service.UpdateArticle(id, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Article updated successfully",
		"data":    updatedArticle,
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
	// Ambil user ID dari JWT context
	userIdRaw, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}
	userIdFloat, ok := userIdRaw.(float64)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user ID type"})
		return
	}
	userId := int(userIdFloat)

	// Ambil ID artikel dari param
	idStr := ctx.Param("article_id")
	articleId, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	// Cek ownership
	article, err := ac.service.FindById(articleId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Article not found"})
		return
	}
	if article.User.Id != userId {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "You do not own this article"})
		return
	}

	// Delete article
	err = ac.service.DeleteArticle(articleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success delete article",
	})
}


func (c *ArticleController) Route() {
	router := c.rg.Group("/article")

	// Public routes
	router.GET("/", c.GetAllArticleHandler)
	router.GET("/:slug", c.GetBySlugHandler)
	router.GET("/u/:user_id", c.GetByUserIdHandler)
	router.GET("/c/:cat_name", c.GetByCategory)

	// Protected routes
	routerAuth := router.Group("/")

	routerAuth.Use(c.md.CheckToken()) // hanya butuh login
	routerAuth.POST("/", c.CreateArticleHandler)
	routerAuth.PUT("/:article_id", c.UpdateArticleHandler)
	routerAuth.DELETE("/:article_id", c.DeleteArticleHandler)
}


func NewArticleController(aS service.ArticleService, md middleware.AuthMiddleware, rg *gin.RouterGroup) *ArticleController {
	return &ArticleController{
		service: aS,
		md:      md,
		rg:      rg,
	}
}
