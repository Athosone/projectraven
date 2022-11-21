package domainDevice

type DeviceRepository interface {
	FindById(id string) (*Device, error)
	Save(device *Device) error
}
