package service

import (
	"database/sql"
	"fmt"
	"log"

	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB
func init() {
    err := godotenv.Load()
    if err != nil {
        log.Println("Error loading .env file:", err) // Changed to log.Println to not fatal
    }
    
    host := os.Getenv("DB_HOST")     // Add this
    port := os.Getenv("DB_PORT")     // Add this
    user := os.Getenv("DB_USER")
    password := os.Getenv("DB_PASSWORD")
    dbname := os.Getenv("DB_NAME")
    sslmode := os.Getenv("DB_SSLMODE")

    connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, password, dbname, sslmode)

    db, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to open DB:", err)
    }

    if err = db.Ping(); err != nil {
        log.Fatal("Failed to ping DB:", err)
    }
}

func execQuery(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func queryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func queryRows(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}
