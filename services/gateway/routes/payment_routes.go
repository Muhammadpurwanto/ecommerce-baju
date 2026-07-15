package routes

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupPaymentRoutes(api fiber.Router, protected fiber.Handler, paymentClientHandler *handler.PaymentClientHandler) {
	payments := api.Group("/payments")
	// Webhook tidak butuh token
	payments.Post("/webhook/midtrans", paymentClientHandler.WebhookMidtrans)

	// Yang lain butuh token
	payments.Post("/", protected, paymentClientHandler.CreatePayment)
	payments.Get("/order/:orderId", protected, paymentClientHandler.GetPaymentByOrderID)
}
