package cfgpromotion

import (
	"context"
	"fmt"

	"reservation/internal/entity"
	"reservation/internal/pkg/apperror"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
)

func (r *cfgPromotionDateRepository) Filter(ctx context.Context, filter *entity.CfgPromotionDate, limit, offset int, isAsc bool) ([]entity.CfgPromotionDate, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	query := db.WithContext(ctx).Model(&entity.CfgPromotionDate{})
	if filter != nil {
		query = query.Where(filter)
	}

	if limit < 0 || offset < 0 {
		return nil, apperror.NewError(40000000, fmt.Errorf("limit and offset must be >= 0"))
	}

	direction := "DESC"
	if isAsc {
		direction = "ASC"
	}
	query = query.Order("reservation_id " + direction)

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	var results []entity.CfgPromotionDate
	if err := query.Find(&results).Error; err != nil {
		return nil, apperror.NewError(50000000, err)
	}
	return results, nil
}
