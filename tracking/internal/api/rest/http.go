package rest

import (
	"net/http"

	gserver "github.com/athosone/golib/pkg/server"
	gmiddleware "github.com/athosone/golib/pkg/server/middleware"
	app "github.com/athosone/projectraven/tracking/internal/application"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type HttpServer struct {
	*chi.Mux
	app    *app.Application
	logger *zap.SugaredLogger
}

type HttpServerConfig struct {
	IsDebug bool
}

func NewHttpServer(application *app.Application, logger *zap.SugaredLogger, cfg HttpServerConfig) HttpServer {

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(gmiddleware.InjectLoggerInRequest(func(r *http.Request) *zap.SugaredLogger {
		return logger.With("request_id", r.Header.Get(middleware.RequestIDHeader))
	}))
	r.Use(gmiddleware.RequestLogger([]string{"/healthy", "/ready"}))
	r.Use(gmiddleware.CompressResponse())
	r.Use(middleware.Heartbeat("/healthy"))
	r.Use(middleware.Recoverer)

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/openapi", func(api chi.Router) {

	})

	if cfg.IsDebug {
		r.Handle("/debug/vars", http.DefaultServeMux)
		r.Handle("/debug/pprof/", http.DefaultServeMux)
		r.Handle("/openapi/", http.StripPrefix("/openapi/", http.FileServer(http.Dir("./docs"))))
	}

	return HttpServer{
		Mux:    r,
		app:    application,
		logger: logger,
	}
}

func (s *HttpServer) Run(addr string) error {
	server := &http.Server{Addr: addr, Handler: s.Mux}
	s.logger.Info("Starting server on ", addr)
	gserver.ListenAndServe(server)
	return nil
}
