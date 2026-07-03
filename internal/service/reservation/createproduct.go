package reservation

import (
	"context"

	"reservation/internal/pkg/apperror"
	inventoryGateway "reservation/internal/gateway/grpc/inventory"
)

func (s *reservationService) CreateProduct(ctx context.Context, request *CreateProductRequest) *apperror.AppError {
	return s.inventoryGateway.CreateProduct(ctx, inventoryGateway.CreateProductRequest{
		ProductID: request.ProductID,
		Name:      request.Name,
		Stock:     uint64(request.Stock),
		BasePrice: request.BasePrice,
	})
}
