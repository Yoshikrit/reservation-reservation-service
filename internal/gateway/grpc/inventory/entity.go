package inventory

type CreateProductRequest struct {
	ProductID string
	Name      string
	Stock     uint64
	BasePrice float64
}

type GetProductByIDResponse struct {
	ProductID string
	Name      string
	BasePrice float64
	Stock     uint64
}

type CheckAndHoldRequest struct {
	ProductID string
	Quantity  uint64
}

type CheckAndHoldResponse struct {
	Available   bool
	ProductID   string
	Name        string
	Description string
	Price       float64
	Quantity    uint64
}
