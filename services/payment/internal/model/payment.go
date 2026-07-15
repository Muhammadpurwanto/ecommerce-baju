package model

import "time"

type Payment struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID       uint       `gorm:"not null;uniqueIndex" json:"order_id"`
	UserID        string     `gorm:"type:varchar(36);not null;index:idx_user_status" json:"user_id"`
	Amount        float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	PaymentType   string     `gorm:"type:varchar(50)" json:"payment_type"`
	TransactionID string     `gorm:"type:varchar(100);index" json:"transaction_id"`
	Status        string     `gorm:"type:enum('pending','settlement','expire','cancel','deny');default:'pending';index:idx_user_status" json:"status"`
	SnapToken     string     `gorm:"type:varchar(255)" json:"snap_token"`
	PaymentURL    string     `gorm:"type:text" json:"payment_url"`
	PaidAt        *time.Time `gorm:"index:idx_paid_at" json:"paid_at,omitempty"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Payment) TableName() string {
	return "payments"
}
