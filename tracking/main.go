package main

import (
	"net/http"
	"os"

	"github.com/athosone/projectraven/tracking/internal/config"
	followposition "github.com/athosone/projectraven/tracking/internal/core/features/followPosition"
	manageusers "github.com/athosone/projectraven/tracking/internal/core/features/manageUsers"
	"github.com/athosone/projectraven/tracking/internal/infrastructure"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	logger *zap.SugaredLogger
)

func init() {
	var l *zap.Logger
	if os.Getenv("IS_DEBUG") == "true" {
		l, _ = zap.NewDevelopment()
	} else {
		l, _ = zap.NewProduction()
	}
	zap.ReplaceGlobals(l)
	logger = zap.S().With("service", "projectraven")
}

func main() {
	fx.New(
		// fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		// 	return &fxevent.ZapLogger{Logger: log}
		// }),
		fx.Provide(
			func() *zap.Logger { return logger.Desugar() },
			config.LoadConfig,
			// Command handlers
			followposition.NewSavePositionCommandHandler,

			// REST part
			newRestServer,
			// cf.: https://uber-go.github.io/fx/get-started/another-handler.html
			AsRoute(manageusers.NewRestUserHandler),
			fx.Annotate(
				newChi,
				fx.ParamTags(`group:"routes"`),
			),
			createServer,

			// MQTT part
			AsMQTTMessageListener(followposition.NewPositionChangedMessageHandler),
			newMQTTListener,
			fx.Annotate(
				subscribeListeners,
				fx.ParamTags(`group:"mqttListeners"`),
			),

			// Repositories
			infrastructure.NewDeviceRepository,

			// Database
			newMongoDB,
		),
		fx.Invoke(func(*http.Server) {}),
		fx.Invoke(func(mqtt.Client) {}),
	).Run()

	// go startHTTPComponents(ctx, cfg)
	// go startMQTTComponents(ctx, cfg)

	// <-ctx.Done()
}
