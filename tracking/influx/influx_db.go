package influx

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxCfg struct {
	Token  string
	Bucket string
	Org    string
	Addr   string
}

func NewInfluxDBClient(cfg *InfluxCfg) (influxdb2.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	url, token := cfg.Addr, cfg.Token
	client := influxdb2.NewClient(url, token)
	h, err := client.Health(ctx)

	if err != nil {
		return nil, err
	}
	if h.Status != "pass" {
		return nil, fmt.Errorf("influxdb health check failed")
	}
	return client, nil
}
