package database

import (
	"fmt"
	"log"
	"os"

	"mini-lms/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "lms_user")
	password := getEnv("DB_PASSWORD", "lms_password")
	dbname := getEnv("DB_NAME", "mini_lms_db")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connection established successfully.")

	// Auto-migrate domain entities
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Course{},
		&domain.Lesson{},
		&domain.Enrollment{},
		&domain.LessonProgress{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database schema: %w", err)
	}

	log.Println("Database auto-migration finished.")
	return db, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
