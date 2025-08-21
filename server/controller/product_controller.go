package controller

import (
	"develapar-server/middleware"
	"develapar-server/model"
	"develapar-server/service"
	"develapar-server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
)

type ProductController struct {
	s              service.ProductService
	rg             *gin.RouterGroup
	mD             middleware.AuthMiddleware
	errorHandler   middleware.ErrorHandler
	responseHelper *utils.ResponseHelper
}

func NewProductController(
	pS service.ProductService,
	rg *gin.RouterGroup,
	mD middleware.AuthMiddleware,
	errorHandler middleware.ErrorHandler,
) *ProductController {
	return &ProductController{
		s:            pS,
		rg:           rg,
		mD:           mD,
		errorHandler: errorHandler,
	}
}

func (c *ProductController) Route() {

	router := c.rg.Group("/products/categories")

	routerAuth := router.Group("/", c.mD.CheckToken())
	routerAuth.POST("/", c.CreateProductCategory)
}

func (c *ProductController) CreateProductCategory(ginCtx *gin.Context) {
	reqCtx, cancel := context.WithTimeout(ginCtx.Request.Context(), 15*time.Second)
	defer cancel()

	var payload model.ProductCategory
	if err := ginCtx.ShouldBindJSON(&payload); err != nil {
		appErr := c.errorHandler.ValidationError(reqCtx, "payload", "Invalid request payload: "+err.Error())
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)

		return
	}

	data, err := c.s.CreateProductCategory(reqCtx, payload)
	if err != nil {
		if reqCtx.Err() == context.DeadlineExceeded {
			appErr := c.errorHandler.TimeoutError(reqCtx, "Create Product Category")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if reqCtx.Err() == context.Canceled {
			appErr := c.errorHandler.CancellationError(reqCtx, "Create Product Category")
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}
		if appErr, ok := err.(*utils.AppError); ok {
			c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
			return
		}

		appErr := c.errorHandler.WrapError(reqCtx, err, utils.ErrInternal, "failed to create category")
		appErr.StatusCode = 500
		c.errorHandler.HandleError(reqCtx, ginCtx, appErr)
		return

	}

	responseData := gin.H{
		"message":          "Product Category created successfully",
		"product_category": data,
	}
	c.responseHelper.SendCreated(ginCtx, responseData)
}
