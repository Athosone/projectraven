package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/athosone/projectraven/tracking/mqttcli"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type MessageHandler func(ctx context.Context, payload []byte, messageId string) error

type MQTTServerConfig struct {
	Broker   string
	ClientID string
}

type MQTTServer struct {
	client mqtt.Client
	rwLock sync.RWMutex
	logger *zap.SugaredLogger

	handlers map[string][]MessageHandler
}

// TODO: Use this module to refacto the http one, also make sure it follows vertical slice
func NewMQTTServer(logger *zap.SugaredLogger, cfg MQTTServerConfig) (*MQTTServer, error) {
	cli, err := mqttcli.NewClient(cfg.Broker, cfg.ClientID)
	if err != nil {
		return nil, fmt.Errorf("error creating MQTT client: %w", err)
	}

	listener := &MQTTServer{client: cli, handlers: make(map[string][]MessageHandler), logger: logger.With("component", "mqtt")}
	return listener, nil
}

func (l *MQTTServer) Subscribe(ctx context.Context, topic string, handler MessageHandler) error {
	l.rwLock.Lock()
	l.handlers[topic] = append(l.handlers[topic], handler)
	l.rwLock.Unlock()
	
  token := l.client.Subscribe(topic, 1, func(client mqtt.Client, msg mqtt.Message) {
		if err := handler(ctx, msg.Payload(), fmt.Sprint(msg.MessageID())); err != nil {
			fmt.Printf("Error handling message: %v", err)
			return
		}
		msg.Ack()
	})
	token.Wait()
	if token.Error() != nil {
		return fmt.Errorf("error subscribing to topic %s: %w", topic, token.Error())
	}
	return nil
}
