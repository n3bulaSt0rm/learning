package user

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// User domain model
type User struct {
	ID        string
	Name      string
	Email     string
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repository interface for user operations
type Repository interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*User, int, error)
}

// InMemoryRepository implements Repository interface using in-memory storage
type InMemoryRepository struct {
	users map[string]*User
	mutex sync.RWMutex
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		users: make(map[string]*User),
	}
}

// Create creates a new user
func (r *InMemoryRepository) Create(ctx context.Context, user *User) (*User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// Check if email already exists
	for _, existingUser := range r.users {
		if existingUser.Email == user.Email {
			return nil, ErrUserAlreadyExists
		}
	}

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Store user
	r.users[user.ID] = user

	return user, nil
}

// GetByID retrieves a user by ID
func (r *InMemoryRepository) GetByID(ctx context.Context, id string) (*User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *InMemoryRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, ErrUserNotFound
}

// Update updates an existing user
func (r *InMemoryRepository) Update(ctx context.Context, user *User) (*User, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	existingUser, exists := r.users[user.ID]
	if !exists {
		return nil, ErrUserNotFound
	}

	// Check if email is being changed and already exists
	if user.Email != existingUser.Email {
		for _, otherUser := range r.users {
			if otherUser.ID != user.ID && otherUser.Email == user.Email {
				return nil, ErrUserAlreadyExists
			}
		}
	}

	// Update fields
	user.CreatedAt = existingUser.CreatedAt
	user.UpdatedAt = time.Now()

	// Store updated user
	r.users[user.ID] = user

	return user, nil
}

// Delete deletes a user by ID
func (r *InMemoryRepository) Delete(ctx context.Context, id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[id]; !exists {
		return ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

// List retrieves users with pagination
func (r *InMemoryRepository) List(ctx context.Context, offset, limit int) ([]*User, int, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// Convert map to slice
	var users []*User
	for _, user := range r.users {
		users = append(users, user)
	}

	total := len(users)

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
		return []*User{}, total, nil
	}

	return users[start:end], total, nil
}
