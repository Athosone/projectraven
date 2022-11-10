package followposition

import (
	"context"
	"encoding/json"
	"fmt"
)

type PositionChangedDeviceHandler struct {
	handler *SavePositionCommandHandler
}

type positionChangedMessage struct {
	DeviceId  string  `json:"deviceId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
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
