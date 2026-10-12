package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
)

func SimpleMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Println("Start func - check from middleware")
		ctx.Writer.Write([]byte("Start func - check from middleware"))

		ctx.Next()

		log.Println("End func - check from middleware")
		ctx.Writer.Write([]byte("End func - check from middleware"))
	}
}
