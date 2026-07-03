package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type EventPublisher struct {
	mock.Mock
}

func (m *EventPublisher) Publish(ctx context.Context, topic, key string, payload []byte, headers map[string]string) error {
	args := m.Called(ctx, topic, key, payload, headers)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (m *EventPublisher) Close() error {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}
