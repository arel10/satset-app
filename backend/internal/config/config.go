package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	App       AppConfig
	DB        DBConfig
	JWT       JWTConfig
	Gemini    GeminiConfig
	AI        AIConfig
	RateLimit RateLimitConfig
	CORS      CORSConfig
}

// AppConfig holds application-level settings.
type AppConfig struct {
	Name     string
	Env      string
	Port     string
	Timezone *time.Location
}

// DBConfig holds database connection settings.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	DSN      string
}

// JWTConfig holds JWT authentication settings.
type JWTConfig struct {
	AccessSecret        string
	RefreshSecret       string
	AccessExpiryMinutes int
	RefreshExpiryDays   int
}

// GeminiConfig holds Gemini AI API settings.
type GeminiConfig struct {
	APIKey string
	Model  string
}

// AIConfig holds AI usage limits.
type AIConfig struct {
	MaxToolCalls       int
	DailyLimitPerUser  int
	RequestTimeoutSecs int
	MaxQuestionLength  int
}

// RateLimitConfig holds rate limiting settings.
type RateLimitConfig struct {
	LoginMaxAttempts    int
	LoginLockoutMinutes int
}

// CORSConfig holds CORS settings.
type CORSConfig struct {
	FrontendURL string
}

// Load reads environment variables and returns a Config struct.
// It loads from .env file if present (for development).
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	timezone, err := time.LoadLocation(getEnv("APP_TIMEZONE", "Asia/Jakarta"))
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %w", err)
	}

	cfg := &Config{
		App: AppConfig{
			Name:     getEnv("APP_NAME", "SatSet"),
			Env:      getEnv("APP_ENV", "development"),
			Port:     getEnv("APP_PORT", "8080"),
			Timezone: timezone,
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "satset"),
		},
		JWT: JWTConfig{
			AccessSecret:        getEnv("JWT_ACCESS_SECRET", ""),
			RefreshSecret:       getEnv("JWT_REFRESH_SECRET", ""),
			AccessExpiryMinutes: getEnvInt("JWT_ACCESS_EXPIRY_MINUTES", 15),
			RefreshExpiryDays:   getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7),
		},
		Gemini: GeminiConfig{
			APIKey: getEnv("GEMINI_API_KEY", ""),
			Model:  getEnv("GEMINI_MODEL", "gemini-2.0-flash"),
		},
		AI: AIConfig{
			MaxToolCalls:       getEnvInt("AI_MAX_TOOL_CALLS", 5),
			DailyLimitPerUser:  getEnvInt("AI_DAILY_LIMIT_PER_USER", 30),
			RequestTimeoutSecs: getEnvInt("AI_REQUEST_TIMEOUT_SECONDS", 30),
			MaxQuestionLength:  getEnvInt("AI_MAX_QUESTION_LENGTH", 2000),
		},
		RateLimit: RateLimitConfig{
			LoginMaxAttempts:    getEnvInt("LOGIN_MAX_ATTEMPTS", 5),
			LoginLockoutMinutes: getEnvInt("LOGIN_LOCKOUT_MINUTES", 15),
		},
		CORS: CORSConfig{
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
	}

	// Build DSN from individual components
	cfg.DB.DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name,
	)

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// IsDevelopment returns true if the app is running in development mode.
func (c *Config) IsDevelopment() bool {
	return c.App.Env == "development"
}

// IsProduction returns true if the app is running in production mode.
func (c *Config) IsProduction() bool {
	return c.App.Env == "production"
}

// AutoMigrate returns true if database auto-migration is enabled.
func (c *Config) AutoMigrate() bool {
	return getEnvBool("AUTO_MIGRATE", false)
}

// validate checks that required configuration values are present.
func (c *Config) validate() error {
	if c.JWT.AccessSecret == "" {
		return fmt.Errorf("JWT_ACCESS_SECRET is required")
	}
	if c.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}
	if c.DB.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}
	return nil
}

// getEnv reads an environment variable or returns a default value.
func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

// getEnvInt reads an integer environment variable or returns a default value.
func getEnvInt(key string, defaultVal int) int {
	if val, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}

// getEnvBool reads a boolean environment variable or returns a default value.
func getEnvBool(key string, defaultVal bool) bool {
	if val, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal
		}
	}
	return defaultVal
}
