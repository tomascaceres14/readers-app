package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Initialize(app *fiber.App, db *gorm.DB) {

	handler := NewHandler()

	app.Get("/register", handler.RegisterView)
	app.Post("/auth/register", handler.Register)
	app.Get("/login", handler.LoginView)
	app.Post("/auth/login", handler.Login)
	app.Post("/auth/logout", handler.Logout)
}
