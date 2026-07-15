package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/util"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/user"
	"github.com/gofiber/fiber/v2"
)

type UserClientHandler struct {
	grpcClient pb.UserServiceClient
}

func NewUserClientHandler(grpcClient pb.UserServiceClient) *UserClientHandler {
	return &UserClientHandler{
		grpcClient: grpcClient,
	}
}

// GET /api/v1/profile
func (h *UserClientHandler) GetProfile(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	req := &pb.GetProfileRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.GetProfile(ctx, req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "message": "Profile retrieved successfully", "data": res})
}

// PUT /api/v1/profile
func (h *UserClientHandler) UpdateProfile(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	var req pb.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}
	req.UserId = userID // Pastikan UserId yang di-update adalah milik user yang sedang login

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.UpdateProfile(ctx, &req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "message": "Profile updated successfully", "data": res})
}

// GET /api/v1/users (admin)
func (h *UserClientHandler) GetAllUsers(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "access denied: require admin role"})
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "20"))

	req := &pb.GetAllUsersRequest{
		Page:    int32(page),
		PerPage: int32(perPage),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.GetAllUsers(ctx, req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "message": "Users retrieved successfully", "data": res.Data, "meta": res.Meta})
}

// GET /api/v1/users/:id (admin)
func (h *UserClientHandler) GetUserByID(c *fiber.Ctx) error {
	if !util.IsAdmin(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "access denied: require admin role"})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid user ID"})
	}

	req := &pb.GetUserByIDRequest{Id: id}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.GetUserByID(ctx, req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "message": "User retrieved successfully", "data": res})
}

// ==========================================
// ADDRESS CLIENT HANDLERS
// ==========================================

// GET /api/v1/users/addresses
func (h *UserClientHandler) GetAddresses(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	req := &pb.GetAddressesRequest{UserId: userID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.GetAddresses(ctx, req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "message": "Addresses retrieved successfully", "data": res.Data})
}

// GET /api/v1/users/addresses/:id
func (h *UserClientHandler) GetAddressByID(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	addressID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid address ID"})
	}

	req := &pb.GetAddressByIDRequest{Id: addressID, UserId: userID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.GetAddressByID(ctx, req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "message": "Address retrieved successfully", "data": res})
}

// POST /api/v1/users/addresses
func (h *UserClientHandler) CreateAddress(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	var reqBody struct {
		Label      string `json:"label"`
		Recipient  string `json:"recipient"`
		Phone      string `json:"phone"`
		Province   string `json:"province"`
		City       string `json:"city"`
		District   string `json:"district"`
		PostalCode string `json:"postal_code"`
		Detail     string `json:"detail"`
		IsDefault  bool   `json:"is_default"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	req := &pb.CreateAddressRequest{
		UserId:     userID,
		Label:      reqBody.Label,
		Recipient:  reqBody.Recipient,
		Phone:      reqBody.Phone,
		Province:   reqBody.Province,
		City:       reqBody.City,
		District:   reqBody.District,
		PostalCode: reqBody.PostalCode,
		Detail:     reqBody.Detail,
		IsDefault:  reqBody.IsDefault,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.CreateAddress(ctx, req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"success": true, "message": "Address created successfully", "data": res})
}

// PUT /api/v1/users/addresses/:id
func (h *UserClientHandler) UpdateAddress(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	addressID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid address ID"})
	}

	var reqBody struct {
		Label      *string `json:"label,omitempty"`
		Recipient  *string `json:"recipient,omitempty"`
		Phone      *string `json:"phone,omitempty"`
		Province   *string `json:"province,omitempty"`
		City       *string `json:"city,omitempty"`
		District   *string `json:"district,omitempty"`
		PostalCode *string `json:"postal_code,omitempty"`
		Detail     *string `json:"detail,omitempty"`
		IsDefault  *bool   `json:"is_default,omitempty"`
	}

	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid request body"})
	}

	req := &pb.UpdateAddressRequest{
		Id:     addressID,
		UserId: userID,
	}
	if reqBody.Label != nil {
		req.Label = reqBody.Label
	}
	if reqBody.Recipient != nil {
		req.Recipient = reqBody.Recipient
	}
	if reqBody.Phone != nil {
		req.Phone = reqBody.Phone
	}
	if reqBody.Province != nil {
		req.Province = reqBody.Province
	}
	if reqBody.City != nil {
		req.City = reqBody.City
	}
	if reqBody.District != nil {
		req.District = reqBody.District
	}
	if reqBody.PostalCode != nil {
		req.PostalCode = reqBody.PostalCode
	}
	if reqBody.Detail != nil {
		req.Detail = reqBody.Detail
	}
	if reqBody.IsDefault != nil {
		req.IsDefault = reqBody.IsDefault
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.UpdateAddress(ctx, req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "message": "Address updated successfully", "data": res})
}

// DELETE /api/v1/users/addresses/:id
func (h *UserClientHandler) DeleteAddress(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	addressID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid address ID"})
	}

	req := &pb.DeleteAddressRequest{Id: addressID, UserId: userID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.DeleteAddress(ctx, req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	if !res.Success {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": res.Message})
	}

	return c.JSON(fiber.Map{"success": true, "message": res.Message})
}

// PUT /api/v1/users/addresses/:id/default
func (h *UserClientHandler) SetDefaultAddress(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "Unauthorized"})
	}

	addressID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "Invalid address ID"})
	}

	req := &pb.SetDefaultAddressRequest{Id: addressID, UserId: userID}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.grpcClient.SetDefaultAddress(ctx, req)
	if err != nil {
		return util.HandleGrpcError(c, err)
	}

	return c.JSON(fiber.Map{"success": true, "message": "Default address set successfully", "data": res})
}
