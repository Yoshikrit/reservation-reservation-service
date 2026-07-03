package reservation

import (
	"context"
	"errors"
	"fmt"

	"reservation/internal/entity"
	"reservation/internal/pkg/apperror"
	"reservation/internal/service/constant"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (r *reservationRepository) FindForUpdate(ctx context.Context, reservationID string) (*entity.Reservation, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	var reservation entity.Reservation
	if err := db.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("reservation_id = ?", reservationID).
		First(&reservation).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NewError(40400000, err, "reservation", reservationID)
		}
		return nil, apperror.NewError(50000000, err)
	}
	return &reservation, nil
}

func (r *reservationRepository) FilterWithStatusOrder(ctx context.Context, filter *entity.Reservation, limit, offset int, isAsc bool) ([]entity.Reservation, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	query := db.WithContext(ctx).Model(&entity.Reservation{})
	if filter != nil {
		query = query.Where(filter)
	}

	isNoFilter := filter == nil || (filter.ProductID == "" && filter.Status == "")
	if isNoFilter {
		query = query.Order(fmt.Sprintf(
			"CASE status WHEN '%s' THEN 1 WHEN '%s' THEN 2 WHEN '%s' THEN 3 ELSE 4 END ASC, updated_at ASC",
			constant.StatusConfirmed, constant.StatusHeld, constant.StatusCancelled,
		))
	} else {
		direction := "DESC"
		if isAsc {
			direction = "ASC"
		}
		query = query.Order("updated_at " + direction)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	var reservations []entity.Reservation
	if err := query.Find(&reservations).Error; err != nil {
		return nil, apperror.NewError(50000000, err)
	}
	return reservations, nil
}

func (r *reservationRepository) SumReservedQuantity(ctx context.Context, productID string) (uint, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)

	var sum uint
	err := db.WithContext(ctx).
		Model(&entity.Reservation{}).
		Where("product_id = ? AND status = ? AND expires_at > NOW()", productID, constant.StatusHeld).
		Select("COALESCE(SUM(quantity), 0)").
		Scan(&sum).Error
	if err != nil {
		return 0, apperror.NewError(50000000, err)
	}
	return sum, nil
}
