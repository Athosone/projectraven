package config

import (
	"flag"
	"os"

	"github.com/athosone/golib/pkg/config"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Service struct {
		Host            string `yaml:"host"`
		Port            int    `yaml:"port"`
		RequestIdHeader string `yaml:"requestIdHeader"`
	} `yaml:"service"`

	MongoDB struct {
		ConnectionString string `yaml:"connectionString"`
		DatabaseName     string `yaml:"databaseName"`
	} `yaml:"database"`

	// IsDebug
	IsDebug bool `yaml:"isDebug"`

	// MQTT
	MQTT struct {
		Broker   string `yaml:"broker"`
		ClientID string `yaml:"clientID"`
	} `yaml:"mqtt"`

	Nats struct {
		URL        string `yaml:"url"`
		StreamName string `yaml:"streamName"`
	} `yaml:"nats"`

	Feature struct {
		FollowPosition struct {
		} `yaml:"followPosition"`
	} `yaml:"feature"`

	InfluxDB struct {
		URL    string `yaml:"url"`
		Token  string `yaml:"token"`
		Org    string `yaml:"org"`
		Bucket string `yaml:"bucket"`
	} `yaml:"influxdb"`
}

// LoadConfig loads the configuration from the environment variable CONFIG_PATH.LoadConfig
// CONFIG_PATH is the path to the configuration folder containing yaml files.
func LoadConfig() (cfg *AppConfig, err error) {
	debug := flag.Bool("debug", false, "debug mode")
	flag.Parse()

	_ = viper.BindEnv("isDebug", "IS_DEBUG")
	_ = viper.BindEnv("service.port", "SERVICE_PORT")
	_ = viper.BindEnv("service.host", "SERVICE_HOST")
	_ = viper.BindEnv("service.requestIdHeader", "SERVICE_REQUEST_ID_HEADER")
	_ = viper.BindEnv("database.connectionString", "DATABASE_CONNECTION_STRING")
	_ = viper.BindEnv("database.databaseName", "DATABASE_DATABASE_NAME")
	_ = viper.BindEnv("mqtt.broker", "MQTT_BROKER")
	_ = viper.BindEnv("mqtt.clientID", "MQTT_CLIENT_ID")
	_ = viper.BindEnv("nats.url", "NATS_URL")
	_ = viper.BindEnv("nats.streamName", "NATS_STREAM_NAME")
	_ = viper.BindEnv("nats.subjects", "NATS_SUBJECTS")
	_ = viper.BindEnv("influxdb.addr", "INFLUXDB_URL")
	_ = viper.BindEnv("influxdb.token", "INFLUXDB_TOKEN")
	_ = viper.BindEnv("influxdb.org", "INFLUXDB_ORG")
	_ = viper.BindEnv("influxdb.bucket", "INFLUXDB_BUCKET")

	viper.SetDefault("isDebug", "true")
	viper.SetDefault("service.port", "5001")
	viper.SetDefault("service.host", "0.0.0.0")
	viper.SetDefault("service.requestIdHeader", middleware.RequestIDHeader)
	viper.SetDefault("database.connectionString", "mongodb://localhost:27017")
	viper.SetDefault("database.databaseName", "projectraven")
	viper.SetDefault("mqtt.broker", "mqtt://localhost:1883")
	viper.SetDefault("mqtt.clientID", "projectraven-tracking")
	viper.SetDefault("nats.url", "nats://localhost:4222")
	viper.SetDefault("nats.streamName", "projectraven-tracking")

	viper.SetDefault("influxdb.addr", "http://localhost:8086")
	viper.SetDefault("influxdb.token", "")
	viper.SetDefault("influxdb.org", "raven")
	viper.SetDefault("influxdb.bucket", "raven")

	cfg, err = config.LoadConfig[AppConfig](os.Getenv("CONFIG_PATH"))
	if *debug {
		cfg.IsDebug = *debug
	}
	return cfg, err
}
