package main

import (
	"context"
	"os"

	"github.com/athosone/projectraven/tracking/internal/config"
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
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	if cfg.IsDebug {
		logger.Info("debug mode enabled")
	}

	// TODO: bind context to signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go startHTTPComponents(ctx, cfg)
	go startMQTTComponents(ctx, cfg)

	<-ctx.Done()
}
