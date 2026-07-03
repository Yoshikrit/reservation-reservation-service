package reservation

import (
	"context"
	"errors"
	"reservation/internal/pkg/json"
	"time"

	"reservation/internal/entity"
	kafkapkg "reservation/internal/gateway/kafka"
	"reservation/internal/gateway/kafka/confirmedreservation"
	"reservation/internal/pkg/apperror"
	"reservation/internal/service/constant"

	"github.com/google/uuid"
)

func (s *reservationService) ConfirmReservation(ctx context.Context, reservationID string) *apperror.AppError {
	var expiredErr *apperror.AppError

	if err := s.trManager.Do(ctx, func(ctx context.Context) error {
		reservation, appErr := s.reservationRepo.FindForUpdate(ctx, reservationID)
		if appErr != nil {
			return appErr
		}

		if reservation.Status != constant.StatusHeld {
			return apperror.NewError(42200000, errors.New("only HELD reservations can be confirmed"))
		}

		if time.Now().After(reservation.ExpiresAt) {
			if patchErr := s.reservationRepo.Patch(ctx, &entity.Reservation{
				ReservationID: reservationID,
				Status:        constant.StatusCancelled,
			}); patchErr != nil {
				return patchErr
			}
			expiredErr = apperror.NewError(42200000, errors.New("reservation has expired"))
			return nil // commit the CANCELLED status
		}

		if patchErr := s.reservationRepo.Patch(ctx, &entity.Reservation{
			ReservationID: reservationID,
			Status:        constant.StatusConfirmed,
		}); patchErr != nil {
			return patchErr
		}

		event := confirmedreservation.ConfirmedReservationEvent{
			ProductID: reservation.ProductID,
			Quantity:  reservation.Quantity,
		}
		payload, err := json.Marshal(event)
		if err != nil {
			return apperror.NewError(50000000, err)
		}

		traceID, _ := ctx.Value(entity.ContextKeyTraceID).(string)
		if appErr := s.outboxRepo.Create(ctx, &entity.Outbox{
			EventID:   uuid.New().String(),
			Topic:     kafkapkg.KafkaConfirmedReservationEvent.Topic(),
			EventType: "ConfirmedReservationEvent",
			Payload:   string(payload),
			AuditModel: entity.AuditModel{
				CreatedByTraceID: traceID,
			},
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

	return expiredErr
}
