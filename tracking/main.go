package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	golibsrv "github.com/athosone/golib/pkg/server"
	golibmiddleware "github.com/athosone/golib/pkg/server/middleware"
	"github.com/athosone/projectraven/tracking/internal/api/rest"
	app "github.com/athosone/projectraven/tracking/internal/application"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	logger = zap.S().With("service", "cluster-sts")
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	if cfg.IsDebug {
		logger.Info("debug mode enabled")
	}

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Service.Port)
	server := &http.Server{Addr: addr, Handler: service(cfg)}
	logger.Info("Starting server on ", addr)
	golibsrv.ListenAndServe(server)
}

func service(cfg *AppConfig) http.Handler {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	application, err := NewApplication(ctx, cfg)
	if err != nil {
		logger.Fatal(err)
	}
	r := configureRouter(cfg, application)

	r.Route("/api", func(api chi.Router) {
		rest.AddUserRoutes(api, rest.NewUserHandler(application.Commands.UserCommands))
	})
	if cfg.IsDebug {
		r.Handle("/debug/vars", http.DefaultServeMux)
	}

	return r
}

func configureRouter(cfg *AppConfig, application *app.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(golibmiddleware.InjectLoggerInRequest(func(r *http.Request) *zap.SugaredLogger {
		return logger.With("request_id", r.Header.Get(middleware.RequestIDHeader))
	}))
	r.Use(golibmiddleware.RequestLogger([]string{"/healthy", "/ready"}))
	r.Use(golibmiddleware.CompressResponse())
	r.Use(middleware.Heartbeat("/healthy"))
	r.Use(middleware.Recoverer)
	r.Use(uiMiddleware)

	r.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Route("/openapi", func(api chi.Router) {

	})

	if cfg.IsDebug {
		r.Handle("/debug/pprof/", http.DefaultServeMux)
		// r.Handle("/openapi/", http.StripPrefix("/", http.FileServer(http.Dir("./internal/api/rest/docs/swagger/swagger-ui"))))
	}

	return r
}

func uiMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Shortcut helpers for swagger-ui
		if r.URL.Path == "/swagger-ui" || r.URL.Path == "/api/help" {
			http.Redirect(w, r, "/swagger-ui/", http.StatusFound)
			return
		}
		// Serving ./swagger-ui/
		if strings.Index(r.URL.Path, "/swagger-ui/") == 0 {
			http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./internal/api/rest/docs/swagger/swagger-ui/"))).ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}
