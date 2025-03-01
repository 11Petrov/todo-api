package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnStr string
	Port      string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file: %v", err)
	}

	return &Config{
		DBConnStr: os.Getenv("DB_CONN_STR"),
		Port:      os.Getenv("PORT"),
	}, nil
}
