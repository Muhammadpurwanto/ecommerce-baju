package routes

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupOrderRoutes(api fiber.Router, protected fiber.Handler, orderClientHandler *handler.OrderClientHandler) {
	orders := api.Group("/orders", protected)
	orders.Post("/", orderClientHandler.CreateOrder)
	orders.Get("/", orderClientHandler.GetMyOrders)
	orders.Get("/all", orderClientHandler.GetAllOrders) // Harus taruh sebelum /:id agar tidak terbaca sebagai parameter :id
	orders.Get("/:id", orderClientHandler.GetOrderByID)
	orders.Put("/:id/cancel", orderClientHandler.CancelOrder)
	orders.Put("/:id/status", orderClientHandler.UpdateOrderStatus)
}
