package product_store

import (
	"context"
	"goCodetest/internal/massage"
	"goCodetest/internal/models"
	"goCodetest/pkg/logger"
	"slices"
	"strings"
	"sync"
)

type Product struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,gte=0"`
	Quantity int64   `json:"quantity" binding:"required,gte=0"`
}

type UpdateProductRequest struct {
	Name     string  `json:"name" binding:"required"`
	Price    float64 `json:"price" binding:"required,gte=0"`
	Quantity int64   `json:"quantity" binding:"required,gte=0"`
}

type ProductStore struct {
	product map[int64]*Product
	id      int64
	mu      sync.RWMutex
}

func NewProductStore() *ProductStore {
	return &ProductStore{
		product: make(map[int64]*Product),
		id:      0,
	}
}

func (ps *ProductStore) AddProduct(ctx context.Context, product *Product) (id int64, ok bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	id = ps.id
	ps.id++
	product.ID = id
	ps.product[product.ID] = product
	logger.Info("TraceId: %s - Product with ID %d added successfully", ctx.Value(models.TraceIDKey), product.ID)
	return product.ID, true
}

func (ps *ProductStore) DeleteProduct(ctx context.Context, id int64) bool {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	if _, exists := ps.product[id]; !exists {
		logger.Error("Product with ID %d does not exist", id)
		return false
	}
	delete(ps.product, id)
	logger.Info("TraceId: %s - Product with ID %d deleted successfully", ctx.Value(models.TraceIDKey), id)
	go massage.EmailNotificationForProductDeletion(ctx, id)
	return true
}

func (ps *ProductStore) UpdateProductByID(ctx context.Context, id int64, newProduct *Product) (*Product, bool) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	existedProduct, exists := ps.product[id]
	if !exists {
		logger.Error("Product with ID %d does not exist", id)
		return nil, false
	}
	existedProduct.Name = newProduct.Name
	existedProduct.Price = newProduct.Price
	existedProduct.Quantity = newProduct.Quantity
	logger.Info("TraceId: %s - Product with ID %d updated successfully", ctx.Value(models.TraceIDKey), id)
	return ps.product[id], true
}

func (ps *ProductStore) GetProductByID(ctx context.Context, id int64) (*Product, bool) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	product, exists := ps.product[id]
	if !exists {
		logger.Error("Product with ID %d does not exist", id)
		return nil, false
	}
	logger.Info("TraceId: %s - Product with ID %d retrieved successfully", ctx.Value(models.TraceIDKey), id)
	return product, true
}

func (ps *ProductStore) productFilter(product *Product, ids []int64, name string, priceRange [2]float64, quantityRange [2]int64) bool {
	if len(ids) > 0 {
		if !slices.Contains(ids, product.ID) {
			return false
		}
	}

	if name != "" && !strings.Contains(strings.ToLower(product.Name), strings.ToLower(name)) {
		return false
	}

	if priceRange != [2]float64{0, 0} && (product.Price <= priceRange[0] || product.Price >= priceRange[1]) {
		return false
	}
	if quantityRange != [2]int64{0, 0} && (product.Quantity <= quantityRange[0] || product.Quantity >= quantityRange[1]) {
		return false
	}
	return true
}

func (ps *ProductStore) GetProductsList(ctx context.Context, page int, pageSize int, ids []int64, name string, priceRange [2]float64, quantityRange [2]int64) []*Product {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var filteredProducts []*Product
	for _, product := range ps.product {
		if ok := ps.productFilter(product, ids, name, priceRange, quantityRange); ok {
			filteredProducts = append(filteredProducts, product)
		}
	}

	slices.SortFunc(filteredProducts, func(a, b *Product) int {
		if a.ID < b.ID {
			return -1
		} else if a.ID > b.ID {
			return 1
		}
		return 0
	})

	start := (page - 1) * pageSize
	end := start + pageSize
	if start >= len(filteredProducts) {
		return nil
	}
	if end > len(filteredProducts) {
		end = len(filteredProducts)
	}
	logger.Info("TraceId: %s - Retrieved %d products for page %d with page size %d", ctx.Value(models.TraceIDKey), len(filteredProducts), page, pageSize)
	return filteredProducts[start:end]
}
