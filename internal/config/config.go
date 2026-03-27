package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	RateLimit RateLimitConfig
	CORS     CORSConfig
	Cache    CacheConfig
	Logging  LoggingConfig
	ReCAPTCHA ReCAPTCHAConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string
	Port string
	Env  string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	PublicKey            string
	PrivateKey           string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	Max      int
	Duration time.Duration
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins string
	AllowedMethods string
	AllowedHeaders string
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	MaxCost     int64
	NumCounters int64
	BufferItems int64
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// ReCAPTCHAConfig holds reCAPTCHA configuration
type ReCAPTCHAConfig struct {
	SecretKey string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "pentsecops"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Auth: AuthConfig{
			PublicKey:            getEnv("PASETO_PUBLIC_KEY", ""),
			PrivateKey:           getEnv("PASETO_PRIVATE_KEY", ""),
			AccessTokenDuration:  getEnvAsDuration("ACCESS_TOKEN_DURATION", 15*time.Minute),
			RefreshTokenDuration: getEnvAsDuration("REFRESH_TOKEN_DURATION", 7*24*time.Hour),
		},
		RateLimit: RateLimitConfig{
			Max:      getEnvAsInt("RATE_LIMIT_MAX", 100),
			Duration: getEnvAsDuration("RATE_LIMIT_DURATION", 1*time.Minute),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000"),
			AllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,DELETE,PATCH,OPTIONS"),
			AllowedHeaders: getEnv("CORS_ALLOWED_HEADERS", "Origin,Content-Type,Accept,Authorization"),
		},
		Cache: CacheConfig{
			MaxCost:     getEnvAsInt64("CACHE_MAX_COST", 100000000),
			NumCounters: getEnvAsInt64("CACHE_NUM_COUNTERS", 1000000),
			BufferItems: getEnvAsInt64("CACHE_BUFFER_ITEMS", 64),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		ReCAPTCHA: ReCAPTCHAConfig{
			SecretKey: getEnv("RECAPTCHA_SECRET_KEY", ""),
		},
	}

	return cfg, nil
}

// GetDSN returns the database connection string
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

