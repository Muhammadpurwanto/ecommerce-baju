package handler

import (
	"context"
	"strconv"
	"time"

	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/cart"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/util"
	"github.com/gofiber/fiber/v2"
)

type CartClientHandler struct {
	grpcClient pb.CartServiceClient
}

func NewCartClientHandler(grpcClient pb.CartServiceClient) *CartClientHandler {
	return &CartClientHandler{grpcClient: grpcClient}
}

func (h *CartClientHandler) GetCart(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.GetCart(ctx, &pb.GetCartRequest{UserId: userID})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *CartClientHandler) AddItem(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	var req pb.AddItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}
	req.UserId = userID

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.AddItem(ctx, &req)
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": res})
}

func (h *CartClientHandler) UpdateItem(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	itemID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid item ID"})
	}

	var req pb.UpdateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}
	req.UserId = userID
	req.ItemId = uint32(itemID)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.UpdateItem(ctx, &req)
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *CartClientHandler) RemoveItem(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	itemID, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid item ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.RemoveItem(ctx, &pb.RemoveItemRequest{
		UserId: userID,
		ItemId: uint32(itemID),
	})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *CartClientHandler) ClearCart(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, errGrpc := h.grpcClient.ClearCart(ctx, &pb.ClearCartRequest{UserId: userID})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "message": "Cart cleared"})
}
