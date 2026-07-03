package outbox

import (
	"context"

	"github.com/Yoshikrit/reservation/internal/entity"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"

	gormtrm "github.com/avito-tech/go-transaction-manager/drivers/gorm/v2"
	"gorm.io/gorm"
)

func (r *outboxRepository) Create(ctx context.Context, outbox *entity.Outbox) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).Create(outbox).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *outboxRepository) FindPending(ctx context.Context, limit int) ([]entity.Outbox, *apperror.AppError) {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	var outboxes []entity.Outbox
	err := db.WithContext(ctx).
		Where("status = ? AND retry_count < ? AND (publish_at IS NULL OR publish_at <= NOW())", StatusPending, MaxRetry).
		Order("created_at ASC").
		Limit(limit).
		Find(&outboxes).Error
	if err != nil {
		return nil, apperror.NewError(50000000, err)
	}
	return outboxes, nil
}

func (r *outboxRepository) UpdateStatus(ctx context.Context, eventID string, status string) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).
		Model(&entity.Outbox{}).
		Where("event_id = ?", eventID).
		Update("status", status).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}

func (r *outboxRepository) IncrRetryCount(ctx context.Context, eventID string) *apperror.AppError {
	db := gormtrm.DefaultCtxGetter.DefaultTrOrDB(ctx, r.db)
	if err := db.WithContext(ctx).
		Model(&entity.Outbox{}).
		Where("event_id = ?", eventID).
		UpdateColumn("retry_count", gorm.Expr("retry_count + 1")).Error; err != nil {
		return apperror.NewError(50000000, err)
	}
	return nil
}
