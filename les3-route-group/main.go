package main

import (
	v1handler "gialong04/go_api/internal/api/v1/handler"
	v2handler "gialong04/go_api/internal/api/v2/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

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

		category := v1.Group("/categories")
		{
			// category group v1
			categoryHandlerV1 := v1handler.NewCategoryHandler()
			category.GET("/:category", categoryHandlerV1.GetCategoryV1)
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

	r.Run(":8080")
}
