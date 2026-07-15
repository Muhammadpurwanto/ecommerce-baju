package dto

import "time"

type AddressResponse struct {
	ID         uint      `json:"id"`
	UserID     string    `json:"user_id"`
	Label      string    `json:"label"`
	Recipient  string    `json:"recipient"`
	Phone      string    `json:"phone"`
	Province   string    `json:"province"`
	City       string    `json:"city"`
	District   string    `json:"district"`
	PostalCode string    `json:"postal_code"`
	Detail     string    `json:"detail"`
	IsDefault  bool      `json:"is_default"`
	CreatedAt  time.Time `json:"created_at"`
}
