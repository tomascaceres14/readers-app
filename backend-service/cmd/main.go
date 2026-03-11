package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/joho/godotenv"
	"github.com/tomascaceres14/readers-app/backend-service/internal/auth"
	"github.com/tomascaceres14/readers-app/backend-service/internal/dashboard"
	"github.com/tomascaceres14/readers-app/backend-service/internal/messaging"
	"github.com/tomascaceres14/readers-app/backend-service/internal/middleware"
	"github.com/tomascaceres14/readers-app/backend-service/internal/resource"
	resourcestatus "github.com/tomascaceres14/readers-app/backend-service/internal/resource_status"
	"github.com/tomascaceres14/readers-app/backend-service/internal/user"
	userresources "github.com/tomascaceres14/readers-app/backend-service/internal/user_resource"
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

	// Middleware
	sessionMiddleware := middleware.FirebaseSessionMiddleware(authClient)

	urRepo := userresources.NewRepository(db)
	statusRepo := resourcestatus.NewRepository(db)

	publisher, err := messaging.NewPublisher()
	if err != nil {
		log.Fatal(err)
	}
	if err := publisher.Setup(); err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()

	responseConsumerCfg := messaging.NewResponseConfig()
	responseConsumer, err := messaging.NewResponseConsumer(responseConsumerCfg, db)
	if err != nil {
		log.Fatal(err)
	}
	if err := responseConsumer.Setup(); err != nil {
		log.Fatal(err)
	}
	go responseConsumer.Listen()
	defer responseConsumer.Close()

	// Users
	usrRepo := user.NewRepository(db)
	usrSvc := user.NewService(usrRepo)
	usrHandler := user.NewHandler(usrSvc)
	user.RegisterRoutes(app, usrHandler)

	// Resources
	rscRepo := resource.NewRepository(db)
	rscSvc := resource.NewService(rscRepo, urRepo, statusRepo, publisher)
	rscHandler := resource.NewHandler(rscSvc)
	resource.RegisterRoutes(app, rscHandler, sessionMiddleware)

	// Auth
	authService := auth.NewService(authClient, usrRepo)
	authHandler := auth.NewHandler(authService)
	auth.RegisterRoutes(app, authHandler)

	// Dashboard
	dashHandler := dashboard.NewHandler(authService, usrSvc)
	dashboard.RegisterRoutes(app, dashHandler, sessionMiddleware)

	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
	}))

	app.Get("/ur", func(c *fiber.Ctx) error {
		urList, _ := urRepo.FindAll()
		return c.JSON(urList)
	})
	app.Static("/static", "./web/static")

	log.Fatal(app.Listen(port))
}
