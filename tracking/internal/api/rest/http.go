package rest

import (
	"net/http"
	"strings"
	"time"

	gmiddleware "github.com/athosone/golib/pkg/server/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type RestServer struct {
	Mux    *chi.Mux
	Addr   string
	logger *zap.SugaredLogger
	Server *http.Server
}

type HttpServerConfig struct {
	Addr    string
	IsDebug bool
}

func NewHttpServer(logger *zap.SugaredLogger, cfg HttpServerConfig) *RestServer {
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

	if cfg.IsDebug {
		r.Handle("/debug/vars", http.DefaultServeMux)
		r.Handle("/debug/pprof/", http.DefaultServeMux)
		r.Handle("/openapi/", http.StripPrefix("/openapi/", http.FileServer(http.Dir("./docs"))))
	}
	srv := &http.Server{Addr: cfg.Addr, Handler: r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return &RestServer{
		Mux:    r,
		Addr:   cfg.Addr,
		logger: logger,
		Server: srv,
	}
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
