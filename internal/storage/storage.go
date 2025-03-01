package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBPool(connStr string) (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Printf("Unable to create connection pool: %v", err)
		return nil, err
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Printf("Unable to ping database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return dbpool, nil
}

func CreateTasksTable(dbpool *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);`

	_, err := dbpool.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Unable to create tasks table: %v", err)
		return err
	}

	log.Println("Tasks table created or already exists")
	return nil
}
