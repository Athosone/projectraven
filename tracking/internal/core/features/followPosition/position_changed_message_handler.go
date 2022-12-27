package followposition

import (
	"context"
	"encoding/json"
	"fmt"

	mqttcli "github.com/athosone/projectraven/tracking/internal/api/mqtt"
	"github.com/athosone/projectraven/tracking/internal/core/features/followPosition/contracts"
)

type MQTTDevicePositionChangedHandler struct {
	handler *SavePositionCommandHandler
}

func NewPositionChangedMessageHandler(topic string, handler *SavePositionCommandHandler) *MQTTDevicePositionChangedHandler {
	return &MQTTDevicePositionChangedHandler{handler: handler}
}

func (m *MQTTDevicePositionChangedHandler) HandleDevicePositionChanged(ctx context.Context, payload []byte, messageId string) error {
	var p contracts.PositionChangedMessage
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("error unmarshaling message: %w", err)
	}
	handler := m.handler
	cmd := SavePositionCommand{
		DeviceId:  p.DeviceId,
		Latitude:  p.Position.Lat,
		Longitude: p.Position.Long,
		Timestamp: p.Timestamp,
	}
	if err := handler.Handle(ctx, cmd); err != nil {
		return fmt.Errorf("error handling command: %w", err)
	}
	return nil
}

func (m *MQTTDevicePositionChangedHandler) SubscribeToTopic(ctx context.Context, server *mqttcli.MQTTServer) error {
	return server.Subscribe(ctx, contracts.MQTTTopicListened, m.HandleDevicePositionChanged)
}
