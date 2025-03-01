package app

import (
	"log"

	"todo-api/internal/config"
	"todo-api/internal/handlers"
	"todo-api/internal/storage/postgres"

	"github.com/gofiber/fiber/v3"
)

func Run(cfg *config.Config) error {
	storage, err := postgres.New(cfg.DBConnStr)
	if err != nil {
		log.Printf("Unable to create storage: %v", err)
		return err
	}
	defer storage.Close()

	app := fiber.New()

	taskHandler := handlers.NewTaskHandler(storage)

	app.Post("/tasks", taskHandler.CreateTask)

	log.Printf("Starting server on :%s", cfg.Port)
	return app.Listen(":" + cfg.Port)
}
