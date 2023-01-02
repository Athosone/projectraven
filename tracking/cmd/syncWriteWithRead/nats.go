package main

import (
	"github.com/athosone/projectraven/tracking/internal/config"
	"github.com/athosone/projectraven/tracking/internal/infrastructure"
	"github.com/athosone/projectraven/tracking/natscli"
	"github.com/nats-io/nats.go"
)

func newNats(cfg *config.AppConfig) (nats.JetStreamContext, error) {
	natsCfg := natscli.NatsConfig{
		URL:        cfg.Nats.URL,
		StreamName: cfg.Nats.StreamName,
		Subjects:   infrastructure.NatsSubjects(),
	}
	err := natscli.InitClient(&natsCfg)
	if err != nil {
		return nil, err
	}
	err = natscli.Configure(&natsCfg)
	if err != nil {
		return nil, err
	}
	jsctx := natscli.JetstreamCtx
	_, err = jsctx.AddConsumer(cfg.Nats.StreamName, &nats.ConsumerConfig{
		Durable:   durableName,
		AckPolicy: nats.AckExplicitPolicy,
	})
	if err != nil {
		// handle error if consumer already exists
		if err.Error() != "consumer already exists" {
			return nil, err
		}
	}

	return natscli.JetstreamCtx, nil
}
