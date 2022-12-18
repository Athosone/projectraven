package main

import (
	"context"
	"fmt"
	"time"

	mqttcli "github.com/athosone/projectraven/tracking/internal/api/mqtt"
	"github.com/athosone/projectraven/tracking/internal/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/fx"
)

type MQTTMessageListener interface {
	SubscribeToTopic(ctx context.Context, server *mqttcli.MQTTServer) error
}

func subscribeListeners(messageListeners []MQTTMessageListener, srv *mqttcli.MQTTServer) (mqtt.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, listener := range messageListeners {
		if err := listener.SubscribeToTopic(ctx, srv); err != nil {
			return nil, fmt.Errorf("error subscribing to topic: %w", err)
		}
	}

	return srv.Client(), nil
}

func newMQTTListener(cfg *config.AppConfig) (*mqttcli.MQTTServer, error) {
	srv, err := mqttcli.NewMQTTServer(logger, mqttcli.MQTTServerConfig{Broker: cfg.MQTT.Broker, ClientID: cfg.MQTT.ClientID})
	if err != nil {
		return nil, fmt.Errorf("error creating MQTT client: %w", err)
	}
	if !srv.Client().IsConnected() {
		if token := srv.Client().Connect(); token.Wait() && token.Error() != nil {
			err := token.Error()
			return nil, fmt.Errorf("error connecting to MQTT broker: %w", err)
		}
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
