package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/components"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/utils"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

type Register struct {
	Username string `json:"username" form:"username"`
	IdToken  string `json:"id-token" form:"token"`
}

type Login struct {
	IdToken string `json:"id-token" form:"token"`
}

func (h *Handler) RegisterView(c *fiber.Ctx) error {
	return utils.Render(c, components.RegisterForm())
}

func (h *Handler) LoginView(c *fiber.Ctx) error {
	return utils.Render(c, components.LoginForm())
}

func (h *Handler) Register(c *fiber.Ctx) error {

	var cred Register
	if err := c.BodyParser(&cred); err != nil {
		fmt.Println(err)
		return c.Status(400).SendString("Invalid request.")
	}

	expiresIn := time.Hour * 24 * 14

	_, err := AuthClient.VerifyIDToken(context.Background(), cred.IdToken)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).SendString("Error validating token.")
	}

	sessionCookie, err := AuthClient.SessionCookie(
		context.Background(),
		cred.IdToken,
		expiresIn,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create session.",
		})
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

func (h *Handler) Login(c *fiber.Ctx) error {

	var cred Login
	if err := c.BodyParser(&cred); err != nil {
		fmt.Println(err)
		return c.Status(400).SendString("Invalid request.")
	}

	expiresIn := time.Hour * 24 * 14

	_, err := AuthClient.VerifyIDToken(context.Background(), cred.IdToken)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).SendString(err.Error())
	}

	sessionCookie, err := AuthClient.SessionCookie(
		context.Background(),
		cred.IdToken,
		expiresIn,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create session.",
		})
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
