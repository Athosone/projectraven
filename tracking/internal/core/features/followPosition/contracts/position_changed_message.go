package contracts

import (
	"fmt"

	domainDevice "github.com/athosone/projectraven/tracking/internal/domain/device"
)

const MQTTTopicListened = domainDevice.RootDeviceTopic + ".position.changed"

func JetStreamDevicePositionChangedSubject(deviceId string) string {
	return fmt.Sprintf("%s.%s.position.changed", domainDevice.RootDeviceTopic, deviceId)
}

type PositionChangedMessage struct {
	MessageId string   `json:"message_id"`
	DeviceId  string   `json:"device_id"`
	Position  Position `json:"position"`
	Timestamp int64    `json:"timestamp"`
}

type Position struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
