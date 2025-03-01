package main

import (
	"log"
	"todo-api/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
