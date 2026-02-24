package middleware

import (
	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
)

func FirebaseSessionMiddleware(auth *auth.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {

		sessionCookie := c.Cookies("session")
		if sessionCookie == "" {
			return c.Status(401).Redirect("/login")
		}
		decoded, err := auth.VerifySessionCookieAndCheckRevoked(
			c.Context(),
			sessionCookie,
		)
		if err != nil {
			c.ClearCookie("session")
			return c.Status(401).Redirect("/login")
		}
		c.Locals("uid", decoded.UID)
		c.Locals("claims", decoded.Claims)

		return c.Next()
	}
}
