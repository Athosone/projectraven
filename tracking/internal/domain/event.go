package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DomainEvent struct {
	EventId      uuid.UUID
	CreatedAtUTC time.Time
}

type EventPublisher func(ctx context.Context, topic string, msg []byte) error
