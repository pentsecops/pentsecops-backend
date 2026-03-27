package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/pentsecops/backend/internal/config"
)

// CORS returns a CORS middleware configured from config
func CORS(cfg *config.CORSConfig) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     cfg.AllowedMethods,
		AllowHeaders:     cfg.AllowedHeaders,
		AllowCredentials: true,
		ExposeHeaders:    "Content-Length,Content-Type",
		MaxAge:           86400, // 24 hours
	})
}
