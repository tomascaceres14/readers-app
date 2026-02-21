package resource

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	Url       string
	CreatedAt time.Time
}
