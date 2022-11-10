package main

import (
	"context"
	"fmt"

	"github.com/athosone/projectraven/tracking/internal/api/rest"
)

func startHTTPComponents(ctx context.Context, cfg *AppConfig) {
	srv := createServer(ctx, cfg)
	addr := fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)
	srv.Run(addr)
}

func createServer(ctx context.Context, cfg *AppConfig) *rest.HttpServer {
	application, err := NewApplication(ctx, cfg)

	if err != nil {
		logger.Fatal(err)
	}
	s := rest.NewHttpServer(application, logger, rest.HttpServerConfig{
		Addr:    fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port),
		IsDebug: cfg.IsDebug,
	})
	return s
}
