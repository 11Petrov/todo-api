package postgres

import (
	"context"
	"fmt"

	"todo-api/internal/model"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(connStr string) (*Storage, error) {
	const op = "storage.postgres.New"

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := RunMigrations(connStr); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{pool: pool}, nil
}

func RunMigrations(connStr string) error {
	const op = "storage.postgres.RunMigrations"

	m, err := migrate.New(
		"file://migrations",
		connStr,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) Close() {
	s.pool.Close()
}

func (s *Storage) CreateTask(ctx context.Context, task *model.Task) error {
	const op = "storage.postgres.CreateTask"

	query := `
	INSERT INTO tasks (title, description, status)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`

	if err := s.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status).Scan(
		&task.ID, &task.CreatedAt, &task.UpdatedAt,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetTasks(ctx context.Context) ([]model.Task, error) {
	const op = "storage.postgres.GetTasks"

	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return tasks, nil
}

func (s *Storage) UpdateTask(ctx context.Context, id string, task *model.Task) error {
	const op = "storage.postgres.UpdateTask"

	query := `
	UPDATE tasks
	SET title = $1, description = $2, status = $3, updated_at = now()
	WHERE id = $4
	RETURNING id, title, description, status, created_at, updated_at`

	if err := s.pool.QueryRow(ctx, query, task.Title, task.Description, task.Status, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteTask(ctx context.Context, id string) error {
	const op = "storage.postgres.DeleteTask"

	query := `DELETE FROM tasks WHERE id = $1`
	result, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("%s: task not found", op)
	}

	return nil
}
