package main

import (
	"log"
	"os"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/components"
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

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Get enviroment variables
	dbUrl := os.Getenv("DB_STRING")
	port := os.Getenv("PORT")

	// Setup services
	app := fiber.New()
	app.Static("/static", "./static")
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Register entities and handlers
	user.Register(app, db)

	app.Get("/login", func(c *fiber.Ctx) error {
		return Render(c, components.LoginForm())
	})
	app.Get("/register", func(c *fiber.Ctx) error {
		return Render(c, components.RegisterForm())
	})

	log.Fatal(app.Listen(port))
}
