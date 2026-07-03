package inventory

import (
	"context"
	"time"

	"github.com/Yoshikrit/reservation/internal/gateway/grpc/inventory/pb"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
	"github.com/Yoshikrit/reservation/internal/pkg/grpcutil"
)

func (g *inventoryGateway) CheckAndHold(ctx context.Context, request CheckAndHoldRequest) (*CheckAndHoldResponse, *apperror.AppError) {
	ctx = grpcutil.AppendTrace(ctx)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := g.breaker.Execute(func() (any, error) {
		return g.inventoryClient.CheckAndHold(ctx, &pb.CheckAndHoldRequest{
			ProductId: request.ProductID,
			Quantity:  request.Quantity,
		})
	})
	if err != nil {
		return nil, grpcutil.ToAppError(err)
	}

	res := result.(*pb.CheckAndHoldResponse)
	return &CheckAndHoldResponse{
		Available:   res.Available,
		ProductID:   res.Product.ProductId,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
		Quantity:    res.Product.Quantity,
	}, nil
}
