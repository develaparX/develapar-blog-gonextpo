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

// @Summary Create a new comment
// @Description Create a new comment on an article
// @Tags Comments
// @Accept json
// @Produce json
// @Param payload body model.Comment true "Comment creation details"
// @Success 200 {object} object{message=string,data=model.Comment} "Comment successfully created"
// @Failure 400 {object} object{message=string} "Invalid payload"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Security BearerAuth
// @Router /comments [post]
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

// @Summary Get comments by article ID
// @Description Get a list of comments for a specific article ID
// @Tags Comments
// @Produce json
// @Param article_id path int true "ID of the article to retrieve comments for"
// @Success 200 {object} object{message=string,data=[]model.Comment} "List of comments for the article"
// @Failure 400 {object} object{error=string} "Invalid article ID"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /comment/article/c{article_id} [get]
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

// @Summary Get comments by user ID
// @Description Get a list of comments by a specific user ID
// @Tags Comments
// @Produce json
// @Param user_id path int true "ID of the user whose comments to retrieve"
// @Success 200 {object} object{message=string,data=[]dto.CommentResponse} "List of comments by the user"
// @Failure 400 {object} object{error=string} "Invalid user ID"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /comment/user/{user_id} [get]
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

// @Summary Update a comment
// @Description Update an existing comment by ID
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path int true "ID of the comment to update"
// @Param payload body object{content=string} true "Comment update details"
// @Success 200 {object} object{message=string} "Comment updated successfully"
// @Failure 400 {object} object{error=string} "Invalid payload"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 403 {object} object{error=string} "Forbidden (user does not own the comment)"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /comment/{id} [put]
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

// @Summary Delete a comment
// @Description Delete a comment by ID
// @Tags Comments
// @Produce json
// @Param id path int true "ID of the comment to delete"
// @Success 200 {object} object{message=string} "Comment deleted successfully"
// @Failure 400 {object} object{error=string} "Invalid comment ID"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Security BearerAuth
// @Router /comment/{id} [delete]
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
