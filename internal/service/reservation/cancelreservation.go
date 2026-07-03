package reservation

import (
	"context"
	"errors"

	"github.com/Yoshikrit/reservation/internal/entity"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
	"github.com/Yoshikrit/reservation/internal/service/constant"
)

func (s *reservationService) CancelReservation(ctx context.Context, reservationID string) *apperror.AppError {
	if err := s.trManager.Do(ctx, func(ctx context.Context) error {
		reservation, appErr := s.reservationRepo.FindForUpdate(ctx, reservationID)
		if appErr != nil {
			return appErr
		}

		if reservation.Status != constant.StatusHeld {
			return apperror.NewError(42200000, errors.New("only HELD reservations can be cancelled"))
		}

		if appErr := s.reservationRepo.Patch(ctx, &entity.Reservation{
			ReservationID: reservationID,
			Status:        constant.StatusCancelled,
		}); appErr != nil {
			return appErr
		}
		return nil
	}); err != nil {
		var appErr *apperror.AppError
		if errors.As(err, &appErr) {
			return appErr
		}
		return apperror.NewError(50000000, err)
	}
	return nil
}
