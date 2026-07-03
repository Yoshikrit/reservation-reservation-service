package kafka

import (
	"context"

	kafkago "github.com/segmentio/kafka-go"
)

type EventPublisher interface {
	Publish(ctx context.Context, topic, key string, payload []byte, headers map[string]string) error
	Close() error
}

type kafkaProducer struct {
	writer *kafkago.Writer
}

func NewKafkaProducer(brokers []string) EventPublisher {
	return &kafkaProducer{
		writer: &kafkago.Writer{
			Addr:         kafkago.TCP(brokers...),
			Balancer:     &kafkago.LeastBytes{},
			RequiredAcks: kafkago.RequireAll,
		},
	}
}

func (p *kafkaProducer) Publish(ctx context.Context, topic, key string, payload []byte, headers map[string]string) error {
	msg := kafkago.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: payload,
	}
	for k, v := range headers {
		msg.Headers = append(msg.Headers, kafkago.Header{Key: k, Value: []byte(v)})
	}
	return p.writer.WriteMessages(ctx, msg)
}

func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}
