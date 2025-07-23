package resolvers

import (
	"learning/internal/graphql/dataloaders"
	"learning/internal/order"
	"learning/internal/product"
	"learning/internal/user"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserRepo    user.Repository
	ProductRepo product.Repository
	OrderRepo   order.Repository
	Loaders     *dataloaders.Loaders
}

// NewResolver creates a new resolver with dependencies
func NewResolver(
	userRepo user.Repository,
	productRepo product.Repository,
	orderRepo order.Repository,
	loaders *dataloaders.Loaders,
) *Resolver {
	return &Resolver{
		UserRepo:    userRepo,
		ProductRepo: productRepo,
		OrderRepo:   orderRepo,
		Loaders:     loaders,
	}
}
