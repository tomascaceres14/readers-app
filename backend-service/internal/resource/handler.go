package resource

import "github.com/gofiber/fiber/v2"

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) Create(c *fiber.Ctx) {
	formUrl := c.FormValue("url")
	if formUrl == "" {
		c.Status(fiber.StatusBadRequest).SendString("No URL found.")
	}

	resource := Resource{
		Url: formUrl,
	}

	if err := h.s.Create(&resource); err != nil {
		c.Status(fiber.StatusInternalServerError).SendString("Error creating resource in db.")
	}

	c.Status(fiber.StatusOK).SendString("resource created.")

}
