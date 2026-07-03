package outbox

import (
	"context"

	"github.com/Yoshikrit/reservation/internal/entity"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"

	"gorm.io/gorm"
)

const (
	StatusPending   = "PENDING"
	StatusPublished = "PUBLISHED"
	StatusFailed    = "FAILED"

	MaxRetry = 3
)

type OutboxRepository interface {
	Create(ctx context.Context, outbox *entity.Outbox) *apperror.AppError
	FindPending(ctx context.Context, limit int) ([]entity.Outbox, *apperror.AppError)
	UpdateStatus(ctx context.Context, eventID string, status string) *apperror.AppError
	IncrRetryCount(ctx context.Context, eventID string) *apperror.AppError
}

type outboxRepository struct {
	db *gorm.DB
}

func NewOutboxRepository(db *gorm.DB) OutboxRepository {
	return &outboxRepository{db}
}
