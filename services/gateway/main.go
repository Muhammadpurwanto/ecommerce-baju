package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberredis "github.com/gofiber/storage/redis/v3"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/middleware"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/routes"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// ============================================================
	// Redis
	// ============================================================
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Warn("Failed to connect to Redis, token blacklist checking will be disabled", zap.Error(err))
		rdb = nil
	} else {
		logger.Info("Connected to Redis for token blacklist checking")
	}

	// ============================================================
	// Fiber App & Middleware
	// ============================================================
	app := fiber.New(fiber.Config{
		AppName:      "API Gateway",
		BodyLimit:    4 * 1024 * 1024, // 4MB max request body
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(middleware.RequestID())
	app.Use(fiberlogger.New())

	// Rate Limiter
	var storage fiber.Storage
	if cfg.RedisHost != "" {
		redisPort, _ := strconv.Atoi(cfg.RedisPort)
		storage = fiberredis.New(fiberredis.Config{
			Host:     cfg.RedisHost,
			Port:     redisPort,
			Password: cfg.RedisPassword,
			Database: 1, // DB terpisah dari cache auth (DB 0)
		})
	}

	app.Use(limiter.New(limiter.Config{
		Max:        cfg.RateLimitMax,
		Expiration: time.Duration(cfg.RateLimitExpiration) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":   "Too Many Requests",
				"message": "Anda telah mencapai batas request, silakan coba beberapa saat lagi.",
			})
		},
		Storage: storage,
	}))

	// CORS — configurable origins, bukan wildcard
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Strip header internal agar tidak bisa dipalsukan dari luar
	app.Use(middleware.StripInternalHeaders())

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		health := fiber.Map{
			"status":  "healthy",
			"service": "api-gateway",
		}

		// Check Redis connectivity
		if rdb != nil {
			if err := rdb.Ping(c.Context()).Err(); err != nil {
				health["status"] = "unhealthy"
				health["redis"] = "disconnected"
			} else {
				health["redis"] = "connected"
			}
		} else {
			health["redis"] = "not configured"
		}

		return c.JSON(health)
	})

	// ============================================================
	// gRPC Clients
	// ============================================================
	grpcClients := InitGrpcClients(cfg, logger)
	defer grpcClients.Close()

	// ============================================================
	// Routes
	// ============================================================
	routes.SetupRoutes(
		app, cfg, rdb,
		grpcClients.AuthHandler,
		grpcClients.UserHandler,
		grpcClients.ProductHandler,
		grpcClients.CategoryHandler,
		grpcClients.CartHandler,
		grpcClients.OrderHandler,
		grpcClients.PaymentHandler,
	)

	// ============================================================
	// Start Server with Graceful Shutdown
	// ============================================================
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
			logger.Fatal("Failed to start gateway", zap.Error(err))
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	logger.Info("Shutting down API Gateway...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error("Failed to gracefully shutdown", zap.Error(err))
	}

	logger.Info("API Gateway stopped gracefully")
}
