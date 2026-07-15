package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/util"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/product"
	"github.com/gofiber/fiber/v2"
)

type ProductClientHandler struct {
	grpcClient pb.ProductServiceClient
}

func NewProductClientHandler(grpcClient pb.ProductServiceClient) *ProductClientHandler {
	return &ProductClientHandler{grpcClient: grpcClient}
}

func (h *ProductClientHandler) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.GetAllProducts(ctx, &pb.GetAllProductsRequest{
		Page:    int32(page),
		PerPage: int32(perPage),
	})
	if err != nil {
		return util.HandleGrpcError(c, err)
	}
	return c.JSON(fiber.Map{
		"success": true,
		"data":    res.Data,
		"meta": fiber.Map{
			"page":        res.Page,
			"per_page":    res.PerPage,
			"total":       res.Total,
			"total_pages": res.TotalPages,
		},
	})
}

func (h *ProductClientHandler) GetBySlug(c *fiber.Ctx) error {
	// Untuk saat ini kita menggunakan ID sementara di protobuf (GetProductByID) karena protobuf menggunakan ID.
	// Tetapi di HTTP kita tetap biarkan URL paramsnya :slug meskipun harus mengirim slug?
	// Oh wait, GetProductByID di gRPC pakai ID, tapi di route `/products/:slug` pakai slug.
	// Ah! Kita sebelumnya tidak buat GetProductBySlug di gRPC. Jadi ini akan error kalau user kirim slug string.
	// Jika gRPC pakai GetByID, kita parsing slug ke uint. Jika gagal berarti slug adalah string beneran.
	
	// Kita akan fallback sementara parsing ID
	id, err := strconv.ParseUint(c.Params("slug"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid product ID. Must be integer for now."})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, errGrpc := h.grpcClient.GetProductByID(ctx, &pb.GetByIDRequest{Id: uint32(id)})
	if errGrpc != nil {
		return util.HandleGrpcError(c, errGrpc)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *ProductClientHandler) Create(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Forbidden"})
	}

	var req pb.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.CreateProduct(ctx, &req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "data": res})
}

func (h *ProductClientHandler) Update(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Forbidden"})
	}

	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	var req pb.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid body"})
	}
	req.Id = uint32(id)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.UpdateProduct(ctx, &req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "data": res})
}

func (h *ProductClientHandler) Delete(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "Forbidden"})
	}

	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := h.grpcClient.DeleteProduct(ctx, &pb.GetByIDRequest{Id: uint32(id)})
	if err != nil {
		return util.HandleGrpcError(c, err)
	}
	return c.JSON(fiber.Map{"success": true, "message": "Product deleted"})
}
