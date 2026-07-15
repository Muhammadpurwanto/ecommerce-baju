package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	cartpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/cart"
	paymentpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/payment"
	productpb "github.com/Muhammadpurwanto/ecommerce-baju/services/common/pb/product"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/broker"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/repository"
)

type OrderService interface {
	CreateOrder(userID string, req *dto.CreateOrderRequest) (*model.Order, error)
	GetOrderByID(id uint, userID string, isAdmin bool) (*model.Order, error)
	GetUserOrders(userID string, page, perPage int) ([]model.Order, int64, error)
	GetAllOrders(page, perPage int) ([]model.Order, int64, error)
	UpdateOrderStatus(id uint, status string) (*model.Order, error)
	CancelOrder(id uint, userID string) (*model.Order, error)
}

type orderService struct {
	repo          repository.OrderRepository
	eventPub      broker.EventPublisher
	productClient productpb.ProductServiceClient
	cartClient    cartpb.CartServiceClient
	paymentClient paymentpb.PaymentServiceClient
}

func NewOrderService(
	repo repository.OrderRepository,
	eventPub broker.EventPublisher,
	productClient productpb.ProductServiceClient,
	cartClient cartpb.CartServiceClient,
	paymentClient paymentpb.PaymentServiceClient,
) OrderService {
	return &orderService{
		repo:          repo,
		eventPub:      eventPub,
		productClient: productClient,
		cartClient:    cartClient,
		paymentClient: paymentClient,
	}
}

func (s *orderService) CreateOrder(userID string, req *dto.CreateOrderRequest) (*model.Order, error) {
	// Generate simple Order Number (INV-TIMESTAMP-USERID_SHORT)
	shortUserID := userID
	if len(userID) > 8 {
		shortUserID = userID[:8]
	}
	orderNumber := fmt.Sprintf("INV-%d-%s", time.Now().Unix(), shortUserID)

	var totalAmount float64
	var orderItems []model.OrderItem

	for _, itemReq := range req.Items {
		subtotal := float64(itemReq.Quantity) * itemReq.Price
		totalAmount += subtotal
		orderItems = append(orderItems, model.OrderItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			Price:     itemReq.Price,
			Subtotal:  subtotal,
		})
	}

	totalAmount += req.ShippingCost

	order := &model.Order{
		OrderNumber:     orderNumber,
		UserID:          userID,
		Status:          "pending",
		TotalAmount:     totalAmount,
		ShippingCost:    req.ShippingCost,
		ShippingAddress: req.ShippingAddress,
		Courier:         req.Courier,
		Notes:           req.Notes,
		Items:           orderItems,
	}

	// Step 1: Save Order to Database in pending status
	if err := s.repo.Create(order); err != nil {
		return nil, err
	}

	// Step 2: Reserve Stock via Product Service (gRPC)
	var stockItems []*productpb.StockItem
	for _, item := range order.Items {
		stockItems = append(stockItems, &productpb.StockItem{
			ProductId: uint32(item.ProductID),
			Quantity:  int32(item.Quantity),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := s.productClient.ReserveStock(ctx, &productpb.ReserveStockRequest{Items: stockItems})
	if err != nil {
		// Compensate Step 1: Mark Order as Cancelled
		order.Status = "cancelled"
		order.Notes = fmt.Sprintf("%s (Auto-cancelled: %v)", order.Notes, err)
		_ = s.repo.Update(order)
		return nil, fmt.Errorf("failed to reserve stock: %w", err)
	}

	// Step 3: Clear Cart via Cart Service (gRPC) - Non-blocking
	go func() {
		cartCtx, cartCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cartCancel()
		_, cartErr := s.cartClient.ClearCart(cartCtx, &cartpb.ClearCartRequest{UserId: userID})
		if cartErr != nil {
			fmt.Printf("Warning: failed to clear cart for user %s: %v\n", userID, cartErr)
		}
	}()

	// Step 4: Create Payment Session via Payment Service (gRPC)
	paymentCtx, paymentCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer paymentCancel()

	paymentRes, err := s.paymentClient.CreatePayment(paymentCtx, &paymentpb.CreatePaymentRequest{
		UserId:  userID,
		OrderId: uint32(order.ID),
		Amount:  order.TotalAmount,
	})
	if err != nil {
		// Compensate Step 2: Release Stock
		releaseCtx, releaseCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer releaseCancel()
		_, _ = s.productClient.ReleaseStock(releaseCtx, &productpb.ReleaseStockRequest{Items: stockItems})

		// Compensate Step 1: Mark Order as Cancelled
		order.Status = "cancelled"
		order.Notes = fmt.Sprintf("%s (Auto-cancelled: failed to create payment session: %v)", order.Notes, err)
		_ = s.repo.Update(order)

		return nil, fmt.Errorf("failed to create payment session: %w", err)
	}

	// Save payment details to order response
	order.SnapToken = paymentRes.GetSnapToken()
	order.PaymentURL = paymentRes.GetPaymentUrl()
	_ = s.repo.Update(order)

	if s.eventPub != nil {
		go s.eventPub.PublishOrderCreated("user@example.com", order.OrderNumber, order.TotalAmount)
	}

	return order, nil
}

func (s *orderService) GetOrderByID(id uint, userID string, isAdmin bool) (*model.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	// Security check: Hanya pemilik asli atau admin yang boleh lihat
	if !isAdmin && order.UserID != userID {
		return nil, errors.New("unauthorized to cancel this order")
	}
	return order, nil
}

func (s *orderService) GetUserOrders(userID string, page, perPage int) ([]model.Order, int64, error) {
	orders, meta, err := s.repo.FindByUserID(userID, page, perPage)
	return orders, meta, err
}

func (s *orderService) GetAllOrders(page, perPage int) ([]model.Order, int64, error) {
	return s.repo.FindAllOrders(page, perPage)
}

func (s *orderService) UpdateOrderStatus(id uint, status string) (*model.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	validStatuses := map[string]bool{
		"pending":    true,
		"paid":       true,
		"processing": true,
		"shipped":    true,
		"completed":  true,
		"cancelled":  true,
	}

	if !validStatuses[status] {
		return nil, errors.New("invalid status")
	}

	order.Status = status
	if err := s.repo.Update(order); err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) CancelOrder(id uint, userID string) (*model.Order, error) {
	order, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if order.Status != "pending" {
		return nil, errors.New("only pending orders can be cancelled")
	}

	order.Status = "cancelled"
	if err := s.repo.Update(order); err != nil {
		return nil, err
	}

	// TODO: Publish event ke RabbitMQ (Order Cancelled -> Return Stock)

	return order, nil
}
