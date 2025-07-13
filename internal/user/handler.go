package user

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "learning/pkg/user/pb"
)

// Handler implements the UserService gRPC server
type Handler struct {
	pb.UnimplementedUserServiceServer
	service *Service
}

// NewHandler creates a new gRPC handler for user service
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateUser creates a new user
func (h *Handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("CreateUser request: %+v", req)

	user, err := h.service.CreateUser(ctx, req.Name, req.Email, req.Phone)
	if err != nil {
		log.Printf("CreateUser error: %v", err)

		// Handle validation errors
		if validationErr, ok := err.(*ValidationError); ok {
			return nil, status.Error(codes.InvalidArgument, validationErr.Message)
		}

		// Handle already exists error
		if err == ErrUserAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}

		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.CreateUserResponse{
		User: h.userToProto(user),
	}, nil
}

// GetUser retrieves a user by ID
func (h *Handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	log.Printf("GetUser request: %+v", req)

	user, err := h.service.GetUser(ctx, req.Id)
	if err != nil {
		log.Printf("GetUser error: %v", err)

		if err == ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to get user")
	}

	return &pb.GetUserResponse{
		User: h.userToProto(user),
	}, nil
}

// UpdateUser updates an existing user
func (h *Handler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Printf("UpdateUser request: %+v", req)

	user, err := h.service.UpdateUser(ctx, req.Id, req.Name, req.Email, req.Phone)
	if err != nil {
		log.Printf("UpdateUser error: %v", err)

		// Handle validation errors
		if validationErr, ok := err.(*ValidationError); ok {
			return nil, status.Error(codes.InvalidArgument, validationErr.Message)
		}

		// Handle not found error
		if err == ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		// Handle already exists error
		if err == ErrUserAlreadyExists {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}

		return nil, status.Error(codes.Internal, "failed to update user")
	}

	return &pb.UpdateUserResponse{
		User: h.userToProto(user),
	}, nil
}

// DeleteUser deletes a user by ID
func (h *Handler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	log.Printf("DeleteUser request: %+v", req)

	err := h.service.DeleteUser(ctx, req.Id)
	if err != nil {
		log.Printf("DeleteUser error: %v", err)

		if err == ErrUserNotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to delete user")
	}

	return &pb.DeleteUserResponse{
		Success: true,
	}, nil
}

// ListUsers retrieves users with pagination
func (h *Handler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	log.Printf("ListUsers request: %+v", req)

	users, total, err := h.service.ListUsers(ctx, int(req.Page), int(req.PageSize))
	if err != nil {
		log.Printf("ListUsers error: %v", err)
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	protoUsers := make([]*pb.User, len(users))
	for i, user := range users {
		protoUsers[i] = h.userToProto(user)
	}

	return &pb.ListUsersResponse{
		Users:    protoUsers,
		Total:    int32(total),
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// userToProto converts domain user to protobuf user
func (h *Handler) userToProto(user *User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
