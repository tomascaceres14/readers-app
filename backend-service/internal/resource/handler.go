package resource

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/tomascaceres14/readers-app/backend-service/internal/errs"
	"github.com/tomascaceres14/readers-app/backend-service/ui"
	"github.com/tomascaceres14/readers-app/backend-service/utils"
)

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
		return c.Status(fiber.StatusOK).SendString("No URL found.")
	}

	resource := Resource{
		Url: formUrl,
	}

	uid := c.Locals("uid").(string)

	if err := h.s.Create(&resource, uid); err != nil {
		if errors.Is(err, errs.ErrAlreadyExists) {
			c.Status(fiber.StatusOK)
			return utils.Render(c, ui.GlobalInfoAlert("URL already exists."))
		}
		return c.Status(fiber.StatusOK).SendString(err.Error())
	}

	return utils.Render(c, ui.SuccessCreateResource(ui.Resource{
		ID:        resource.ID,
		Url:       resource.Url,
		CreatedAt: resource.CreatedAt,
	}))
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	resources, err := h.s.GetAll()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(resources)
}

func (h *Handler) DeleteAll(c *fiber.Ctx) error {
	h.s.DeleteAll()
	return c.Status(200).SendString("ok")
}
