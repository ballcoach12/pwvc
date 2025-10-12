package migrate
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run migrate.go [up|down|version|force <version>]")
	}

	command := os.Args[1]

	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:password@localhost:5432/pwvc?sslmode=disable"
	}

	// Open database connection
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create postgres driver: %v", err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations completed successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		fmt.Println("Rollback completed successfully")
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get version: %v", err)
		}
		fmt.Printf("Version: %d, Dirty: %t\n", version, dirty)
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("Force command requires version number")
		}
		var version int
		if _, err := fmt.Sscanf(os.Args[2], "%d", &version); err != nil {
			log.Fatalf("Invalid version number: %v", err)
		}
		if err := m.Force(version); err != nil {
			log.Fatalf("Failed to force version: %v", err)
		}
		fmt.Printf("Forced to version %d\n", version)
	default:
		log.Fatal("Unknown command. Use: up, down, version, or force <version>")
	}
}