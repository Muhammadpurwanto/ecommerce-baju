package repository

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/order/internal/model"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *model.Order) error
	FindByID(id uint) (*model.Order, error)
	FindByUserID(userID string, page, perPage int) ([]model.Order, int64, error)
	FindAllOrders(page, perPage int) ([]model.Order, int64, error)
	Update(order *model.Order) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *model.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := r.db.Preload("Items").Where("id = ?", id).First(&order).Error
	return &order, err
}

func (r *orderRepository) FindByUserID(userID string, page, perPage int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.Model(&model.Order{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * perPage
	err := query.Preload("Items").Order("created_at desc").Offset(offset).Limit(perPage).Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) FindAllOrders(page, perPage int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := r.db.Model(&model.Order{})
	query.Count(&total)

	offset := (page - 1) * perPage
	err := query.Preload("Items").Order("created_at desc").Offset(offset).Limit(perPage).Find(&orders).Error

	return orders, total, err
}

func (r *orderRepository) Update(order *model.Order) error {
	return r.db.Save(order).Error
}
