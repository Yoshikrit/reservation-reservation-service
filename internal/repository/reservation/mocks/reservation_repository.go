package mocks

import (
	"context"

	"github.com/Yoshikrit/reservation/internal/entity"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"

	"github.com/stretchr/testify/mock"
)

type ReservationRepository struct {
	mock.Mock
}

func (m *ReservationRepository) Create(ctx context.Context, reservation *entity.Reservation) *apperror.AppError {
	args := m.Called(ctx, reservation)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *ReservationRepository) Patch(ctx context.Context, reservation *entity.Reservation) *apperror.AppError {
	args := m.Called(ctx, reservation)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *ReservationRepository) Find(ctx context.Context, filter *entity.Reservation) (*entity.Reservation, *apperror.AppError) {
	args := m.Called(ctx, filter)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).(*entity.Reservation), nil
}

func (m *ReservationRepository) Filter(ctx context.Context, filter *entity.Reservation, limit, offset int, isAsc bool) ([]entity.Reservation, *apperror.AppError) {
	args := m.Called(ctx, filter, limit, offset, isAsc)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).([]entity.Reservation), nil
}

func (m *ReservationRepository) FindAll(ctx context.Context) ([]entity.Reservation, *apperror.AppError) {
	args := m.Called(ctx)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).([]entity.Reservation), nil
}

func (m *ReservationRepository) Count(ctx context.Context, filter *entity.Reservation) (int64, *apperror.AppError) {
	args := m.Called(ctx, filter)
	if v := args.Get(1); v != nil {
		return 0, v.(*apperror.AppError)
	}
	return args.Get(0).(int64), nil
}

func (m *ReservationRepository) SumReservedQuantity(ctx context.Context, productID string) (uint, *apperror.AppError) {
	args := m.Called(ctx, productID)
	if v := args.Get(1); v != nil {
		return 0, v.(*apperror.AppError)
	}
	return args.Get(0).(uint), nil
}

func (m *ReservationRepository) FindForUpdate(ctx context.Context, reservationID string) (*entity.Reservation, *apperror.AppError) {
	args := m.Called(ctx, reservationID)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).(*entity.Reservation), nil
}

func (m *ReservationRepository) FilterWithStatusOrder(ctx context.Context, filter *entity.Reservation, limit, offset int, isAsc bool) ([]entity.Reservation, *apperror.AppError) {
	args := m.Called(ctx, filter, limit, offset, isAsc)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).([]entity.Reservation), nil
}
