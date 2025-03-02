package app

import (
	"log/slog"

	"todo-api/internal/config"
	"todo-api/internal/handlers"
	"todo-api/internal/storage/postgres"

	"github.com/gofiber/fiber/v3"
)

func Run(cfg *config.Config) error {
	storage, err := postgres.New(cfg.DBConnStr)
	if err != nil {
		slog.Error("Unable to create storage", "error", err)
		return err
	}
	defer storage.Close()

	app := fiber.New()

	taskHandler := handlers.NewTaskHandler(storage)

	app.Post("/tasks", taskHandler.CreateTask)
	app.Get("/tasks", taskHandler.GetTasks)
	app.Put("/tasks/:id", taskHandler.UpdateTask)
	app.Delete("/tasks/:id", taskHandler.DeleteTask)

	slog.Info("Starting server", "port", cfg.Port)
	return app.Listen(":" + cfg.Port)
}
