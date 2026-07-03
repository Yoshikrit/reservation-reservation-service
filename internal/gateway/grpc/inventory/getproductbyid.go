package inventory

import (
	"context"
	"time"

	"reservation/internal/gateway/grpc/inventory/pb"
	"reservation/internal/pkg/apperror"
	"reservation/internal/pkg/grpcutil"
)

func (g *inventoryGateway) GetProductByID(ctx context.Context, productID string) (*GetProductByIDResponse, *apperror.AppError) {
	ctx = grpcutil.AppendTrace(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := g.breaker.Execute(func() (any, error) {
		return g.inventoryClient.GetProductByID(ctx, &pb.GetProductByIDRequest{
			ProductId: productID,
		})
	})
	if err != nil {
		return nil, grpcutil.ToAppError(err)
	}

	res := result.(*pb.ProductResponse)
	return &GetProductByIDResponse{
		ProductID: res.ProductId,
		Name:      res.Name,
		BasePrice: res.Price,
		Stock:     res.Quantity,
	}, nil
}
