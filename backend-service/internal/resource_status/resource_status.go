package resourcestatus

import (
	"time"

	"github.com/google/uuid"
)

const (
	OK      = "OK"
	PENDING = "PENDING"
	FAILED  = "FAILED"
	ERROR   = "ERROR"
)

type ResourceStatus struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();column:id" json:"id"`
	Name      string    `gorm:"uniqueIndex;not null;column:name" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

func (ResourceStatus) TableName() string {
	return "resources_status"
}
