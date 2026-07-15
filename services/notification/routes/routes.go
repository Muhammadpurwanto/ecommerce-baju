package routes

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, notifHandler *handler.NotificationHandler) {
	api := app.Group("/api/v1/notifications")

	// Endpoint ini sengaja dibuat public DULU untuk mempermudah testing.
	// Nanti service ini hanya akan subscribe ke RabbitMQ queue.
	api.Post("/send-email", notifHandler.SendEmail)
}
