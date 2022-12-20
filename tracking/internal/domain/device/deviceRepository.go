package domainDevice

import "context"

type DeviceRepository interface {
	FindById(ctx context.Context, id string) (*Device, error)
	CreateOrUpdate(ctx context.Context, device *Device) error
}
