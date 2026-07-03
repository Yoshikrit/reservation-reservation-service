package mocks

import (
	"context"

	"github.com/Yoshikrit/reservation/internal/entity"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"

	"github.com/stretchr/testify/mock"
)

type CfgPromotionDateRepository struct {
	mock.Mock
}

func (m *CfgPromotionDateRepository) Filter(ctx context.Context, filter *entity.CfgPromotionDate, limit, offset int, isAsc bool) ([]entity.CfgPromotionDate, *apperror.AppError) {
	args := m.Called(ctx, filter, limit, offset, isAsc)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).([]entity.CfgPromotionDate), nil
}
