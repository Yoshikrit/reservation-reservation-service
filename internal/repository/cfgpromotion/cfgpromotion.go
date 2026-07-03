package cfgpromotion

import (
	"context"

	"reservation/internal/entity"
	"reservation/internal/pkg/apperror"

	"gorm.io/gorm"
)

type CfgPromotionDateRepository interface {
	Filter(ctx context.Context, filter *entity.CfgPromotionDate, limit, offset int, isAsc bool) ([]entity.CfgPromotionDate, *apperror.AppError)
}

type cfgPromotionDateRepository struct {
	db *gorm.DB
}

func NewCfgPromotionDateRepository(db *gorm.DB) CfgPromotionDateRepository {
	return &cfgPromotionDateRepository{db}
}
