package dataloaders

import (
	"context"
	"fmt"
	"sync"
	"time"

	"learning/internal/product"
)

// ProductLoader provides batched loading of products
type ProductLoader struct {
	repo  product.Repository
	cache map[string]*product.Product
	mutex sync.RWMutex

	// Batching fields
	batch      map[string][]chan productResult
	batchMutex sync.Mutex
	maxBatch   int
	wait       time.Duration
}

type productResult struct {
	product *product.Product
	err     error
}

// NewProductLoader creates a new ProductLoader
func NewProductLoader(repo product.Repository) *ProductLoader {
	return &ProductLoader{
		repo:     repo,
		cache:    make(map[string]*product.Product),
		batch:    make(map[string][]chan productResult),
		maxBatch: 100,
		wait:     1 * time.Millisecond,
	}
}

// Load loads a single product by ID with batching
func (pl *ProductLoader) Load(ctx context.Context, id string) (*product.Product, error) {
	// Check cache first
	pl.mutex.RLock()
	if cached, exists := pl.cache[id]; exists {
		pl.mutex.RUnlock()
		return cached, nil
	}
	pl.mutex.RUnlock()

	// Create result channel
	result := make(chan productResult, 1)

	// Add to batch
	pl.batchMutex.Lock()
	pl.batch[id] = append(pl.batch[id], result)
	batch := pl.batch[id]

	// If this is the first request for this ID, start the batch timer
	if len(batch) == 1 {
		go pl.startBatch(id)
	}
	pl.batchMutex.Unlock()

	// Wait for result
	res := <-result
	return res.product, res.err
}

// LoadMany loads multiple products by IDs
func (pl *ProductLoader) LoadMany(ctx context.Context, ids []string) ([]*product.Product, []error) {
	products := make([]*product.Product, len(ids))
	errors := make([]error, len(ids))

	for i, id := range ids {
		products[i], errors[i] = pl.Load(ctx, id)
	}

	return products, errors
}

// Prime adds a product to the cache
func (pl *ProductLoader) Prime(ctx context.Context, id string, product *product.Product) {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	pl.cache[id] = product
}

// Clear removes a product from the cache
func (pl *ProductLoader) Clear(ctx context.Context, id string) {
	pl.mutex.Lock()
	defer pl.mutex.Unlock()
	delete(pl.cache, id)
}

// startBatch starts a batch operation for the given ID
func (pl *ProductLoader) startBatch(id string) {
	time.Sleep(pl.wait)

	pl.batchMutex.Lock()
	batch := pl.batch[id]
	delete(pl.batch, id)
	pl.batchMutex.Unlock()

	// Fetch the product
	product, err := pl.repo.GetByID(context.Background(), id)

	// Cache successful results
	if err == nil {
		pl.Prime(context.Background(), id, product)
	}

	// Send results to all waiters
	result := productResult{product: product, err: err}
	for _, ch := range batch {
		ch <- result
		close(ch)
	}
}

// batchGetProducts fetches multiple products in a single operation
func batchGetProducts(ctx context.Context, repo product.Repository, ids []string) ([]*product.Product, []error) {
	// Map to store results
	productMap := make(map[string]*product.Product)
	errors := make([]error, len(ids))

	// Fetch each product
	for _, id := range ids {
		product, err := repo.GetByID(ctx, id)
		if err != nil {
			productMap[id] = nil
		} else {
			productMap[id] = product
		}
	}

	// Build result slice in the same order as input IDs
	results := make([]*product.Product, len(ids))
	for i, id := range ids {
		results[i] = productMap[id]
		if results[i] == nil && errors[i] == nil {
			errors[i] = fmt.Errorf("product not found: %s", id)
		}
	}

	return results, errors
}
