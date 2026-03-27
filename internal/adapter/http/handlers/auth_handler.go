package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/internal/core/domain"
	"github.com/pentsecops/backend/internal/core/domain/dto"
	"github.com/pentsecops/backend/pkg/auth/logger"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	useCase domain.AuthUseCase
	logger  *logger.Logger
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(useCase domain.AuthUseCase, logger *logger.Logger) *AuthHandler {
	return &AuthHandler{
		useCase: useCase,
		logger:  logger,
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate user and return access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse login request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Log login attempt (without password)
	h.logger.Info("Login attempt", "email", req.Email, "ip", c.IP())

	resp, err := h.useCase.Login(c.Context(), &req)
	if err != nil {
		// Log failed login attempt
		h.logger.Warn("Login failed", "email", req.Email, "ip", c.IP(), "error", err.Error())
		
		// Return generic error message for security
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials or account locked",
		})
	}

	// Log successful login
	h.logger.Info("Login successful", "email", req.Email, "user_id", resp.UserID, "ip", c.IP())

	return c.Status(fiber.StatusOK).JSON(resp)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Generate new access and refresh tokens using a valid refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.RefreshTokenResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req dto.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse refresh token request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	h.logger.Info("Token refresh attempt", "ip", c.IP())

	resp, err := h.useCase.RefreshToken(c.Context(), &req)
	if err != nil {
		h.logger.Warn("Token refresh failed", "ip", c.IP(), "error", err.Error())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired refresh token",
		})
	}

	h.logger.Info("Token refresh successful", "ip", c.IP())

	return c.Status(fiber.StatusOK).JSON(resp)
}

// ChangePassword handles password change
// @Summary Change user password
// @Description Change the password for the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Password change request"
// @Success 200 {object} dto.ChangePasswordResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/change-password [post]
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	// Get user ID from context (set by auth middleware)
	userID, ok := c.Locals("user_id").(string)
	if !ok || userID == "" {
		h.logger.Error("User ID not found in context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req dto.ChangePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse change password request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	h.logger.Info("Password change attempt", "user_id", userID, "ip", c.IP())

	resp, err := h.useCase.ChangePassword(c.Context(), userID, &req)
	if err != nil {
		h.logger.Warn("Password change failed", "user_id", userID, "ip", c.IP(), "error", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	h.logger.Info("Password change successful", "user_id", userID, "ip", c.IP())

	return c.Status(fiber.StatusOK).JSON(resp)
}

// Logout handles user logout
// @Summary User logout
// @Description Invalidate the refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LogoutRequest true "Logout request"
// @Success 200 {object} dto.LogoutResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req dto.LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse logout request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Try to get user ID from context (optional for logout)
	userID, _ := c.Locals("user_id").(string)
	h.logger.Info("Logout attempt", "user_id", userID, "ip", c.IP())

	resp, err := h.useCase.Logout(c.Context(), &req)
	if err != nil {
		h.logger.Warn("Logout failed", "user_id", userID, "ip", c.IP(), "error", err.Error())
		// Don't fail logout even if token deletion fails
	}

	h.logger.Info("Logout successful", "user_id", userID, "ip", c.IP())

	return c.Status(fiber.StatusOK).JSON(resp)
}

