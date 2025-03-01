package main

import (
	"log"

	"todo-api/internal/app"
	"todo-api/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := app.Run(cfg); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}
