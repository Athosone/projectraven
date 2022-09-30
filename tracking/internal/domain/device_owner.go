package domain

import (
	"github.com/google/uuid"
)

type DeviceOwner struct {
	ID   uuid.UUID
	Name string
}
