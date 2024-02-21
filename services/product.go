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

var MerchantProductsMap map[string][]Product

func Initialize() {
	MerchantProductsMap = make(map[string][]Product)
}

func CreateProductHandler(c *gin.Context) {
	var createRequest CreateProductRequest
	if err := c.BindJSON(&createRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	merchantID := createRequest.MerchantID
	product := createRequest.Product
	product.CreatedAt = time.Now()

	MerchantProductsMap[merchantID] = append(MerchantProductsMap[merchantID], product)

	c.JSON(http.StatusCreated, product)
}

func DisplayProductsHandler(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	products, found := MerchantProductsMap[merchantID]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func EditProductHandler(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	skuID := c.Query("sku_id")

	products, found := MerchantProductsMap[merchantID]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}

	var updatedProduct Product
	if err := c.BindJSON(&updatedProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Find and update the product with the matching SKU ID
	for i, p := range products {
		if p.SKUID == skuID {
			updatedProduct.CreatedAt = p.CreatedAt
			products[i] = updatedProduct
			c.JSON(http.StatusOK, products[i])
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}

func DeleteProductHandler(c *gin.Context) {
	merchantID := c.Query("merchant_id")
	skuID := c.Query("sku_id")

	products, found := MerchantProductsMap[merchantID]
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Merchant not found"})
		return
	}

	// Find and remove the product with the matching SKU ID
	for i, p := range products {
		if p.SKUID == skuID {
			MerchantProductsMap[merchantID] = append(products[:i], products[i+1:]...)
			c.Status(http.StatusNoContent)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
}