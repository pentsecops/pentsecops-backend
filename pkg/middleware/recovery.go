package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Recovery returns a panic recovery middleware
func Recovery() fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
			// Log panic (can be enhanced with proper logger later)
			fmt.Printf("Panic recovered: %v, path: %s, method: %s\n", e, c.Path(), c.Method())
		},
	})
}
