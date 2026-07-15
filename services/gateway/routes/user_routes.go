package routes

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupUserRoutes(api fiber.Router, protected fiber.Handler, userClientHandler *handler.UserClientHandler) {
	userGroup := api.Group("/users", protected)
	
	// Profile routes
	userGroup.Get("/profile", userClientHandler.GetProfile)
	userGroup.Put("/profile", userClientHandler.UpdateProfile)
	
	// Address routes
	userGroup.Get("/addresses", userClientHandler.GetAddresses)
	userGroup.Get("/addresses/:id", userClientHandler.GetAddressByID)
	userGroup.Post("/addresses", userClientHandler.CreateAddress)
	userGroup.Put("/addresses/:id", userClientHandler.UpdateAddress)
	userGroup.Delete("/addresses/:id", userClientHandler.DeleteAddress)
	userGroup.Put("/addresses/:id/default", userClientHandler.SetDefaultAddress)

	// // Admin routes
	userGroup.Get("/", userClientHandler.GetAllUsers)
	userGroup.Get("/:id", userClientHandler.GetUserByID)
}
