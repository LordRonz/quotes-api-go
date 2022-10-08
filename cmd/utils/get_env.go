package utils

import (
	"log"
	"os"

	_ "backend-2/api/docs"

	"github.com/joho/godotenv"
)

func GetEnv(key string, def ...string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	env := os.Getenv(key)
	if env != "" {
		return env
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}