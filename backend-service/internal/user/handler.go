package user

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/resource"
)

type Handler struct {
	svc *Service
}

type CreateUserReq struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	Username string `json:"username"`
}

type UserResponse struct {
	ID        string              `json:"id"`
	Username  string              `json:"username"`
	Resources []resource.Resource `json:"scrapped_resources"`
	CreatedAt time.Time           `json:"created_at"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{svc: service}
}

func (h *Handler) Create(c *fiber.Ctx) error {
	req := new(CreateUserReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user := User{
		Username: req.Username,
	}

	if err := h.svc.Register(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Status(fiber.StatusCreated)

	return c.JSON(fiber.Map{
		"user_id": user.ID,
	})
}

func (h *Handler) GetAll(c *fiber.Ctx) error {
	users, err := h.svc.GetAll()
	if err != nil {
		return c.SendString(err.Error())
	}

	fmt.Println(users)

	usersResponse := make([]UserResponse, 0)

	for _, u := range users {
		usersResponse = append(usersResponse, UserResponse{
			ID:        u.ID,
			Username:  u.Username,
			CreatedAt: u.CreatedAt,
			Resources: u.Resources,
		})
	}

	return c.JSON(usersResponse)
}

func (h *Handler) FindById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id in path not found"})
	}

	user, err := h.svc.FindById(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(UserResponse{
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}

func (h *Handler) UpdateById(c *fiber.Ctx) error {

	req := new(UpdateUserReq)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id in path not found"})
	}

	if err := h.svc.UpdateById(id, req.Username); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Status(fiber.StatusCreated)

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("User id: %s successfully updated.", id),
	})
}

func (h *Handler) DeleteById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Id in path not found"})
	}

	h.svc.DeleteById(id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted.",
	})
}
