package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv               string
	Port                 string
	DatabaseURL          string
	JWTSecret            string
	TokenTTLMinutes      int
	DefaultAdminName     string
	DefaultAdminEmail    string
	DefaultAdminPassword string
	FrontendURL          string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		AppEnv:               getEnv("APP_ENV", "development"),
		Port:                 getEnv("PORT", "8080"),
		DatabaseURL:          getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/techbrein_pm?sslmode=disable"),
		JWTSecret:            getEnv("JWT_SECRET", "change-me"),
		TokenTTLMinutes:      getEnvAsInt("TOKEN_TTL_MINUTES", 60),
		DefaultAdminName:     getEnv("DEFAULT_ADMIN_NAME", "Admin User"),
		DefaultAdminEmail:    getEnv("DEFAULT_ADMIN_EMAIL", "admin@techbrein.local"),
		DefaultAdminPassword: getEnv("DEFAULT_ADMIN_PASSWORD", "Admin@123"),
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return fallback
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("invalid int env for %s, using fallback", key)
		return fallback
	}

	return value
}
