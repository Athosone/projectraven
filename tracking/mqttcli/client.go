package mqttcli

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var client mqtt.Client

func NewClient(broker string, clientID string) (mqtt.Client, error) {
	if client != nil {
		return client, nil
	}
	opts := mqtt.
		NewClientOptions().
		AddBroker(broker).
		SetClientID(clientID).
		SetAutoAckDisabled(true)

	client = mqtt.NewClient(opts)
	err := client.Connect().Error()
	if err != nil {
		return nil, fmt.Errorf("error connecting to MQTT broker: %w", err)
	}
	return client, nil
}

func GetClient() mqtt.Client {
	return client
}
