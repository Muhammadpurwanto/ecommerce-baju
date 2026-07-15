package handler

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/cart/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/cart/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/cart/internal/service"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/cart"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartGrpcHandler struct {
	pb.UnimplementedCartServiceServer
	cartService service.CartService
	validate    *validator.Validate
}

func NewCartHandler(cartService service.CartService) *CartGrpcHandler {
	return &CartGrpcHandler{
		cartService: cartService,
		validate:    validator.New(),
	}
}

func toPbCartResponse(c *model.Cart) *pb.CartResponse {
	var pbItems []*pb.CartItem
	for _, item := range c.Items {
		pbItems = append(pbItems, &pb.CartItem{
			Id:        uint32(item.ID),
			CartId:    uint32(item.CartID),
			ProductId: uint32(item.ProductID),
			Quantity:  int32(item.Quantity),
			CreatedAt: item.CreatedAt.Format(time.RFC3339),
			UpdatedAt: item.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &pb.CartResponse{
		Id:        uint32(c.ID),
		UserId:    c.UserID,
		CreatedAt: c.CreatedAt.Format(time.RFC3339),
		UpdatedAt: c.UpdatedAt.Format(time.RFC3339),
		Items:     pbItems,
	}
}

func (h *CartGrpcHandler) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.CartResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	cart, err := h.cartService.GetCart(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get cart: %v", err)
	}

	return toPbCartResponse(cart), nil
}

func (h *CartGrpcHandler) AddItem(ctx context.Context, req *pb.AddItemRequest) (*pb.CartResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	dtoReq := &dto.AddItemRequest{
		ProductID: uint(req.GetProductId()),
		Quantity:  int(req.GetQuantity()),
	}

	// Validasi input: Pastikan Quantity positif dan ID ada
	if dtoReq.ProductID == 0 {
		return nil, status.Error(codes.InvalidArgument, "product_id is required")
	}
	if dtoReq.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than zero")
	}

	cart, err := h.cartService.AddItem(req.GetUserId(), dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to add item: %v", err)
	}

	return toPbCartResponse(cart), nil
}

func (h *CartGrpcHandler) UpdateItem(ctx context.Context, req *pb.UpdateItemRequest) (*pb.CartResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	dtoReq := &dto.UpdateItemRequest{
		Quantity: int(req.GetQuantity()),
	}

	// Validasi input: Pastikan Quantity positif
	if dtoReq.Quantity <= 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity must be greater than zero")
	}

	cart, err := h.cartService.UpdateItem(req.GetUserId(), uint(req.GetItemId()), dtoReq)
	if err != nil {
		if err.Error() == "item not found in cart" {
			return nil, status.Errorf(codes.NotFound, "item not found in cart")
		}
		return nil, status.Errorf(codes.Internal, "failed to update item: %v", err)
	}

	return toPbCartResponse(cart), nil
}

func (h *CartGrpcHandler) RemoveItem(ctx context.Context, req *pb.RemoveItemRequest) (*pb.CartResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetItemId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "item_id is required")
	}

	cart, err := h.cartService.RemoveItem(req.GetUserId(), uint(req.GetItemId()))
	if err != nil {
		if err.Error() == "item not found in cart" {
			return nil, status.Errorf(codes.NotFound, "item not found in cart")
		}
		return nil, status.Errorf(codes.Internal, "failed to remove item: %v", err)
	}

	return toPbCartResponse(cart), nil
}

func (h *CartGrpcHandler) ClearCart(ctx context.Context, req *pb.ClearCartRequest) (*pb.CartEmptyResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	err := h.cartService.ClearCart(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to clear cart: %v", err)
	}

	return &pb.CartEmptyResponse{Success: true}, nil
}
