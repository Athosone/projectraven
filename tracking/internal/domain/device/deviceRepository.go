package domainDevice

import "context"

type DeviceRepository interface {
	FindById(ctx context.Context, id string) (*Device, error)
	Save(ctx context.Context, device *Device) error
}
