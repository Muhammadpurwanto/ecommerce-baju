package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/broker"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/database"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/handler"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/repository"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/service"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/payment"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewMySQL(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	publisher, err := broker.NewRabbitMQPublisher(cfg.RabbitMQURL)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ Publisher", zap.Error(err))
	}

	// Dependency Injection
	paymentRepo := repository.NewPaymentRepository(db)
	paymentSvc := service.NewPaymentService(paymentRepo, cfg, publisher)
	paymentHandler := handler.NewPaymentHandler(paymentSvc)

	app := fiber.New(fiber.Config{
		AppName:      "Payment Service",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(fiberlogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowOrigins, // Menggunakan konfigurasi dinamis (bukan wildcard)
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-User-ID,X-User-Role",
	}))

	// Health check komprehensif
	app.Get("/health", func(c *fiber.Ctx) error {
		health := fiber.Map{
			"status":  "healthy",
			"service": "payment-service",
		}

		// Cek DB MySQL
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			health["mysql"] = "disconnected"
		} else {
			health["mysql"] = "connected"
		}

		return c.JSON(health)
	})

	// Inisialisasi gRPC Server
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	go func() {
		grpcPort := "50056"
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			logger.Fatal("Failed to listen on gRPC port", zap.Error(err))
		}

		logger.Info("gRPC Server is running on port " + grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", cfg.AppPort)); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	
	// Graceful shutdown gRPC
	grpcServer.GracefulStop()

	// Graceful shutdown HTTP Fiber
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server stopped gracefully")
}
