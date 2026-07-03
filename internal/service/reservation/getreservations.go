package reservation

import (
	"context"
	"math"

	"github.com/Yoshikrit/reservation/internal/entity"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
)

func (s *reservationService) GetReservations(ctx context.Context, body *ListReservationBody) (*ListReservationResponse, *apperror.AppError) {
	return s.getFilteredReservations(ctx, body)
}

func (s *reservationService) getFilteredReservations(ctx context.Context, body *ListReservationBody) (*ListReservationResponse, *apperror.AppError) {
	limit := int64(DefaultLimit)
	if body.Limit > 0 {
		if body.Limit > MaxLimit {
			limit = MaxLimit
		} else {
			limit = body.Limit
		}
	}

	offset := int64(0)
	if body.Offset > 0 {
		offset = body.Offset
	}

	filter := &entity.Reservation{
		ProductID: body.ProductID,
		Status:    body.Status,
	}

	total, appErr := s.reservationRepo.Count(ctx, filter)
	if appErr != nil {
		return nil, appErr
	}

	items, appErr := s.reservationRepo.FilterWithStatusOrder(ctx, filter, int(limit), int(offset), true)
	if appErr != nil {
		return nil, appErr
	}

	hasMore := offset+int64(len(items)) < total
	currentPage := int64(1)
	if limit > 0 {
		currentPage = offset/limit + 1
	}

	return &ListReservationResponse{
		Reservations: mapReservationItems(items),
		Pagination: PaginationInfo{
			Limit:       limit,
			Offset:      offset,
			Count:       int64(len(items)),
			HasMore:     hasMore,
			CurrentPage: currentPage,
			PerPage:     limit,
		},
	}, nil
}

func mapReservationItems(reservations []entity.Reservation) []ReservationItem {
	items := make([]ReservationItem, len(reservations))
	for i, r := range reservations {
		items[i] = ReservationItem{
			ReservationID: r.ReservationID,
			ProductID:     r.ProductID,
			Qty:           r.Quantity,
			Status:        r.Status,
			ExpiresAt:     r.ExpiresAt,
			BasePrice:     r.BasePrice,
			DiscountRate:  int(math.Round(r.DiscountRate * 100)),
			Price:         r.Price,
		}
	}
	return items
}
