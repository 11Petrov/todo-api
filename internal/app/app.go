package app

import (
	"log"

	"todo-api/internal/config"
	"todo-api/internal/storage"

	"github.com/gofiber/fiber/v3"
)

func Run(cfg *config.Config) error {
	dbpool, err := storage.NewDBPool(cfg.DBConnStr)
	if err != nil {
		log.Printf("Unable to create connection pool: %v", err)
		return err
	}
	defer dbpool.Close()

	if err := storage.CreateTasksTable(dbpool); err != nil {
		log.Printf("Unable to create tasks table: %v", err)
		return err
	}

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, Todo API!")
	})

	log.Printf("Starting server on :%s", cfg.Port)
	return app.Listen(":" + cfg.Port)
}
