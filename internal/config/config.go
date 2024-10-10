package config

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	ServerAddress string
	DatabaseURL   string
	JWTSecret     string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No '.env' file found, using environment variable.")
	}

	serverAddress := getEnv("SERVER_ADDRESS", ":3000")
	databaseurl := getEnv("DATABASE_URL", "")
	jwtSecret := getEnv("JWT_SECRET", "")

	return &Config{
		ServerAddress: serverAddress,
		DatabaseURL:   databaseurl,
		JWTSecret:     jwtSecret,
	}

}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
