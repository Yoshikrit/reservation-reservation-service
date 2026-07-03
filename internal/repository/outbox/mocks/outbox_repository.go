package mocks

import (
	"context"

	"reservation/internal/entity"
	"reservation/internal/pkg/apperror"

	"github.com/stretchr/testify/mock"
)

type OutboxRepository struct {
	mock.Mock
}

func (m *OutboxRepository) Create(ctx context.Context, outbox *entity.Outbox) *apperror.AppError {
	args := m.Called(ctx, outbox)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *OutboxRepository) FindPending(ctx context.Context, limit int) ([]entity.Outbox, *apperror.AppError) {
	args := m.Called(ctx, limit)
	if v := args.Get(1); v != nil {
		return nil, v.(*apperror.AppError)
	}
	return args.Get(0).([]entity.Outbox), nil
}

func (m *OutboxRepository) UpdateStatus(ctx context.Context, eventID string, status string) *apperror.AppError {
	args := m.Called(ctx, eventID, status)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}

func (m *OutboxRepository) IncrRetryCount(ctx context.Context, eventID string) *apperror.AppError {
	args := m.Called(ctx, eventID)
	if v := args.Get(0); v != nil {
		return v.(*apperror.AppError)
	}
	return nil
}
