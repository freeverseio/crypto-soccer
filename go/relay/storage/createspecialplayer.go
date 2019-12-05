package storage

import (
	"time"

	"github.com/google/uuid"
)

type CreateSpecialPlayer struct {
	UUID      uuid.UUID
	Timestamp time.Time
	Name      string
}
