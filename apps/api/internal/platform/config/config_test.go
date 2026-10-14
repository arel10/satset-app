package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Load_DevelopmentDefaults(t *testing.T) {
	os.Clearenv()
	_ = os.Setenv("APP_ENV", "development")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "development", cfg.AppEnv)
	assert.Equal(t, "8080", cfg.AppPort)
	assert.Equal(t, 25, cfg.DBMaxOpenConns)
	assert.Equal(t, 10, cfg.DBMaxIdleConns)
	assert.Contains(t, cfg.DSN(), "satset_db")
	assert.False(t, cfg.IsProduction())
}

func TestConfig_Validate_InvalidPort(t *testing.T) {
	cfg := &Config{
		AppEnv:         "development",
		AppPort:        "invalid_port",
		DBMaxOpenConns: 25,
		DBMaxIdleConns: 10,
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid PORT")
}

func TestConfig_Validate_InvalidPooling(t *testing.T) {
	cfg := &Config{
		AppEnv:         "development",
		AppPort:        "8080",
		DBMaxOpenConns: 5,
		DBMaxIdleConns: 10, // Idle > Open is invalid
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be greater than DB_MAX_OPEN_CONNS")
}

func TestConfig_Validate_ProductionFailFast(t *testing.T) {
	cfg := &Config{
		AppEnv:          "production",
		AppPort:         "8080",
		DBMaxOpenConns:  25,
		DBMaxIdleConns:  10,
		JWTAccessSecret: "change_me_default",
	}

	err := cfg.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_ACCESS_SECRET must be set to a secure secret")
	assert.True(t, cfg.IsProduction())
}
