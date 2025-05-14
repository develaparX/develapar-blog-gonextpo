package controller

import (
	"develapar-server/model"
	"develapar-server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	service service.CommentService
	rg      *gin.RouterGroup
}

func (c *CommentController) CreateCommentHandler(ctx *gin.Context) {
	var payload model.Comment
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
	commentId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var payload struct {
		Content string `json:"content"`
	}

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.EditComment(commentId, payload.Content); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Comment updated successfully"})
}

func (c *CommentController) DeleteCommentHandler(ctx *gin.Context) {
	commentId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	if err := c.service.DeleteComment(commentId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

func (c *CommentController) Route() {
	router := c.rg.Group("/comment")
	router.POST("/", c.CreateCommentHandler)
	router.GET("/article/c:article_id", c.FindCommentByArticleIdHandler)
	router.GET("/user/:user_id", c.FindCommentByUserIdHandler)
	router.PUT("/:id", c.UpdateCommentHandler)
	router.DELETE("/:id", c.DeleteCommentHandler)
}

func NewCommentController(cS service.CommentService, rg *gin.RouterGroup) *CommentController {
	return &CommentController{
		service: cS,
		rg:      rg,
	}
}
