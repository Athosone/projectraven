package infrastructure

import (
	"context"

	domainDevice "github.com/athosone/projectraven/tracking/internal/domain/device"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type deviceRepository struct {
	collection *mongo.Collection
}

func NewDeviceRepository(database *mongo.Database) (domainDevice.DeviceRepository, error) {
	return &deviceRepository{database: database.Collection("device")}, nil
}

func (r *deviceRepository) FindById(ctx context.Context, id string) (*domainDevice.Device, error) {
	// use the database to find the device
	c := r.collection.FindOne(ctx, bson.M{"_id": id})
	if c.Err() != nil {
		return nil, c.Err()
	}
	return c.Decode()
}

func (r *deviceRepository) Save(device *domainDevice.Device) error {
	return nil
}
