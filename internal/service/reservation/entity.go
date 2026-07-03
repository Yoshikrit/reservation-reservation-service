package reservation

import "time"

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

type CreateProductRequest struct {
	ProductID string
	Name      string
	Stock     uint
	BasePrice float64
}

type GetProductResponse struct {
	ProductID      string
	Name           string
	BasePrice      float64
	StockTotal     uint
	StockReserved  uint
	StockAvailable uint
}

type CreateReservationRequest struct {
	ProductID string
	Quantity  uint
	TtlSecond uint
}

type ListReservationBody struct {
	ProductID string
	Status    string
	Limit     int64
	Offset    int64
}

type ListReservationResponse struct {
	Reservations []ReservationItem `json:"reservations"`
	Pagination   PaginationInfo    `json:"pagination"`
}

type ReservationItem struct {
	ReservationID string    `json:"reservation_id"`
	ProductID     string    `json:"product_id"`
	Qty           uint      `json:"qty"`
	Status        string    `json:"status"`
	ExpiresAt     time.Time `json:"expires_at"`
	BasePrice     float64   `json:"base_price"`
	DiscountRate  int       `json:"discount_rate"`
	Price         float64   `json:"price"`
}

type PaginationInfo struct {
	Limit       int64 `json:"limit"`
	Offset      int64 `json:"offset"`
	Count       int64 `json:"count"`
	HasMore     bool  `json:"has_more"`
	CurrentPage int64 `json:"current_page"`
	PerPage     int64 `json:"per_page"`
}
