package controller

import (
	"develapar-server/model/dto"
	"develapar-server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleTagController struct {
	service service.ArticleTagService
	rg      *gin.RouterGroup
}

type AssignTagRequest struct {
	ArticleID int   `json:"article_id"`
	TagIDs    []int `json:"tag_ids"`
}

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
	router.POST("/", at.AssignTagToArticleByNameHandler)
	router.GET("/tags/:article_id", at.GetTagsByArticleIDHandler)
	router.GET("/article/:tag_id", at.GetArticlesByTagIDHandler)
	router.DELETE("/articles/:article_id/tags/:tag_id", at.RemoveTagFromArticleHandler)

}

func NewArticleTagController(s service.ArticleTagService, rg *gin.RouterGroup) *ArticleTagController {
	return &ArticleTagController{service: s, rg: rg}
}
