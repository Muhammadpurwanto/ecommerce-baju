package model

import "time"

type Order struct {
	ID              uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderNumber     string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	UserID          string      `gorm:"type:varchar(36);not null;index:idx_user_status" json:"user_id"`
	Status          string      `gorm:"type:enum('pending','paid','processing','shipped','completed','cancelled');default:'pending';index:idx_user_status" json:"status"`
	TotalAmount     float64     `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	ShippingCost    float64     `gorm:"type:decimal(10,2);not null" json:"shipping_cost"`
	ShippingAddress string      `gorm:"type:text;not null" json:"shipping_address"`
	Courier         string      `gorm:"type:varchar(50)" json:"courier"`
	TrackingNumber  string      `gorm:"type:varchar(100)" json:"tracking_number"`
	Notes           string      `gorm:"type:text" json:"notes"`
	PaidAt          *time.Time  `gorm:"index:idx_paid_at" json:"paid_at,omitempty"`
	ShippedAt       *time.Time  `gorm:"index:idx_shipped_at" json:"shipped_at,omitempty"`
	CreatedAt       time.Time   `gorm:"autoCreateTime;index:idx_created_at" json:"created_at"`
	UpdatedAt       time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
	Items           []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
	SnapToken       string      `gorm:"type:varchar(255)" json:"snap_token,omitempty"`
	PaymentURL      string      `gorm:"type:text" json:"payment_url,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
