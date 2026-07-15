package dto

type CreatePaymentRequest struct {
	OrderID uint    `json:"order_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required,gt=0"`
}

type MidtransWebhookPayload map[string]interface{}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}
