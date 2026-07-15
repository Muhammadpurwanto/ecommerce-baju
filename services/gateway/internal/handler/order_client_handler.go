package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/util"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/order"
	"github.com/gofiber/fiber/v2"
)

type OrderClientHandler struct {
	grpcClient pb.OrderServiceClient
}

func NewOrderClientHandler(grpcClient pb.OrderServiceClient) *OrderClientHandler {
	return &OrderClientHandler{grpcClient: grpcClient}
}

func (h *OrderClientHandler) CreateOrder(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	var req pb.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}
	req.UserId = userID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.CreateOrder(ctx, &req)
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": res})
}

func (h *OrderClientHandler) GetOrderByID(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	isAdmin := util.IsAdmin(c)
	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid order ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.GetOrderByID(ctx, &pb.GetOrderByIDRequest{
		Id:      uint32(orderID),
		UserId:  userID,
		IsAdmin: isAdmin,
	})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *OrderClientHandler) GetMyOrders(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.GetMyOrders(ctx, &pb.GetOrdersRequest{
		UserId:  userID,
		Page:    int32(page),
		PerPage: int32(perPage),
	})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "data": res.Data, "meta": res.Meta})
}

func (h *OrderClientHandler) GetAllOrders(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Forbidden"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "50"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.GetAllOrders(ctx, &pb.GetOrdersRequest{
		Page:    int32(page),
		PerPage: int32(perPage),
	})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "data": res.Data, "meta": res.Meta})
}

func (h *OrderClientHandler) UpdateOrderStatus(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Forbidden"})
	}

	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid order ID"})
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.UpdateOrderStatus(ctx, &pb.UpdateOrderStatusRequest{
		Id:     uint32(orderID),
		Status: req.Status,
	})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *OrderClientHandler) CancelOrder(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	orderID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid order ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.CancelOrder(ctx, &pb.CancelOrderRequest{
		Id:     uint32(orderID),
		UserId: userID,
	})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}
