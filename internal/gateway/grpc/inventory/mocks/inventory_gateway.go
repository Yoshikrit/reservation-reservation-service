package mocks

import (
	"context"

	inventoryGateway "github.com/Yoshikrit/reservation/internal/gateway/grpc/inventory"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"

	"github.com/stretchr/testify/mock"
)

type InventoryGateway struct {
	mock.Mock
}

func (m *InventoryGateway) CreateProduct(ctx context.Context, request inventoryGateway.CreateProductRequest) *apperror.AppError {
	args := m.Called(ctx, request)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *InventoryGateway) GetProductByID(ctx context.Context, productID string) (*inventoryGateway.GetProductByIDResponse, *apperror.AppError) {
	args := m.Called(ctx, productID)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).(*inventoryGateway.GetProductByIDResponse), nil
}

func (m *InventoryGateway) CheckAndHold(ctx context.Context, request inventoryGateway.CheckAndHoldRequest) (*inventoryGateway.CheckAndHoldResponse, *apperror.AppError) {
	args := m.Called(ctx, request)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).(*inventoryGateway.CheckAndHoldResponse), nil
}
