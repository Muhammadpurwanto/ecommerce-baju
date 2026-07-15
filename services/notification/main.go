package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/broker"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/internal/handler"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/internal/service"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/notification/routes"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Dependency Injection
	emailSvc := service.NewEmailService(cfg, logger)
	notifHandler := handler.NewNotificationHandler(emailSvc)

	rmqConsumer, err := broker.NewRabbitMQConsumer(cfg.RabbitMQURL, emailSvc)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", zap.Error(err))
	} else {
		go rmqConsumer.Listen() // Jalankan listener di goroutine terpisah
	}

	app := fiber.New(fiber.Config{
		AppName:      "Notification Service",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(fiberlogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowOrigins, // Menggunakan konfigurasi dinamis (bukan wildcard)
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy", "service": "notification-service"})
	})

	routes.SetupRoutes(app, notifHandler)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	
	// Graceful shutdown RabbitMQ Consumer
	if rmqConsumer != nil {
		rmqConsumer.Close()
	}

	// Graceful shutdown HTTP Fiber
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server stopped gracefully")
}
