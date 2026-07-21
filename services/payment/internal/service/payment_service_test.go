package service

import (
	"testing"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/model"
)

// Mock PaymentRepository
type mockPaymentRepository struct {
	CreateFunc              func(payment *model.Payment) error
	FindByOrderIDFunc       func(orderID uint) (*model.Payment, error)
	FindByTransactionIDFunc func(trxID string) (*model.Payment, error)
	UpdateFunc              func(payment *model.Payment) error
}

func (m *mockPaymentRepository) Create(payment *model.Payment) error {
	return m.CreateFunc(payment)
}
func (m *mockPaymentRepository) FindByOrderID(orderID uint) (*model.Payment, error) {
	return m.FindByOrderIDFunc(orderID)
}
func (m *mockPaymentRepository) FindByTransactionID(trxID string) (*model.Payment, error) {
	return m.FindByTransactionIDFunc(trxID)
}
func (m *mockPaymentRepository) Update(payment *model.Payment) error {
	return m.UpdateFunc(payment)
}

// Mock EventPublisher
type mockEventPublisher struct {
	PublishPaymentSuccessFunc func(email string, orderID string) error
}

func (m *mockEventPublisher) PublishPaymentSuccess(email string, orderID string) error {
	if m.PublishPaymentSuccessFunc != nil {
		return m.PublishPaymentSuccessFunc(email, orderID)
	}
	return nil
}

func TestPaymentService_GetPaymentByOrderID(t *testing.T) {
	cfg := &config.Config{
		MidtransIsProd: false,
		MidtransServer: "dummy-key",
	}

	t.Run("Success Get Payment", func(t *testing.T) {
		mockRepo := &mockPaymentRepository{
			FindByOrderIDFunc: func(orderID uint) (*model.Payment, error) {
				return &model.Payment{OrderID: 1, UserID: "user-1", Status: "pending"}, nil
			},
		}

		srv := NewPaymentService(mockRepo, cfg, nil)
		resp, err := srv.GetPaymentByOrderID(1, "user-1")

		if err != nil {
			t.Fatalf("diharapkan tidak ada error, mendapat: %v", err)
		}
		if resp.OrderID != 1 {
			t.Errorf("diharapkan OrderID 1, mendapat: %d", resp.OrderID)
		}
	})

	t.Run("Unauthorized Get Payment", func(t *testing.T) {
		mockRepo := &mockPaymentRepository{
			FindByOrderIDFunc: func(orderID uint) (*model.Payment, error) {
				return &model.Payment{OrderID: 1, UserID: "user-1"}, nil
			},
		}

		srv := NewPaymentService(mockRepo, cfg, nil)
		resp, err := srv.GetPaymentByOrderID(1, "user-2")

		if err == nil {
			t.Fatal("diharapkan terjadi error otorisasi, mendapat nil")
		}
		if resp != nil {
			t.Errorf("diharapkan response nil, mendapat: %+v", resp)
		}
	})
}

func TestPaymentService_HandleMidtransWebhook(t *testing.T) {
	cfg := &config.Config{
		MidtransIsProd: false,
		MidtransServer: "dummy-key",
	}

	t.Run("Success Webhook Settlement", func(t *testing.T) {
		updated := false
		mockRepo := &mockPaymentRepository{
			FindByOrderIDFunc: func(orderID uint) (*model.Payment, error) {
				return &model.Payment{OrderID: 123, Status: "pending"}, nil
			},
			UpdateFunc: func(payment *model.Payment) error {
				if payment.Status != "settlement" {
					t.Errorf("diharapkan status 'settlement', mendapat: %s", payment.Status)
				}
				updated = true
				return nil
			},
		}

		srv := NewPaymentService(mockRepo, cfg, nil)
		payload := map[string]interface{}{
			"order_id":           "ORDER-123-1620000000",
			"transaction_status": "settlement",
			"payment_type":       "credit_card",
			"transaction_id":     "trx-456",
		}

		err := srv.HandleMidtransWebhook(payload)
		if err != nil {
			t.Fatalf("diharapkan sukses, mendapat error: %v", err)
		}
		if !updated {
			t.Error("diharapkan repository Update terpanggil")
		}
	})

	t.Run("Invalid Payload Webhook", func(t *testing.T) {
		mockRepo := &mockPaymentRepository{}
		srv := NewPaymentService(mockRepo, cfg, nil)

		err := srv.HandleMidtransWebhook(map[string]interface{}{})
		if err == nil {
			t.Fatal("diharapkan terjadi error untuk payload kosong")
		}
	})
}
