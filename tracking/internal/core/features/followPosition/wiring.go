package followposition

import (
	"context"

	api "github.com/athosone/projectraven/tracking/internal/api/mqtt"
	"github.com/athosone/projectraven/tracking/internal/config"
	domainDevice "github.com/athosone/projectraven/tracking/internal/domain/device"
)

func WireFollowPositionFeature(ctx context.Context, cfg config.AppConfig,
	deviceRepository domainDevice.DeviceRepository,
	server *api.MQTTServer) error {
	savePositionHandler, err := NewSavePositionCommandHandler(deviceRepository)
	if err != nil {
		return err
	}
	p := &PositionChangedDeviceHandler{handler: savePositionHandler}
	return server.Subscribe(ctx, "device.position.changed", p.HandleDevicePositionChanged)
}
