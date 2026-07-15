package model

import "time"

type Address struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     string    `gorm:"type:varchar(36);not null;index:idx_user_id" json:"user_id"`
	Label      string    `gorm:"type:varchar(50);not null" json:"label"`
	Recipient  string    `gorm:"type:varchar(255);not null" json:"recipient"`
	Phone      string    `gorm:"type:varchar(20);not null" json:"phone"`
	Province   string    `gorm:"type:varchar(100);not null" json:"province"`
	City       string    `gorm:"type:varchar(100);not null" json:"city"`
	District   string    `gorm:"type:varchar(100);not null" json:"district"`
	PostalCode string    `gorm:"type:varchar(10);not null;column:postal_code" json:"postal_code"`
	Detail     string    `gorm:"type:text;not null" json:"detail"`
	IsDefault  bool      `gorm:"default:false;index:idx_is_default" json:"is_default"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Address) TableName() string {
	return "addresses"
}
