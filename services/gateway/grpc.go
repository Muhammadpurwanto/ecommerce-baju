package main

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	authpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/auth"
	cartpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/cart"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/internal/handler"
	orderpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/order"
	paymentpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/payment"
	productpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/product"
	userpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/user"
)

// GrpcClients menampung semua koneksi dan handler gRPC ke backend services.
type GrpcClients struct {
	conns []*grpc.ClientConn

	AuthHandler     *handler.AuthClientHandler
	UserHandler     *handler.UserClientHandler
	ProductHandler  *handler.ProductClientHandler
	CategoryHandler *handler.CategoryClientHandler
	CartHandler     *handler.CartClientHandler
	OrderHandler    *handler.OrderClientHandler
	PaymentHandler  *handler.PaymentClientHandler
}

// Close menutup semua koneksi gRPC yang aktif.
func (gc *GrpcClients) Close() {
	for _, conn := range gc.conns {
		conn.Close()
	}
}

// dialGrpc membuat koneksi gRPC ke service tertentu.
// Menggunakan passthrough resolver agar kompatibel dengan Docker Compose DNS.
func dialGrpc(target, serviceName string, logger *zap.Logger) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		"passthrough:///"+target,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		logger.Fatal("Gagal koneksi ke gRPC service", zap.String("service", serviceName), zap.Error(err))
	}
	logger.Info("Terhubung ke gRPC service", zap.String("service", serviceName), zap.String("target", target))
	return conn
}

// InitGrpcClients menginisialisasi semua koneksi gRPC berdasarkan konfigurasi.
func InitGrpcClients(cfg *config.Config, logger *zap.Logger) *GrpcClients {
	gc := &GrpcClients{}

	// Auth Service
	authConn := dialGrpc(cfg.AuthServiceURL, "auth", logger)
	gc.conns = append(gc.conns, authConn)
	gc.AuthHandler = handler.NewAuthClientHandler(authpb.NewAuthServiceClient(authConn))

	// User Service
	userConn := dialGrpc(cfg.UserServiceURL, "user", logger)
	gc.conns = append(gc.conns, userConn)
	gc.UserHandler = handler.NewUserClientHandler(userpb.NewUserServiceClient(userConn))

	// Product Service (satu koneksi untuk Product, Category, dan Variant)
	productConn := dialGrpc(cfg.ProductServiceURL, "product", logger)
	gc.conns = append(gc.conns, productConn)
	gc.ProductHandler = handler.NewProductClientHandler(productpb.NewProductServiceClient(productConn))
	gc.CategoryHandler = handler.NewCategoryClientHandler(productpb.NewCategoryServiceClient(productConn))

	// Cart Service
	cartConn := dialGrpc(cfg.CartServiceURL, "cart", logger)
	gc.conns = append(gc.conns, cartConn)
	gc.CartHandler = handler.NewCartClientHandler(cartpb.NewCartServiceClient(cartConn))

	// Order Service
	orderConn := dialGrpc(cfg.OrderServiceURL, "order", logger)
	gc.conns = append(gc.conns, orderConn)
	gc.OrderHandler = handler.NewOrderClientHandler(orderpb.NewOrderServiceClient(orderConn))

	// Payment Service
	paymentConn := dialGrpc(cfg.PaymentServiceURL, "payment", logger)
	gc.conns = append(gc.conns, paymentConn)
	gc.PaymentHandler = handler.NewPaymentClientHandler(paymentpb.NewPaymentServiceClient(paymentConn))

	return gc
}
