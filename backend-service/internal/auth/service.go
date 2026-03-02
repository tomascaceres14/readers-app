package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/tomascaceres14/readers-app/backend-service/internal/user"
)

type Service struct {
	UsrRepo *user.Repository
	Auth    *auth.Client
}

func NewService(auth *auth.Client, usrRepo *user.Repository) *Service {
	return &Service{
		UsrRepo: usrRepo,
		Auth:    auth,
	}
}

func (s *Service) Register(ctx *fiber.Ctx, credentials Register) error {

	token, err := s.Auth.VerifyIDToken(context.Background(), credentials.IdToken)
	if err != nil {
		fmt.Println(err)
		return errors.New("Unable to verify ID Token")
	}

	username := credentials.Username
	if username == "" {
		username = token.Claims["name"].(string)
	}

	userID := token.UID
	user := &user.User{ID: userID, Username: username}
	if err := s.UsrRepo.Register(user); err != nil {
		return errors.New("Error registering user")
	}

	expiresIn := s.GetCookieExpiration()
	sessionCookie, err := s.Auth.SessionCookie(
		context.Background(),
		credentials.IdToken,
		expiresIn,
	)

	if err != nil {
		return errors.New("Failed to create new session.")
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    sessionCookie,
		Path:     "/",
		MaxAge:   int(expiresIn.Seconds()),
		HTTPOnly: true,
		Secure:   true,
	})

	return nil
}

func (s *Service) GetCookieExpiration() time.Duration {
	val := os.Getenv("SESSION_EXPIRES_IN")
	duration, err := time.ParseDuration(val)
	if err != nil || val == "" {
		return time.Hour * 24 * 7
	}

	return duration
}

func (s *Service) ValidateSessionCookie(ctx context.Context, cookie string) (*auth.Token, error) {
	return s.Auth.VerifySessionCookieAndCheckRevoked(ctx, cookie)
}
