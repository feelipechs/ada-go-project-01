package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	Database    DatabaseConfig
	Port        string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	db := DatabaseConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		User:     getEnv("POSTGRES_USER", "app"),
		Password: getEnv("POSTGRES_PASSWORD", "app"),
		Name:     getEnv("POSTGRES_DB", "app"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
	}

	return Config{
		DatabaseURL: getEnv("DATABASE_URL", db.URL()),
		Database:    db,
		Port:        getEnv("PORT", "8080"),
	}
}

func (d DatabaseConfig) URL() string {
	return "postgres://" + d.User + ":" + d.Password + "@" + d.Host + ":" + d.Port + "/" + d.Name + "?sslmode=" + d.SSLMode
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		d, err := time.ParseDuration(value)
		if err == nil {
			return d
		}
	}
	return fallback
}
