package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables from .env
	currentDir, _ := os.Getwd()
	envPath := currentDir + "/.env"
	if err := godotenv.Load(envPath); err != nil {
		log.Fatal("Error loading .env file from", envPath)
	}

	// Read database credentials from environment variables
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Build connection string (no schema)
	connStr := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=%s",
		user, password, dbname, sslmode)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	// Setup migration driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	// Load migration files
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations", // Correct relative path to migrations
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrate instance: %v", err)
	}

	// Handle CLI command
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration command: up | down | drop")
	}

	var migrationErr error
	switch os.Args[1] {
	case "up":
		migrationErr = m.Up()
	case "down":
		migrationErr = m.Steps(-1)
	// Uncomment for drop command if needed
	// case "drop":
	// 	migrationErr = m.Drop()
	default:
		log.Fatalf("Unknown command: %s", os.Args[1])
	}

	if migrationErr != nil && migrationErr != migrate.ErrNoChange {
		log.Fatalf("Migration error: %v", migrationErr)
	}

	// Log the current migration version
	version, _, err := m.Version()
	if err != nil {
		log.Fatalf("Error fetching migration version: %v", err)
	}
	log.Printf("Migration %s executed successfully. Current version: %d", os.Args[1], version)

	// Get the current migration version applied in the database
	currentVersion, _, err := m.Version()
	if err != nil {
		log.Fatalf("Error fetching current migration version: %v", err)
	}
	log.Printf("Currently applied migration version: %d", currentVersion)

	// Optional: Log output to a file
	logFile, err := os.OpenFile("migration_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	// Set output to log file as well
	log.SetOutput(logFile)
	log.Println("Logging migration details to migration_log.txt")
	log.Printf("Migration %s executed successfully. Current version: %d", os.Args[1], currentVersion)
}