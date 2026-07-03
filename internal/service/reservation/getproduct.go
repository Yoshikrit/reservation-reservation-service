package reservation

import (
	"context"

	"reservation/internal/pkg/apperror"
)

func (s *reservationService) GetProduct(ctx context.Context, productID string) (*GetProductResponse, *apperror.AppError) {
	product, appErr := s.inventoryGateway.GetProductByID(ctx, productID)
	if appErr != nil {
		return nil, appErr
	}

	reserved, appErr := s.reservationRepo.SumReservedQuantity(ctx, productID)
	if appErr != nil {
		return nil, appErr
	}

	stockTotal := uint(product.Stock)
	available := uint(0)
	if stockTotal > reserved {
		available = stockTotal - reserved
	}

	return &GetProductResponse{
		ProductID:      product.ProductID,
		Name:           product.Name,
		BasePrice:      product.BasePrice,
		StockTotal:     stockTotal,
		StockReserved:  reserved,
		StockAvailable: available,
	}, nil
}
