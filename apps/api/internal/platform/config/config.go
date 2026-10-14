package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// App
	AppEnv     string
	AppPort    string
	AppBaseURL string

	// Database (MySQL)
	DBDSN                    string
	DBHost                   string
	DBPort                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	DBMaxOpenConns           int
	DBMaxIdleConns           int
	DBConnMaxLifetimeMinutes int

	// Redis
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Auth & JWT
	JWTAccessSecret     string
	JWTRefreshSecret    string
	JWTAccessTTLMinutes int
	JWTRefreshTTLDays   int

	// Gemini AI
	GeminiAPIKey             string
	GeminiModel              string
	AIRequestTimeoutSeconds  int
	AIMaxToolCalls           int
	AIMaxInputChars          int
	AIDailyLimitPerUser      int
	AIGlobalRPMLimit         int
	AICacheTTLSeconds        int

	// S3 / MinIO
	S3Endpoint  string
	S3Bucket    string
	S3AccessKey string
	S3SecretKey string
}

func Load() (*Config, error) {
	dbMaxOpenConns := getEnvInt("DB_MAX_OPEN_CONNS", 25)
	dbMaxIdleConns := getEnvInt("DB_MAX_IDLE_CONNS", 10)
	dbConnMaxLifetime := getEnvInt("DB_CONN_MAX_LIFETIME_MINUTES", 30)
	redisDB := getEnvInt("REDIS_DB", 0)

	jwtAccessTTL := getEnvInt("JWT_ACCESS_TTL_MINUTES", 15)
	jwtRefreshTTL := getEnvInt("JWT_REFRESH_TTL_DAYS", 14)

	aiTimeout := getEnvInt("AI_REQUEST_TIMEOUT_SECONDS", 30)
	aiMaxToolCalls := getEnvInt("AI_MAX_TOOL_CALLS", 5)
	aiMaxInputChars := getEnvInt("AI_MAX_INPUT_CHARS", 2000)
	aiDailyLimit := getEnvInt("AI_DAILY_LIMIT_PER_USER", 20)
	aiRPMLimit := getEnvInt("AI_GLOBAL_RPM_LIMIT", 10)
	aiCacheTTL := getEnvInt("AI_CACHE_TTL_SECONDS", 300)

	cfg := &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("PORT", getEnv("APP_PORT", "8080")),
		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:3000"),

		DBDSN:                    getEnv("DB_DSN", ""),
		DBHost:                   getEnv("DB_HOST", "localhost"),
		DBPort:                   getEnv("DB_PORT", "3306"),
		DBUser:                   getEnv("DB_USER", "satset"),
		DBPassword:               getEnv("DB_PASSWORD", "satset_secret"),
		DBName:                   getEnv("DB_NAME", "satset_db"),
		DBMaxOpenConns:           dbMaxOpenConns,
		DBMaxIdleConns:           dbMaxIdleConns,
		DBConnMaxLifetimeMinutes: dbConnMaxLifetime,

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		JWTAccessSecret:     getEnv("JWT_ACCESS_SECRET", "change_me_access_secret_satset_development"),
		JWTRefreshSecret:    getEnv("JWT_REFRESH_SECRET", "change_me_refresh_secret_satset_development"),
		JWTAccessTTLMinutes: jwtAccessTTL,
		JWTRefreshTTLDays:   jwtRefreshTTL,

		GeminiAPIKey:            getEnv("GEMINI_API_KEY", ""),
		GeminiModel:             getEnv("GEMINI_MODEL", "gemini-1.5-flash"),
		AIRequestTimeoutSeconds: aiTimeout,
		AIMaxToolCalls:          aiMaxToolCalls,
		AIMaxInputChars:         aiMaxInputChars,
		AIDailyLimitPerUser:     aiDailyLimit,
		AIGlobalRPMLimit:        aiRPMLimit,
		AICacheTTLSeconds:       aiCacheTTL,

		S3Endpoint:  getEnv("S3_ENDPOINT", "localhost:9000"),
		S3Bucket:    getEnv("S3_BUCKET", "satset-assets"),
		S3AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
		S3SecretKey: getEnv("S3_SECRET_KEY", "minioadminpassword"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	var errs []string

	// Validate Port
	portNum, err := strconv.Atoi(c.AppPort)
	if err != nil || portNum < 1 || portNum > 65535 {
		errs = append(errs, fmt.Sprintf("invalid PORT/APP_PORT: %s", c.AppPort))
	}

	// Validate DB Pooling
	if c.DBMaxIdleConns > c.DBMaxOpenConns {
		errs = append(errs, fmt.Sprintf("DB_MAX_IDLE_CONNS (%d) cannot be greater than DB_MAX_OPEN_CONNS (%d)", c.DBMaxIdleConns, c.DBMaxOpenConns))
	}

	// Production strict fail-fast validation
	if c.IsProduction() {
		if c.JWTAccessSecret == "" || strings.Contains(c.JWTAccessSecret, "change_me") {
			errs = append(errs, "JWT_ACCESS_SECRET must be set to a secure secret in production")
		}
		if c.JWTRefreshSecret == "" || strings.Contains(c.JWTRefreshSecret, "change_me") {
			errs = append(errs, "JWT_REFRESH_SECRET must be set to a secure secret in production")
		}
		if c.DBPassword == "" || c.DBPassword == "satset_secret" {
			errs = append(errs, "DB_PASSWORD cannot be default or empty in production")
		}
		if c.DBDSN == "" && (c.DBHost == "" || c.DBUser == "" || c.DBName == "") {
			errs = append(errs, "DB_DSN or database connection credentials must be fully provided in production")
		}
	}

	if len(errs) > 0 {
		return errors.New("configuration validation failed: " + strings.Join(errs, "; "))
	}

	return nil
}

func (c *Config) DSN() string {
	if c.DBDSN != "" {
		return c.DBDSN
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func (c *Config) DBConnMaxLifetime() time.Duration {
	return time.Duration(c.DBConnMaxLifetimeMinutes) * time.Minute
}

func (c *Config) IsProduction() bool {
	return strings.ToLower(c.AppEnv) == "production"
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return defaultVal
}
