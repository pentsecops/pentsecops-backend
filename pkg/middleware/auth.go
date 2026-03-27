package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/pentsecops/backend/pkg/auth"
)

// AuthMiddleware creates authentication middleware with PASETO validation
type AuthMiddleware struct {
	pasetoService *auth.PasetoService
}

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(pasetoService *auth.PasetoService) *AuthMiddleware {
	return &AuthMiddleware{
		pasetoService: pasetoService,
	}
}

// RequireAuth returns a middleware that requires authentication
func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := m.pasetoService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Check token type (should be access token)
		if claims.TokenType != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token type",
			})
		}

		// Set user data in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}

// RequireRole returns a middleware that checks if the user has the required role
func (m *AuthMiddleware) RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid authorization header format",
			})
		}

		token := parts[1]

		// Validate token
		claims, err := m.pasetoService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// Check token type (should be access token)
		if claims.TokenType != "access" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token type",
			})
		}

		// Set user data in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)

		// Check role
		if strings.ToLower(claims.Role) != strings.ToLower(role) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireAnyRole returns a middleware that checks if the user has any of the required roles
func (m *AuthMiddleware) RequireAnyRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First, authenticate the user
		if err := m.RequireAuth()(c); err != nil {
			return err
		}

		// Check role
		userRole, ok := c.Locals("user_role").(string)
		if !ok || userRole == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User role not found",
			})
		}

		// Check if user has any of the required roles
		for _, role := range roles {
			if strings.ToLower(userRole) == strings.ToLower(role) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}
