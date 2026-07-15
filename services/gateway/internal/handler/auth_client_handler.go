package handler

import (
	"context"
	"time"

	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/auth"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/util"
	"github.com/gofiber/fiber/v2"
)

type AuthClientHandler struct {
	grpcClient pb.AuthServiceClient
}

func NewAuthClientHandler(grpcClient pb.AuthServiceClient) *AuthClientHandler {
	return &AuthClientHandler{
		grpcClient: grpcClient,
	}
}

func (h *AuthClientHandler) Register(c *fiber.Ctx) error {
	var req pb.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.Register(ctx, &req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": res})
}

func (h *AuthClientHandler) Login(c *fiber.Ctx) error {
	var req pb.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.Login(ctx, &req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *AuthClientHandler) Logout(c *fiber.Ctx) error {
	var req pb.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.Logout(ctx, &req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *AuthClientHandler) RefreshToken(c *fiber.Ctx) error {
	var req pb.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.RefreshToken(ctx, &req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "data": res})
}
