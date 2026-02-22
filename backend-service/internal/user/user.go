package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/resource"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Username  string
	Resources []resource.Resource `gorm:"many2many;user_resources"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

func RegisterRoutes(app *fiber.App, h *Handler) {
	app.Post("/users", h.Create)
	app.Get("/users", h.GetAll)
	app.Get("/users/:id", h.FindById)
	app.Put("/users/:id", h.UpdateById)
	app.Delete("/users/:id", h.DeleteById)
}
