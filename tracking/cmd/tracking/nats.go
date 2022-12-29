package main

import (
	"context"

	"github.com/athosone/projectraven/tracking/internal/config"
	"github.com/athosone/projectraven/tracking/internal/infrastructure"
	"github.com/athosone/projectraven/tracking/natscli"
	"github.com/nats-io/nats.go"
	"go.uber.org/fx"
)

func newNats(cfg *config.AppConfig, lc fx.Lifecycle) (nats.JetStreamContext, error) {
	natsCfg := natscli.NatsConfig{
		URL:        cfg.Nats.URL,
		StreamName: cfg.Nats.StreamName,
		Subjects:   infrastructure.NatsSubjects(),
	}
	err := natscli.InitClient(&natsCfg)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return natscli.Configure(&natsCfg)
		},
	})

	return natscli.JetstreamCtx, nil
}
