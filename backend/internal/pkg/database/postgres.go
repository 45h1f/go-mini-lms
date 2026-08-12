package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"mini-lms/internal/domain"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dbPath := getEnv("DB_PATH", "mini_lms.db")

	// Ensure directory exists if dbPath includes a folder
	dir := filepath.Dir(dbPath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	log.Printf("SQLite database connected successfully at: %s", dbPath)

	// Enable Foreign Key support for SQLite
	db.Exec("PRAGMA foreign_keys = ON;")

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

	log.Println("SQLite database auto-migration & seeding finished.")
	return db, nil
}

func seedAdminUser(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Where("email = ?", "admin@lms.com").Count(&count)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
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
