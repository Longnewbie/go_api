package v1handler

import (
	"gialong04/go_api/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
}

type PostCategoryV1Param struct {
	Name   string `form:"name" binding:"required"`
	Status string `form:"status" binding:"required,oneof=1 2"`
}

type GetCategoryV1Param struct {
	Category string `uri:"category" binding:"oneof=php python golang"`
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

func (c *CategoryHandler) GetCategoryV1(ctx *gin.Context) {
	var param GetCategoryV1Param
	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	log.Printf("into GetCategoryV1")

	val, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Missing username"})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Category found",
		"category": param.Category,
		"username": val,
	})
}

func (c *CategoryHandler) PostCategoryV1(ctx *gin.Context) {
	var param PostCategoryV1Param
	if err := ctx.ShouldBind(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Post category (V1)",
		"name":    param.Name,
		"status":  param.Status,
	})
}
