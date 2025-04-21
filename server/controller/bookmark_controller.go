package controller

import (
	"develapar-server/model"
	"develapar-server/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookmarkController struct {
	service service.BookmarkService
	rg      *gin.RouterGroup
}

func (b *BookmarkController) CreateBookmarkHandler(ctx *gin.Context) {
	var payload model.Bookmark
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid payload: " + err.Error(),
		})
		return
	}

	data, err := b.service.CreateBookmark(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create bookmark: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Create New Category",
		"data":    data,
	})
}

func (b *BookmarkController) GetBookmarkByUserId(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	// Mendapatkan daftar bookmark
	bookmarks, err := b.service.FindByUserId(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get Bookmark",
		"data":    bookmarks,
	})
}

func (c *BookmarkController) Route() {
	router := c.rg.Group("/bookmark")
	router.GET("/:user_id", c.GetBookmarkByUserId)
	router.POST("/", c.CreateBookmarkHandler)
}

func NewBookmarkController(bS service.BookmarkService, rg *gin.RouterGroup) *BookmarkController {
	return &BookmarkController{
		service: bS,
		rg:      rg,
	}
}
