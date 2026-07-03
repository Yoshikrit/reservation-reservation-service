package reservation

import (
	"context"
	"math"
	"time"

	"github.com/Yoshikrit/reservation/internal/entity"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
)

func (s *reservationService) isPromotionDay(ctx context.Context, t time.Time) (bool, *apperror.AppError) {
	list, appErr := s.cfgPromotionRepo.Filter(ctx, &entity.CfgPromotionDate{
		Date:      uint32(t.Day()),
		IsEnabled: true,
	}, 0, 0, true)
	if appErr != nil {
		return false, appErr
	}

	for _, cfg := range list {
		if cfg.EndDate.After(t) {
			return true, nil
		}
	}
	return false, nil
}

func (s *reservationService) calculateDiscountRate(ctx context.Context, qty uint, t time.Time) (int, *apperror.AppError) {
	isPromo, appErr := s.isPromotionDay(ctx, t)
	if appErr != nil {
		return 0, appErr
	}
	if !isPromo {
		return 0, nil
	}

	switch {
	case qty >= 301:
		return 15, nil
	case qty >= 101:
		return 10, nil
	case qty >= 10:
		return 5, nil
	default:
		return 0, nil
	}
}

// calculatePrice returns base_price * qty * (1 - discount_rate/100), rounded to 2 decimal places.
func calculatePrice(basePrice float64, qty uint, discountRate int) float64 {
	total := basePrice * float64(qty)
	if discountRate == 0 {
		return math.Round(total*100) / 100
	}
	return math.Round(total*(1.0-float64(discountRate)/100.0)*100) / 100
}
