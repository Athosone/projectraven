package infrastructure

import domainDevice "github.com/athosone/projectraven/tracking/internal/domain/device"

func NatsSubjects() []string {
	sub := []string{
		domainDevice.RootDeviceTopic,
	}
	natsSub := []string{}
	for _, s := range sub {
		natsSub = append(natsSub, s+".>")
	}
	return natsSub
}
