package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/utils"
)

// UsersHandler handles HTTP requests for users management
type UsersHandler struct {
	useCase domain.UsersUseCase
}

// NewUsersHandler creates a new UsersHandler
func NewUsersHandler(useCase domain.UsersUseCase) *UsersHandler {
	return &UsersHandler{
		useCase: useCase,
	}
}

// CreateUser handles user creation
// POST /api/admin/users
// UC13: Open Add User Dialog
// UC14: Create New User with Email and Role
func (h *UsersHandler) CreateUser(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	// Create user
	response, err := h.useCase.CreateUser(ctx, &req)
	if err != nil {
		if err.Error() == "email already exists" {
			return utils.BadRequest(c, "Email already exists", nil)
		}
		if err.Error() == "validation failed: Key: 'CreateUserRequest.FullName' Error:Field validation for 'FullName' failed on the 'required' tag" {
			return utils.BadRequest(c, "Full Name is required", nil)
		}
		if err.Error() == "validation failed: Key: 'CreateUserRequest.Email' Error:Field validation for 'Email' failed on the 'required' tag" {
			return utils.BadRequest(c, "Email is required", nil)
		}
		if err.Error() == "validation failed: Key: 'CreateUserRequest.Role' Error:Field validation for 'Role' failed on the 'required' tag" {
			return utils.BadRequest(c, "Role is required", nil)
		}
		return utils.InternalServerError(c, "Failed to create user")
	}

	return utils.Success(c, response, "User created successfully")
}

// ListUsers handles listing users with pagination
// GET /api/admin/users?page=1&per_page=5
// UC17: Display All Users with Pagination
// UC18: Navigate Users Table Pages
func (h *UsersHandler) ListUsers(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "5"))

	// Get users
	response, err := h.useCase.ListUsers(ctx, page, perPage)
	if err != nil {
		fmt.Printf("ListUsers Handler Error: %v\n", err)
		return utils.InternalServerError(c, "Failed to fetch users")
	}

	return utils.Success(c, response, "Users retrieved successfully")
}

// RefreshUsers handles refreshing the users list
// GET /api/admin/users/refresh?page=1&per_page=5
// UC12: Refresh Users List
func (h *UsersHandler) RefreshUsers(c *fiber.Ctx) error {
	ctx := c.Context()

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "5"))

	// Refresh users (same as list)
	response, err := h.useCase.RefreshUsers(ctx, page, perPage)
	if err != nil {
		return utils.InternalServerError(c, "Failed to refresh users")
	}

	return utils.Success(c, response, "Users refreshed successfully")
}

// DeleteUser handles user deletion
// DELETE /api/admin/users/:id
// UC19: Delete User from Table
func (h *UsersHandler) DeleteUser(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := c.Params("id")
	if userID == "" {
		return utils.BadRequest(c, "User ID is required", nil)
	}

	// Delete user
	err := h.useCase.DeleteUser(ctx, userID)
	if err != nil {
		if err.Error() == "invalid user ID format" {
			return utils.BadRequest(c, "Invalid user ID format", nil)
		}
		return utils.InternalServerError(c, "Failed to delete user")
	}

	return utils.Success(c, nil, "User deleted successfully")
}

// UpdateUser handles user updates
// PUT /api/admin/users/:id
func (h *UsersHandler) UpdateUser(c *fiber.Ctx) error {
	ctx := c.Context()

	userID := c.Params("id")
	if userID == "" {
		return utils.BadRequest(c, "User ID is required", nil)
	}

	var req dto.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body", nil)
	}

	// Update user
	response, err := h.useCase.UpdateUser(ctx, userID, &req)
	if err != nil {
		if err.Error() == "invalid user ID format" {
			return utils.BadRequest(c, "Invalid user ID format", nil)
		}
		if err.Error() == "user not found" {
			return utils.NotFound(c, "User not found")
		}
		return utils.InternalServerError(c, "Failed to update user")
	}

	return utils.Success(c, response, "User updated successfully")
}

// GetUserStats handles fetching user statistics
// GET /api/admin/users/stats
// UC24: Fetch and Display User Statistics
func (h *UsersHandler) GetUserStats(c *fiber.Ctx) error {
	ctx := c.Context()

	// Get stats
	stats, err := h.useCase.GetUserStats(ctx)
	if err != nil {
		return utils.InternalServerError(c, "Failed to fetch user statistics")
	}

	return utils.Success(c, stats, "User statistics retrieved successfully")
}

// ExportUsersToCSV handles exporting users to CSV
// GET /api/admin/users/export
// UC25: Export Users to CSV
func (h *UsersHandler) ExportUsersToCSV(c *fiber.Ctx) error {
	ctx := c.Context()

	// Export users
	csvData, err := h.useCase.ExportUsersToCSV(ctx)
	if err != nil {
		fmt.Printf("ExportUsersToCSV Handler Error: %v\n", err)
		return utils.InternalServerError(c, "Failed to export users")
	}

	// Set headers for CSV download
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("users_export_%s.csv", timestamp)

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	return c.Send(csvData)
}
