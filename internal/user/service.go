package user

import (
	"context"
	"strings"
)

// Service handles business logic for user operations
type Service struct {
	repo Repository
}

// NewService creates a new user service
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// CreateUser creates a new user with validation
func (s *Service) CreateUser(ctx context.Context, name, email, phone string) (*User, error) {
	// Validate input
	if err := s.validateUser(name, email, phone); err != nil {
		return nil, err
	}

	user := &User{
		Name:  strings.TrimSpace(name),
		Email: strings.ToLower(strings.TrimSpace(email)),
		Phone: strings.TrimSpace(phone),
	}

	return s.repo.Create(ctx, user)
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	if id == "" {
		return nil, ErrUserNotFound
	}
	return s.repo.GetByID(ctx, id)
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(ctx context.Context, id, name, email, phone string) (*User, error) {
	if id == "" {
		return nil, ErrUserNotFound
	}

	// Validate input
	if err := s.validateUser(name, email, phone); err != nil {
		return nil, err
	}

	user := &User{
		ID:    id,
		Name:  strings.TrimSpace(name),
		Email: strings.ToLower(strings.TrimSpace(email)),
		Phone: strings.TrimSpace(phone),
	}

	return s.repo.Update(ctx, user)
}

// DeleteUser deletes a user by ID
func (s *Service) DeleteUser(ctx context.Context, id string) error {
	if id == "" {
		return ErrUserNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ListUsers retrieves users with pagination
func (s *Service) ListUsers(ctx context.Context, page, pageSize int) ([]*User, int, error) {
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
	return s.repo.List(ctx, offset, pageSize)
}

// validateUser validates user input
func (s *Service) validateUser(name, email, phone string) error {
	if strings.TrimSpace(name) == "" {
		return NewValidationError("name is required")
	}

	if strings.TrimSpace(email) == "" {
		return NewValidationError("email is required")
	}

	if !isValidEmail(email) {
		return NewValidationError("invalid email format")
	}

	if strings.TrimSpace(phone) == "" {
		return NewValidationError("phone is required")
	}

	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}

	// Basic email validation
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}

	if !strings.Contains(parts[1], ".") {
		return false
	}

	return true
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
