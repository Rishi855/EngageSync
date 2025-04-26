package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
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

// Main function
func main() {
	// Create a new router
	r := mux.NewRouter()

	// Prefix `/api/` for all API routes
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Define routes under `/api/` and apply CORS middleware
	apiRouter.HandleFunc("/signup", enableCORS(signupHandler)).Methods("POST")
	apiRouter.HandleFunc("/login", enableCORS(loginHandler)).Methods("POST")

	apiRouter.HandleFunc("/ideas", enableCORS(authMiddleware(GetAllIdeasHandler))).Methods("GET")
	apiRouter.HandleFunc("/idea", enableCORS(authMiddleware(CreateIdeaHandler))).Methods("POST")
	apiRouter.HandleFunc("/idea/delete", enableCORS(authMiddleware(DeleteIdeaHandler))).Methods("DELETE")

	apiRouter.HandleFunc("/comments", enableCORS(authMiddleware(GetAllCommentsByIdeaIDHandler))).Methods("GET")
	apiRouter.HandleFunc("/comment", enableCORS(authMiddleware(CreateCommentHandler))).Methods("POST")
	apiRouter.HandleFunc("/comment/delete", enableCORS(authMiddleware(DeleteCommentHandler))).Methods("DELETE")

	apiRouter.HandleFunc("/projects", enableCORS(authMiddleware(GetAllProjectsHandler))).Methods("GET")
	apiRouter.HandleFunc("/project", enableCORS(authMiddleware(CreateProjectHandler))).Methods("POST")
	apiRouter.HandleFunc("/project/delete", enableCORS(authMiddleware(DeleteProjectHandler))).Methods("DELETE")

	// Start the server
	log.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
