package main

import (
	"flag"
	"os"

	"github.com/athosone/golib/pkg/config"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Service struct {
		Port            int    `yaml:"port"`
		RequestIdHeader string `yaml:"requestIdHeader"`
	} `yaml:"service"`

	Database struct {
		ConnectionString string `yaml:"connectionString"`
		DatabaseName     string `yaml:"databaseName"`
	} `yaml:"database"`

	// IsDebug
	IsDebug bool `yaml:"isDebug"`
}

// LoadConfig loads the configuration from the environment variable CONFIG_PATH.LoadConfig
// CONFIG_PATH is the path to the configuration folder containing yaml files.
func LoadConfig() (cfg *AppConfig, err error) {
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	_ = viper.BindEnv("isDebug", "IS_DEBUG")
	_ = viper.BindEnv("service.port", "SERVICE_PORT")
	_ = viper.BindEnv("service.requestIdHeader", "SERVICE_REQUEST_ID_HEADER")
	_ = viper.BindEnv("database.connectionString", "DATABASE_CONNECTION_STRING")
	_ = viper.BindEnv("database.databaseName", "DATABASE_DATABASE_NAME")

	viper.SetDefault("isDebug", "true")
	viper.SetDefault("service.port", "5000")
	viper.SetDefault("service.requestIdHeader", middleware.RequestIDHeader)
	viper.SetDefault("database.connectionString", "mongodb://localhost:27017")
	viper.SetDefault("database.databaseName", "projectraven")

	cfg, err = config.LoadConfig[AppConfig](os.Getenv("CONFIG_PATH"))
	if err == nil {
		middleware.RequestIDHeader = cfg.Service.RequestIdHeader
	}
	if *debug {
		cfg.IsDebug = *debug
	}
	return cfg, err
}
