package kafka

import (
	"strings"

	"github.com/rs/zerolog/log"
)

// GatewayName is a typed key for Kafka topic registration.
// Call .Topic() to resolve the topic string from the registry.
type GatewayName string

func (n GatewayName) Topic() string {
	if current == nil {
		return ""
	}
	return current.topics[n]
}

const (
	KafkaConfirmedReservationEvent GatewayName = "kafka_confirmed_reservation_event"
)

type Config struct {
	Brokers                   string
	TopicConfirmedReservation string
}

type gatewayInstance struct {
	topics map[GatewayName]string
}

var current *gatewayInstance

// InitGateway builds the topic registry and initialises the Kafka producer.
// Call once at startup before serving any requests.
func InitGateway(cfg Config) EventPublisher {
	current = &gatewayInstance{
		topics: map[GatewayName]string{
			KafkaConfirmedReservationEvent: cfg.TopicConfirmedReservation,
		},
	}

	for name, topic := range current.topics {
		if topic == "" {
			log.Fatal().Str("gateway_name", string(name)).Msg("gateway: topic is empty — check KAFKA_TOPIC_* env vars")
		}
	}

	if cfg.Brokers == "" {
		return nil
	}
	return NewKafkaProducer(strings.Split(cfg.Brokers, ","))
}

// InitForTesting sets up the topic registry for unit tests without a real Kafka connection.
func InitForTesting(topics map[GatewayName]string) {
	current = &gatewayInstance{topics: topics}
}
