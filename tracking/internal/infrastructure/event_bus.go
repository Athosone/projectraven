package infrastructure

import (
	"context"

	"github.com/athosone/projectraven/tracking/internal/domain"
	"github.com/nats-io/nats.go"
)

var jetstreamCtx nats.JetStreamContext

func NewEventPublisher(streamName string, jsCtx nats.JetStreamContext) (domain.EventPublisher, error) {
	return func(ctx context.Context, topic string, msg []byte) error {
		_, err := jsCtx.Publish(topic, msg)
		return err
	}, nil
}
