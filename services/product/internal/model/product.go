package model

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID  uint      `gorm:"not null;index" json:"category_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	Brand       string    `gorm:"type:varchar(100);index:idx_brand" json:"brand"`
	Gender      string    `gorm:"type:enum('men','women','unisex');default:'unisex'" json:"gender"`
	BasePrice   float64   `gorm:"type:decimal(10,2);not null;index:idx_base_price" json:"base_price"`
	Weight      float64   `gorm:"type:decimal(10,2);not null" json:"weight"` // in grams
	Stock       int       `gorm:"type:int;not null;default:0" json:"stock"`
	ImageURL    string    `gorm:"type:varchar(512)" json:"image_url"`
	IsActive    bool      `gorm:"default:true;index:idx_is_active" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index:idx_created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Product) TableName() string {
	return "products"
}
