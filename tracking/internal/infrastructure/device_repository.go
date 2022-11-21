package infrastructure

import (
	"context"

	domainDevice "github.com/athosone/projectraven/tracking/internal/domain/device"
	"go.mongodb.org/mongo-driver/mongo"
)

type deviceRepository struct {
	database *mongo.Database
}

func NewDeviceRepository(ctx context.Context, database *mongo.Database) (domainDevice.DeviceRepository, error) {
	return &deviceRepository{database: database}, nil
}

func (r *deviceRepository) FindById(id string) (*domainDevice.Device, error) {
	return nil, nil
}

func (r *deviceRepository) Save(device *domainDevice.Device) error {
	return nil
}
