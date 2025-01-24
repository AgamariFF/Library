package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DBDSN      string
}

func LoadConfig() Config {
	// Загрузка .env файла
	if err := godotenv.Load("./configs/.env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Чтение переменных из окружения
	config := Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBDSN:      getEnv("DB_DSN", "localhost"),
	}

	return config
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
