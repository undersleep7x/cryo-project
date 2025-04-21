package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host string
	Port string
	User string
	Password string
	Name string
}
type AppConfig struct {
	Env       string
	Port      string
	RedisHost string
	RedisPort string
	LoggingPath string
	LoggingPerms string
	DB DBConfig
}

func LoadConfig() *AppConfig {
	env := os.Getenv("ENV")
	if env == ""{
		log.Println("No application environment found; loading with local env config.")
		env = "dev"
	}

	if env == "dev" {
        if err := godotenv.Load("/app/.env.dev"); err != nil {
            log.Printf("No .env.dev file found — relying on system environment: Error - %v", err)
        }
	}

	return &AppConfig{
		Env:       env,
		Port:      getEnv("PORT", "8080"),
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),
		LoggingPath: getEnv("LOGGING_PATH", "logs/apps.log"),
		LoggingPerms: getEnv("LOGGING_PERMS", "0666"),
		DB: DBConfig{
			Host: getEnv("DB_HOST", "postgres"),
			Port: getEnv("DB_PORT", "5432"),
			User: getEnv("DB_USER", "cryouser"),
			Password: getEnv("DB_PASSWORD", "cryopass"),
			Name: getEnv("DB_NAME", "cryo"),
		},
	}
}


func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}