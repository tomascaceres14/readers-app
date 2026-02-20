package auth

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler) {
	app.Get("/register", h.RegisterView)
	app.Get("/login", h.LoginView)
	app.Post("/auth/register", h.Register)
	app.Post("/auth/login", h.Login)
	app.Post("/auth/logout", h.Logout)
}
