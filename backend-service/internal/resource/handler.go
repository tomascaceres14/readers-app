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

func (h *Handler) Create(c *fiber.Ctx) error {
	formUrl := c.FormValue("url")
	if formUrl == "" {
		c.Status(fiber.StatusBadRequest).SendString("No URL found.")
	}

	resource := Resource{
		Url: formUrl,
	}

	if err := h.s.Create(&resource); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).Next()
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	resources, err := h.s.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(resources)
}
