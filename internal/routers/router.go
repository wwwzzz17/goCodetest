package routers

import (
	"context"
	"errors"
	"goCodetest/internal/models"
	"goCodetest/pkg/logger"
	"goCodetest/pkg/product_store"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var pStore *product_store.ProductStore

func Setup(app *gin.Engine, store *product_store.ProductStore) {
	pStore = store
	app.Use(gin.Recovery())
	// Debug for gin
	if gin.Mode() == gin.DebugMode {
		logger.Debug("Running in debug mode")
	}
	SetRoutes(app) // Set up all API routes.
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// simple auth middleware example
		// auth_token := c.Request.Header.Get("Authorization")
		// if auth_token == "" {
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		// 	c.Abort()
		// 	return
		// }
		c.Next() // 继续处理请求
	}
}

// use context to record trace ID
func replaceCtxWithTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		originalCtx := c.Request.Context()
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			logger.Warn("X-Trace-ID header not found, using default trace ID")
			traceID = uuid.New().String()
		}
		newCtx := context.WithValue(originalCtx, models.TraceIDKey, traceID)
		c.Request = c.Request.WithContext(newCtx)
		c.Next()
	}
}

func SetRoutes(app *gin.Engine) {
	app.Use(Cors())
	app.Use(AuthMiddleware())
	app.Use(replaceCtxWithTraceID())

	app.GET("/health", func(c *gin.Context) {
		respondWithSuccess(c, "success", map[string]interface{}{
			"status": "healthy",
		})
	})

	app.POST("/products", addProduct)
	app.GET("/products/:id", getProductByID)
	app.PUT("/products/:id", updateProductByID)
	app.DELETE("/products/:id", deleteProductByID)
	app.GET("/products", getProductsList)

}

func respondWithError(c *gin.Context, code int, message string) {
	logger.Error("%s", message)
	c.JSON(code, models.HTTPResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func respondWithSuccess(c *gin.Context, message string, data map[string]interface{}) {
	c.JSON(http.StatusOK, models.HTTPResponse{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func parseIDParam(c *gin.Context) (int64, error) {
	idStr := c.Param("id")
	if idStr == "" {
		return 0, errors.New("ID is required")
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid ID format: " + err.Error())
	}

	return id, nil
}

func addProduct(c *gin.Context) {
	ctx := c.Request.Context()

	var addProductReq product_store.Product
	if err := c.ShouldBindJSON(&addProductReq); err != nil {
		logger.Error("addProduct: bind json err %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	id, ok := pStore.AddProduct(ctx, &addProductReq)

	if !ok {
		logger.Error("addProduct: add product failed")
		respondWithError(c, http.StatusInternalServerError, "Failed to add product")
		return
	}

	respondWithSuccess(c, "Product added successfully", map[string]interface{}{"id": id})

	logger.Info("Product added successfully with ID: %d", id)
}

func getProductByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := parseIDParam(c)
	if err != nil {
		logger.Error("getProductByID: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	product, exists := pStore.GetProductByID(ctx, id)
	if !exists {
		logger.Error("getProductByID: product with ID %d does not exist", id)
		respondWithError(c, http.StatusNotFound, "Product not found")
		return
	}

	respondWithSuccess(c, "Product retrieved successfully", map[string]interface{}{
		"result": &product,
	})
	logger.Info("Product retrieved successfully with ID: %d", id)
}

func deleteProductByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := parseIDParam(c)
	if err != nil {
		logger.Error("deleteProductByID: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	if !pStore.DeleteProduct(ctx, id) {
		logger.Error("deleteProductByID: product with ID %d does not exist", id)
		respondWithError(c, http.StatusNotFound, "Product not found")
		return
	}

	respondWithSuccess(c, "Product deleted successfully", nil)
}

func updateProductByID(c *gin.Context) {
	ctx := c.Request.Context()
	id, err := parseIDParam(c)
	if err != nil {
		logger.Error("updateProductByID: %s", err.Error())
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}

	var productRequest *product_store.UpdateProductRequest
	if err := c.ShouldBindJSON(&productRequest); err != nil {
		logger.Error("updateProductByID: bind json err %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Invalid request data: "+err.Error())
		return
	}

	product, ok := pStore.UpdateProductByID(ctx, id, &product_store.Product{
		Name:     productRequest.Name,
		Price:    productRequest.Price,
		Quantity: productRequest.Quantity,
	})

	if !ok {
		logger.Error("updateProductByID: product with ID %d does not exist", id)
		respondWithError(c, http.StatusNotFound, "Product not found")
		return
	}

	respondWithSuccess(c, "Product updated successfully", map[string]interface{}{
		"result": &product,
	})
	logger.Info("Product updated successfully with ID: %d", id)
}

// GET /products?page=1&limit=10&ids[]=1&ids[]=2&ids[]=3&name=premium&price_min=5.0&price_max=50.0&quantity_min=10&quantity_max=1000
func getProductsList(c *gin.Context) {
	ctx := c.Request.Context()

	var query models.ProductListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		logger.Error("getProductsList: bind query err %s", err.Error())
		respondWithError(c, http.StatusBadRequest, "Invalid query parameters: "+err.Error())
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	if query.PriceMax == 0 {
		query.PriceMax = math.MaxFloat64
	}
	if query.QuantityMax == 0 {
		query.QuantityMax = math.MaxInt64
	}

	products, total := pStore.GetProductsList(ctx, query.Page, query.Limit, query.IDs, query.Name, [2]float64{query.PriceMin, query.PriceMax}, [2]int64{query.QuantityMin, query.QuantityMax})
	if len(products) == 0 {
		logger.Info("getProductsList: no products found")
		respondWithSuccess(c, "No products found", map[string]interface{}{
			"results": nil,
			"total":   0,
		})
		return
	}

	respondWithSuccess(c, "Products retrieved successfully", map[string]interface{}{
		"results": products,
		"total":   total,
	})
}
