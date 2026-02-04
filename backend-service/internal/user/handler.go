package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	return c.JSON([]fiber.Map{
		{
			"username": "user_teasdasdasst",
			"email":    "who_cares@gmail.com",
		},
		{
			"username": "user_test2",
			"email":    "mymail@gmail.com",
		},
	})
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var user User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Status(fiber.StatusCreated)

	return c.JSON(user)
}
