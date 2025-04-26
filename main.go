package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	// "golang.org/x/crypto/bcrypt"
)

var SCHEMA string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	SCHEMA = os.Getenv("DB_SCHEMA")

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", user, password, dbname, sslmode)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("Server crashed with error: %v", r)
		}
	}()

	mux := http.NewServeMux()

	// Serve Description.md at root "/"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "Description.md")
	})

	// API routes under /api/
	apiMux := http.NewServeMux()

	apiMux.HandleFunc("/signup", enableCORS(signupHandler))
	apiMux.HandleFunc("/login", enableCORS(loginHandler))

	apiMux.HandleFunc("/ideas", enableCORS(authMiddleware(GetAllIdeasHandler)))
	apiMux.HandleFunc("/idea", enableCORS(authMiddleware(CreateIdeaHandler)))
	apiMux.HandleFunc("/idea/delete", enableCORS(authMiddleware(DeleteIdeaHandler)))

	apiMux.HandleFunc("/comments", enableCORS(authMiddleware(GetAllCommentsByIdeaIDHandler)))
	apiMux.HandleFunc("/comment", enableCORS(authMiddleware(CreateCommentHandler)))
	apiMux.HandleFunc("/comment/delete", enableCORS(authMiddleware(DeleteCommentHandler)))

	// Prefix /api routes
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	log.Println("Server started at :8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
