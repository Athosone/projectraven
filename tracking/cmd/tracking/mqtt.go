package main

import (
	"context"
	"fmt"

	mqttcli "github.com/athosone/projectraven/tracking/internal/api/mqtt"
	"github.com/athosone/projectraven/tracking/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/fx"
)

type MQTTMessageListener interface {
	SubscribeToTopic(ctx context.Context, server *mqttcli.MQTTServer) error
}

func subscribeListeners(messageListeners []MQTTMessageListener, srv *mqttcli.MQTTServer, lc fx.Lifecycle) (mqtt.Client, error) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if token := srv.Client().Connect(); token.Wait() && token.Error() != nil {
				err := token.Error()
				return fmt.Errorf("error connecting to MQTT broker: %w", err)
			}

			for _, listener := range messageListeners {
				if err := listener.SubscribeToTopic(ctx, srv); err != nil {
					return fmt.Errorf("error subscribing to topic: %w", err)
				}
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			srv.Client().Disconnect(2000)
			return nil
		},
	})

	return srv.Client(), nil
}

func newMQTTListener(cfg *config.AppConfig) (*mqttcli.MQTTServer, error) {
	srv, err := mqttcli.NewMQTTServer(logger, mqttcli.MQTTServerConfig{Broker: cfg.MQTT.Broker, ClientID: cfg.MQTT.ClientID})
	if err != nil {
		return nil, fmt.Errorf("error creating MQTT client: %w", err)
	}

	return srv, nil
}

func AsMQTTMessageListener(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(MQTTMessageListener)),
		fx.ResultTags(`group:"mqttListeners"`),
	)
}
