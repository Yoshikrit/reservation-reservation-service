package reservation

import (
	"context"

	"reservation/internal/entity"
	"reservation/internal/pkg/apperror"

	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(ctx context.Context, reservation *entity.Reservation) *apperror.AppError
	Patch(ctx context.Context, reservation *entity.Reservation) *apperror.AppError
	Find(ctx context.Context, filter *entity.Reservation) (*entity.Reservation, *apperror.AppError)
	Filter(ctx context.Context, filter *entity.Reservation, limit, offset int, isAsc bool) ([]entity.Reservation, *apperror.AppError)
	FindAll(ctx context.Context) ([]entity.Reservation, *apperror.AppError)
	Count(ctx context.Context, filter *entity.Reservation) (int64, *apperror.AppError)

	// custom
	SumReservedQuantity(ctx context.Context, productID string) (uint, *apperror.AppError)
	FindForUpdate(ctx context.Context, reservationID string) (*entity.Reservation, *apperror.AppError)
	FilterWithStatusOrder(ctx context.Context, filter *entity.Reservation, limit, offset int, isAsc bool) ([]entity.Reservation, *apperror.AppError)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db}
}
