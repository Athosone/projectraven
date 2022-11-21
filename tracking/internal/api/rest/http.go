package rest

import (
	"net/http"
	"strings"

	gserver "github.com/athosone/golib/pkg/server"
	gmiddleware "github.com/athosone/golib/pkg/server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type HttpServer struct {
	*chi.Mux
	logger *zap.SugaredLogger
}

type HttpServerConfig struct {
	Addr    string
	IsDebug bool
}

func NewHttpServer(logger *zap.SugaredLogger, cfg HttpServerConfig) *HttpServer {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(gmiddleware.InjectLoggerInRequest(func(r *http.Request) *zap.SugaredLogger {
		return logger.With("request_id", r.Header.Get(middleware.RequestIDHeader))
	}))
	r.Use(gmiddleware.RequestLogger([]string{"/healthy", "/ready"}))
	r.Use(gmiddleware.CompressResponse())
	r.Use(middleware.Heartbeat("/healthy"))
	r.Use(middleware.Recoverer)
	r.Use(openApi)

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/api", func(api chi.Router) {
		// AddUserRoutes(api, NewUserHandler(application.Commands.UserCommands))
	})

	if cfg.IsDebug {
		r.Handle("/debug/vars", http.DefaultServeMux)
		r.Handle("/debug/pprof/", http.DefaultServeMux)
		r.Handle("/openapi/", http.StripPrefix("/openapi/", http.FileServer(http.Dir("./docs"))))
	}

	return &HttpServer{
		Mux:    r,
		logger: logger,
	}
}

func (s *HttpServer) AddRoute(pattern string, router func(r chi.Router)) {
	s.Mux.Route(pattern, router)
}

func (s *HttpServer) Run(addr string) error {
	server := &http.Server{Addr: addr, Handler: s.Mux}
	s.logger.Info("Starting server on ", addr)
	gserver.ListenAndServe(server)
	return nil
}

func openApi(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Index(r.URL.Path, "/openapi") == 0 || strings.Index(r.URL.Path, "/openapi/") == 0 {
			http.StripPrefix("/openapi", http.FileServer(http.Dir("./internal/api/rest/docs/swagger/swagger-ui/"))).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
