package controller

import (
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LikeController struct {
	service service.LikeService
	rg      *gin.RouterGroup
	md      middleware.AuthMiddleware
}

// @Summary Add a like to an article
// @Description Add a like to a specific article by the authenticated user
// @Tags Likes
// @Accept json
// @Produce json
// @Param payload body model.Likes true "Like creation details"
// @Success 200 {object} object{message=string,data=model.Likes} "Like successfully added"
// @Failure 400 {object} object{message=string} "Invalid payload"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Security BearerAuth
// @Router /likes [post]
func (l *LikeController) AddLikeHandler(ctx *gin.Context) {
	var payload model.Likes

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

// @Summary Get likes by article ID
// @Description Get a list of likes for a specific article ID
// @Tags Likes
// @Produce json
// @Param article_id path int true "ID of the article to retrieve likes for"
// @Success 200 {object} object{message=string,data=[]model.Likes} "List of likes for the article"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /likes/article/{article_id} [get]
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

// @Summary Get likes by user ID
// @Description Get a list of likes by a specific user ID
// @Tags Likes
// @Produce json
// @Param user_id path int true "ID of the user whose likes to retrieve"
// @Success 200 {object} object{message=string,data=[]model.Likes} "List of likes by the user"
// @Failure 400 {object} object{error=string} "Invalid user ID"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /likes/user/{user_id} [get]
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

// @Summary Remove a like from an article
// @Description Remove a like from a specific article by the authenticated user
// @Tags Likes
// @Accept json
// @Produce json
// @Param payload body object{article_id=int} true "Article ID to unlike"
// @Success 200 {object} object{message=string} "Like deleted successfully"
// @Failure 400 {object} object{error=string} "Invalid article ID"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /likes [delete]
func (l *LikeController) DeleteLikeHandler(ctx *gin.Context) {
	var payload model.Likes

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

	err := l.service.DeleteLike(payload.User.Id, payload.Article.Id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Delete Like",
	})
}

// @Summary Check if an article is liked by the current user
// @Description Check if a specific article is liked by the authenticated user
// @Tags Likes
// @Produce json
// @Param article_id query int true "ID of the article to check"
// @Success 200 {object} object{liked=bool} "Like status"
// @Failure 400 {object} object{error=string} "Invalid article ID"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /likes/check [get]
func (c *LikeController) CheckLikeHandler(ctx *gin.Context) {
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

	liked, err := c.service.IsLiked(userId, articleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"liked": liked})
}

func (l *LikeController) Route() {
	router := l.rg.Group("/likes")
	router.GET("/article/:article_id", l.GetLikeByArticleIdHandler)
	router.GET("/user/:user_id", l.GetLikeByUserIdHandler)

	routerAuth := router.Group("/", l.md.CheckToken())
	routerAuth.POST("/", l.AddLikeHandler)
	routerAuth.DELETE("/", l.DeleteLikeHandler)
	routerAuth.GET("/check", l.CheckLikeHandler)
}

func NewLikeController(lS service.LikeService, rg *gin.RouterGroup, md middleware.AuthMiddleware) *LikeController {
	return &LikeController{
		service: lS,
		rg:      rg,
		md:      md,
	}
}
