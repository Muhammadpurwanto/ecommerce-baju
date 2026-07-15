package routes

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func SetupAuthRoutes(api fiber.Router, cfg *config.Config, authClientHandler *handler.AuthClientHandler) {
	authGroup := api.Group("/auth")
	authGroup.Post("/register", authClientHandler.Register)
	authGroup.Post("/login", authClientHandler.Login)
	authGroup.Post("/logout", authClientHandler.Logout)
	authGroup.Post("/refresh", authClientHandler.RefreshToken)
	
	// Proxy khusus untuk OAuth (tetap HTTP karena redirect browser)
	authGroup.All("/google*", func(c *fiber.Ctx) error {
		target := cfg.AuthServiceHttpURL + c.OriginalURL()
		return proxy.Do(c, target)
	})
}
