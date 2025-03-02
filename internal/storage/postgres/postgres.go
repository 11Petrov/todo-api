package postgres

import (
	"context"
	"log"

	"todo-api/internal/model"
	"todo-api/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(connStr string) (*Storage, error) {
	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Printf("Unable to create connection pool: %v", err)
		return nil, storage.ErrDBConnection
	}

	if err := dbpool.Ping(context.Background()); err != nil {
		log.Printf("Unable to ping database: %v", err)
		return nil, storage.ErrDBConnection
	}

	log.Println("Successfully connected to the database")

	if err := createTasksTable(dbpool); err != nil {
		log.Printf("Unable to create tasks table: %v", err)
		return nil, err
	}

	return &Storage{db: dbpool}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func createTasksTable(db *pgxpool.Pool) error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
		created_at TIMESTAMP DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	);`

	_, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Printf("Unable to create tasks table: %v", err)
		return err
	}

	log.Println("Tasks table created or already exists")
	return nil
}

func (s *Storage) CreateTask(ctx context.Context, task *model.Task) error {
	query := `
	INSERT INTO tasks (title, description, status)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`

	return s.db.QueryRow(ctx, query, task.Title, task.Description, task.Status).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
}

func (s *Storage) GetTasks(ctx context.Context) ([]model.Task, error) {
	query := `SELECT id, title, description, status, created_at, updated_at FROM tasks`
	rows, err := s.db.Query(ctx, query)
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

	err := s.db.QueryRow(ctx, query, task.Title, task.Description, task.Status, id).Scan(
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
