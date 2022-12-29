package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/athosone/projectraven/tracking/internal/core/features/followPosition/contracts"
	"github.com/athosone/projectraven/tracking/mqttcli"
	"github.com/google/uuid"
)

func main() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	cli, err := mqttcli.NewClient("tcp://localhost:1883", "tracking")
	if err != nil {
		panic(err)
	}
	cli.Connect().Wait()
  defer cli.Disconnect(250)
	go func() {
		for {
			// generate random position
			pos := generateRandomPosition()
			// send mqtt message
			j, _ := json.Marshal(pos)
			token := cli.Publish("device.position.changed", 1, false, j)
			token.Wait()
			if token.Error() != nil {
				panic(token.Error())
			}
			fmt.Println("Sent message: ", pos)
			select {
			case <-shutdown:
				return
			case <-time.After(1 * time.Second):
			}
		}
	}()
	<-shutdown
}

func generateRandomPosition() contracts.PositionChangedMessage {
	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	return contracts.PositionChangedMessage{
		MessageId: uuid.New().String(),
		DeviceId:  fmt.Sprintf("device-%d", r.Intn(5)),
		Position: contracts.Position{
			Lat:  float64(r.Intn(100)),
			Long: float64(r.Intn(100)),
		},
		Timestamp: time.Now().Unix(),
	}
}
