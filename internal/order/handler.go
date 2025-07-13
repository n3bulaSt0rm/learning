package order

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "learning/pkg/order/pb"
)

// Handler implements the OrderService gRPC server
type Handler struct {
	pb.UnimplementedOrderServiceServer
	service *Service
}

// NewHandler creates a new gRPC handler for order service
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateOrder creates a new order
func (h *Handler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.CreateOrderResponse, error) {
	log.Printf("CreateOrder request: %+v", req)

	// Convert proto items to domain items
	var items []*OrderItemRequest
	for _, item := range req.Items {
		items = append(items, &OrderItemRequest{
			ProductID: item.ProductId,
			Quantity:  item.Quantity,
		})
	}

	order, err := h.service.CreateOrder(ctx, req.UserId, items)
	if err != nil {
		log.Printf("CreateOrder error: %v", err)

		// Handle validation errors
		if validationErr, ok := err.(*ValidationError); ok {
			return nil, status.Error(codes.InvalidArgument, validationErr.Message)
		}

		return nil, status.Error(codes.Internal, "failed to create order")
	}

	return &pb.CreateOrderResponse{
		Order: h.orderToProto(order),
	}, nil
}

// GetOrder retrieves an order by ID
func (h *Handler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.GetOrderResponse, error) {
	log.Printf("GetOrder request: %+v", req)

	order, err := h.service.GetOrder(ctx, req.Id)
	if err != nil {
		log.Printf("GetOrder error: %v", err)

		if err == ErrOrderNotFound {
			return nil, status.Error(codes.NotFound, "order not found")
		}

		return nil, status.Error(codes.Internal, "failed to get order")
	}

	return &pb.GetOrderResponse{
		Order: h.orderToProto(order),
	}, nil
}

// UpdateOrderStatus updates the status of an order
func (h *Handler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.UpdateOrderStatusResponse, error) {
	log.Printf("UpdateOrderStatus request: %+v", req)

	order, err := h.service.UpdateOrderStatus(ctx, req.Id, OrderStatus(req.Status))
	if err != nil {
		log.Printf("UpdateOrderStatus error: %v", err)

		// Handle validation errors
		if validationErr, ok := err.(*ValidationError); ok {
			return nil, status.Error(codes.InvalidArgument, validationErr.Message)
		}

		if err == ErrOrderNotFound {
			return nil, status.Error(codes.NotFound, "order not found")
		}

		return nil, status.Error(codes.Internal, "failed to update order status")
	}

	return &pb.UpdateOrderStatusResponse{
		Order: h.orderToProto(order),
	}, nil
}

// ListOrdersByUser retrieves orders for a specific user
func (h *Handler) ListOrdersByUser(ctx context.Context, req *pb.ListOrdersByUserRequest) (*pb.ListOrdersByUserResponse, error) {
	log.Printf("ListOrdersByUser request: %+v", req)

	orders, total, err := h.service.ListOrdersByUser(ctx, req.UserId, int(req.Page), int(req.PageSize))
	if err != nil {
		log.Printf("ListOrdersByUser error: %v", err)

		// Handle validation errors
		if validationErr, ok := err.(*ValidationError); ok {
			return nil, status.Error(codes.InvalidArgument, validationErr.Message)
		}

		return nil, status.Error(codes.Internal, "failed to list orders")
	}

	protoOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		protoOrders[i] = h.orderToProto(order)
	}

	return &pb.ListOrdersByUserResponse{
		Orders:   protoOrders,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// ListOrders retrieves all orders with pagination and optional status filter
func (h *Handler) ListOrders(ctx context.Context, req *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	log.Printf("ListOrders request: %+v", req)

	orders, total, err := h.service.ListOrders(ctx, int(req.Page), int(req.PageSize), OrderStatus(req.Status))
	if err != nil {
		log.Printf("ListOrders error: %v", err)
		return nil, status.Error(codes.Internal, "failed to list orders")
	}

	protoOrders := make([]*pb.Order, len(orders))
	for i, order := range orders {
		protoOrders[i] = h.orderToProto(order)
	}

	return &pb.ListOrdersResponse{
		Orders:   protoOrders,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// orderToProto converts domain order to protobuf order
func (h *Handler) orderToProto(order *Order) *pb.Order {
	items := make([]*pb.OrderItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = &pb.OrderItem{
			Id:           item.ID,
			ProductId:    item.ProductID,
			ProductName:  item.ProductName,
			ProductPrice: item.ProductPrice,
			Quantity:     item.Quantity,
			Total:        item.Total,
		}
	}

	return &pb.Order{
		Id:          order.ID,
		UserId:      order.UserID,
		Items:       items,
		TotalAmount: order.TotalAmount,
		Status:      pb.OrderStatus(order.Status),
		CreatedAt:   timestamppb.New(order.CreatedAt),
		UpdatedAt:   timestamppb.New(order.UpdatedAt),
	}
}
