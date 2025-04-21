package controller

import (
	"develapar-server/model"
	"develapar-server/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagController struct {
	service service.TagService
	rg      *gin.RouterGroup
}

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
	router.POST("/", t.CreateTagHandler)
}

func NewTagController(tS service.TagService, rg *gin.RouterGroup) *TagController {
	return &TagController{
		service: tS,
		rg:      rg,
	}
}
