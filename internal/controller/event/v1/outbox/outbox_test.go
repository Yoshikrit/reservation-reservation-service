package outbox_test

import (
	"context"
	"errors"
	"testing"

	outboxCtrl "github.com/Yoshikrit/reservation/internal/controller/event/v1/outbox"
	kafkaMocks "github.com/Yoshikrit/reservation/internal/gateway/kafka/mocks"
	"github.com/Yoshikrit/reservation/internal/entity"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"
	outboxRepo "github.com/Yoshikrit/reservation/internal/repository/outbox"
	outboxMocks "github.com/Yoshikrit/reservation/internal/repository/outbox/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProcess_Success(t *testing.T) {
	repo := new(outboxMocks.OutboxRepository)
	pub := new(kafkaMocks.EventPublisher)

	repo.On("FindPending", mock.Anything, 100).Return([]entity.Outbox{
		{EventID: "evt-001", Topic: "reservation-confirmed", Payload: `{"product_id":"prod-001"}`, RetryCount: 0},
	}, nil)
	pub.On("Publish", mock.Anything, "reservation-confirmed", "evt-001", mock.Anything, mock.Anything).
		Return(nil)
	repo.On("UpdateStatus", mock.Anything, "evt-001", outboxRepo.StatusPublished).Return(nil)

	h := outboxCtrl.NewOutboxRelay(repo, pub)
	h.Process(context.Background())

	repo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestProcess_WithTraceID(t *testing.T) {
	repo := new(outboxMocks.OutboxRepository)
	pub := new(kafkaMocks.EventPublisher)

	repo.On("FindPending", mock.Anything, 100).Return([]entity.Outbox{
		{
			EventID: "evt-002", Topic: "reservation-confirmed", Payload: `{}`, RetryCount: 0,
			AuditModel: entity.AuditModel{CreatedByTraceID: "trace-abc"},
		},
	}, nil)
	pub.On("Publish", mock.Anything, "reservation-confirmed", "evt-002", mock.Anything,
		mock.MatchedBy(func(h map[string]string) bool {
			return h["X-Request-ID"] == "trace-abc"
		})).Return(nil)
	repo.On("UpdateStatus", mock.Anything, "evt-002", outboxRepo.StatusPublished).Return(nil)

	h := outboxCtrl.NewOutboxRelay(repo, pub)
	h.Process(context.Background())

	pub.AssertExpectations(t)
}

func TestProcess_FindPendingError(t *testing.T) {
	repo := new(outboxMocks.OutboxRepository)
	pub := new(kafkaMocks.EventPublisher)

	repo.On("FindPending", mock.Anything, 100).
		Return(nil, apperror.NewError(50000000, errors.New("db error")))

	h := outboxCtrl.NewOutboxRelay(repo, pub)
	h.Process(context.Background())

	pub.AssertNotCalled(t, "Publish")
	repo.AssertNotCalled(t, "UpdateStatus")
}

func TestProcess_PublishFail_BelowMaxRetry(t *testing.T) {
	repo := new(outboxMocks.OutboxRepository)
	pub := new(kafkaMocks.EventPublisher)

	// RetryCount=0 → after this attempt: 1 < MaxRetry(3), so only IncrRetryCount, no UpdateStatus
	repo.On("FindPending", mock.Anything, 100).Return([]entity.Outbox{
		{EventID: "evt-003", Topic: "reservation-confirmed", Payload: `{}`, RetryCount: 0},
	}, nil)
	pub.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("kafka timeout"))
	repo.On("IncrRetryCount", mock.Anything, "evt-003").Return(nil)

	h := outboxCtrl.NewOutboxRelay(repo, pub)
	h.Process(context.Background())

	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "UpdateStatus")
}

func TestProcess_PublishFail_MaxRetryReached(t *testing.T) {
	repo := new(outboxMocks.OutboxRepository)
	pub := new(kafkaMocks.EventPublisher)

	// RetryCount=MaxRetry-1 (2) → after this attempt: 3 >= MaxRetry(3), mark FAILED
	repo.On("FindPending", mock.Anything, 100).Return([]entity.Outbox{
		{EventID: "evt-004", Topic: "reservation-confirmed", Payload: `{}`, RetryCount: outboxRepo.MaxRetry - 1},
	}, nil)
	pub.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("kafka timeout"))
	repo.On("IncrRetryCount", mock.Anything, "evt-004").Return(nil)
	repo.On("UpdateStatus", mock.Anything, "evt-004", outboxRepo.StatusFailed).Return(nil)

	h := outboxCtrl.NewOutboxRelay(repo, pub)
	h.Process(context.Background())

	repo.AssertExpectations(t)
	assert.True(t, true) // reached here means no panic
}
