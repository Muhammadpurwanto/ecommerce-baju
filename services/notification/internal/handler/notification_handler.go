package handler

import (
	"regexp"

	"github.com/gofiber/fiber/v2"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/internal/service"
)

type NotificationHandler struct {
	emailService service.EmailService
}

func NewNotificationHandler(emailSvc service.EmailService) *NotificationHandler {
	return &NotificationHandler{emailService: emailSvc}
}

// Regex sederhana untuk memvalidasi format email
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// Endpoint HTTP sementara untuk ngetest kirim email
func (h *NotificationHandler) SendEmail(c *fiber.Ctx) error {
	var req dto.SendEmailRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
		})
	}

	if req.To == "" || req.Subject == "" || req.Body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "To, Subject, and Body fields are required",
		})
	}

	if !emailRegex.MatchString(req.To) {
		return c.Status(fiber.StatusBadRequest).JSON(dto.APIResponse{
			Success: false,
			Message: "Invalid email format",
		})
	}

	// Proses kirim email dibuat asynchronous (background)
	go func() {
		_ = h.emailService.SendEmail(&req)
	}()

	return c.JSON(dto.APIResponse{
		Success: true,
		Message: "Email is being processed",
	})
}
