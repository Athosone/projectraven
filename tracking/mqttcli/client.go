package mqttcli

import (
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
	return client, nil
}

func GetClient() mqtt.Client {
	return client
}
