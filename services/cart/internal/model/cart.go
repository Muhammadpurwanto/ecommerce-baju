package model

import "time"

type Cart struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    string     `gorm:"type:varchar(36);not null;uniqueIndex" json:"user_id"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	Items     []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE" json:"items,omitempty"`
}

func (Cart) TableName() string {
	return "carts"
}
