package config_test

import (
	"os"
	"testing"

	 "booking_system/booking_service/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestLoadConfig_Success(t *testing.T) {
	// Set required env vars
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	os.Setenv("PORT", "9090")
	os.Setenv("PROJECT_ROOT", "/test/root")
	os.Setenv("BASE_URL", "http://localhost:9090")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("LOG_DIR", "./logs")
	os.Setenv("LOG_FILE", "test.log")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("PORT")
		os.Unsetenv("PROJECT_ROOT")
		os.Unsetenv("BASE_URL")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("LOG_DIR")
		os.Unsetenv("LOG_FILE")
	}()

	cfg, err := config.LoadConfig("booking-service", "nonexistent-env")

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "postgres://user:pass@localhost:5432/testdb", cfg.DatabaseURL)
	assert.Equal(t, "/test/root", cfg.ProjectRoot)
	assert.Equal(t, "http://localhost:9090", cfg.BaseUrl)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, "./logs", cfg.Log.Directory)
	assert.Equal(t, "test.log", cfg.Log.FileName)
}

func TestLoadConfig_DefaultValues(t *testing.T) {
	// Only set required DATABASE_URL
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/testdb")
	defer os.Unsetenv("DATABASE_URL")

	// Clear optional vars to test defaults
	os.Unsetenv("PORT")
	os.Unsetenv("PROJECT_ROOT")
	os.Unsetenv("BASE_URL")
	os.Unsetenv("LOG_LEVEL")
	os.Unsetenv("LOG_DIR")
	os.Unsetenv("LOG_FILE")

	cfg, err := config.LoadConfig("booking-service", "nonexistent-env")

	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Port)            // default
	assert.Equal(t, "info", cfg.Log.Level)        // default
	assert.Equal(t, "./logs", cfg.Log.Directory)  // default
	assert.Equal(t, "booking-service.log", cfg.Log.FileName) // default = serviceName+".log"
}

func TestLoadConfig_MissingDatabaseURL(t *testing.T) {
	os.Unsetenv("DATABASE_URL")

	cfg, err := config.LoadConfig("booking-service", "nonexistent-env")

	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "DATABASE_URL")
}

