package userresources

import (
	"time"

	"github.com/google/uuid"
)

type UserResource struct {
	UserID     string    `gorm:"primaryKey"`
	ResourceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt  time.Time
}
