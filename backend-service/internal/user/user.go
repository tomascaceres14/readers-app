package user

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Register(app *fiber.App, db *gorm.DB) {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	app.Post("/users", handler.Register)
	app.Get("/users", handler.GetAll)
	app.Get("/users/:id", handler.FindById)
	app.Put("/users/:id", handler.UpdateById)
	app.Delete("/users/:id", handler.DeleteById)
}
