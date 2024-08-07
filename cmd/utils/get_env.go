package utils

import (
	"os"

	"github.com/rs/zerolog/log"

	_ "backend-2/api/docs"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Print(".env file not found")
	}
}

func GetEnv(key string, def ...string) string {
	env := os.Getenv(key)
	if env != "" {
		return env
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}
