package dashboard

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, h *Handler, authMiddleware fiber.Handler) {
	app.Get("/dashboard", authMiddleware, h.Dashboard)
}
