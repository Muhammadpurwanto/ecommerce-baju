package repository

import (
	"github.com/Muhammadpurwanto/ecommerce-baju/services/cart/internal/model"
	"gorm.io/gorm"
)

type CartRepository interface {
	GetCartByUserID(userID string) (*model.Cart, error)
	CreateCart(cart *model.Cart) error
	
	FindItem(cartID, productID uint) (*model.CartItem, error)
	AddItem(item *model.CartItem) error
	UpdateItem(item *model.CartItem) error
	RemoveItem(itemID uint) error
	ClearCart(cartID uint) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) GetCartByUserID(userID string) (*model.Cart, error) {
	var cart model.Cart
	err := r.db.Preload("Items").Where("user_id = ?", userID).First(&cart).Error
	return &cart, err
}

func (r *cartRepository) CreateCart(cart *model.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) FindItem(cartID, productID uint) (*model.CartItem, error) {
	var item model.CartItem
	err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&item).Error
	return &item, err
}

func (r *cartRepository) AddItem(item *model.CartItem) error {
	return r.db.Create(item).Error
}

func (r *cartRepository) UpdateItem(item *model.CartItem) error {
	return r.db.Model(item).Update("quantity", item.Quantity).Error
}

func (r *cartRepository) RemoveItem(itemID uint) error {
	return r.db.Delete(&model.CartItem{}, itemID).Error
}

func (r *cartRepository) ClearCart(cartID uint) error {
	return r.db.Where("cart_id = ?", cartID).Delete(&model.CartItem{}).Error
}
