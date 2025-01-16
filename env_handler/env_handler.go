package env_handler

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func InitialEnvParams() {
	// Загружаем переменные из файла .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Env params is loaded!")
}

func GetEnvParam(key string) string {
	param := os.Getenv(key)

	if param == "" {
		log.Fatalf("Missing required environment variable: %s!", key)
	}

	return param
}
