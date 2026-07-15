package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/util"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/payment"
	"github.com/gofiber/fiber/v2"
)

type PaymentClientHandler struct {
	grpcClient pb.PaymentServiceClient
}

func NewPaymentClientHandler(grpcClient pb.PaymentServiceClient) *PaymentClientHandler {
	return &PaymentClientHandler{grpcClient: grpcClient}
}

func (h *PaymentClientHandler) CreatePayment(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	var req pb.CreatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}
	req.UserId = userID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.CreatePayment(ctx, &req)
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": res})
}

func (h *PaymentClientHandler) GetPaymentByOrderID(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	orderID, err := strconv.ParseUint(c.Params("orderId"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid order ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.GetPaymentByOrderID(ctx, &pb.GetPaymentByOrderIDRequest{
		OrderId: uint32(orderID),
		UserId:  userID,
	})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *PaymentClientHandler) WebhookMidtrans(c *fiber.Ctx) error {
	payloadBytes := c.Body()
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.WebhookMidtrans(ctx, &pb.WebhookMidtransRequest{
		PayloadJson: string(payloadBytes),
	})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}

	return c.JSON(fiber.Map{"status": res.GetStatus(), "message": res.GetMessage()})
}
