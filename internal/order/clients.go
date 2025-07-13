package order

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	productpb "learning/pkg/product/pb"
	userpb "learning/pkg/user/pb"
)

// UserServiceClient wraps the user service gRPC client
type UserServiceClient struct {
	client userpb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserServiceClient creates a new user service client
func NewUserServiceClient(address string) (*UserServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	client := userpb.NewUserServiceClient(conn)

	return &UserServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

// GetUser retrieves a user by ID
func (c *UserServiceClient) GetUser(ctx context.Context, userID string) (*userpb.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := c.client.GetUser(ctx, &userpb.GetUserRequest{
		Id: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return resp.User, nil
}

// Close closes the connection
func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}

// ProductServiceClient wraps the product service gRPC client
type ProductServiceClient struct {
	client productpb.ProductServiceClient
	conn   *grpc.ClientConn
}

// NewProductServiceClient creates a new product service client
func NewProductServiceClient(address string) (*ProductServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	client := productpb.NewProductServiceClient(conn)

	return &ProductServiceClient{
		client: client,
		conn:   conn,
	}, nil
}

// GetProduct retrieves a product by ID
func (c *ProductServiceClient) GetProduct(ctx context.Context, productID string) (*productpb.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := c.client.GetProduct(ctx, &productpb.GetProductRequest{
		Id: productID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return resp.Product, nil
}

// UpdateStock updates product stock
func (c *ProductServiceClient) UpdateStock(ctx context.Context, productID string, quantity int32) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := c.client.UpdateStock(ctx, &productpb.UpdateStockRequest{
		ProductId: productID,
		Quantity:  quantity, // negative to reduce stock
	})
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	log.Printf("Updated stock for product %s by %d", productID, quantity)
	return nil
}

// Close closes the connection
func (c *ProductServiceClient) Close() error {
	return c.conn.Close()
}
