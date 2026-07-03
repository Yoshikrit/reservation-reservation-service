package inventory

import (
	"context"
	"time"

	"github.com/Yoshikrit/reservation/internal/gateway/grpc/inventory/pb"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"

	"github.com/sony/gobreaker/v2"
)

type InventoryGateway interface {
	CreateProduct(ctx context.Context, request CreateProductRequest) *apperror.AppError
	GetProductByID(ctx context.Context, productID string) (*GetProductByIDResponse, *apperror.AppError)
	CheckAndHold(ctx context.Context, request CheckAndHoldRequest) (*CheckAndHoldResponse, *apperror.AppError)
}

type inventoryGateway struct {
	inventoryClient pb.InventoryServiceClient
	breaker         *gobreaker.CircuitBreaker[any]
}

func NewInventoryGateway(inventoryClient pb.InventoryServiceClient) InventoryGateway {
	breaker := gobreaker.NewCircuitBreaker[any](gobreaker.Settings{
		Name:        "inventory-grpc",
		MaxRequests: 1,
		Interval:    30 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
	})
	return &inventoryGateway{inventoryClient: inventoryClient, breaker: breaker}
}
