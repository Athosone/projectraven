package domainDevice

import (
	"time"

	"github.com/athosone/projectraven/tracking/internal/domain"
	"github.com/google/uuid"
)

type Device struct {
	ID                string
	Name              string
	Position          DevicePosition
	UncommittedEvents []any
}

type DevicePosition struct {
	Latitude  float64
	Longitude float64
}

func (d *Device) UpdatePosition(latitude float64, longitude float64) {
	newPos := DevicePosition{Latitude: latitude, Longitude: longitude}
	oldPos := d.Position

	d.Position = newPos
	d.UncommittedEvents = append(d.UncommittedEvents, newDevicePositionChanged(*d, oldPos, newPos))
}

func newDevicePositionChanged(d Device, oldPosition, newPosition DevicePosition) any {
	return DevicePositionChangedEvent{
		DomainEvent: domain.DomainEvent{
			EventId:      uuid.New(),
			CreatedAtUTC: time.Now().UTC(),
		},
		DeviceId:          d.ID,
		OldDevicePosition: oldPosition,
		NewDevicePosition: newPosition,
	}
}
