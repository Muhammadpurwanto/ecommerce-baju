package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/broker"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/config"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/dto"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/model"
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/repository"
)

type PaymentService interface {
	CreatePayment(userID string, req *dto.CreatePaymentRequest) (*model.Payment, error)
	GetPaymentByOrderID(orderID uint, userID string) (*model.Payment, error)
	HandleMidtransWebhook(payload map[string]interface{}) error
}

type paymentService struct {
	repo repository.PaymentRepository
	cfg  *config.Config
	snap snap.Client
	eventPub broker.EventPublisher
}

func NewPaymentService(repo repository.PaymentRepository, cfg *config.Config, eventPub broker.EventPublisher) PaymentService {
	var s snap.Client
	env := midtrans.Sandbox
	if cfg.MidtransIsProd {
		env = midtrans.Production
	}
	s.New(cfg.MidtransServer, env)

	return &paymentService{repo: repo, cfg: cfg, snap: s, eventPub: eventPub}
}

func (s *paymentService) CreatePayment(userID string, req *dto.CreatePaymentRequest) (*model.Payment, error) {
	// Cek jika payment sudah ada untuk order ini
	existing, err := s.repo.FindByOrderID(req.OrderID)
	if err == nil {
		return existing, nil // Return exist jika sudah ada
	}

	orderIDStr := fmt.Sprintf("ORDER-%d-%d", req.OrderID, time.Now().Unix())

	// Persiapkan req ke Midtrans
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderIDStr,
			GrossAmt: int64(req.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			// Akan dikembangkan nanti dengan call ke user-service
			FName: "User",
			LName: userID,
		},
		Callbacks: &snap.Callbacks{
			Finish: fmt.Sprintf("%s/orders", s.cfg.FrontendURL),
		},
	}

	var snapToken string
	var paymentURL string

	if s.cfg.MidtransServer == "SB-Mid-server-xxx" || s.cfg.MidtransServer == "" {
		// Mock response untuk testing / CI/CD
		snapToken = "mock-snap-token-12345"
		paymentURL = "https://app.sandbox.midtrans.com/snap/v2/vtweb/mock-snap-token-12345"
	} else {
		// Buat transaksi Midtrans asli
		snapResp, midErr := s.snap.CreateTransaction(snapReq)
		if midErr != nil {
			return nil, fmt.Errorf("midtrans error: %v", midErr)
		}
		snapToken = snapResp.Token
		paymentURL = snapResp.RedirectURL
	}

	payment := &model.Payment{
		OrderID:    req.OrderID,
		UserID:     userID,
		Amount:     req.Amount,
		Status:     "pending",
		SnapToken:  snapToken,
		PaymentURL: paymentURL,
	}

	if err := s.repo.Create(payment); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *paymentService) GetPaymentByOrderID(orderID uint, userID string) (*model.Payment, error) {
	payment, err := s.repo.FindByOrderID(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}

	if payment.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return payment, nil
}

func (s *paymentService) HandleMidtransWebhook(payload map[string]interface{}) error {
	orderIDStr, ok := payload["order_id"].(string)
	if !ok {
		return errors.New("invalid order_id in payload")
	}

	transactionStatus, ok := payload["transaction_status"].(string)
	if !ok {
		return errors.New("invalid transaction_status in payload")
	}

	var orderID uint
	fmt.Sscanf(orderIDStr, "ORDER-%d", &orderID)

	payment, err := s.repo.FindByOrderID(orderID)
	if err != nil {
		return err
	}

	switch transactionStatus {
	case "capture", "settlement":
		payment.Status = "settlement"
		now := time.Now()
		payment.PaidAt = &now

		if s.eventPub != nil {
			go s.eventPub.PublishPaymentSuccess("user@example.com", fmt.Sprintf("%d", payment.OrderID))
		}
		if pType, ok := payload["payment_type"].(string); ok {
			payment.PaymentType = pType
		}
		if trxID, ok := payload["transaction_id"].(string); ok {
			payment.TransactionID = trxID
		}
	case "deny", "cancel", "expire":
		payment.Status = "cancel"
	case "pending":
		payment.Status = "pending"
	}

	if err := s.repo.Update(payment); err != nil {
		return err
	}

	// TODO: Publish event (RabbitMQ) => Payment Success ke Order Service

	return nil
}
