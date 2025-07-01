package controller

import (
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	service service.TagService
	rg      *gin.RouterGroup
	md      middleware.AuthMiddleware
}

// @Summary Create a new tag
// @Description Create a new tag with a given name
// @Tags Tags
// @Accept json
// @Produce json
// @Param payload body model.Tags true "Tag creation details"
// @Success 200 {object} object{message=string,data=model.Tags} "Tag successfully created"
// @Failure 400 {object} object{message=string} "Invalid payload"
// @Failure 401 {object} object{message=string} "Unauthorized"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Security BearerAuth
// @Router /tags [post]
func (t *TagController) CreateTagHandler(ctx *gin.Context) {
	var payload model.Tags
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid payload: " + err.Error(),
		})
		return
	}

	data, err := t.service.CreateTag(payload)
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

// @Summary Get all tags
// @Description Get a list of all tags
// @Tags Tags
// @Produce json
// @Success 200 {object} object{message=string,data=[]model.Tags} "List of tags"
// @Failure 500 {object} object{message=string} "Internal server error"
// @Router /tags [get]
func (t *TagController) GetAllTagHandler(ctx *gin.Context) {
	data, err := t.service.FindAll()
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

// @Summary Get tag by ID
// @Description Get tag details by its ID
// @Tags Tags
// @Produce json
// @Param tags_id path int true "ID of the tag to retrieve"
// @Success 200 {object} object{message=string,data=model.Tags} "Tag details"
// @Failure 400 {object} object{error=string} "Invalid tag ID"
// @Failure 500 {object} object{error=string} "Internal server error"
// @Router /tags/{tags_id} [get]
func (t *TagController) GetByTagIdHandler(ctx *gin.Context) {
	tagId, err := strconv.Atoi(ctx.Param("tags_id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	tags, err := t.service.FindById(tagId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Success Get Tags",
		"data":    tags,
	})

}

func (t *TagController) Route() {
	router := t.rg.Group("/tags")
	router.GET("/:tags_id", t.GetByTagIdHandler)
	router.GET("/", t.GetAllTagHandler)

	routerAuth := router.Group("/", t.md.CheckToken())
	routerAuth.POST("/", t.CreateTagHandler)
}

func NewTagController(tS service.TagService, rg *gin.RouterGroup, md middleware.AuthMiddleware) *TagController {
	return &TagController{
		service: tS,
		rg:      rg,
		md:      md,
	}
}
