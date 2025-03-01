package app

import (
	"log"

	"github.com/gofiber/fiber/v3"
)

func Run() error {
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, Todo API!")
	})

	log.Println("Starting server on :8082")
	return app.Listen(":8082")
}
