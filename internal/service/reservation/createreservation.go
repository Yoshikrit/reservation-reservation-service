package reservation

import (
	"context"
	"errors"
	"time"

	"github.com/Yoshikrit/reservation/internal/entity"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
	"github.com/Yoshikrit/reservation/internal/service/constant"
	"github.com/Yoshikrit/reservation/internal/service/utility"
	inventoryGateway "github.com/Yoshikrit/reservation/internal/gateway/grpc/inventory"
)

func (s *reservationService) CreateReservation(ctx context.Context, request *CreateReservationRequest) (appErr *apperror.AppError) {
	lockKey := "lock:rsv:create:" + request.ProductID
	var acquired bool
	for range 5 {
		var err error
		acquired, err = s.redis.SetNX(ctx, lockKey, "1", 10*time.Second).Result()
		if err != nil {
			return apperror.NewError(50000000, err)
		}
		if acquired {
			break
		}
		select {
		case <-ctx.Done():
			return apperror.NewError(50000000, ctx.Err())
		case <-time.After(300 * time.Millisecond):
		}
	}
	if !acquired {
		return apperror.NewError(42900000, errors.New("product is being processed, please retry"))
	}
	defer s.redis.Del(context.WithoutCancel(ctx), lockKey)

	sumHeldQuantity, appErr := s.reservationRepo.SumReservedQuantity(ctx, request.ProductID)
	if appErr != nil {
		return appErr
	}

	checkAndHoldResult, appErr := s.inventoryGateway.CheckAndHold(ctx, inventoryGateway.CheckAndHoldRequest{
		ProductID: request.ProductID,
		Quantity:  uint64(sumHeldQuantity + request.Quantity),
	})
	if appErr != nil {
		return appErr
	}

	if !checkAndHoldResult.Available {
		return apperror.NewError(42200000, errors.New("insufficient stock for product "+request.ProductID))
	}

	now := time.Now()
	discountRate, appErr := s.calculateDiscountRate(ctx, request.Quantity, now)
	if appErr != nil {
		return appErr
	}

	reserveID, appErr := utility.CreateUUID("rsv")
	if appErr != nil {
		return appErr
	}

	reservation := &entity.Reservation{
		ReservationID: reserveID,
		ProductID:     request.ProductID,
		Quantity:      request.Quantity,
		Status:        constant.StatusHeld,
		BasePrice:     checkAndHoldResult.Price,
		DiscountRate:  float64(discountRate) / 100.0,
		Price:         calculatePrice(checkAndHoldResult.Price, request.Quantity, discountRate),
		ExpiresAt:     now.Add(time.Duration(request.TtlSecond) * time.Second),
	}

	if appErr = s.reservationRepo.Create(ctx, reservation); appErr != nil {
		return appErr
	}

	defer func() {
		if appErr != nil {
			_ = s.reservationRepo.Patch(context.WithoutCancel(ctx), &entity.Reservation{
				ReservationID: reservation.ReservationID,
				Status:        constant.StatusCancelled,
			})
		}
	}()

	return nil
}
