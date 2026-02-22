package resource

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Resource struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Url       string
	CreatedAt time.Time
}

func RegisterRoutes(app *fiber.App, h *Handler) {
	app.Post("/resource", h.Create)
	app.Get("/api/resource", h.GetAll)
	app.Delete("/api/resource", h.DeleteAll)
}
