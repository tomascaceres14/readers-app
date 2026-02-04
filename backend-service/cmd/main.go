package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Link struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Url    string
	UserID string
}

type apiConfig struct {
	port string
	db   *gorm.DB
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// Get enviroment variables
	dbUrl := os.Getenv("DB_STRING")
	port := os.Getenv("API_PORT")

	// Setup services
	app := fiber.New()
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Register entities and handlers
	user.Register(app, db)

	log.Fatal(app.Listen(port))
}
