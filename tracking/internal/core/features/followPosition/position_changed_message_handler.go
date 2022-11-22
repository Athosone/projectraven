package followposition

import (
	"context"
	"encoding/json"
	"fmt"

	mqttcli "github.com/athosone/projectraven/tracking/internal/api/mqtt"
)

type PositionChangedDeviceHandler struct {
	handler *SavePositionCommandHandler
}

func NewPositionChangedMessageHandler(handler *SavePositionCommandHandler) *PositionChangedDeviceHandler {
	return &PositionChangedDeviceHandler{handler: handler}
}

type positionChangedMessage struct {
	DeviceId  string  `json:"deviceId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timestamp int64   `json:"timestamp"`
}

func (m *PositionChangedDeviceHandler) HandleDevicePositionChanged(ctx context.Context, payload []byte, messageId string) error {
	var p positionChangedMessage
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("error unmarshaling message: %w", err)
	}
	handler := m.handler
	cmd := SavePositionCommand(p)
	if err := handler.Handle(cmd); err != nil {
		return fmt.Errorf("error handling command: %w", err)
	}
	return nil
}

func (m *PositionChangedDeviceHandler) SubscribeToTopic(ctx context.Context, server *mqttcli.MQTTServer) error {
	return server.Subscribe(ctx, "device.position.changed", m.HandleDevicePositionChanged)
}
