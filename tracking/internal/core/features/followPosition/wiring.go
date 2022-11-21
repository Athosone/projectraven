package followposition

import (
	"context"

	api "github.com/athosone/projectraven/tracking/internal/api/mqtt"
	"github.com/athosone/projectraven/tracking/internal/config"
	"github.com/athosone/projectraven/tracking/internal/infrastructure"
	"github.com/athosone/projectraven/tracking/mongodb"
)

func WireFollowPositionFeature(ctx context.Context, cfg config.AppConfig, server *api.MQTTServer) error {
	deviceRepository, err := infrastructure.NewDeviceRepository(ctx, mongodb.Database)
	if err != nil {
		return err
	}
	savePositionHandler, err := NewSavePositionCommandHandler(deviceRepository)
	if err != nil {
		return err
	}
	p := &PositionChangedDeviceHandler{handler: savePositionHandler}
	return server.Subscribe(ctx, "device.position.changed", p.HandleDevicePositionChanged)
}
