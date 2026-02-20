package dashboard

import "github.com/gofiber/fiber/v2"

func RegisterRoutes(app *fiber.App, h *Handler) {
	app.Get("/dashboard", h.Dashboard)
}
