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

// @Summary Create a new article
// @Description Create a new blog article
// @Tags Articles
// @Accept json
// @Produce json
// @Param payload body model.Article true "Article creation details"
// @Success 200 {object} object{message=string,data=model.Article} "Article successfully created"
// @Failure 400 {object} object{message=string} "Invalid payload"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Security BearerAuth
// @Router /article [post]
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

// @Summary Get all articles
// @Description Get a list of all blog articles
// @Tags Articles
// @Produce json
// @Success 200 {object} object{message=string,data=[]model.Article} "List of articles"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Router /article [get]
func (c *ArticleController) GetAllArticleHandler(ctx *gin.Context) {
	data, err := c.service.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get articles: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get All Articles",
		"data":    data,
	})
}

// @Summary Update an article
// @Description Update an existing article by ID
// @Tags Articles
// @Accept json
// @Produce json
// @Param article_id path int true "ID of the article to update"
// @Param payload body dto.UpdateArticleRequest true "Article update details"
// @Success 200 {object} object{message=string,data=model.Article} "Article updated successfully"
// @Failure 400 {object} object{error=string} "Invalid article ID or payload"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 403 {object} object{message=string} "Forbidden (user does not own the article)"
// @Failure 404 {object} object{message=string} "Article not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /article/{article_id} [put]
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


// @Summary Get article by slug
// @Description Get article details by its slug
// @Tags Articles
// @Produce json
// @Param slug path string true "Slug of the article to retrieve"
// @Success 200 {object} object{message=string,data=model.Article} "Article details"
// @Failure 404 {object} object{error=string} "Article not found"
// @Router /article/{slug} [get]
func (c *ArticleController) GetBySlugHandler(ctx *gin.Context) {
	slug := ctx.Param("slug")

	article, err := c.service.FindBySlug(slug)
	if err != nil {
		if err.Error() == "not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get article by slug: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success get article by slug",
		"data":    article,
	})
}

// @Summary Get articles by user ID
// @Description Get a list of articles by a specific user ID
// @Tags Articles
// @Produce json
// @Param user_id path int true "ID of the user whose articles to retrieve"
// @Success 200 {object} object{message=string,data=[]model.Article} "List of articles by user"
// @Failure 400 {object} object{error=string} "Invalid user ID"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /article/u/{user_id} [get]
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

// @Summary Get articles by category name
// @Description Get a list of articles by category name
// @Tags Articles
// @Produce json
// @Param cat_name path string true "Name of the category to retrieve articles from"
// @Success 200 {object} object{message=string,data=[]model.Article} "List of articles by category"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /article/c/{cat_name} [get]
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

// @Summary Delete an article
// @Description Delete an article by ID
// @Tags Articles
// @Produce json
// @Param article_id path int true "ID of the article to delete"
// @Success 200 {object} object{message=string} "Article deleted successfully"
// @Failure 400 {object} object{error=string} "Invalid article ID"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 403 {object} object{message=string} "Forbidden (user does not own the article)"
// @Failure 404 {object} object{message=string} "Article not found"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /article/{article_id} [delete]
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
	// Public routes
	publicRoutes := c.rg.Group("/article")
	publicRoutes.GET("/", c.GetAllArticleHandler)
	publicRoutes.GET("/:slug", c.GetBySlugHandler)
	publicRoutes.GET("/u/:user_id", c.GetByUserIdHandler)
	publicRoutes.GET("/c/:cat_name", c.GetByCategory)

	// Protected routes
	protectedRoutes := c.rg.Group("/article")
	protectedRoutes.Use(c.md.CheckToken()) // hanya butuh login
	protectedRoutes.POST("/", c.CreateArticleHandler)
	protectedRoutes.PUT("/:article_id", c.UpdateArticleHandler)
	protectedRoutes.DELETE("/:article_id", c.DeleteArticleHandler)
}


func NewArticleController(aS service.ArticleService, md middleware.AuthMiddleware, rg *gin.RouterGroup) *ArticleController {
	return &ArticleController{
		service: aS,
		md:      md,
		rg:      rg,
	}
}
