package product

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "learning/pkg/product/pb"
)

// Handler implements the ProductService gRPC server
type Handler struct {
	pb.UnimplementedProductServiceServer
	service *Service
}

// NewHandler creates a new gRPC handler for product service
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateProduct creates a new product
func (h *Handler) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	log.Printf("CreateProduct request: %+v", req)

	product, err := h.service.CreateProduct(ctx, req.Name, req.Description, req.Category, req.Price, req.Stock)
	if err != nil {
		log.Printf("CreateProduct error: %v", err)

		// Handle validation errors
		if validationErr, ok := err.(*ValidationError); ok {
			return nil, status.Error(codes.InvalidArgument, validationErr.Message)
		}

		// Handle already exists error
		if err == ErrProductAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "product already exists")
		}

		return nil, status.Error(codes.Internal, "failed to create product")
	}

	return &pb.CreateProductResponse{
		Product: h.productToProto(product),
	}, nil
}

// GetProduct retrieves a product by ID
func (h *Handler) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	log.Printf("GetProduct request: %+v", req)

	product, err := h.service.GetProduct(ctx, req.Id)
	if err != nil {
		log.Printf("GetProduct error: %v", err)

		if err == ErrProductNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}

		return nil, status.Error(codes.Internal, "failed to get product")
	}

	return &pb.GetProductResponse{
		Product: h.productToProto(product),
	}, nil
}

// UpdateProduct updates an existing product
func (h *Handler) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	log.Printf("UpdateProduct request: %+v", req)

	product, err := h.service.UpdateProduct(ctx, req.Id, req.Name, req.Description, req.Category, req.Price, req.Stock)
	if err != nil {
		log.Printf("UpdateProduct error: %v", err)

		// Handle validation errors
		if validationErr, ok := err.(*ValidationError); ok {
			return nil, status.Error(codes.InvalidArgument, validationErr.Message)
		}

		// Handle not found error
		if err == ErrProductNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}

		return nil, status.Error(codes.Internal, "failed to update product")
	}

	return &pb.UpdateProductResponse{
		Product: h.productToProto(product),
	}, nil
}

// DeleteProduct deletes a product by ID
func (h *Handler) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	log.Printf("DeleteProduct request: %+v", req)

	err := h.service.DeleteProduct(ctx, req.Id)
	if err != nil {
		log.Printf("DeleteProduct error: %v", err)

		if err == ErrProductNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}

		return nil, status.Error(codes.Internal, "failed to delete product")
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

// ListProducts retrieves products with pagination and optional category filter
func (h *Handler) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	log.Printf("ListProducts request: %+v", req)

	products, total, err := h.service.ListProducts(ctx, int(req.Page), int(req.PageSize), req.Category)
	if err != nil {
		log.Printf("ListProducts error: %v", err)
		return nil, status.Error(codes.Internal, "failed to list products")
	}

	protoProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		protoProducts[i] = h.productToProto(product)
	}

	return &pb.ListProductsResponse{
		Products: protoProducts,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// UpdateStock updates product stock
func (h *Handler) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	log.Printf("UpdateStock request: %+v", req)

	product, err := h.service.UpdateStock(ctx, req.ProductId, req.Quantity)
	if err != nil {
		log.Printf("UpdateStock error: %v", err)

		if err == ErrProductNotFound {
			return nil, status.Error(codes.NotFound, "product not found")
		}

		if err == ErrInsufficientStock {
			return nil, status.Error(codes.InvalidArgument, "insufficient stock")
		}

		return nil, status.Error(codes.Internal, "failed to update stock")
	}

	return &pb.UpdateStockResponse{
		Success:  true,
		NewStock: product.Stock,
	}, nil
}

// productToProto converts domain product to protobuf product
func (h *Handler) productToProto(product *Product) *pb.Product {
	return &pb.Product{
		Id:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       product.Stock,
		Category:    product.Category,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
	}
}
