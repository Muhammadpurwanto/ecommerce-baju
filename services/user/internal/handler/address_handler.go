package handler

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/service"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/util"
)

type AddressHandler struct {
	service  service.AddressService
	validate *validator.Validate
}

func NewAddressHandler(svc service.AddressService) *AddressHandler {
	return &AddressHandler{
		service:  svc,
		validate: validator.New(),
	}
}

// GET /api/v1/users/me/addresses
func (h *AddressHandler) GetAddresses(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	addresses, err := h.service.GetAddresses(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.APIResponse{
		Success: true,
		Message: "Addresses retrieved successfully",
		Data:    addresses,
	})
}

// POST /api/v1/users/me/addresses
func (h *AddressHandler) CreateAddress(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	var req dto.CreateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  util.FormatValidationErrors(err),
		})
	}

	address, err := h.service.CreateAddress(userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(dto.APIResponse{
		Success: true,
		Message: "Address created successfully",
		Data:    address,
	})
}

// PUT /api/v1/users/me/addresses/:id
func (h *AddressHandler) UpdateAddress(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	addressID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid address ID",
		})
	}

	var req dto.UpdateAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Validation failed",
			Errors:  util.FormatValidationErrors(err),
		})
	}

	address, err := h.service.UpdateAddress(uint(addressID), userID, &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.APIResponse{
		Success: true,
		Message: "Address updated successfully",
		Data:    address,
	})
}

// DELETE /api/v1/users/me/addresses/:id
func (h *AddressHandler) DeleteAddress(c *fiber.Ctx) error {
	userID, err := util.GetUserIDFromContext(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
	}

	addressID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid address ID",
		})
	}

	if err := h.service.DeleteAddress(uint(addressID), userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return c.JSON(dto.APIResponse{
		Success: true,
		Message: "Address deleted successfully",
	})
}
