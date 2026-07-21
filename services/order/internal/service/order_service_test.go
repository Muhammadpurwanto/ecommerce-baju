package service

import (
	"context"
	"testing"

	"google.golang.org/grpc"

	cartpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/cart"
	paymentpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/payment"
	productpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/product"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/model"
)

// Mock OrderRepository
type mockOrderRepository struct {
	CreateFunc        func(order *model.Order) error
	FindByIDFunc      func(id uint) (*model.Order, error)
	FindByUserIDFunc  func(userID string, page, perPage int) ([]model.Order, int64, error)
	FindAllOrdersFunc func(page, perPage int) ([]model.Order, int64, error)
	UpdateFunc        func(order *model.Order) error
}

func (m *mockOrderRepository) Create(order *model.Order) error {
	return m.CreateFunc(order)
}
func (m *mockOrderRepository) FindByID(id uint) (*model.Order, error) {
	return m.FindByIDFunc(id)
}
func (m *mockOrderRepository) FindByUserID(userID string, page, perPage int) ([]model.Order, int64, error) {
	return m.FindByUserIDFunc(userID, page, perPage)
}
func (m *mockOrderRepository) FindAllOrders(page, perPage int) ([]model.Order, int64, error) {
	return m.FindAllOrdersFunc(page, perPage)
}
func (m *mockOrderRepository) Update(order *model.Order) error {
	return m.UpdateFunc(order)
}

// Mock EventPublisher
type mockEventPublisher struct {
	PublishOrderCreatedFunc func(email string, orderNumber string, amount float64) error
}

func (m *mockEventPublisher) PublishOrderCreated(email string, orderNumber string, amount float64) error {
	if m.PublishOrderCreatedFunc != nil {
		return m.PublishOrderCreatedFunc(email, orderNumber, amount)
	}
	return nil
}

// Mock gRPC clients
type mockProductServiceClient struct {
	productpb.ProductServiceClient
	ReserveStockFunc func(ctx context.Context, in *productpb.ReserveStockRequest, opts ...grpc.CallOption) (*productpb.StockResponse, error)
	ReleaseStockFunc func(ctx context.Context, in *productpb.ReleaseStockRequest, opts ...grpc.CallOption) (*productpb.StockResponse, error)
}

func (m *mockProductServiceClient) ReserveStock(ctx context.Context, in *productpb.ReserveStockRequest, opts ...grpc.CallOption) (*productpb.StockResponse, error) {
	return m.ReserveStockFunc(ctx, in, opts...)
}
func (m *mockProductServiceClient) ReleaseStock(ctx context.Context, in *productpb.ReleaseStockRequest, opts ...grpc.CallOption) (*productpb.StockResponse, error) {
	return m.ReleaseStockFunc(ctx, in, opts...)
}

type mockCartServiceClient struct {
	cartpb.CartServiceClient
	ClearCartFunc func(ctx context.Context, in *cartpb.ClearCartRequest, opts ...grpc.CallOption) (*cartpb.CartEmptyResponse, error)
}

func (m *mockCartServiceClient) ClearCart(ctx context.Context, in *cartpb.ClearCartRequest, opts ...grpc.CallOption) (*cartpb.CartEmptyResponse, error) {
	return m.ClearCartFunc(ctx, in, opts...)
}

type mockPaymentServiceClient struct {
	paymentpb.PaymentServiceClient
	CreatePaymentFunc func(ctx context.Context, in *paymentpb.CreatePaymentRequest, opts ...grpc.CallOption) (*paymentpb.Payment, error)
}

func (m *mockPaymentServiceClient) CreatePayment(ctx context.Context, in *paymentpb.CreatePaymentRequest, opts ...grpc.CallOption) (*paymentpb.Payment, error) {
	return m.CreatePaymentFunc(ctx, in, opts...)
}

func TestOrderService_GetOrderByID(t *testing.T) {
	t.Run("Success Get Order By ID", func(t *testing.T) {
		mockRepo := &mockOrderRepository{
			FindByIDFunc: func(id uint) (*model.Order, error) {
				return &model.Order{ID: 1, UserID: "user-1", Status: "pending"}, nil
			},
		}

		srv := NewOrderService(mockRepo, nil, nil, nil, nil)
		resp, err := srv.GetOrderByID(1, "user-1", false)

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, mendapat: %v", err)
		}
		if resp.ID != 1 {
			t.Errorf("diharapkan ID order 1, mendapat: %d", resp.ID)
		}
	})

	t.Run("Unauthorized Order Access", func(t *testing.T) {
		mockRepo := &mockOrderRepository{
			FindByIDFunc: func(id uint) (*model.Order, error) {
				return &model.Order{ID: 1, UserID: "user-1"}, nil
			},
		}

		srv := NewOrderService(mockRepo, nil, nil, nil, nil)
		resp, err := srv.GetOrderByID(1, "user-2", false)

		if err == nil {
			t.Fatal("diharapkan terjadi error otorisasi, mendapat nil")
		}
		if resp != nil {
			t.Errorf("diharapkan response nil, mendapat: %+v", resp)
		}
	})
}

func TestOrderService_CreateOrder(t *testing.T) {
	t.Run("Success Create Order", func(t *testing.T) {
		mockRepo := &mockOrderRepository{
			CreateFunc: func(order *model.Order) error {
				order.ID = 1
				return nil
			},
			UpdateFunc: func(order *model.Order) error {
				return nil
			},
		}

		mockProduct := &mockProductServiceClient{
			ReserveStockFunc: func(ctx context.Context, in *productpb.ReserveStockRequest, opts ...grpc.CallOption) (*productpb.StockResponse, error) {
				return &productpb.StockResponse{Success: true}, nil
			},
		}

		mockCart := &mockCartServiceClient{
			ClearCartFunc: func(ctx context.Context, in *cartpb.ClearCartRequest, opts ...grpc.CallOption) (*cartpb.CartEmptyResponse, error) {
				return &cartpb.CartEmptyResponse{Success: true}, nil
			},
		}

		mockPayment := &mockPaymentServiceClient{
			CreatePaymentFunc: func(ctx context.Context, in *paymentpb.CreatePaymentRequest, opts ...grpc.CallOption) (*paymentpb.Payment, error) {
				return &paymentpb.Payment{
					SnapToken:  "snap-123",
					PaymentUrl: "http://snap.url",
				}, nil
			},
		}

		srv := NewOrderService(mockRepo, nil, mockProduct, mockCart, mockPayment)
		req := &dto.CreateOrderRequest{
			ShippingCost:    15000,
			ShippingAddress: "Jl. Test",
			Items: []dto.OrderItemRequest{
				{ProductID: 1, Quantity: 2, Price: 50000},
			},
		}

		resp, err := srv.CreateOrder("user-1", req)
		if err != nil {
			t.Fatalf("diharapkan tidak ada error, mendapat: %v", err)
		}
		if resp.TotalAmount != 115000 {
			t.Errorf("diharapkan total amount 115000, mendapat: %f", resp.TotalAmount)
		}
		if resp.SnapToken != "snap-123" {
			t.Errorf("diharapkan snap token 'snap-123', mendapat: %s", resp.SnapToken)
		}
	})
}
