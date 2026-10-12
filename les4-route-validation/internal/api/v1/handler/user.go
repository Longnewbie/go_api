package v1handler

import (
	"gialong04/go_api/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
}

type GetUsersByIdV1Param struct {
	ID int `uri:"id" binding:"gt=0"`
}

type GetUsersByUUIDV1Param struct {
	Uuid int `uri:"uuid" binding:"uuid"`
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (u *UserHandler) GetUsersV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "Get all Users (V1)"})
}

func (u *UserHandler) GetUsersByIdV1(ctx *gin.Context) {
	var param GetUsersByIdV1Param

	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	log.Printf("into GetUsersByIdV1")

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get User by ID (V1)",
		"user_id": param.ID,
	})
}

func (u *UserHandler) GetUsersByUUIDV1(ctx *gin.Context) {
	var param GetUsersByUUIDV1Param
	if err := ctx.ShouldBindUri(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.HandleValidationErrors(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "Get user by UUID (V1)",
		"user_uuid": param.Uuid,
	})
}

func (u *UserHandler) PostUsersV1(ctx *gin.Context) {
	ctx.JSON(http.StatusCreated, gin.H{"message": "User created (V1)"})
}

func (u *UserHandler) PutUsersV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "User updated (V1)"})
}

func (u *UserHandler) DeleteUsersV1(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted (V1)"})
}
