package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/handler"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/middleware"
)

func SetupRoutes(
	app *fiber.App,
	cfg *config.Config,
	rdb *redis.Client,
	authClientHandler *handler.AuthClientHandler,
	userClientHandler *handler.UserClientHandler,
	productClientHandler *handler.ProductClientHandler,
	categoryClientHandler *handler.CategoryClientHandler,
	cartClientHandler *handler.CartClientHandler,
	orderClientHandler *handler.OrderClientHandler,
	paymentClientHandler *handler.PaymentClientHandler,
) {
	api := app.Group("/api/v1")
	protected := middleware.Protected(cfg, rdb)

	// ==========================================
	// ROUTING PROXY KE MICROSERVICES
	// ==========================================

	SetupAuthRoutes(api, cfg, authClientHandler)
	SetupUserRoutes(api, protected, userClientHandler)
	SetupProductRoutes(api, protected, cfg, productClientHandler, categoryClientHandler)
	SetupCartRoutes(api, protected, cartClientHandler)
	SetupOrderRoutes(api, protected, orderClientHandler)
	SetupPaymentRoutes(api, protected, paymentClientHandler)
}
