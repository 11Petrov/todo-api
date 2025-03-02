package postgres

import (
	"context"
	"log"

	"todo-api/internal/model"
	"todo-api/internal/storage"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(connStr string) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Printf("Unable to create connection pool: %v", err)
		return nil, storage.ErrDBConnection
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Printf("Unable to ping database: %v", err)
		return nil, storage.ErrDBConnection
	}

	log.Println("Successfully connected to the database")

	if err := RunMigrations(connStr); err != nil {
		log.Printf("Failed to run migrations: %v", err)
		return nil, err
	}
	return &Storage{pool: pool}, nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func RunMigrations(connStr string) error {
	m, err := migrate.New(
		"file://migrations",
		connStr,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}

func (s *Storage) CreateTask(ctx context.Context, task *model.Task) error {
	query := `
	INSERT INTO tasks (title, description, status)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`

	return s.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (s *Storage) GetTasks(ctx context.Context) ([]model.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks`
	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		log.Printf("Failed to fetch tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			log.Printf("Failed to scan task: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Storage) UpdateTask(ctx context.Context, id string, task *model.Task) error {
	query := `
	UPDATE tasks
	SET title = $1, description = $2, status = $3, updated_at = now()
	WHERE id = $4
	RETURNING id, title, description, status, created_at, updated_at`

	err := s.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		log.Printf("Failed to update task: %v", err)
		return err
	}

	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Failed to delete task: %v", err)
		return err
	}

	return nil
}
