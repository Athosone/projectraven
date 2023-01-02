package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/athosone/projectraven/tracking/internal/config"
	"github.com/athosone/projectraven/tracking/internal/core/features/followPosition/contracts"
	"github.com/nats-io/nats.go"
)

const (
	anyDevice   = "*"
	durableName = "syncwritetoread"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	jsctx, err := newNats(cfg)
	if err != nil {
		panic(err)
	}
	// add pull consumer
	fetcher, err := jsctx.PullSubscribe(contracts.JetStreamDevicePositionChangedSubject(anyDevice), durableName,
		nats.AckExplicit(),
		nats.ManualAck(),
		nats.AckWait(1))
	if err != nil {
		panic(err)
	}
	// pull messages
	for {
		select {
		case <-stop:
			return
		default:
			msg, err := fetcher.Fetch(1, nats.MaxWait(1*time.Second))
			if err != nil {
				// if err is deadline exceeded, then no messages were received
				if err != context.DeadlineExceeded {
					panic(err)
				}
			}
			for _, m := range msg {
				fmt.Println("Received message: ", string(m.Data))
				m.Ack()
			}
		}
	}
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
