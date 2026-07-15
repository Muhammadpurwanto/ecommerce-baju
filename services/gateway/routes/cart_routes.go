package routes

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupCartRoutes(api fiber.Router, protected fiber.Handler, cartClientHandler *handler.CartClientHandler) {
	carts := api.Group("/carts", protected)
	carts.Get("/", cartClientHandler.GetCart)
	carts.Post("/items", cartClientHandler.AddItem)
	carts.Put("/items/:id", cartClientHandler.UpdateItem)
	carts.Delete("/items/:id", cartClientHandler.RemoveItem)
	carts.Delete("/", cartClientHandler.ClearCart)
}
