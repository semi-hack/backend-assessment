package routes

import (
	"github.com/gin-gonic/gin"

	"backend-assessment/services"
)

func Initialize() {
	r := gin.Default()

	r.GET("/products", services.DisplayProductsHandler)
	r.POST("/product/create", services.CreateProductHandler)
	r.PUT("/product/edit", services.EditProductHandler)
	r.DELETE("/product/delete", services.DeleteProductHandler)

	r.Run()
}