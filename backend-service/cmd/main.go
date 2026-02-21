package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/joho/godotenv"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/auth"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/dashboard"
	"github.com/tomascaceres14/readers-app/app-server/backend-service/internal/resource"
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
	usrRepo := user.NewRepository(db)
	usrSvc := user.NewService(usrRepo)
	usrHandler := user.NewHandler(usrSvc)
	user.RegisterRoutes(app, usrHandler)

	// Resources
	rscRepo := resource.NewRepository(db)
	rscSvc := resource.NewService(rscRepo)
	rscHandler := resource.NewHandler(rscSvc)
	resource.RegisterRoutes(app, rscHandler)

	// Auth
	authService := auth.NewService(authClient, usrRepo)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(app, authHandler)

	// Dashboard
	dashHandler := dashboard.NewHandler(authService, usrSvc)
	dashboard.RegisterRoutes(app, dashHandler)

	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
	}))
	app.Static("/static", "./web/static")

	log.Fatal(app.Listen(port))
}
