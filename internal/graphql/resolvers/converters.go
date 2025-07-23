package resolvers

import (
	"strconv"

	"learning/internal/graphql/models"
	"learning/internal/order"
	"learning/internal/product"
	"learning/internal/user"
)

// User converters
func domainUserToGraphQL(u *user.User) *models.User {
	if u == nil {
		return nil
	}

	return &models.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Phone:     u.Phone,
		CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: u.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func domainUsersToGraphQL(users []*user.User) []*models.User {
	result := make([]*models.User, len(users))
	for i, u := range users {
		result[i] = domainUserToGraphQL(u)
	}
	return result
}

// Product converters
func domainProductToGraphQL(p *product.Product) *models.Product {
	if p == nil {
		return nil
	}

	return &models.Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       int(p.Stock),
		Category:    p.Category,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		IsInStock:   p.Stock > 0,
		StockStatus: getStockStatus(p.Stock),
	}
}

func domainProductsToGraphQL(products []*product.Product) []*models.Product {
	result := make([]*models.Product, len(products))
	for i, p := range products {
		result[i] = domainProductToGraphQL(p)
	}
	return result
}

func getStockStatus(stock int32) models.StockStatus {
	switch {
	case stock == 0:
		return models.StockStatusOutOfStock
	case stock <= 10:
		return models.StockStatusLowStock
	default:
		return models.StockStatusInStock
	}
}

// Order converters
func domainOrderToGraphQL(o *order.Order) *models.Order {
	if o == nil {
		return nil
	}

	return &models.Order{
		ID:             o.ID,
		TotalAmount:    o.TotalAmount,
		Status:         domainOrderStatusToGraphQL(o.Status),
		CreatedAt:      o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:      o.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		ItemCount:      len(o.Items),
		CanBeCancelled: o.Status == order.OrderStatusPending || o.Status == order.OrderStatusConfirmed,
	}
}

func domainOrdersToGraphQL(orders []*order.Order) []*models.Order {
	result := make([]*models.Order, len(orders))
	for i, o := range orders {
		result[i] = domainOrderToGraphQL(o)
	}
	return result
}

func domainOrderItemToGraphQL(item *order.OrderItem) *models.OrderItem {
	if item == nil {
		return nil
	}

	return &models.OrderItem{
		ID:           item.ID,
		ProductName:  item.ProductName,
		ProductPrice: item.ProductPrice,
		Quantity:     int(item.Quantity),
		Total:        item.Total,
	}
}

func domainOrderItemsToGraphQL(items []*order.OrderItem) []*models.OrderItem {
	result := make([]*models.OrderItem, len(items))
	for i, item := range items {
		result[i] = domainOrderItemToGraphQL(item)
	}
	return result
}

func domainOrderStatusToGraphQL(status order.OrderStatus) models.OrderStatus {
	switch status {
	case order.OrderStatusPending:
		return models.OrderStatusPending
	case order.OrderStatusConfirmed:
		return models.OrderStatusConfirmed
	case order.OrderStatusProcessing:
		return models.OrderStatusProcessing
	case order.OrderStatusShipped:
		return models.OrderStatusShipped
	case order.OrderStatusDelivered:
		return models.OrderStatusDelivered
	case order.OrderStatusCancelled:
		return models.OrderStatusCancelled
	default:
		return models.OrderStatusPending
	}
}

func graphQLOrderStatusToDomain(status models.OrderStatus) order.OrderStatus {
	switch status {
	case models.OrderStatusPending:
		return order.OrderStatusPending
	case models.OrderStatusConfirmed:
		return order.OrderStatusConfirmed
	case models.OrderStatusProcessing:
		return order.OrderStatusProcessing
	case models.OrderStatusShipped:
		return order.OrderStatusShipped
	case models.OrderStatusDelivered:
		return order.OrderStatusDelivered
	case models.OrderStatusCancelled:
		return order.OrderStatusCancelled
	default:
		return order.OrderStatusPending
	}
}

// Error helpers
func createUserError(field string, message string, code models.ErrorCode) *models.UserError {
	return &models.UserError{
		Field:   &field,
		Message: message,
		Code:    code,
	}
}

func createProductError(field string, message string, code models.ErrorCode) *models.ProductError {
	return &models.ProductError{
		Field:   &field,
		Message: message,
		Code:    code,
	}
}

func createOrderError(field string, message string, code models.ErrorCode) *models.OrderError {
	return &models.OrderError{
		Field:   &field,
		Message: message,
		Code:    code,
	}
}

// Pagination helpers
func createUserConnection(users []*models.User, total int, page int, pageSize int) *models.UserConnection {
	edges := make([]*models.UserEdge, len(users))
	for i, user := range users {
		cursor := strconv.Itoa((page-1)*pageSize + i + 1)
		edges[i] = &models.UserEdge{
			Node:   user,
			Cursor: cursor,
		}
	}

	hasNextPage := (page * pageSize) < total
	hasPreviousPage := page > 1

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}

	return &models.UserConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
		TotalCount: total,
	}
}

func createProductConnection(products []*models.Product, total int, page int, pageSize int) *models.ProductConnection {
	edges := make([]*models.ProductEdge, len(products))
	for i, product := range products {
		cursor := strconv.Itoa((page-1)*pageSize + i + 1)
		edges[i] = &models.ProductEdge{
			Node:   product,
			Cursor: cursor,
		}
	}

	hasNextPage := (page * pageSize) < total
	hasPreviousPage := page > 1

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}

	return &models.ProductConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
		TotalCount: total,
	}
}

func createOrderConnection(orders []*models.Order, total int, page int, pageSize int) *models.OrderConnection {
	edges := make([]*models.OrderEdge, len(orders))
	for i, order := range orders {
		cursor := strconv.Itoa((page-1)*pageSize + i + 1)
		edges[i] = &models.OrderEdge{
			Node:   order,
			Cursor: cursor,
		}
	}

	hasNextPage := (page * pageSize) < total
	hasPreviousPage := page > 1

	var startCursor, endCursor *string
	if len(edges) > 0 {
		startCursor = &edges[0].Cursor
		endCursor = &edges[len(edges)-1].Cursor
	}

	return &models.OrderConnection{
		Edges: edges,
		PageInfo: &models.PageInfo{
			HasNextPage:     hasNextPage,
			HasPreviousPage: hasPreviousPage,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
		TotalCount: total,
	}
}
