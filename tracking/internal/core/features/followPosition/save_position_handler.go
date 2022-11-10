package followposition

import (
	domainDevice "github.com/athosone/projectraven/tracking/internal/domain/device"
)

type SavePositionCommand struct {
	DeviceId  string
	Latitude  float64
	Longitude float64
}

type SavePositionCommandHandler struct {
	deviceRepository domainDevice.DeviceRepository
}

func NewSavePositionCommandHandler(deviceRepository domainDevice.DeviceRepository) (*SavePositionCommandHandler, error) {
  return &SavePositionCommandHandler{deviceRepository: deviceRepository}, nil
}

func (h *SavePositionCommandHandler) Handle(command SavePositionCommand) error {
  device, err := h.deviceRepository.FindById(command.DeviceId)
  if err != nil {
    return err
  }
  device.UpdatePosition(command.Latitude, command.Longitude)
  return h.deviceRepository.Save(device)
}
