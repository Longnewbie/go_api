package v1handler

import (
	"fmt"
	"gialong04/go_api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductHandler struct {
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

type GetProductsBySlugV1Param struct {
	Slug string `uri:"slug" binding:"slug,min=5,max=100"`
}

type GetProductV1Param struct {
	Search string `form:"search" binding:"required,min=3,max=50,search"`
	Limit  int    `form:"limit" binding:"omitempty,gte=1,lte=100"`
	Email  string `form:"email" binding:"omitempty,email"`
	Date   string `form:"date" binding:"omitempty,datetime=2006-01-02"`
}

type ProductImage struct {
	ImageName string `json:"image_name" binding:"required"`
	ImageLink string `json:"image_link" binding:"required,file_ext=jpg png jpeg gif"`
}

type ProductAttribute struct {
	AttributeName  string `json:"attribute_name" binding:"required"`
	AttributeValue string `json:"attribute_value" binding:"required"`
}

type ProductInfo struct {
	InfoKey   string `json:"info_key" binding:"required"`
	InfoValue string `json:"info_value" binding:"required"`
}

type PostProductsV1Param struct {
	Name             string                 `json:"name" binding:"required,min=3,max=100"`
	Price            int                    `json:"price" binding:"required,min_int=10000"`
	Display          *bool                  `json:"display" binding:"omitempty"`
	ProductImage     ProductImage           `json:"product_image" binding:"required"`
	Tags             []string               `json:"tags" binding:"required,gt=3,lt=5"`
	ProductAttribute []ProductAttribute     `json:"product_attribute" binding:"required,gt=0,dive"`
	ProductInfo      map[string]ProductInfo `json:"product_info" binding:"required,gt=0,dive"`
	ProductMetadata  map[string]any         `json:"product_metadata" binding:"omitempty"`
}

func (p *ProductHandler) GetProductsV1(ctx *gin.Context) {
	var param GetProductV1Param
	if err := ctx.ShouldBindQuery(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	if param.Limit == 0 {
		param.Limit = 1
	}

	if param.Email == "" {
		param.Email = "No Email"
	}

	if param.Date == "" {
		param.Date = time.Now().Format("2006-01-02")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get all Products (V1)",
		"search":  param.Search,
		"limit":   param.Limit,
		"email":   param.Email,
		"date":    param.Date,
	})
}

func (p *ProductHandler) GetProductsByIdV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Get Product by ID (V1)"})
}

func (p *ProductHandler) GetProductsBySlugV1(ctx *gin.Context) {
	var param GetProductsBySlugV1Param

	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get Product by Slug (V1)",
		"slug":    param.Slug,
	})
}

func (p *ProductHandler) PostProductsV1(ctx *gin.Context) {
	var param PostProductsV1Param

	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	for key := range param.ProductInfo {
		if _, err := uuid.Parse(key); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"errors": gin.H{
					"product_info": fmt.Sprintf("Key '%s' trong product_info khong phai la UUID hop le", key),
				},
			})
			return
		}
	}

	if param.Display == nil {
		defaultDisplay := true
		param.Display = &defaultDisplay
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":           "Product created (V1)",
		"name":              param.Name,
		"price":             param.Price,
		"display":           param.Display,
		"product_image":     param.ProductImage,
		"tags":              param.Tags,
		"product_attribute": param.ProductAttribute,
		"product_info":      param.ProductInfo,
		"product_metadata":  param.ProductMetadata,
	})
}

func (p *ProductHandler) PutProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Product updated (V1)"})
}

func (p *ProductHandler) DeleteProductsV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Product deleted (V1)"})
}
