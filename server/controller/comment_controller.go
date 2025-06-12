package controller

import (
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	service service.CommentService
	rg      *gin.RouterGroup
	md      middleware.AuthMiddleware
}

func (c *CommentController) CreateCommentHandler(ctx *gin.Context) {
	var payload model.Comment

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

	data, err := c.service.CreateComment(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create comment: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Create New Comment",
		"data":    data,
	})
}

func (c *CommentController) FindCommentByArticleIdHandler(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("article_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comments, err := c.service.FindCommentByArticleId(articleId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get Comments",
		"data":    comments,
	})
}

func (c *CommentController) FindCommentByUserIdHandler(ctx *gin.Context) {
	user_id, err := strconv.Atoi(ctx.Param("user_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comments, err := c.service.FindCommentByUserId(user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get Comments",
		"data":    comments,
	})
}

func (c *CommentController) UpdateCommentHandler(ctx *gin.Context) {
	userId := ctx.GetInt("userId")
	commentId, _ := strconv.Atoi(ctx.Param("id"))

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.EditComment(commentId, req.Content, userId)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update comment"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (c *CommentController) DeleteCommentHandler(ctx *gin.Context) {
	userId := ctx.GetInt("userId")

	commentId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if err := c.service.DeleteComment(commentId, userId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func (c *CommentController) Route() {
	router := c.rg.Group("/comment")
	router.GET("/article/c:article_id", c.FindCommentByArticleIdHandler)
	router.GET("/user/:user_id", c.FindCommentByUserIdHandler)

	routerAuth := router.Group("/", c.md.CheckToken())

	routerAuth.POST("/", c.CreateCommentHandler)
	routerAuth.PUT("/:id", c.UpdateCommentHandler)
	routerAuth.DELETE("/:id", c.DeleteCommentHandler)
}

func NewCommentController(cS service.CommentService, rg *gin.RouterGroup, md middleware.AuthMiddleware) *CommentController {
	return &CommentController{
		service: cS,
		rg:      rg,
		md:      md,
	}
}
