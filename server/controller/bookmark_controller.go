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

// @Summary Create a new bookmark
// @Description Create a new bookmark for an article
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param payload body model.Bookmark true "Bookmark creation details"
// @Success 200 {object} object{message=string,data=model.Bookmark} "Bookmark successfully created"
// @Failure 400 {object} object{message=string} "Invalid payload"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Security BearerAuth
// @Router /bookmark [post]
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

// @Summary Get bookmarks by user ID
// @Description Get a list of bookmarks for a specific user ID
// @Tags Bookmarks
// @Produce json
// @Param user_id path int true "ID of the user whose bookmarks to retrieve"
// @Success 200 {object} object{message=string,data=[]model.Bookmark} "List of bookmarks for the user"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /bookmark/{user_id} [get]
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

// @Summary Delete a bookmark
// @Description Delete a bookmark for an article by article ID
// @Tags Bookmarks
// @Accept json
// @Produce json
// @Param article_id body object{article_id=int} true "Article ID to unbookmark"
// @Success 200 {object} object{message=string} "Bookmark deleted successfully"
// @Failure 400 {object} object{error=string} "Invalid article ID"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /bookmark [delete]
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

// @Summary Check if an article is bookmarked by the current user
// @Description Check if a specific article is bookmarked by the authenticated user
// @Tags Bookmarks
// @Produce json
// @Param article_id query int true "ID of the article to check"
// @Success 200 {object} object{bookmarked=bool} "Bookmark status"
// @Failure 400 {object} object{error=string} "Invalid article ID"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /bookmark/check [get]
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
