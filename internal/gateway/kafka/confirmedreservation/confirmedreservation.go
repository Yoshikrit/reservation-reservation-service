package confirmedreservation

import (
	"context"
	"reservation/internal/pkg/json"

	"reservation/internal/entity"
	kafkapkg "reservation/internal/gateway/kafka"
)

type ConfirmedReservationEvent struct {
	ProductID string `json:"product_id"`
	Quantity  uint   `json:"quantity"`
}

type Publisher interface {
	Publish(ctx context.Context, event ConfirmedReservationEvent) error
}

type publisher struct {
	producer kafkapkg.EventPublisher
}

func NewPublisher(producer kafkapkg.EventPublisher) Publisher {
	return publisher{producer: producer}
}

func (p publisher) Publish(ctx context.Context, event ConfirmedReservationEvent) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	headers := map[string]string{}
	if traceID, _ := ctx.Value(entity.ContextKeyTraceID).(string); traceID != "" {
		headers["X-Request-ID"] = traceID
	}
	return p.producer.Publish(ctx, kafkapkg.KafkaConfirmedReservationEvent.Topic(), event.ProductID, payload, headers)
}
