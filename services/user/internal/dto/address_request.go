package dto

type CreateAddressRequest struct {
	Label      string `json:"label" validate:"required,min=1,max=50"`
	Recipient  string `json:"recipient" validate:"required,min=2,max=255"`
	Phone      string `json:"phone" validate:"required,min=8,max=20"`
	Province   string `json:"province" validate:"required,min=2,max=100"`
	City       string `json:"city" validate:"required,min=2,max=100"`
	District   string `json:"district" validate:"required,min=2,max=100"`
	PostalCode string `json:"postal_code" validate:"required,min=3,max=10"`
	Detail     string `json:"detail" validate:"required,min=5"`
	IsDefault  bool   `json:"is_default"`
}

type UpdateAddressRequest struct {
	Label      string `json:"label" validate:"omitempty,min=1,max=50"`
	Recipient  string `json:"recipient" validate:"omitempty,min=2,max=255"`
	Phone      string `json:"phone" validate:"omitempty,min=8,max=20"`
	Province   string `json:"province" validate:"omitempty,min=2,max=100"`
	City       string `json:"city" validate:"omitempty,min=2,max=100"`
	District   string `json:"district" validate:"omitempty,min=2,max=100"`
	PostalCode string `json:"postal_code" validate:"omitempty,min=3,max=10"`
	Detail     string `json:"detail" validate:"omitempty,min=5"`
	IsDefault  *bool  `json:"is_default"`
}
