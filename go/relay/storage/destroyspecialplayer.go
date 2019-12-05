package storage

import (
	"time"

	"github.com/google/uuid"
)

type DestroySpecialPlayer struct {
	UUID      uuid.UUID
	Timestamp time.Time
}
