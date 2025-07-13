package product

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrProductAlreadyExists = errors.New("product already exists")
	ErrInsufficientStock    = errors.New("insufficient stock")
)

// Product domain model
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int32
	Category    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Repository interface for product operations
type Repository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	GetByID(ctx context.Context, id string) (*Product, error)
	Update(ctx context.Context, product *Product) (*Product, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int, category string) ([]*Product, int, error)
	UpdateStock(ctx context.Context, productID string, quantity int32) (*Product, error)
}

// InMemoryRepository implements Repository interface using in-memory storage
type InMemoryRepository struct {
	products map[string]*Product
	mutex    sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		products: make(map[string]*Product),
	}
}

// Create creates a new product
func (r *InMemoryRepository) Create(ctx context.Context, product *Product) (*Product, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Generate ID if not provided
	if product.ID == "" {
		product.ID = uuid.New().String()
	}

	// Check if product already exists
	if _, exists := r.products[product.ID]; exists {
		return nil, ErrProductAlreadyExists
	}

	// Set timestamps
	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	// Store product
	r.products[product.ID] = product

	return product, nil
}

// GetByID retrieves a product by ID
func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Product, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	product, exists := r.products[id]
	if !exists {
		return nil, ErrProductNotFound
	}

	return product, nil
}

// Update updates an existing product
func (r *InMemoryRepository) Update(ctx context.Context, product *Product) (*Product, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	existingProduct, exists := r.products[product.ID]
	if !exists {
		return nil, ErrProductNotFound
	}

	// Update fields
	product.CreatedAt = existingProduct.CreatedAt
	product.UpdatedAt = time.Now()

	// Store updated product
	r.products[product.ID] = product

	return product, nil
}

// Delete deletes a product by ID
func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.products[id]; !exists {
		return ErrProductNotFound
	}

	delete(r.products, id)
	return nil
}

// List retrieves products with pagination and optional category filter
func (r *InMemoryRepository) List(ctx context.Context, offset, limit int, category string) ([]*Product, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Convert map to slice and filter by category if specified
	var products []*Product
	for _, product := range r.products {
		if category == "" || product.Category == category {
			products = append(products, product)
		}
	}

	total := len(products)

	// Apply pagination
	start := offset
	if start > total {
		start = total
	}

	end := start + limit
	if end > total {
		end = total
	}

	if start >= total {
		return []*Product{}, total, nil
	}

	return products[start:end], total, nil
}

// UpdateStock updates product stock
func (r *InMemoryRepository) UpdateStock(ctx context.Context, productID string, quantity int32) (*Product, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	product, exists := r.products[productID]
	if !exists {
		return nil, ErrProductNotFound
	}

	newStock := product.Stock + quantity
	if newStock < 0 {
		return nil, ErrInsufficientStock
	}

	// Create a copy and update stock
	updatedProduct := *product
	updatedProduct.Stock = newStock
	updatedProduct.UpdatedAt = time.Now()

	// Store updated product
	r.products[productID] = &updatedProduct

	return &updatedProduct, nil
}
