package db

import (
	"backend-2/api/cmd/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDB() (*gorm.DB, error) {
	host := utils.GetEnv("DB_HOST", "localhost")
	user := utils.GetEnv("DB_USER", "postgres")
	pass := utils.GetEnv("DB_PASS", "urmomgae69420")
	name := utils.GetEnv("DB_NAME", "backend")
	port := utils.GetEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, pass, name, port)

	// log.Fatal(dsn)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}