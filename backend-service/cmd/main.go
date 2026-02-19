package main

import (
	"context"
	"log"
	"os"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/components"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/auth"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Link struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Url    string
	UserID string
}

func Render(c *fiber.Ctx, component templ.Component) error {
	c.Set("Content-Type", "text/html")
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func main() {

	// Get enviroment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}
	dbUrl := os.Getenv("DB_STRING")
	port := os.Getenv("PORT")

	// Setup services
	app := fiber.New()
	app.Static("/static", "./web/static")

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	auth.InitFirebase()

	// Register entities and handlers
	user.Initialize(app, db)
	auth.Initialize(app, db)

	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return Render(c, components.Welcome())
	})

	app.Get("/dashboard", func(c *fiber.Ctx) error {
		sessionCookie := c.Cookies("session")

		if sessionCookie == "" {
			return c.Redirect("/login")
		}

		decoded, err := auth.AuthClient.VerifySessionCookie(
			context.Background(),
			sessionCookie,
		)

		if err != nil {
			c.ClearCookie("session")
			return c.Redirect("/login")
		}
		return Render(c, components.Dashboard(decoded.Claims["email"].(string)))
	})

	log.Fatal(app.Listen(port))
}
