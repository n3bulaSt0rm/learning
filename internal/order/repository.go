package order

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
)

// OrderStatus represents the status of an order
type OrderStatus int32

const (
	OrderStatusUnspecified OrderStatus = 0
	OrderStatusPending     OrderStatus = 1
	OrderStatusConfirmed   OrderStatus = 2
	OrderStatusProcessing  OrderStatus = 3
	OrderStatusShipped     OrderStatus = 4
	OrderStatusDelivered   OrderStatus = 5
	OrderStatusCancelled   OrderStatus = 6
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID           string
	ProductID    string
	ProductName  string
	ProductPrice float64
	Quantity     int32
	Total        float64
}

// Order domain model
type Order struct {
	ID          string
	UserID      string
	Items       []*OrderItem
	TotalAmount float64
	Status      OrderStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Repository interface for order operations
type Repository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	GetByID(ctx context.Context, id string) (*Order, error)
	UpdateStatus(ctx context.Context, id string, status OrderStatus) (*Order, error)
	ListByUser(ctx context.Context, userID string, offset, limit int) ([]*Order, int, error)
	List(ctx context.Context, offset, limit int, status OrderStatus) ([]*Order, int, error)
}

// InMemoryRepository implements Repository interface using in-memory storage
type InMemoryRepository struct {
	orders map[string]*Order
	mutex  sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		orders: make(map[string]*Order),
	}
}

// Create creates a new order
func (r *InMemoryRepository) Create(ctx context.Context, order *Order) (*Order, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Generate ID if not provided
	if order.ID == "" {
		order.ID = uuid.New().String()
	}

	// Generate IDs for order items
	for _, item := range order.Items {
		if item.ID == "" {
			item.ID = uuid.New().String()
		}
	}

	// Check if order already exists
	if _, exists := r.orders[order.ID]; exists {
		return nil, ErrOrderAlreadyExists
	}

	// Set timestamps
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	// Store order
	r.orders[order.ID] = order

	return order, nil
}

// GetByID retrieves an order by ID
func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

// UpdateStatus updates the status of an order
func (r *InMemoryRepository) UpdateStatus(ctx context.Context, id string, status OrderStatus) (*Order, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	order, exists := r.orders[id]
	if !exists {
		return nil, ErrOrderNotFound
	}

	// Create a copy and update status
	updatedOrder := *order
	updatedOrder.Status = status
	updatedOrder.UpdatedAt = time.Now()

	// Store updated order
	r.orders[id] = &updatedOrder

	return &updatedOrder, nil
}

// ListByUser retrieves orders for a specific user with pagination
func (r *InMemoryRepository) ListByUser(ctx context.Context, userID string, offset, limit int) ([]*Order, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Filter orders by user ID
	var userOrders []*Order
	for _, order := range r.orders {
		if order.UserID == userID {
			userOrders = append(userOrders, order)
		}
	}

	total := len(userOrders)

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
		return []*Order{}, total, nil
	}

	return userOrders[start:end], total, nil
}

// List retrieves orders with pagination and optional status filter
func (r *InMemoryRepository) List(ctx context.Context, offset, limit int, status OrderStatus) ([]*Order, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Convert map to slice and filter by status if specified
	var orders []*Order
	for _, order := range r.orders {
		if status == OrderStatusUnspecified || order.Status == status {
			orders = append(orders, order)
		}
	}

	total := len(orders)

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
		return []*Order{}, total, nil
	}

	return orders[start:end], total, nil
}
