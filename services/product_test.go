package services

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)


func TestCreateProductHandler(t *testing.T) {

	Initialize()

	// Create a test router
	router := gin.Default()
	router.POST("/products/create", CreateProductHandler)

	testCase := CreateProductRequest{
		MerchantID: "123",
		Product: Product{
			SKUID:      "ABC123",
			Name:       "Test Product",
			Description: "This is a test product.",
			Price:      19.99,
		},
	}

	payload, err := json.Marshal(testCase)

	assert.NoError(t, err)

	req, err := http.NewRequest("POST", "/products/create?merchant_id="+testCase.MerchantID, strings.NewReader(string(payload)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var createdProduct Product
	err = json.Unmarshal(w.Body.Bytes(), &createdProduct)
	assert.NoError(t, err)

	assert.Equal(t, testCase.Product.SKUID, createdProduct.SKUID)
	assert.Equal(t, testCase.Product.Name, createdProduct.Name)
	assert.Equal(t, testCase.Product.Description, createdProduct.Description)
	assert.Equal(t, testCase.Product.Price, createdProduct.Price)

}

func TestEditProductHandler(t *testing.T) {
	Initialize()

	router := gin.Default()
	router.PUT("/products/edit", EditProductHandler)

	MerchantProductsMap = map[string]map[string]Product{
		"merchant1": {
			"sku123": {SKUID: "sku123", Name: "Product1", Description: "Description1", Price: 19.99, CreatedAt: time.Now()},
		},
	}

	// Test
	editRequest := Product{
		SKUID:      "sku123",
		Name:       "UpdatedProduct",
		Description: "UpdatedDescription",
		Price:      49.99,
	}

	reqBody, err := json.Marshal(editRequest)
	assert.NoError(t, err)

	req, _ := http.NewRequest("PUT", "/products/edit?merchant_id=merchant1&sku_id=sku123", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	updatedProduct := MerchantProductsMap["merchant1"]["sku123"]

	assert.Equal(t, editRequest.Name, updatedProduct.Name)
	assert.Equal(t, editRequest.Description, updatedProduct.Description)
	assert.Equal(t, editRequest.Price, updatedProduct.Price)
}

func TestDisplayProductsHandler(t *testing.T) {
	Initialize()

	router := gin.Default()
	router.GET("/products", DisplayProductsHandler)

	MerchantProductsMap = map[string]map[string]Product{
		"merchant1": {
			"sku123": {SKUID: "sku123", Name: "Product1", Description: "Description1", Price: 19.99, CreatedAt: time.Now()},
			"sku456": {SKUID: "sku456", Name: "Product2", Description: "Description2", Price: 29.99, CreatedAt: time.Now()},
		},
	}

	t.Run("Retrieve Single Product", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products?merchant_id=merchant1&sku_id=sku123", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		var product Product
		err := json.Unmarshal(w.Body.Bytes(), &product)
		assert.NoError(t, err)
		assert.Equal(t, "sku123", product.SKUID)
	})

	t.Run("Retrieve All Products for Merchant", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products?merchant_id=merchant1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusOK, w.Code)
		var products []Product
		err := json.Unmarshal(w.Body.Bytes(), &products)
		assert.NoError(t, err)
		assert.Len(t, products, 2)
	})

	t.Run("Product Not Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/products?merchant_id=merchant1&sku_id=sku_not_found", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assertions
		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Product not found", response["error"])
	})
}

func TestDeleteProductHandler(t *testing.T) {
	// Setup
	Initialize()

	router := gin.Default()
	router.DELETE("/product/delete", DeleteProductHandler)

	MerchantProductsMap = map[string]map[string]Product{
		"merchant1": {
			"sku123": {SKUID: "sku123", Name: "Product1", Description: "Description1", Price: 19.99, CreatedAt: time.Now()},
		},
	}

	// Test
	req, _ := http.NewRequest("DELETE", "/product/delete?merchant_id=merchant1&sku_id=sku123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusNoContent, w.Code)

	// Additional Assertions
	_, productExists := MerchantProductsMap["merchant1"]["sku123"]
	assert.False(t, productExists)
}