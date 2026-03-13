package resource

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	resourcestatus "github.com/tomascaceres14/readers-app/backend-service/internal/resource_status"
	"github.com/tomascaceres14/readers-app/backend-service/ui"
)

type Resource struct {
	ID        uuid.UUID                     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Url       string                        `json:"url"`
	Title     string                        `json:"title"`
	Excerpt   string                        `json:"excerpt"`
	Language  string                        `gorm:"type:varchar(10)" json:"language"`
	StatusID  uuid.UUID                     `gorm:"type:uuid;column:status_id" json:"status_id"`
	Status    resourcestatus.ResourceStatus `json:"status" gorm:"foreignKey:StatusID;references:ID"`
	CreatedAt time.Time                     `json:"created_at"`
	UpdatedAt time.Time                     `json:"updated_at"`
}

func (r Resource) toUI() ui.Resource {
	return ui.Resource{
		ID:        r.ID,
		Url:       r.Url,
		CreatedAt: r.CreatedAt,
	}
}

func ManyToUI(rs []Resource) []ui.Resource {
	uis := make([]ui.Resource, len(rs))

	for i, r := range rs {
		uis[i] = r.toUI()
	}

	return uis
}

func RegisterRoutes(app *fiber.App, h *Handler, authMiddleware fiber.Handler) {
	// public
	app.Get("/api/resource", h.GetAll)
	app.Delete("/api/resource", h.DeleteAll)

	// protected
	api := app.Group("/api/resource", authMiddleware)
	api.Post("", h.Create)
}
