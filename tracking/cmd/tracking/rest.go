package main

// TODO: after the mqtt part
import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/athosone/projectraven/tracking/internal/api/rest"
	"github.com/athosone/projectraven/tracking/internal/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/fx"
)

//TODO: try to check the modules

func createServer(cfg *config.AppConfig, r chi.Router, lc fx.Lifecycle) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)
	srv := &http.Server{Addr: addr, Handler: r}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", addr)
			if err != nil {
				return err
			}
			logger.Infow("Starting HTTP server at: " + addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return srv
}

// Create base router including base middleware such as logging, recovery, etc.
func newRestServer(cfg *config.AppConfig) *rest.RestServer {
	middleware.RequestIDHeader = cfg.Service.RequestIdHeader
	addr := fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)
	s := rest.NewHttpServer(logger, rest.HttpServerConfig{
		Addr:    addr,
		IsDebug: cfg.IsDebug,
	})
	return s
}

type ChiRoute interface {
	AddRoutes(r chi.Router)
}

// Add routes wiring here
// looks like routes must be first due to value group
func newChi(routes []ChiRoute, s *rest.RestServer) chi.Router {
	r := s.Mux
	r.Route("/api", func(api chi.Router) {
		for _, route := range routes {
			route.AddRoutes(api)
		}
	})
	return r
}

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(ChiRoute)),
		fx.ResultTags(`group:"routes"`),
	)
}
