package services

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)


type Product struct {
	SKUID      string    `json:"sku_id"`
	Name       string    `json:"name"`
	Description string   `json:"description"`
	Price      float64   `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateProductRequest struct {
	MerchantID string `json:"merchant_id"`
	Product    Product `json:"product"`
}

var MerchantProductsMap map[string]map[string]Product

func Initialize() {
	MerchantProductsMap = make(map[string]map[string]Product)
}

func CreateProductHandler(c *gin.Context) {
	var createRequest CreateProductRequest
	

	if err := c.BindJSON(&createRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	merchantID := createRequest.MerchantID

	if _, exists := MerchantProductsMap[merchantID]; !exists {
		MerchantProductsMap[merchantID] = make(map[string]Product)
	}

	product := createRequest.Product
	product.CreatedAt = time.Now()

	MerchantProductsMap[merchantID][product.SKUID] = product

	c.JSON(http.StatusCreated, product)
}

func DisplayProductsHandler(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	skuID := c.Query("sku_id")

	if skuID != "" {
		product, found := MerchantProductsMap[merchantID][skuID]
		if !found {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	} else {
		products := make([]Product, 0, len(MerchantProductsMap[merchantID]))
		for _, product := range MerchantProductsMap[merchantID] {
			products = append(products, product)
		}
		c.JSON(http.StatusOK, products)
	}
}

func EditProductHandler(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	skuID := c.Query("sku_id")

	existingProduct, found := MerchantProductsMap[merchantID][skuID]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}

	var updatedProduct Product
	if err := c.BindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	existingProduct.Name = updatedProduct.Name
	existingProduct.Description = updatedProduct.Description
	existingProduct.Price = updatedProduct.Price

	MerchantProductsMap[merchantID][skuID] = existingProduct

	c.JSON(http.StatusOK, existingProduct)
}

func DeleteProductHandler(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	skuID := c.Query("sku_id")

	if products, merchantExists := MerchantProductsMap[merchantID]; merchantExists {
		if _, productExists := products[skuID]; productExists {
			// Delete the product
			delete(products, skuID)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}