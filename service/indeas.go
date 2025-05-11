package service

import (
	"encoding/json"
	// "log"
	"net/http"
	// "os"
	"time"

	"github.com/google/uuid"
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
	rows, err := db.Query(`
		SELECT 
			ideaid, title, description, submittedby, submittedat,
			isanonymous, isapproved, approvedby, approvedat
		FROM kanaka.ideas
	`)
	if err != nil {
		http.Error(w, "Failed to fetch ideas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Idea struct {
		IdeaID      string     `json:"idea_id"`
		Title       string     `json:"title"`
		Description string     `json:"description"`
		SubmittedBy string     `json:"submitted_by"`
		SubmittedAt time.Time  `json:"submitted_at"`
		IsAnonymous bool       `json:"is_anonymous"`
		IsApproved  *bool      `json:"is_approved"` // Nullable
		ApprovedBy  *string    `json:"approved_by"` // Nullable
		ApprovedAt  *time.Time `json:"approved_at"` // Nullable
	}

	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    []Idea `json:"data,omitempty"`
	}

	var ideas []Idea

	for rows.Next() {
		var idea Idea
		err := rows.Scan(
			&idea.IdeaID,
			&idea.Title,
			&idea.Description,
			&idea.SubmittedBy,
			&idea.SubmittedAt,
			&idea.IsAnonymous,
			&idea.IsApproved,
			&idea.ApprovedBy,
			&idea.ApprovedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning idea data", http.StatusInternalServerError)
			return
		}
		ideas = append(ideas, idea)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Projects fetched successfully",
		Data:    ideas,
	})
}

// CreateIdeaHandler creates a new idea
func CreateIdeaHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	var idea struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		SubmittedBy string `json:"submitted_by"`
		IsAnonymous bool   `json:"is_anonymous"`
	}

	if err := json.NewDecoder(r.Body).Decode(&idea); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Invalid request body"})
		return
	}

	// Generate a new UUID for IdeaID
	ideaID := uuid.New()

	// Insert into the database
	_, err := db.Exec(`
		INSERT INTO kanaka.ideas (
			ideaid, title, description, submittedby, isanonymous
		) VALUES ($1, $2, $3, $4, $5)`,
		ideaID, idea.Title, idea.Description, idea.SubmittedBy, idea.IsAnonymous,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Failed to create idea"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Status: http.StatusOK, Message: "Idea created successfully"})
}

// DeleteIdeaHandler deletes an idea by its ID
func DeleteIdeaHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	// Get IdeaID from query parameters
	ideaID := r.URL.Query().Get("id")
	if ideaID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Missing idea ID"})
		return
	}

	// Execute delete query
	result, err := db.Exec(`DELETE FROM kanaka.ideas WHERE ideaid = $1`, ideaID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Failed to delete idea"})
		return
	}

	// Check if any row was actually deleted
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Status: http.StatusNotFound, Message: "Idea not found"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Status: http.StatusOK, Message: "Idea deleted successfully"})
}
