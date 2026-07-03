package event

import (
	"context"

	outboxRelay "github.com/Yoshikrit/reservation/internal/controller/event/v1/outbox"
	kafkapkg "github.com/Yoshikrit/reservation/internal/gateway/kafka"
	outboxRepo "github.com/Yoshikrit/reservation/internal/repository/outbox"
	"github.com/Yoshikrit/reservation/internal/entity"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func StartEvent(ctx context.Context, db *gorm.DB, producer kafkapkg.EventPublisher) {
	if producer == nil {
		log.Warn().Msg("event: kafka producer not configured, relay disabled")
		return
	}

	relays := []relayConfig{
		{
			name:    entity.Outbox{}.TableName(),
			handler: outboxRelay.NewOutboxRelay(outboxRepo.NewOutboxRepository(db), producer),
		},
	}

	for _, r := range relays {
		startRelay(ctx, r)
	}
}
