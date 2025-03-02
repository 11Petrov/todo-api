package main

import (
	"log/slog"

	"todo-api/internal/app"
	"todo-api/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
	}

	if err := app.Run(cfg); err != nil {
		slog.Error("Failed to start application", "error", err)
	}
}
