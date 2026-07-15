package repository

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/payment/internal/model"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *model.Payment) error
	FindByOrderID(orderID uint) (*model.Payment, error)
	FindByTransactionID(trxID string) (*model.Payment, error)
	Update(payment *model.Payment) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *model.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByOrderID(orderID uint) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}

func (r *paymentRepository) FindByTransactionID(trxID string) (*model.Payment, error) {
	var payment model.Payment
	err := r.db.Where("transaction_id = ?", trxID).First(&payment).Error
	return &payment, err
}

func (r *paymentRepository) Update(payment *model.Payment) error {
	return r.db.Save(payment).Error
}
