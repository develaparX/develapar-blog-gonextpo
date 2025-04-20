package controller

import (
	"develapar-server/model"
	"develapar-server/service"
	"net/http"

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

func (c *ArticleController) Route() {
	router := c.rg.Group("/article")
	router.GET("/", c.GetAllArticleHandler)
	router.POST("/", c.CreateArticleHandler)
}

func NewArticleController(aS service.ArticleService, rg *gin.RouterGroup) *ArticleController {
	return &ArticleController{
		service: aS,
		rg:      rg,
	}
}
