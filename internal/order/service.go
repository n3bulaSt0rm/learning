package order

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// OrderItemRequest represents a request to add an item to an order
type OrderItemRequest struct {
	ProductID string
	Quantity  int32
}

// Service handles business logic for order operations
type Service struct {
	repo          Repository
	userClient    *UserServiceClient
	productClient *ProductServiceClient
}

// NewService creates a new order service
func NewService(repo Repository, userClient *UserServiceClient, productClient *ProductServiceClient) *Service {
	return &Service{
		repo:          repo,
		userClient:    userClient,
		productClient: productClient,
	}
}

// CreateOrder creates a new order with validation and stock management
func (s *Service) CreateOrder(ctx context.Context, userID string, items []*OrderItemRequest) (*Order, error) {
	log.Printf("Creating order for user %s with %d items", userID, len(items))

	// Validate input
	if userID == "" {
		return nil, NewValidationError("user ID is required")
	}

	if len(items) == 0 {
		return nil, NewValidationError("at least one item is required")
	}

	// Verify user exists
	_, err := s.userClient.GetUser(ctx, userID)
	if err != nil {
		log.Printf("Failed to verify user %s: %v", userID, err)
		if status.Code(err) == codes.NotFound {
			return nil, NewValidationError("user not found")
		}
		return nil, fmt.Errorf("failed to verify user: %w", err)
	}

	// Process order items
	var orderItems []*OrderItem
	var totalAmount float64

	for _, itemReq := range items {
		// Validate item
		if itemReq.ProductID == "" {
			return nil, NewValidationError("product ID is required for all items")
		}
		if itemReq.Quantity <= 0 {
			return nil, NewValidationError("quantity must be greater than 0")
		}

		// Get product details
		product, err := s.productClient.GetProduct(ctx, itemReq.ProductID)
		if err != nil {
			log.Printf("Failed to get product %s: %v", itemReq.ProductID, err)
			if status.Code(err) == codes.NotFound {
				return nil, NewValidationError(fmt.Sprintf("product %s not found", itemReq.ProductID))
			}
			return nil, fmt.Errorf("failed to get product: %w", err)
		}

		// Check stock availability
		if product.Stock < itemReq.Quantity {
			return nil, NewValidationError(fmt.Sprintf("insufficient stock for product %s", product.Name))
		}

		// Calculate item total
		itemTotal := product.Price * float64(itemReq.Quantity)
		totalAmount += itemTotal

		// Create order item
		orderItem := &OrderItem{
			ProductID:    product.Id,
			ProductName:  product.Name,
			ProductPrice: product.Price,
			Quantity:     itemReq.Quantity,
			Total:        itemTotal,
		}

		orderItems = append(orderItems, orderItem)
	}

	// Create order
	order := &Order{
		UserID:      userID,
		Items:       orderItems,
		TotalAmount: totalAmount,
		Status:      OrderStatusPending,
	}

	// Save order
	savedOrder, err := s.repo.Create(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	// Update stock for all items (reduce stock)
	for _, item := range orderItems {
		err := s.productClient.UpdateStock(ctx, item.ProductID, -item.Quantity)
		if err != nil {
			log.Printf("Failed to update stock for product %s: %v", item.ProductID, err)
			// In a real system, we would need compensation logic here
			// For simplicity, we'll just log the error
		}
	}

	log.Printf("Order %s created successfully", savedOrder.ID)
	return savedOrder, nil
}

// GetOrder retrieves an order by ID
func (s *Service) GetOrder(ctx context.Context, id string) (*Order, error) {
	if id == "" {
		return nil, ErrOrderNotFound
	}
	return s.repo.GetByID(ctx, id)
}

// UpdateOrderStatus updates the status of an order
func (s *Service) UpdateOrderStatus(ctx context.Context, id string, status OrderStatus) (*Order, error) {
	if id == "" {
		return nil, ErrOrderNotFound
	}

	// Validate status transition (simplified)
	if status == OrderStatusUnspecified {
		return nil, NewValidationError("invalid order status")
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

// ListOrdersByUser retrieves orders for a specific user with pagination
func (s *Service) ListOrdersByUser(ctx context.Context, userID string, page, pageSize int) ([]*Order, int, error) {
	if userID == "" {
		return nil, 0, NewValidationError("user ID is required")
	}

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
	return s.repo.ListByUser(ctx, userID, offset, pageSize)
}

// ListOrders retrieves orders with pagination and optional status filter
func (s *Service) ListOrders(ctx context.Context, page, pageSize int, status OrderStatus) ([]*Order, int, error) {
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
	return s.repo.List(ctx, offset, pageSize, status)
}

// Close closes external service connections
func (s *Service) Close() error {
	var lastErr error

	if err := s.userClient.Close(); err != nil {
		log.Printf("Failed to close user client: %v", err)
		lastErr = err
	}

	if err := s.productClient.Close(); err != nil {
		log.Printf("Failed to close product client: %v", err)
		lastErr = err
	}

	return lastErr
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
