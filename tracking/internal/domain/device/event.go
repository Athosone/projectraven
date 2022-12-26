package domainDevice

import "github.com/athosone/projectraven/tracking/internal/domain"

type DevicePositionChangedEvent struct {
	domain.DomainEvent
	DeviceId          string
	OldDevicePosition DevicePosition
	NewDevicePosition DevicePosition
}
