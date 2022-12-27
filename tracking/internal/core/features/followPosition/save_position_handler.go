package followposition

import (
	"context"
	"encoding/json"
	"fmt"

	domain "github.com/athosone/projectraven/tracking/internal/domain"
	domainDevice "github.com/athosone/projectraven/tracking/internal/domain/device"
)

type SavePositionCommand struct {
	DeviceId  string
	Latitude  float64
	Longitude float64
	Timestamp int64
}

type SavePositionCommandHandler struct {
	deviceRepository domainDevice.DeviceRepository
	eventPublisher   domain.EventPublisher
}

// TODO: Add jetstream event publisher
func NewSavePositionCommandHandler(deviceRepository domainDevice.DeviceRepository, publisher domain.EventPublisher) (*SavePositionCommandHandler, error) {
	return &SavePositionCommandHandler{deviceRepository: deviceRepository, eventPublisher: publisher}, nil
}

// TODO: Publish event using jetstream
func (h *SavePositionCommandHandler) Handle(ctx context.Context, command SavePositionCommand) error {
	device, err := h.deviceRepository.FindById(ctx, command.DeviceId)
	// if err not nil and err is not of type ErrDeviceNotFound, return err
	if err != nil && !domainDevice.IsErrDeviceNotFound(err) {
		return err
	}
	// if err is of type ErrDeviceNotFound, create a new device
	if domainDevice.IsErrDeviceNotFound(err) {
		device = &domainDevice.Device{
			ID: command.DeviceId,
		}
	}
	// update the position of the device and save it
	device.UpdatePosition(command.Latitude, command.Longitude)
	ue := device.UncommittedEvents
	for _, e := range ue {
		topic, msg := formatMsg(domainDevice.RootDeviceTopic, device.ID, e)
		if err := h.eventPublisher(ctx, topic, msg); err != nil {
			return err
		}
	}
	// TODO: use outbox pattern

	return h.deviceRepository.CreateOrUpdate(ctx, device)
}

func formatMsg(rootTopic string, deviceId string, event any) (string, []byte) {
	topic := fmt.Sprintf("%s.%s.positionChanged", rootTopic, deviceId)
	data, _ := json.Marshal(event)
	return topic, data
}
