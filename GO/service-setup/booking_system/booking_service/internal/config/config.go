package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type LogConfig struct {
	Level     string
	Directory string
	FileName  string
}

type Config struct {
	Port        string
	DatabaseURL string
	ProjectRoot string
	BaseUrl     string
	Log         LogConfig
}

func LoadConfig(serviceName, env string) (*Config, error) {
	envFile := fmt.Sprintf("envs/.env.%s", env)

	if err := godotenv.Load(envFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("ENV file isn't present, continuing without it:", envFile)
		} else {
			return nil, fmt.Errorf("failed to load %s: %w", envFile, err)
		}
	}

	config := &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		ProjectRoot: getEnv("PROJECT_ROOT", ""),
		BaseUrl:     getEnv("BASE_URL", ""),

		// Logging config
		Log: LogConfig{
			Level:     getEnv("LOG_LEVEL", "info"),
			Directory: getEnv("LOG_DIR", "./logs"),
			FileName:  getEnv("LOG_FILE", serviceName+".log"),
		},
	}

	if config.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
