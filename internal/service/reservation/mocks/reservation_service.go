package mocks

import (
	"context"

	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
	svc "github.com/Yoshikrit/reservation/internal/service/reservation"

	"github.com/stretchr/testify/mock"
)

type ReservationService struct {
	mock.Mock
}

func (m *ReservationService) CreateProduct(ctx context.Context, request *svc.CreateProductRequest) *apperror.AppError {
	args := m.Called(ctx, request)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *ReservationService) GetProduct(ctx context.Context, productID string) (*svc.GetProductResponse, *apperror.AppError) {
	args := m.Called(ctx, productID)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).(*svc.GetProductResponse), nil
}

func (m *ReservationService) CreateReservation(ctx context.Context, request *svc.CreateReservationRequest) *apperror.AppError {
	args := m.Called(ctx, request)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *ReservationService) GetReservations(ctx context.Context, body *svc.ListReservationBody) (*svc.ListReservationResponse, *apperror.AppError) {
	args := m.Called(ctx, body)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).(*svc.ListReservationResponse), nil
}

func (m *ReservationService) CancelReservation(ctx context.Context, reservationID string) *apperror.AppError {
	args := m.Called(ctx, reservationID)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *ReservationService) ConfirmReservation(ctx context.Context, reservationID string) *apperror.AppError {
	args := m.Called(ctx, reservationID)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}
