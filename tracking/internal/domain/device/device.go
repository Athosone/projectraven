package domainDevice

type Device struct {
	ID       string
	Name     string
	Position DevicePosition
}

type DevicePosition struct {
	Latitude  float64
	Longitude float64
}

func (d *Device) UpdatePosition(latitude float64, longitude float64) {
	d.Position = DevicePosition{Latitude: latitude, Longitude: longitude}
}
