package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Success sends a successful response
func Success(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 created response
func Created(c *fiber.Ctx, data interface{}, message string) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// NoContent sends a 204 no content response
func NoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

// BadRequest sends a 400 bad request response
func BadRequest(c *fiber.Ctx, message string, details interface{}) error {
	return c.Status(fiber.StatusBadRequest).JSON(Response{
		Success: false,
		Error: ErrorDetail{
			Code:    "BAD_REQUEST",
			Message: message,
			Details: details,
		},
	})
}

// Unauthorized sends a 401 unauthorized response
func Unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(Response{
		Success: false,
		Error: ErrorDetail{
			Code:    "UNAUTHORIZED",
			Message: message,
		},
	})
}

// Forbidden sends a 403 forbidden response
func Forbidden(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusForbidden).JSON(Response{
		Success: false,
		Error: ErrorDetail{
			Code:    "FORBIDDEN",
			Message: message,
		},
	})
}

// NotFound sends a 404 not found response
func NotFound(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusNotFound).JSON(Response{
		Success: false,
		Error: ErrorDetail{
			Code:    "NOT_FOUND",
			Message: message,
		},
	})
}

// Conflict sends a 409 conflict response
func Conflict(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusConflict).JSON(Response{
		Success: false,
		Error: ErrorDetail{
			Code:    "CONFLICT",
			Message: message,
		},
	})
}

// ValidationError sends a 422 validation error response
func ValidationError(c *fiber.Ctx, message string, details interface{}) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(Response{
		Success: false,
		Error: ErrorDetail{
			Code:    "VALIDATION_ERROR",
			Message: message,
			Details: details,
		},
	})
}

// TooManyRequests sends a 429 too many requests response
func TooManyRequests(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusTooManyRequests).JSON(Response{
		Success: false,
		Error: ErrorDetail{
			Code:    "TOO_MANY_REQUESTS",
			Message: message,
		},
	})
}

// InternalServerError sends a 500 internal server error response
func InternalServerError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Error: ErrorDetail{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: message,
		},
	})
}

// ServiceUnavailable sends a 503 service unavailable response
func ServiceUnavailable(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusServiceUnavailable).JSON(Response{
		Success: false,
		Error: ErrorDetail{
			Code:    "SERVICE_UNAVAILABLE",
			Message: message,
		},
	})
}

// Error sends a custom error response with status code, error code, and message
func Error(c *fiber.Ctx, statusCode int, errorCode string, message string) error {
	return c.Status(statusCode).JSON(Response{
		Success: false,
		Error: ErrorDetail{
			Code:    errorCode,
			Message: message,
		},
	})
}
