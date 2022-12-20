package infrastructure

import (
	"context"

	domainDevice "github.com/athosone/projectraven/tracking/internal/domain/device"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type deviceDB struct {
	Id        string  `bson:"_id"`
	Name      string  `bson:"name"`
	Latitude  float64 `bson:"position_lat"`
	Longitude float64 `bson:"position_lon"`
}

type DevicePosition struct {
	Latitude  float64
	Longitude float64
}

type deviceRepository struct {
	collection *mongo.Collection
}

func NewDeviceRepository(database *mongo.Database) (domainDevice.DeviceRepository, error) {
	return &deviceRepository{collection: database.Collection("device")}, nil
}

func (r *deviceRepository) FindById(ctx context.Context, id string) (*domainDevice.Device, error) {
	// use the database to find the device
	c := r.collection.FindOne(ctx, bson.M{"_id": id})
	if c.Err() != nil {
		if c.Err() == mongo.ErrNoDocuments {
			return nil, &domainDevice.ErrDeviceNotFound{ID: id}
		}
		return nil, c.Err()
	}
	deviceDB := deviceDB{}
	err := c.Decode(&deviceDB)
	if err != nil {
		return nil, err
	}
	device := domainDevice.Device{
		ID:   deviceDB.Id,
		Name: deviceDB.Name,
		Position: domainDevice.DevicePosition{
			Latitude:  deviceDB.Latitude,
			Longitude: deviceDB.Longitude,
		},
	}
	return &device, nil
}

func (r *deviceRepository) CreateOrUpdate(ctx context.Context, device *domainDevice.Device) error {
	deviceDB := deviceDB{
		Id:        device.ID,
		Name:      device.Name,
		Latitude:  device.Position.Latitude,
		Longitude: device.Position.Longitude,
	}
	// insert or replace the device
	opt := options.Replace().SetUpsert(true)
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": device.ID}, deviceDB, opt)

	if err != nil {
		return err
	}
	return nil
}
