package routes

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func SetupProductRoutes(
	api fiber.Router, 
	protected fiber.Handler, 
	cfg *config.Config, 
	productClientHandler *handler.ProductClientHandler,
	categoryClientHandler *handler.CategoryClientHandler,
) {
	// === Categories ===
	categories := api.Group("/categories")
	categories.Get("/", categoryClientHandler.GetAll)
	categories.Post("/", protected, categoryClientHandler.Create)
	categories.Put("/:id", protected, categoryClientHandler.Update)
	categories.Delete("/:id", protected, categoryClientHandler.Delete)

	// === Products ===
	products := api.Group("/products")
	products.Get("/", productClientHandler.GetAll)
	products.Get("/:slug", productClientHandler.GetBySlug)
	products.Post("/", protected, productClientHandler.Create)
	products.Put("/:id", protected, productClientHandler.Update)
	products.Delete("/:id", protected, productClientHandler.Delete)

	// Proxy upload image to product-service HTTP server
	products.Post("/upload-image", protected, func(c *fiber.Ctx) error {
		target := cfg.ProductServiceHttpURL + c.OriginalURL()
		return proxy.Do(c, target)
	})
}
