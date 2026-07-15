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

	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/database"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/handler"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/repository"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/service"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/storage"
	pb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/product"
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

	rdb, err := database.NewRedis(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	storageSvc, err := storage.NewMinioStorage(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	if err := storageSvc.CreateBucketIfNotExist(context.Background()); err != nil {
		logger.Warn("Failed to ensure bucket exists", zap.Error(err))
	}

	// Dependency Injection
	productRepo := repository.NewProductRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	productCache := service.NewProductCacheService(rdb)
	productSvc := service.NewProductService(productRepo, productCache)
	categorySvc := service.NewCategoryService(categoryRepo)

	// HTTP Handlers
	uploadHandler := handler.NewUploadHandler(storageSvc)

	// gRPC Handlers
	productGrpcHandler := handler.NewProductGrpcHandler(productSvc)
	categoryGrpcHandler := handler.NewCategoryGrpcHandler(categorySvc)

	app := fiber.New(fiber.Config{
		AppName:      "Product Service",
		ErrorHandler: customErrorHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	app.Use(recover.New())
	app.Use(fiberlogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-User-ID,X-User-Role",
	}))

	// Endpoint health check
	app.Get("/health", func(c *fiber.Ctx) error {
		health := fiber.Map{
			"status":  "healthy",
			"service": "product-service",
		}

		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			health["mysql"] = "disconnected"
		} else {
			health["mysql"] = "connected"
		}

		if err := rdb.Ping(context.Background()).Err(); err != nil {
			health["redis"] = "disconnected"
		} else {
			health["redis"] = "connected"
		}

		return c.JSON(health)
	})

	// Image Upload Route
	app.Post("/api/v1/products/upload-image", uploadHandler.UploadImage)

	// Inisialisasi gRPC Server
	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, productGrpcHandler)
	pb.RegisterCategoryServiceServer(grpcServer, categoryGrpcHandler)

	go func() {
		grpcPort := "50053"
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
	grpcServer.GracefulStop()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	
	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}
	logger.Info("Server stopped gracefully")
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "An internal error occurred"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"success": false,
		"message": message,
	})
}
