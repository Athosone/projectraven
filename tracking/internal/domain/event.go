package domain

import (
	"time"

	"github.com/google/uuid"
)

type DomainEvent struct {
	EventId      uuid.UUID
	CreatedAtUTC time.Time
}

type EventRepository interface {
	Save(event any) error
}
