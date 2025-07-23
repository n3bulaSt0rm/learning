package dataloaders

import (
	"context"
	"net/http"

	"learning/internal/order"
	"learning/internal/product"
	"learning/internal/user"
)

type contextKey string

const loadersKey = contextKey("dataloaders")

// Loaders holds all DataLoaders
type Loaders struct {
	UserLoader    *UserLoader
	ProductLoader *ProductLoader
	// OrderLoader could be added here if needed
}

// NewLoaders creates a new set of DataLoaders
func NewLoaders(
	userRepo user.Repository,
	productRepo product.Repository,
	orderRepo order.Repository,
) *Loaders {
	return &Loaders{
		UserLoader:    NewUserLoader(userRepo),
		ProductLoader: NewProductLoader(productRepo),
	}
}

// Middleware adds DataLoaders to the request context
func Middleware(loaders *Loaders) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), loadersKey, loaders)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// FromContext retrieves DataLoaders from context
func FromContext(ctx context.Context) *Loaders {
	loaders, ok := ctx.Value(loadersKey).(*Loaders)
	if !ok {
		panic("dataloaders not found in context")
	}
	return loaders
}

// GetUser is a convenience function to get a user via DataLoader
func GetUser(ctx context.Context, id string) (*user.User, error) {
	loaders := FromContext(ctx)
	return loaders.UserLoader.Load(ctx, id)
}

// GetProduct is a convenience function to get a product via DataLoader
func GetProduct(ctx context.Context, id string) (*product.Product, error) {
	loaders := FromContext(ctx)
	return loaders.ProductLoader.Load(ctx, id)
}

// GetUsers is a convenience function to get multiple users via DataLoader
func GetUsers(ctx context.Context, ids []string) ([]*user.User, []error) {
	loaders := FromContext(ctx)
	return loaders.UserLoader.LoadMany(ctx, ids)
}

// GetProducts is a convenience function to get multiple products via DataLoader
func GetProducts(ctx context.Context, ids []string) ([]*product.Product, []error) {
	loaders := FromContext(ctx)
	return loaders.ProductLoader.LoadMany(ctx, ids)
}
