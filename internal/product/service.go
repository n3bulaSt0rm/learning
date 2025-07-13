package product

import (
	"context"
	"strings"
)

// Service handles business logic for product operations
type Service struct {
	repo Repository
}

// NewService creates a new product service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateProduct creates a new product with validation
func (s *Service) CreateProduct(ctx context.Context, name, description, category string, price float64, stock int32) (*Product, error) {
	// Validate input
	if err := s.validateProduct(name, description, category, price, stock); err != nil {
		return nil, err
	}

	product := &Product{
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Category:    strings.TrimSpace(category),
		Price:       price,
		Stock:       stock,
	}

	return s.repo.Create(ctx, product)
}

// GetProduct retrieves a product by ID
func (s *Service) GetProduct(ctx context.Context, id string) (*Product, error) {
	if id == "" {
		return nil, ErrProductNotFound
	}
	return s.repo.GetByID(ctx, id)
}

// UpdateProduct updates an existing product
func (s *Service) UpdateProduct(ctx context.Context, id, name, description, category string, price float64, stock int32) (*Product, error) {
	if id == "" {
		return nil, ErrProductNotFound
	}

	// Validate input
	if err := s.validateProduct(name, description, category, price, stock); err != nil {
		return nil, err
	}

	product := &Product{
		ID:          id,
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Category:    strings.TrimSpace(category),
		Price:       price,
		Stock:       stock,
	}

	return s.repo.Update(ctx, product)
}

// DeleteProduct deletes a product by ID
func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return ErrProductNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ListProducts retrieves products with pagination and optional category filter
func (s *Service) ListProducts(ctx context.Context, page, pageSize int, category string) ([]*Product, int, error) {
	// Set default page size if not provided
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	// Set default page if not provided
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * pageSize
	return s.repo.List(ctx, offset, pageSize, category)
}

// UpdateStock updates product stock
func (s *Service) UpdateStock(ctx context.Context, productID string, quantity int32) (*Product, error) {
	if productID == "" {
		return nil, ErrProductNotFound
	}
	return s.repo.UpdateStock(ctx, productID, quantity)
}

// validateProduct validates product input
func (s *Service) validateProduct(name, description, category string, price float64, stock int32) error {
	if strings.TrimSpace(name) == "" {
		return NewValidationError("name is required")
	}

	if strings.TrimSpace(description) == "" {
		return NewValidationError("description is required")
	}

	if strings.TrimSpace(category) == "" {
		return NewValidationError("category is required")
	}

	if price <= 0 {
		return NewValidationError("price must be greater than 0")
	}

	if stock < 0 {
		return NewValidationError("stock cannot be negative")
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}
