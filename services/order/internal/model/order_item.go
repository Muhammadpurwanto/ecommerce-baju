package model

type OrderItem struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID   uint    `gorm:"not null;index:idx_order_id" json:"order_id"`
	ProductID uint    `gorm:"not null;index:idx_product_id" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	Subtotal  float64 `gorm:"type:decimal(10,2);not null" json:"subtotal"`
}

func (OrderItem) TableName() string {
	return "order_items"
}
