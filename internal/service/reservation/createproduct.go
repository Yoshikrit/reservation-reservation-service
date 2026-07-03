package reservation

import (
	"context"

	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
	inventoryGateway "github.com/Yoshikrit/reservation/internal/gateway/grpc/inventory"
)

func (s *reservationService) CreateProduct(ctx context.Context, request *CreateProductRequest) *apperror.AppError {
	return s.inventoryGateway.CreateProduct(ctx, inventoryGateway.CreateProductRequest{
		ProductID: request.ProductID,
		Name:      request.Name,
		Stock:     uint64(request.Stock),
		BasePrice: request.BasePrice,
	})
}
