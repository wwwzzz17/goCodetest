package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"goCodetest/internal/routers"
	"goCodetest/pkg/product_store"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupTestProductStore() *product_store.ProductStore {
	return product_store.NewProductStore()
}

func setupTestRouter() *gin.Engine {
	store := setupTestProductStore()
	router := gin.New()
	routers.Setup(router, store)
	return router
}

func TestAddProduct_Success(t *testing.T) {
	router := setupTestRouter()

	product := map[string]interface{}{
		"name":     "Test",
		"price":    4.99,
		"quantity": 100,
	}

	jsonData, _ := json.Marshal(product)

	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonData))
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if !assert.Equal(t, http.StatusOK, w.Code) {
		t.Fatal("Expected status OK, got", w.Code)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	if !assert.NoError(t, err) {
		t.Fatal(err)
	}

	if !assert.NotNil(t, responseBody) {
		t.Fatal("Expected response body to be not nil")
	}
}
