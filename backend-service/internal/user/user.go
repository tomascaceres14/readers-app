package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        string `gorm:"primaryKey"`
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func RegisterRoutes(app *fiber.App, h *Handler) {
	app.Post("/users", h.Create)
	app.Get("/users", h.GetAll)
	app.Get("/users/:id", h.FindById)
	app.Put("/users/:id", h.UpdateById)
	app.Delete("/users/:id", h.DeleteById)
}
