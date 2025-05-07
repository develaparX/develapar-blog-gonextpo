package controller

import (
	"develapar-server/model"
	"develapar-server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LikeController struct {
	service service.LikeService
	rg      *gin.RouterGroup
}

func (l *LikeController) AddLikeHandler(ctx *gin.Context) {
	var payload model.Likes

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid payload: " + err.Error(),
		})
		return
	}

	data, err := l.service.CreateLike(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to add Like: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Add New Like",
		"data":    data,
	})

}

func (l *LikeController) GetLikeByArticleIdHandler(ctx *gin.Context) {
	articleID, err := strconv.Atoi(ctx.Param("article_id"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mendapatkan daftar bookmark
	bookmarks, err := l.service.FindLikeByArticleId(articleID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get Bookmark",
		"data":    bookmarks,
	})
}

func (l *LikeController) GetLikeByUserIdHandler(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mendapatkan daftar bookmark
	bookmarks, err := l.service.FindLikeByUserId(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get Bookmark",
		"data":    bookmarks,
	})
}

func (l *LikeController) DeleteLikeHandler(ctx *gin.Context) {
	var payload model.Likes

	err := l.service.DeleteLike(payload.User.Id, payload.Article.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Delete Like",
	})
}

func (l *LikeController) Route() {
	router := l.rg.Group("/likes")
	router.POST("/", l.AddLikeHandler)
	router.GET("/article/:article_id", l.GetLikeByArticleIdHandler)
	router.GET("/user/:user_id", l.GetLikeByUserIdHandler)
	router.DELETE("/del", l.DeleteLikeHandler)
}

func NewLikeController(lS service.LikeService, rg *gin.RouterGroup) *LikeController {
	return &LikeController{
		service: lS,
		rg:      rg,
	}
}
