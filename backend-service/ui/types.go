package ui

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Username  string
	Resources []Resource
}

type Resource struct {
	ID        uuid.UUID
	Url       string
	CreatedAt time.Time
}
