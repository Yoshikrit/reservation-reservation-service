package product

type CreateProductRequest struct {
	ProductID string  `json:"product_id" validate:"required"`
	Name      string  `json:"name"       validate:"required"`
	Stock     uint    `json:"stock"      validate:"gte=0"`
	BasePrice float64 `json:"base_price" validate:"gt=0"`
}

type GetProductResponse struct {
	ProductID      string  `json:"product_id"`
	Name           string  `json:"name"`
	BasePrice      float64 `json:"base_price"`
	StockTotal     uint    `json:"stock_total"`
	StockReserved  uint    `json:"stock_reserved"`
	StockAvailable uint    `json:"stock_available"`
}
