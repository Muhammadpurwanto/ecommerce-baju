package handler

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/service"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/payment"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentGrpcHandler struct {
	pb.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

func NewPaymentHandler(svc service.PaymentService) *PaymentGrpcHandler {
	return &PaymentGrpcHandler{paymentService: svc}
}

func toPbPayment(p *model.Payment) *pb.Payment {
	paidAt := ""
	if p.PaidAt != nil {
		paidAt = p.PaidAt.Format(time.RFC3339)
	}

	return &pb.Payment{
		Id:            uint32(p.ID),
		OrderId:       uint32(p.OrderID),
		UserId:        p.UserID,
		Amount:        p.Amount,
		PaymentType:   p.PaymentType,
		TransactionId: p.TransactionID,
		Status:        p.Status,
		SnapToken:     p.SnapToken,
		PaymentUrl:    p.PaymentURL,
		PaidAt:        paidAt,
		CreatedAt:     p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     p.UpdatedAt.Format(time.RFC3339),
	}
}

func (h *PaymentGrpcHandler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.Payment, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	if req.GetOrderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	// Validasi keamanan fatal: Pastikan Amount lebih dari 0
	if req.GetAmount() <= 0 {
		return nil, status.Error(codes.InvalidArgument, "payment amount must be greater than zero")
	}

	dtoReq := &dto.CreatePaymentRequest{
		OrderID: uint(req.GetOrderId()),
		Amount:  req.GetAmount(),
	}

	payment, err := h.paymentService.CreatePayment(req.GetUserId(), dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment: %v", err)
	}

	return toPbPayment(payment), nil
}

func (h *PaymentGrpcHandler) GetPaymentByOrderID(ctx context.Context, req *pb.GetPaymentByOrderIDRequest) (*pb.Payment, error) {
	if req.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}
	
	if req.GetOrderId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	payment, err := h.paymentService.GetPaymentByOrderID(uint(req.GetOrderId()), req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "payment not found: %v", err)
	}

	return toPbPayment(payment), nil
}

func (h *PaymentGrpcHandler) WebhookMidtrans(ctx context.Context, req *pb.WebhookMidtransRequest) (*pb.WebhookMidtransResponse, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(req.GetPayloadJson()), &payload); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid JSON payload: %v", err)
	}

	if err := h.paymentService.HandleMidtransWebhook(payload); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to handle webhook: %v", err)
	}

	return &pb.WebhookMidtransResponse{
		Status:  "success",
		Message: "webhook handled",
	}, nil
}
