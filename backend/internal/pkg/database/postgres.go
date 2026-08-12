package database

import (
	"fmt"
	"log"
	"os"

	"mini-lms/internal/domain"

	"golang.org/x/crypto/bcrypt"
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

	// Seed default Admin user if not exists
	seedAdminUser(db)

	log.Println("Database auto-migration & seeding finished.")
	return db, nil
}

func seedAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Where("email = ?", "admin@lms.com").Count(&count)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash admin password: %v", err)
			return
		}

		admin := domain.User{
			FullName: "System Administrator",
			Email:    "admin@lms.com",
			Password: string(hashedPassword),
			Role:     domain.RoleAdmin,
		}

		if err := db.Create(&admin).Error; err != nil {
			log.Printf("Failed to seed admin user: %v", err)
		} else {
			log.Println("Default Admin user seeded successfully (admin@lms.com / admin123)")
		}
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
