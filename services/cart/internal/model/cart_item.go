package model

import "time"

type CartItem struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CartID    uint      `gorm:"not null;index:idx_cart_id" json:"cart_id"`
	ProductID uint      `gorm:"not null;index:idx_product_id" json:"product_id"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (CartItem) TableName() string {
	return "cart_items"
}
