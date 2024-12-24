package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	LoadEnv()
}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func Get(key string, defalutValue string) string {
	cache := make(map[string]string)
	if val, ok := cache[key]; ok {
		return val
	}
	newVal := os.Getenv(key)
	if newVal == "" {
		return defalutValue
	}
	cache[key] = newVal

	return newVal
}
