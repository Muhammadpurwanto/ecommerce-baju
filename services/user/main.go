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

	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/broker"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/database"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/handler"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/repository"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/service"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/user/internal/storage"
	authpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/auth"
	userpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/user"
)

func main() {
	// Logger
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()

	// Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Database
	db, err := database.NewMySQL(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Redis
	rdb, err := database.NewRedis(cfg, logger)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Storage
	storageSvc, err := storage.NewMinioStorage(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to MinIO", zap.Error(err))
	}
	
	if err := storageSvc.InitBucket(context.Background()); err != nil {
		logger.Fatal("Failed to initialize MinIO bucket", zap.Error(err))
	}

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	addressRepo := repository.NewAddressRepository(db)
	addressSvc := service.NewAddressService(addressRepo)

	// RabbitMQ Publisher
	eventPublisher, err := broker.NewRabbitMQPublisher(cfg.RabbitMQURL)
	if err != nil {
		logger.Warn("Failed to initialize RabbitMQ Event Publisher", zap.Error(err))
	}

	// Auth and JWT Services
	jwtSvc := service.NewJWTService(cfg)
	tokenCacheSvc := service.NewTokenCacheService(rdb)
	authSvc := service.NewAuthService(userRepo, jwtSvc, tokenCacheSvc, cfg, eventPublisher)

	// OAuth2 Handler
	oauthHandler := handler.NewOAuthHandler(cfg, authSvc)

	// Fiber App
	app := fiber.New(fiber.Config{
		AppName:      "User Service",
		ErrorHandler: customErrorHandler,
		BodyLimit:    5 * 1024 * 1024, // 5MB limit
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	})

	// Middleware
	app.Use(recover.New())
	app.Use(fiberlogger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowOrigins,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-User-ID,X-User-Role",
	}))

	// Health check (Komprehensif: MySQL & Redis)
	app.Get("/health", func(c *fiber.Ctx) error {
		health := fiber.Map{
			"status":  "healthy",
			"service": "user-service",
		}

		// Cek DB
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			health["mysql"] = "disconnected"
		} else {
			health["mysql"] = "connected"
		}

		// Cek Redis
		if err := rdb.Ping(context.Background()).Err(); err != nil {
			health["redis"] = "disconnected"
		} else {
			health["redis"] = "connected"
		}

		return c.JSON(health)
	})

	// Google OAuth2 Routes
	app.Get("/api/v1/auth/google", oauthHandler.HandleLogin)
	app.Get("/api/v1/auth/google/login", oauthHandler.HandleLogin)
	app.Get("/api/v1/auth/google/callback", oauthHandler.HandleCallback)

	grpcServer := grpc.NewServer()
	
	userHandler := handler.NewUserGrpcHandler(userSvc, addressSvc)
	authGrpcHandler := handler.NewAuthGrpcHandler(authSvc)

	userpb.RegisterUserServiceServer(grpcServer, userHandler)
	authpb.RegisterAuthServiceServer(grpcServer, authGrpcHandler)

	go func() {
		grpcPort := "50052"
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			logger.Fatal("Failed to listen on gRPC port", zap.Error(err))
		}
		
		logger.Info("gRPC Server is running on port " + grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// Menjalankan Fiber
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

	// Graceful shutdown Fiber
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
