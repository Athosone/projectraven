package followposition

import (
	"context"

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
}

// TODO: Add jetstream event publisher
func NewSavePositionCommandHandler(deviceRepository domainDevice.DeviceRepository) (*SavePositionCommandHandler, error) {
	return &SavePositionCommandHandler{deviceRepository: deviceRepository}, nil
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
	return h.deviceRepository.CreateOrUpdate(ctx, device)
}
