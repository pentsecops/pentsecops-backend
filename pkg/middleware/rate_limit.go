package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/pentsecops/backend/internal/config"
	"github.com/pentsecops/backend/pkg/utils"
)

// RateLimit returns a rate limiting middleware
func RateLimit(cfg *config.RateLimitConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cfg.Max,
		Expiration: cfg.Duration,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return utils.TooManyRequests(c, "Rate limit exceeded. Please try again later.")
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		Storage:                nil, // Use in-memory storage
	})
}

// AuthRateLimit returns a rate limiting middleware specifically for auth endpoints
// Limits to 15 requests per minute per IP as per security requirements
func AuthRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        15,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use IP address as key
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":       "Too many authentication attempts",
				"retry_after": 60, // seconds
			})
		},
		SkipFailedRequests:     false,
		SkipSuccessfulRequests: false,
		Storage:                nil, // Use in-memory storage
	})
}
