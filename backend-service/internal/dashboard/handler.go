package dashboard

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/components"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/auth"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/user"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/utils"
)

type Handler struct {
	authService *auth.Service
	userRepo    *user.Service
}

func NewHandler(authService *auth.Service, userService *user.Service) *Handler {
	return &Handler{
		authService: authService,
		userRepo:    userService,
	}
}

func (h *Handler) Dashboard(c *fiber.Ctx) error {

	// All of this should be a middleware

	sessionCookie := c.Cookies("session")

	if sessionCookie == "" {
		return c.Redirect("/login")
	}
	decoded, err := h.authService.ValidateSessionCookie(
		c.Context(),
		sessionCookie,
	)

	if err != nil {
		c.ClearCookie("session")
		return c.Redirect("/login")
	}
	//
	usr, err := h.userRepo.FindById(decoded.UID)
	if err != nil {
		return err
	}

	return utils.Render(c, components.Dashboard(usr.Username))
}
