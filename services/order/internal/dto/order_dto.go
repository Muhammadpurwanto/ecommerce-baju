package dto

type CreateOrderRequest struct {
	ShippingCost    float64            `json:"shipping_cost" validate:"required,min=0"`
	ShippingAddress string             `json:"shipping_address" validate:"required"`
	Courier         string             `json:"courier" validate:"required"`
	Notes           string             `json:"notes"`
	Items           []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductID uint    `json:"product_id" validate:"required"`
	Quantity  int     `json:"quantity" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,min=0"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}
