package main

import (
	"encoding/json"
	// "log"
	"net/http"
	// "os"
	"time"
	// STATIC "websocket-demo/VAR" // Assuming this is where TANENT is declared

	// "github.com/joho/godotenv"
)

// var SCHEMA string

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	SCHEMA = os.Getenv("DB_SCHEMA")
// }

type Idea struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
	Likes     int       `json:"likes"`
	Comments  int       `json:"comments"`
}

// GetAllIdeasHandler retrieves all ideas
func GetAllIdeasHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, title, user_id, created_at, content, total_likes, total_comments 
		FROM ` + SCHEMA + `.ideas
	`
	rows, err := queryRows(query)
	if err != nil {
		http.Error(w, "Error fetching ideas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ideas []Idea
	for rows.Next() {
		var idea Idea
		if err := rows.Scan(&idea.ID, &idea.Title, &idea.UserID, &idea.CreatedAt, &idea.Text, &idea.Likes, &idea.Comments); err != nil {
			http.Error(w, "Error scanning ideas", http.StatusInternalServerError)
			return
		}
		ideas = append(ideas, idea)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ideas)
}

// CreateIdeaHandler creates a new idea
func CreateIdeaHandler(w http.ResponseWriter, r *http.Request) {
	var idea Idea
	if err := json.NewDecoder(r.Body).Decode(&idea); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO ` + SCHEMA + `.ideas 
	(title, user_id, created_at, content, total_likes, total_comments) 
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

	_, err := execQuery(query, idea.Title, idea.UserID, time.Now(), idea.Text, 0, 0)
	if err != nil {
		http.Error(w, "Error creating idea", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Idea created")
}

// DeleteIdeaHandler deletes an idea by its ID
func DeleteIdeaHandler(w http.ResponseWriter, r *http.Request) {
	ideaID := r.URL.Query().Get("id")
	if ideaID == "" {
		http.Error(w, "Idea ID is required", http.StatusBadRequest)
		return
	}

	query := `
		DELETE FROM ` + SCHEMA + `.ideas 
		WHERE id = $1
	`
	_, err := execQuery(query, ideaID)
	if err != nil {
		http.Error(w, "Error deleting idea", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Idea deleted")
}
