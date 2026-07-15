package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/util"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/product"
	"github.com/gofiber/fiber/v2"
)

type CategoryClientHandler struct {
	grpcClient pb.CategoryServiceClient
}

func NewCategoryClientHandler(grpcClient pb.CategoryServiceClient) *CategoryClientHandler {
	return &CategoryClientHandler{grpcClient: grpcClient}
}

func (h *CategoryClientHandler) GetAll(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.GetAllCategories(ctx, &pb.Empty{})
	if err != nil {
		return util.HandleGrpcError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": res.Data})
}

func (h *CategoryClientHandler) Create(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Forbidden"})
	}

	var req pb.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.CreateCategory(ctx, &req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": res})
}

func (h *CategoryClientHandler) Update(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Forbidden"})
	}

	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var req pb.UpdateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}
	req.Id = uint32(id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.UpdateCategory(ctx, &req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *CategoryClientHandler) Delete(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Forbidden"})
	}

	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	req := &pb.GetByIDRequest{Id: uint32(id)}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.grpcClient.DeleteCategory(ctx, req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "message": "Category deleted"})
}
