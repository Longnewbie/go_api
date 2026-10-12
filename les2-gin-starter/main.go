package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() // create router

	// Define a simple GET endpoint
	r.GET("/demo", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Hello, My name is Long!"})
	})

	r.GET("/users", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"Message": "Get all Users"})
	})

	r.GET("/user/:userId", func(ctx *gin.Context) {
		userId := ctx.Param("userId")

		ctx.JSON(200, gin.H{
			"data": "Get user with ID: " + userId,
		})
	})

	r.GET("/products", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"Message": "Get all Products"})
	})

	r.GET("/product/detail/:productName", func(ctx *gin.Context) {
		productName := ctx.Param("productName")
		price := ctx.Query("price")
		color := ctx.Query("color")
		ctx.JSON(200, gin.H{
			"Message":       "Get all Products",
			"data":          "Get product with name: " + productName,
			"product_price": price,
			"product_color": color,
		})
	})

	r.Run(":8080") // init server - localhost:8080
}
