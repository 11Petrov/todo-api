package handlers

import (
	"context"
	"net/http"

	"todo-api/internal/model"

	"github.com/gofiber/fiber/v3"
)

type TaskStore interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetTasks(ctx context.Context) ([]model.Task, error)
	UpdateTask(ctx context.Context, id string, task *model.Task) error
	DeleteTask(ctx context.Context, id string) error
}

type TaskHandler struct {
	store TaskStore
}

func NewTaskHandler(store TaskStore) *TaskHandler {
	return &TaskHandler{store: store}
}

func (h *TaskHandler) CreateTask(c fiber.Ctx) error {
	var task model.Task
	if err := c.Bind().Body(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := h.store.CreateTask(c.Context(), &task); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create task"})
	}

	return c.Status(http.StatusCreated).JSON(task)
}

func (h *TaskHandler) GetTasks(c fiber.Ctx) error {
	tasks, err := h.store.GetTasks(c.Context())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch tasks"})
	}

	return c.JSON(tasks)
}

func (h *TaskHandler) UpdateTask(c fiber.Ctx) error {
	id := c.Params("id")

	var task model.Task
	if err := c.Bind().Body(&task); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := h.store.UpdateTask(c.Context(), id, &task); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update task"})
	}

	return c.JSON(task)
}

func (h *TaskHandler) DeleteTask(c fiber.Ctx) error {

	id := c.Params("id")

	if err := h.store.DeleteTask(c.Context(), id); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete task"})
	}

	return c.SendStatus(http.StatusNoContent)
}
