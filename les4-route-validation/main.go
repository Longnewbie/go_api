package main

import (
	v1handler "gialong04/go_api/internal/api/v1/handler"
	v2handler "gialong04/go_api/internal/api/v2/handler"
	"gialong04/go_api/middleware"
	"gialong04/go_api/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := utils.RegisterValidators(); err != nil {
		panic(err)
	}

	err := godotenv.Load()

	if err != nil {
		log.Printf("No .env file found")
	}

	r := gin.Default()

	go middleware.CleanupClients() // run parallel with main

	// 1. middleware is used for all routes
	r.Use(
		middleware.LoggerMiddleware(),
		middleware.ApiKeyMiddleware(),
		middleware.RateLimitingMiddleware(),
	)

	v1 := r.Group("/api/v1")
	{
		user := v1.Group("/users")
		{
			// users group v1
			userHandlerV1 := v1handler.NewUserHandler()
			user.GET("", userHandlerV1.GetUsersV1)
			user.GET("/:id", userHandlerV1.GetUsersByIdV1)
			user.GET("/admin/:uuid", userHandlerV1.GetUsersByUUIDV1)
			user.POST("", userHandlerV1.PostUsersV1)
			user.PUT("/:id", userHandlerV1.PutUsersV1)
			user.DELETE("/:id", userHandlerV1.DeleteUsersV1)
		}

		product := v1.Group("/products")
		{
			// products group v1
			productHandlerV1 := v1handler.NewProductHandler()
			product.GET("", productHandlerV1.GetProductsV1)
			product.GET("/:slug", productHandlerV1.GetProductsBySlugV1)
			// product.GET("/:id", productHandlerV1.GetProductsByIdV1)
			product.POST("", productHandlerV1.PostProductsV1)
			product.PUT("/:id", productHandlerV1.PutProductsV1)
			product.DELETE("/:id", productHandlerV1.DeleteProductsV1)
		}

		// 2. middleware is used for a group route
		category := v1.Group("/categories")
		{
			// category group v1
			categoryHandlerV1 := v1handler.NewCategoryHandler()
			category.GET("/:category", categoryHandlerV1.GetCategoryV1)
			category.POST("", categoryHandlerV1.PostCategoryV1)
		}

		news := v1.Group("/news")
		{
			// category group v1
			newsHandlerV1 := v1handler.NewNewsHandler()
			news.GET("/", newsHandlerV1.GetNewsV1)
			news.GET("/:slug", middleware.SimpleMiddleware(), newsHandlerV1.GetNewsV1) // 3. middleware is used for a route
			news.POST("/", newsHandlerV1.PostNewsV1)
			news.POST("/upload-file", newsHandlerV1.PostUploadFileNewsV1)
			news.POST("/upload-multiple-files", newsHandlerV1.PostUploadMultipleFileNewsV1)
		}
	}

	v2 := r.Group("/api/v2")
	{
		user := v2.Group("/users")
		{
			// users group v2
			userHandlerV2 := v2handler.NewUserHandler()
			user.GET("", userHandlerV2.GetUsersV1)
			user.GET("/:id", userHandlerV2.GetUsersByIdV1)
			user.POST("", userHandlerV2.PostUsersV1)
			user.PUT("/:id", userHandlerV2.PutUsersV1)
			user.DELETE("/:id", userHandlerV2.DeleteUsersV1)
		}
	}

	r.StaticFS("/images", gin.Dir("./uploads", false))

	r.Run(":8080")
}
