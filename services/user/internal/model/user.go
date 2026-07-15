package model

import "time"

type User struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Password  string    `gorm:"type:varchar(255)" json:"-"`
	Provider  string    `gorm:"type:enum('local','google');default:'local';index:idx_provider" json:"provider"`
	ProviderID *string   `gorm:"type:varchar(255);index:idx_provider_id" json:"provider_id,omitempty"`
	Phone     *string   `gorm:"type:varchar(20)" json:"phone,omitempty"`
	AvatarURL *string   `gorm:"type:varchar(500);column:avatar_url" json:"avatar_url,omitempty"`
	Role      string    `gorm:"type:enum('customer','admin');default:'customer'" json:"role"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Addresses []Address `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"addresses,omitempty"`
}

func (User) TableName() string {
	return "users"
}
