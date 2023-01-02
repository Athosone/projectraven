package influx

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxCfg struct {
	Token  string
	Bucket string
	Org    string
	Addr   string
}

var Client influxdb2.Client

func InitInfluxDBClient(cfg *InfluxCfg) (influxdb2.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	url, token := cfg.Addr, cfg.Token
	client := influxdb2.NewClient(url, token)
	pinged, err := client.Ping(ctx)

	if err != nil {
		return nil, err
	}
	if !pinged {
		return nil, err
	}
	Client = client
	return client, nil
}
