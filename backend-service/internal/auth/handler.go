package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
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

type Register struct {
	Username string `json:"username" form:"username"`
	IdToken  string `json:"id-token" form:"token"`
}

type Login struct {
	IdToken string `json:"id-token" form:"token"`
}

func (h *Handler) RegisterView(c *fiber.Ctx) error {
	return utils.Render(c, ui.RegisterForm())
}

func (h *Handler) LoginView(c *fiber.Ctx) error {
	return utils.Render(c, ui.LoginForm())
}

func (h *Handler) Register(c *fiber.Ctx) error {

	var cred Register
	if err := c.BodyParser(&cred); err != nil {
		fmt.Println(err)
		return c.Status(400).SendString("Invalid request.")
	}

	if err := h.s.Register(c, cred); err != nil {
		return c.Status(400).SendString(err.Error())
	}
	return c.Redirect("/dashboard")
}

func (h *Handler) Login(c *fiber.Ctx) error {

	var cred Login
	if err := c.BodyParser(&cred); err != nil {
		fmt.Println(err)
		return utils.Render(c, ui.GlobalAlert(ui.AlertError, "Invalid request."))
	}

	expiresIn := time.Hour * 24 * 14

	_, err := h.s.Auth.VerifyIDToken(context.Background(), cred.IdToken)
	if err != nil {
		return utils.Render(c, ui.GlobalAlert(ui.AlertError, "Ups! An error occured. Try again later."))
	}

	sessionCookie, err := h.s.Auth.SessionCookie(
		context.Background(),
		cred.IdToken,
		expiresIn,
	)
	if err != nil {
		return utils.Render(c, ui.GlobalAlert(ui.AlertError, "Failed to create session."))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    sessionCookie,
		Path:     "/",
		MaxAge:   int(expiresIn.Seconds()),
		HTTPOnly: true,
		Secure:   true,
	})
	return c.Redirect("/dashboard")
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    "",
		HTTPOnly: true,
		Expires:  time.Now().Add(-time.Second),
	})

	return c.Redirect("/login")
}
