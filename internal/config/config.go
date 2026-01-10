package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

// Config содержит всю конфигурацию приложения
type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

// AppConfig содержит конфигурацию приложения
type AppConfig struct {
	Name    string
	Version string
	Port    string
	Debug   bool
}

// DatabaseConfig содержит конфигурацию подключения к базе данных
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int // в минутах
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	log.Println("Loading application configuration...")

	cfg := &Config{
		App: AppConfig{
			Name:    getEnv("APP_NAME", "GoToDoList"),
			Version: getEnv("APP_VERSION", "1.0.0"),
			Port:    getEnv("APP_PORT", "8080"),
			Debug:   getEnvAsBool("APP_DEBUG", true),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "todolist"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: getEnvAsInt("DB_CONN_MAX_LIFETIME", 5),
		},
	}

	// Валидация обязательных параметров
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	log.Println("Configuration loaded successfully")
	return cfg, nil
}

// validate проверяет корректность конфигурации
func (c *Config) validate() error {
	if c.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}

	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	return nil
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает значение переменной окружения как int или возвращает значение по умолчанию
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool получает значение переменной окружения как bool или возвращает значение по умолчанию
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
