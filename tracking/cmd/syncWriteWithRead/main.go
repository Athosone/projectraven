package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/athosone/projectraven/tracking/influx"
	api "github.com/athosone/projectraven/tracking/internal/api/jetstream"
	"github.com/athosone/projectraven/tracking/internal/config"
	"github.com/athosone/projectraven/tracking/internal/core/features/followPosition/contracts"
	domainDevice "github.com/athosone/projectraven/tracking/internal/domain/device"
	influxdb2api "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

const (
	anyDevice   = "*"
	durableName = "syncwritetoread"
)

func init() {
	logger, _ := zap.NewProduction()
	logger = logger.With(zap.String("service", durableName))
	zap.ReplaceGlobals(logger)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	jsctx, err := newNats(cfg)
	if err != nil {
		panic(err)
	}
	srv, err := api.NewJSServer(zap.S(), jsctx)
	if err != nil {
		panic(err)
	}

	writeAPI, err := newInfluxWriteAPI(cfg)
	if err != nil {
		panic(err)
	}
	topic := contracts.JetStreamDevicePositionChangedSubject(anyDevice)
	err = srv.Subscribe(ctx, topic, durableName, func(ctx context.Context, payload []byte, messageId string) error {
		var msg domainDevice.DevicePositionChangedEvent
		if err := json.Unmarshal(payload, &msg); err != nil {
			panic(err)
			return err
		}
		return influxDbHandler(ctx, writeAPI, msg)
	})
	if err != nil {
		panic(err)
	}
	<-stop
	cancel()
	ctxWithTO, cc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cc()
	srv.Stop(ctxWithTO)
}

func influxDbHandler(ctx context.Context, writeAPI influxdb2api.WriteAPIBlocking, msg domainDevice.DevicePositionChangedEvent) error {
	zap.S().Infof("Received message: %v", msg)
	tags := map[string]string{
		"device_id": msg.DeviceId,
	}
	fields := map[string]interface{}{
		"lat": msg.NewDevicePosition.Latitude,
		"lon": msg.NewDevicePosition.Longitude,
	}
	point := write.NewPoint("device_position", tags, fields, msg.CreatedAtUTC)

	err := writeAPI.WritePoint(ctx, point)
	if err != nil {
		return err
	}
	return nil
}

func newInfluxWriteAPI(cfg *config.AppConfig) (influxdb2api.WriteAPIBlocking, error) {
	client, err := influx.InitInfluxDBClient(&influx.InfluxCfg{
		Token:  cfg.InfluxDb.Token,
		Bucket: cfg.InfluxDb.Bucket,
		Org:    cfg.InfluxDb.Org,
		Addr:   cfg.InfluxDb.URL,
	})
	if err != nil {
		return nil, err
	}
	writeAPI := client.WriteAPIBlocking(cfg.InfluxDb.Org, cfg.InfluxDb.Bucket)
	return writeAPI, nil
}

// org := "raven"
// bucket := "raven_position"
// writeAPI := client.WriteAPIBlocking(org, bucket)
// for value := 0; value < 5; value++ {
// 	tags := map[string]string{
// 		"tagname1": "tagvalue1",
// 	}
// 	fields := map[string]interface{}{
// 		"field1": value,
// 	}
// 	point := write.NewPoint("measurement1", tags, fields, time.Now())
// 	time.Sleep(1 * time.Second) // separate points by 1 second

// 	if err := writeAPI.WritePoint(context.Background(), point); err != nil {
// 		log.Fatal(err)
// 	}
// }
