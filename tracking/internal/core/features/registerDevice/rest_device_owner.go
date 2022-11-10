package registerdevice

import (
	"net/http"
)

const (
	ProduceDeviceOwnerV1         = "application/vnd.athosone.projectraven.deviceowner+*; v=1"
	ConsumeRegisterDeviceInputV1 = "application/vnd.athosone.projectraven.registerDevice+json; v=1"
)

type DeviceOwnerHandler interface {
	RegisterDevice(w http.ResponseWriter, r *http.Request)
}

type deviceOwnerHandler struct{}

func NewDeviceOwnerHandler() DeviceOwnerHandler {
	return &deviceOwnerHandler{}
}

func (h *deviceOwnerHandler) RegisterDevice(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
