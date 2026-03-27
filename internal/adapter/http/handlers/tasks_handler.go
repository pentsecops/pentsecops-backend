package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/utils"
)

type TasksHandler struct {
	usecase domain.TasksUseCase
}

func NewTasksHandler(usecase domain.TasksUseCase) *TasksHandler {
	return &TasksHandler{
		usecase: usecase,
	}
}

// CreateTask - UC37: Create New Task
func (h *TasksHandler) CreateTask(c *fiber.Ctx) error {
	var req dto.CreateTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	response, err := h.usecase.CreateTask(c.Context(), &req)
	if err != nil {
		return utils.InternalServerError(c, "Failed to create task")
	}

	return utils.Success(c, response, "Task created successfully")
}

// ListTasksByProject - UC36: List tasks for a specific project
func (h *TasksHandler) ListTasksByProject(c *fiber.Ctx) error {
	projectID := c.Params("project_id")
	if projectID == "" {
		return utils.BadRequest(c, "Project ID is required", nil)
	}

	response, err := h.usecase.ListTasksByProject(c.Context(), projectID)
	if err != nil {
		return utils.InternalServerError(c, "Failed to retrieve tasks")
	}

	return utils.Success(c, response, "Tasks retrieved successfully")
}

// ListAllTasks - UC36: Display Task Board Interface
func (h *TasksHandler) ListAllTasks(c *fiber.Ctx) error {
	response, err := h.usecase.ListAllTasks(c.Context())
	if err != nil {
		return utils.InternalServerError(c, "Failed to retrieve task board")
	}

	return utils.Success(c, response, "Task board retrieved successfully")
}

// UpdateTaskStatus - UC38: Update Task Status
func (h *TasksHandler) UpdateTaskStatus(c *fiber.Ctx) error {
	taskID := c.Params("id")
	if taskID == "" {
		return utils.BadRequest(c, "Task ID is required", nil)
	}

	var req dto.UpdateTaskStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	response, err := h.usecase.UpdateTaskStatus(c.Context(), taskID, &req)
	if err != nil {
		return utils.InternalServerError(c, "Failed to update task status")
	}

	return utils.Success(c, response, "Task status updated successfully")
}

// DeleteTask - Delete a task
func (h *TasksHandler) DeleteTask(c *fiber.Ctx) error {
	taskID := c.Params("id")
	if taskID == "" {
		return utils.BadRequest(c, "Task ID is required", nil)
	}

	if err := h.usecase.DeleteTask(c.Context(), taskID); err != nil {
		return utils.InternalServerError(c, "Failed to delete task")
	}

	return utils.Success(c, nil, "Task deleted successfully")
}
