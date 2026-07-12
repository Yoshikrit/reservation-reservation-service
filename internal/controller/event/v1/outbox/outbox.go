package outbox

import (
	"context"

	kafkapkg "github.com/Yoshikrit/reservation/internal/gateway/kafka"
	outboxRepo "github.com/Yoshikrit/reservation/internal/repository/outbox"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

var tracer = otel.Tracer("outbox.relay")

type handler struct {
	repo     outboxRepo.OutboxRepository
	producer kafkapkg.EventPublisher
}

func NewOutboxRelay(repo outboxRepo.OutboxRepository, producer kafkapkg.EventPublisher) *handler {
	return &handler{repo: repo, producer: producer}
}

func (h *handler) Process(ctx context.Context) {
	records, err := h.repo.FindPending(ctx, 100)
	if err != nil {
		log.Error().Err(err).Msg("outbox: failed to find pending records")
		return
	}

	for _, rec := range records {
		recCtx, span := tracer.Start(ctx, "outbox.relay "+rec.Topic)

		headers := map[string]string{}
		if rec.CreatedByTraceID != "" {
			headers["X-Request-ID"] = rec.CreatedByTraceID
		}
		carrier := propagation.MapCarrier{}
		otel.GetTextMapPropagator().Inject(recCtx, carrier)
		for k, v := range carrier {
			headers[k] = v
		}

		if pubErr := h.producer.Publish(recCtx, rec.Topic, rec.EventID, []byte(rec.Payload), headers); pubErr != nil {
			span.End()
			if incrErr := h.repo.IncrRetryCount(ctx, rec.EventID); incrErr != nil {
				log.Error().Err(incrErr).Str("event_id", rec.EventID).Msg("outbox: failed to increment retry_count")
			}
			if rec.RetryCount+1 >= outboxRepo.MaxRetry {
				log.Error().
					Str("event_id", rec.EventID).
					Str("topic", rec.Topic).
					Int("retry_count", rec.RetryCount+1).
					Str("payload", rec.Payload).
					Msg("outbox: max retries reached, marking as dead-letter")
				if statusErr := h.repo.UpdateStatus(ctx, rec.EventID, outboxRepo.StatusFailed); statusErr != nil {
					log.Error().Err(statusErr).Str("event_id", rec.EventID).Msg("outbox: failed to update status to failed")
				}
			} else {
				log.Warn().
					Str("event_id", rec.EventID).
					Int("retry_count", rec.RetryCount+1).
					Int("max_retry", outboxRepo.MaxRetry).
					Err(pubErr).
					Msg("outbox: publish failed, will retry")
			}
			continue
		}

		span.End()
		log.Info().Str("event_id", rec.EventID).Str("topic", rec.Topic).Msg("outbox: event published")
		if statusErr := h.repo.UpdateStatus(ctx, rec.EventID, outboxRepo.StatusPublished); statusErr != nil {
			log.Error().Err(statusErr).Str("event_id", rec.EventID).Msg("outbox: failed to mark event as published")
		}
	}
}
