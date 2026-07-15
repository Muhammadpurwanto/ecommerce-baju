package handler

import (
	"context"
	"time"
	"math"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/service"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderGrpcHandler struct {
	pb.UnimplementedOrderServiceServer
	orderService service.OrderService
}

func NewOrderHandler(svc service.OrderService) *OrderGrpcHandler {
	return &OrderGrpcHandler{orderService: svc}
}

func toPbOrder(o *model.Order) *pb.Order {
	var items []*pb.OrderItem
	for _, item := range o.Items {
		items = append(items, &pb.OrderItem{
			Id:        uint32(item.ID),
			ProductId: uint32(item.ProductID),
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
			Subtotal:  item.Subtotal,
		})
	}

	paidAt := ""
	if o.PaidAt != nil {
		paidAt = o.PaidAt.Format(time.RFC3339)
	}

	shippedAt := ""
	if o.ShippedAt != nil {
		shippedAt = o.ShippedAt.Format(time.RFC3339)
	}

	return &pb.Order{
		Id:              uint32(o.ID),
		OrderNumber:     o.OrderNumber,
		UserId:          o.UserID,
		Status:          o.Status,
		TotalAmount:     o.TotalAmount,
		ShippingCost:    o.ShippingCost,
		ShippingAddress: o.ShippingAddress,
		Courier:         o.Courier,
		TrackingNumber:  o.TrackingNumber,
		Notes:           o.Notes,
		PaidAt:          paidAt,
		ShippedAt:       shippedAt,
		CreatedAt:       o.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       o.UpdatedAt.Format(time.RFC3339),
		Items:           items,
		SnapToken:       o.SnapToken,
		PaymentUrl:      o.PaymentURL,
	}
}

func (h *OrderGrpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if len(req.GetItems()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "order must contain at least one item")
	}

	if req.GetShippingCost() < 0 {
		return nil, status.Error(codes.InvalidArgument, "shipping cost cannot be negative")
	}

	var itemsReq []dto.OrderItemRequest
	for _, item := range req.GetItems() {
		// Validasi masing-masing item
		if item.GetQuantity() <= 0 {
			return nil, status.Errorf(codes.InvalidArgument, "quantity for product_id %d must be greater than zero", item.GetProductId())
		}
		if item.GetPrice() < 0 {
			return nil, status.Errorf(codes.InvalidArgument, "price for product_id %d cannot be negative", item.GetProductId())
		}
		
		itemsReq = append(itemsReq, dto.OrderItemRequest{
			ProductID: uint(item.GetProductId()),
			Quantity:  int(item.GetQuantity()),
			Price:     item.GetPrice(),
		})
	}

	dtoReq := &dto.CreateOrderRequest{
		ShippingCost:    req.GetShippingCost(),
		ShippingAddress: req.GetShippingAddress(),
		Courier:         req.GetCourier(),
		Notes:           req.GetNotes(),
		Items:           itemsReq,
	}

	order, err := h.orderService.CreateOrder(req.GetUserId(), dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	return toPbOrder(order), nil
}

func (h *OrderGrpcHandler) GetOrderByID(ctx context.Context, req *pb.GetOrderByIDRequest) (*pb.Order, error) {
	if req.GetUserId() == "" && !req.GetIsAdmin() {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	order, err := h.orderService.GetOrderByID(uint(req.GetId()), req.GetUserId(), req.GetIsAdmin())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	return toPbOrder(order), nil
}

func (h *OrderGrpcHandler) GetMyOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.OrderListResponse, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	page := int(req.GetPage())
	if page < 1 {
		page = 1
	}
	perPage := int(req.GetPerPage())
	if perPage < 1 {
		perPage = 10
	}

	orders, totalItems, err := h.orderService.GetUserOrders(req.GetUserId(), page, perPage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get orders: %v", err)
	}

	var pbData []*pb.Order
	for _, o := range orders {
		val := o
		pbData = append(pbData, toPbOrder(&val))
	}

	totalPages := int32(math.Ceil(float64(totalItems) / float64(perPage)))

	return &pb.OrderListResponse{
		Data: pbData,
		Meta: &pb.Meta{
			CurrentPage: int32(page),
			PerPage:     int32(perPage),
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

func (h *OrderGrpcHandler) GetAllOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.OrderListResponse, error) {
	page := int(req.GetPage())
	if page < 1 {
		page = 1
	}
	perPage := int(req.GetPerPage())
	if perPage < 1 {
		perPage = 50
	}

	orders, totalItems, err := h.orderService.GetAllOrders(page, perPage)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all orders: %v", err)
	}

	var pbData []*pb.Order
	for _, o := range orders {
		val := o
		pbData = append(pbData, toPbOrder(&val))
	}

	totalPages := int32(math.Ceil(float64(totalItems) / float64(perPage)))

	return &pb.OrderListResponse{
		Data: pbData,
		Meta: &pb.Meta{
			CurrentPage: int32(page),
			PerPage:     int32(perPage),
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

func (h *OrderGrpcHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}
	if req.GetStatus() == "" {
		return nil, status.Error(codes.InvalidArgument, "status is required")
	}

	order, err := h.orderService.UpdateOrderStatus(uint(req.GetId()), req.GetStatus())
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	return toPbOrder(order), nil
}

func (h *OrderGrpcHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderRequest) (*pb.Order, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	if req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	order, err := h.orderService.CancelOrder(uint(req.GetId()), req.GetUserId())
	if err != nil {
		if err.Error() == "record not found" {
			return nil, status.Errorf(codes.NotFound, "order not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to cancel order: %v", err)
	}

	return toPbOrder(order), nil
}
