package event

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type RelayHandler interface {
	Process(ctx context.Context)
}

type relayConfig struct {
	name    string
	handler RelayHandler
}

func startRelay(ctx context.Context, cfg relayConfig) {
	log.Info().Str("relay", cfg.name).Msg("event: relay started")
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				log.Info().Str("relay", cfg.name).Msg("event: relay stopped")
				return
			case <-ticker.C:
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Error().
								Str("relay", cfg.name).
								Str("panic", fmt.Sprintf("%v", r)).
								Msg("event: relay panic recovered")
						}
					}()
					cfg.handler.Process(ctx)
				}()
			}
		}
	}()
}
