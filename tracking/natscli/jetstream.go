package natscli

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type NatsConfig struct {
	URL        string
	StreamName string
	Subjects   []string
}

var JetstreamCtx nats.JetStreamContext

func InitClient(cfg *NatsConfig) error {
	nc, err := nats.Connect(cfg.URL)

	if err != nil {
		return fmt.Errorf("error connecting to NATS: %w", err)
	}
	jsCtx, err := nc.JetStream()

	if err != nil {
		return fmt.Errorf("error connecting to NATS JetStream: %w", err)
	}
	JetstreamCtx = jsCtx
	return nil
}

func Configure(cfg *NatsConfig) error {
	// check if stream exists
	_, err := JetstreamCtx.StreamInfo(cfg.StreamName)
	if err == nil {
		_, err = JetstreamCtx.UpdateStream(&nats.StreamConfig{
			Name:     cfg.StreamName,
			Subjects: cfg.Subjects,
		})
		return err
	}

	// create stream
	// Eventually replace with NACK controller and crds
	_, err = JetstreamCtx.AddStream(&nats.StreamConfig{
		Name:      cfg.StreamName,
		Subjects:  cfg.Subjects,
		Retention: nats.LimitsPolicy,
		MaxBytes:  1024 * 1024 * 1024 * 10,
	})

	if err != nil {
		return fmt.Errorf("error creating stream: %w", err)
	}
	return nil
}
