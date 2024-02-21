package services

import (
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

	initialProduct := Product{
		SKUID:      "ABC123",
		Name:       "Initial Product",
		Description: "This is the initial product.",
		Price:      29.99,
		CreatedAt:  time.Now(),
	}

	MerchantProductsMap["123"] = append(MerchantProductsMap["123"], initialProduct)

	// Define a test case for the edit request
	editRequest := Product{
		Name:        "Updated Product",
		Description: "This is the updated product.",
		Price:       39.99,
	}

	payload, err := json.Marshal(editRequest)
	assert.NoError(t, err)

	req, err := http.NewRequest("PUT", "/products/edit?merchant_id=123&sku_id=ABC123", strings.NewReader(string(payload)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var updatedProduct Product
	err = json.Unmarshal(w.Body.Bytes(), &updatedProduct)
	assert.NoError(t, err)

	assert.Equal(t, editRequest.Name, updatedProduct.Name)
	assert.Equal(t, editRequest.Description, updatedProduct.Description)
	assert.Equal(t, editRequest.Price, updatedProduct.Price)
}

func TestDisplayProductsHandler(t *testing.T) {
	Initialize()

	router := gin.Default()
	router.GET("/products", DisplayProductsHandler)

	MerchantProductsMap = map[string][]Product{
		"merchant1": {
			{SKUID: "123", Name: "Product1", Description: "Description1", Price: 19.99, CreatedAt: time.Now()},
			{SKUID: "456", Name: "Product2", Description: "Description2", Price: 29.99, CreatedAt: time.Now()},
		},
	}

	// Test
	req, _ := http.NewRequest("GET", "/products?merchant_id=merchant1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	var products []Product
	err := json.Unmarshal(w.Body.Bytes(), &products)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
}
