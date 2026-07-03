package reservation

type CreateReservationRequest struct {
	ProductID string `json:"product_id" validate:"required"`
	Quantity  uint   `json:"quantity"   validate:"required,gt=0"`
	TtlSecond uint   `json:"ttl_second" validate:"required,gte=30,lte=300"`
}

type GetReservationsQuery struct {
	ProductID string `query:"product_id"`
	Status    string `query:"status"`
	Limit     int64  `query:"limit"`
	Offset    int64  `query:"offset"`
}

