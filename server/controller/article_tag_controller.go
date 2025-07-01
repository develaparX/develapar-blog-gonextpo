package controller

import (
	"develapar-server/middleware"
	"develapar-server/model/dto"
	"develapar-server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleTagController struct {
	service service.ArticleTagService
	rg      *gin.RouterGroup
	md      middleware.AuthMiddleware
}

type AssignTagRequest struct {
	ArticleID int   `json:"article_id"`
	TagIDs    []int `json:"tag_ids"`
}

// @Summary Assign tags to an article by tag names
// @Description Assigns a list of tags (by name) to a specific article
// @Tags Tags
// @Accept json
// @Produce json
// @Param payload body dto.AssignTagsByNameDTO true "Article ID and list of tag names"
// @Success 200 {object} object{message=string} "Tags assigned successfully"
// @Failure 400 {object} object{error=string} "Invalid payload"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /article-to-tag [post]
func (c *ArticleTagController) AssignTagToArticleByNameHandler(ctx *gin.Context) {
	var req dto.AssignTagsByNameDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.AsignTagsByName(req.ArticleID, req.Tags)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tags assigned successfully"})
}

func (c *ArticleTagController) AssignTagToArticleByIdHandler(ctx *gin.Context) {
	var req AssignTagRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.AssignTags(req.ArticleID, req.TagIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Tags assigned successfully"})
}

// @Summary Get tags by article ID
// @Description Get a list of tags associated with a specific article ID
// @Tags Tags
// @Produce json
// @Param article_id path int true "ID of the article to retrieve tags for"
// @Success 200 {object} object{data=[]model.Tags} "List of tags for the article"
// @Failure 400 {object} object{error=string} "Invalid article ID"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /article-to-tag/tags/{article_id} [get]
func (c *ArticleTagController) GetTagsByArticleIDHandler(ctx *gin.Context) {
	articleID, err := strconv.Atoi(ctx.Param("article_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	tags, err := c.service.FindTagByArticleId(articleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": tags})
}

// @Summary Get articles by tag ID
// @Description Get a list of articles associated with a specific tag ID
// @Tags Articles
// @Produce json
// @Param tag_id path int true "ID of the tag to retrieve articles for"
// @Success 200 {object} object{data=[]model.Article} "List of articles with the specified tag"
// @Failure 400 {object} object{error=string} "Invalid tag ID"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /article-to-tag/article/{tag_id} [get]
func (c *ArticleTagController) GetArticlesByTagIDHandler(ctx *gin.Context) {
	tagID, err := strconv.Atoi(ctx.Param("tag_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	articles, err := c.service.FindArticleByTagId(tagID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": articles})
}

// @Summary Remove a tag from an article
// @Description Remove a specific tag from an article by their IDs
// @Tags Tags
// @Produce json
// @Param article_id path int true "ID of the article"
// @Param tag_id path int true "ID of the tag to remove"
// @Success 200 {object} object{message=string} "Tag removed from article successfully"
// @Failure 400 {object} object{error=string} "Invalid article ID or tag ID"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /article-to-tag/articles/{article_id}/tags/{tag_id} [delete]
func (c *ArticleTagController) RemoveTagFromArticleHandler(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("article_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Article ID"})
		return
	}

	tagId, err := strconv.Atoi(ctx.Param("tag_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag Id"})
		return
	}

	err = c.service.RemoveTagFromArticle(articleId, tagId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove tag from article"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Tag removed from article successfully",
	})
}

func (at *ArticleTagController) Route() {
	router := at.rg.Group("/article-to-tag")
	router.GET("/tags/:article_id", at.GetTagsByArticleIDHandler)
	router.GET("/article/:tag_id", at.GetArticlesByTagIDHandler)

	routerAuth := router.Group("/")
	routerAuth.Use(at.md.CheckToken())
	routerAuth.POST("/", at.AssignTagToArticleByNameHandler)
	routerAuth.DELETE("/articles/:article_id/tags/:tag_id", at.RemoveTagFromArticleHandler)

}

func NewArticleTagController(s service.ArticleTagService, rg *gin.RouterGroup, md middleware.AuthMiddleware) *ArticleTagController {
	return &ArticleTagController{service: s, rg: rg, md: md}
}
