package config

import (
	kafkapkg "github.com/Yoshikrit/reservation/internal/gateway/kafka"
)

type KafkaConfig struct {
	Brokers                   string `env:"KAFKA_BROKERS,required"`
	TopicConfirmedReservation string `env:"KAFKA_TOPIC_CONFIRMED_RESERVATION,required"`
}

// InitKafkaProducer initialises the topic registry and returns the Kafka producer.
// Returns nil if KAFKA_BROKERS is empty (disables the event relay).
func InitKafkaProducer(cfg KafkaConfig) kafkapkg.EventPublisher {
	return kafkapkg.InitGateway(kafkapkg.Config{
		Brokers:                   cfg.Brokers,
		TopicConfirmedReservation: cfg.TopicConfirmedReservation,
	})
}
