package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	direction := flag.String("direction", "up", "Migration direction: 'up' or 'down'")
	dbPath := flag.String("db", "mini_lms.db", "Path to SQLite database file")
	flag.Parse()

	migrationsPath := "file://migrations"

	dbURL := fmt.Sprintf("sqlite3://%s", *dbPath)

	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}

	if *direction == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration UP failed: %v", err)
		}
		fmt.Println("Successfully applied migrations UP.")
	} else if *direction == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Migration DOWN failed: %v", err)
		}
		fmt.Println("Successfully rolled back migrations DOWN.")
	} else {
		log.Fatalf("Invalid direction parameter. Use -direction=up or -direction=down")
	}
}
