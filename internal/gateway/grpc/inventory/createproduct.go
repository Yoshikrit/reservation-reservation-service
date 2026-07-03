package inventory

import (
	"context"
	"time"

	"github.com/Yoshikrit/reservation/internal/gateway/grpc/inventory/pb"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
	"github.com/Yoshikrit/reservation/internal/pkg/grpcutil"
)

func (g *inventoryGateway) CreateProduct(ctx context.Context, request CreateProductRequest) *apperror.AppError {
	ctx = grpcutil.AppendTrace(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := g.breaker.Execute(func() (any, error) {
		return g.inventoryClient.CreateProduct(ctx, &pb.CreateProductRequest{
			ProductId: request.ProductID,
			Name:      request.Name,
			Price:     request.BasePrice,
			Quantity:  request.Stock,
		})
	})
	if err != nil {
		return grpcutil.ToAppError(err)
	}
	return nil
}
