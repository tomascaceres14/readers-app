package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/joho/godotenv"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/auth"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/dashboard"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Get enviroment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}
	dbUrl := os.Getenv("DB_STRING")
	port := os.Getenv("PORT")

	// Setup services
	app := fiber.New()

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	authClient, err := auth.InitFirebase()
	if err != nil {
		log.Fatal(err)
	}

	// Users
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	user.RegisterRoutes(app, userHandler)

	// Auth
	authService := auth.NewService(authClient, userRepo)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(app, authHandler)

	// Dashboard
	dashboardHandler := dashboard.NewHandler(authService, userService)
	dashboard.RegisterRoutes(app, dashboardHandler)

	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
	}))
	app.Static("/static", "./web/static")

	log.Fatal(app.Listen(port))
}
