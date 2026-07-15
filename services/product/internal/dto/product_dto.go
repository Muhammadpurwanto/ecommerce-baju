package dto

type ProductRequest struct {
	CategoryID  uint    `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required,min=3"`
	Description string  `json:"description"`
	Brand       string  `json:"brand" validate:"required"`
	Gender      string  `json:"gender" validate:"required,oneof=men women unisex"`
	BasePrice   float64 `json:"base_price" validate:"required,gt=0"`
	Weight      float64 `json:"weight" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"min=0"`
	ImageURL    string  `json:"image_url"`
	IsActive    bool    `json:"is_active"`
}

type CategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2"`
}

type StockItem struct {
	ProductID uint `json:"product_id" validate:"required"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}
