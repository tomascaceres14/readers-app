package dashboard

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/components"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/auth"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/user"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/utils"
)

type Handler struct {
	userSvc *user.Service
}

func NewHandler(authService *auth.Service, userService *user.Service) *Handler {
	return &Handler{
		userSvc: userService,
	}
}

func (h *Handler) Dashboard(c *fiber.Ctx) error {
	uid := c.Locals("uid").(string)
	usr, err := h.userSvc.FindById(uid)
	if err != nil {
		return err
	}
	return utils.Render(c, components.Dashboard(usr.Username))
}
