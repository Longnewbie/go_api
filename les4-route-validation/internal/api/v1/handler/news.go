package v1handler

import (
	"fmt"
	"gialong04/go_api/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
}

type PostNewsV1Param struct {
	Title  string `form:"title" binding:"required"`
	Status string `form:"status" binding:"required,oneof=1 2"`
}

func NewNewsHandler() *NewsHandler {
	return &NewsHandler{}
}

func (c *NewsHandler) GetNewsV1(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if slug == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "News found",
			"slug":    "No News",
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "News found",
			"slug":    slug,
		})
	}
}

func (c *NewsHandler) PostNewsV1(ctx *gin.Context) {
	var param PostNewsV1Param
	if err := ctx.ShouldBind(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// giới hạn file 5MB
	if image.Size > 5<<20 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File too large (5MB)"})
		return
	}

	// os.ModePerm = 0777 (octal)
	// Có quyền đọc, ghi, thực thi (read, write, execute) chot ất cả mọi người (owner, group, others)
	err = os.MkdirAll("./uploads", os.ModePerm) // nếu không có folder thì tạo folder
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create upload folder"})
		return
	}

	dst := fmt.Sprintf("./uploads/%s", filepath.Base(image.Filename))

	if err := ctx.SaveUploadedFile(image, dst); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot save file"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Post news (V1)",
		"title":   param.Title,
		"status":  param.Status,
		"image":   image.Filename,
		"path":    dst,
	})
}

func (c *NewsHandler) PostUploadFileNewsV1(ctx *gin.Context) {
	var param PostNewsV1Param

	if err := ctx.ShouldBind(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	image, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	filename, err := utils.ValidateAndSaveFile(image, "./uploads")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Post news (V1)",
		"title":   param.Title,
		"status":  param.Status,
		"image":   filename,
		"path":    "./upload/" + filename,
	})
}

func (c *NewsHandler) PostUploadMultipleFileNewsV1(ctx *gin.Context) {
	var param PostNewsV1Param
	const publicURL = "http://localhost:8080/images/"

	if err := ctx.ShouldBind(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
		return
	}

	images := form.File["images"]
	if len(images) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	var successFiles []string           // data accepted
	var failedFiles []map[string]string // data failed
	for _, img := range images {
		filename, err := utils.ValidateAndSaveFile(img, "./uploads")
		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"filename": img.Filename,
				"error":    err.Error(),
			})

			continue
		}

		publicImageURL := publicURL + filename
		successFiles = append(successFiles, publicImageURL)
	}

	response := gin.H{
		"message":       "Post news (V1)",
		"title":         param.Title,
		"status":        param.Status,
		"success_files": successFiles,
	}

	if len(failedFiles) > 0 {
		response["message"] = "Upload completed with partial error"
		response["error_files"] = failedFiles
	}

	ctx.JSON(http.StatusCreated, response)
}
