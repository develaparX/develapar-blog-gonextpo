package controller

import (
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookmarkController struct {
	service service.BookmarkService
	rg      *gin.RouterGroup
	md      middleware.AuthMiddleware
}

func (b *BookmarkController) CreateBookmarkHandler(ctx *gin.Context) {
	var payload model.Bookmark
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

	payload.User.Id = userId

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

func (b *BookmarkController) DeleteBookmarkHandler(ctx *gin.Context) {
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

	articleIdParam := ctx.Param("article_id")
	articleId, err := strconv.Atoi(articleIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	err = b.service.DeleteBookmark(userId, articleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Bookmark deleted successfully",
	})
}

func (c *BookmarkController) CheckBookmarkHandler(ctx *gin.Context) {
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

	articleId, err := strconv.Atoi(ctx.Query("article_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article_id"})
		return
	}

	bookmarked, err := c.service.IsBookmarked(userId, articleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"bookmarked": bookmarked})
}

func (c *BookmarkController) Route() {
	router := c.rg.Group("/bookmark")
	router.GET("/:user_id", c.GetBookmarkByUserId)

	routerAuth := router.Group("/")
	routerAuth.Use(c.md.CheckToken())
	routerAuth.POST("/", c.CreateBookmarkHandler)
	routerAuth.DELETE("/", c.DeleteBookmarkHandler)
	routerAuth.GET("/check", c.CheckBookmarkHandler)
}

func NewBookmarkController(bS service.BookmarkService, rg *gin.RouterGroup, md middleware.AuthMiddleware) *BookmarkController {
	return &BookmarkController{
		service: bS,
		rg:      rg,
		md:      md,
	}
}
