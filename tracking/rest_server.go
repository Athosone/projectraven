package main

// TODO: after the mqtt part
import (
	"context"
	"fmt"

	"github.com/athosone/projectraven/tracking/internal/api/rest"
	"github.com/athosone/projectraven/tracking/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func startHTTPComponents(ctx context.Context, cfg *config.AppConfig) {
	srv := createServer(ctx, cfg)
	addr := fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)
	srv.Run(addr)
}

func createServer(ctx context.Context, cfg *config.AppConfig) *rest.HttpServer {
	middleware.RequestIDHeader = cfg.Service.RequestIdHeader

	s := rest.NewHttpServer(logger, rest.HttpServerConfig{
		Addr:    fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port),
		IsDebug: cfg.IsDebug,
	})
	// TODO: add endpoints
	s.AddRoute("/api", func(r chi.Router) {
		// AddUserRoutes(r, userHandler)
	})
	return s
}
