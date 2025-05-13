package service

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type GuessImage struct {
	ImageID     uuid.UUID `json:"image_id"`
	ImageURL    string    `json:"image_url"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	UploadedAt  time.Time `json:"uploaded_at"`
	UploadedBy  uuid.UUID `json:"uploaded_by"`
	IsCompleted bool      `json:"is_completed"`
}

func GetGuessImageHandler(w http.ResponseWriter, r *http.Request) {

	userDetails, err := ExtractUserFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if userDetails.Role != "Admin" {
		http.Error(w, "Only admin can add users", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	category := strings.ToLower(vars["category"])

	// Validate category
	validCategories := map[string]bool{
		"movies":          true,
		"songs":           true,
		"actressandactor": true,
	}

	if !validCategories[category] {
		http.Error(w, "Invalid category", http.StatusBadRequest)
		return
	}

	var img GuessImage

	err = db.QueryRow(`
		SELECT 
			imageid, imageurl, title, description, category, uploadedat, uploadedby, iscompleted
		FROM kanaka.guessimages 
		WHERE LOWER(category) = $1 AND iscompleted = FALSE 
		ORDER BY RANDOM() 
		LIMIT 1
	`, category).Scan(
		&img.ImageID,
		&img.ImageURL,
		&img.Title,
		&img.Description,
		&img.Category,
		&img.UploadedAt,
		&img.UploadedBy,
		&img.IsCompleted,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "No image found for this category", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error querying image", http.StatusInternalServerError)
		return
	}

	// Respond with the image data
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(img)

	// Mark the image as completed
	go func(imageID uuid.UUID) {
		_, err := db.Exec(`
			UPDATE kanaka.guessimages 
			SET iscompleted = TRUE 
			WHERE imageid = $1
		`, imageID)
		if err != nil {
			// Log error if needed
		}
	}(img.ImageID)
}
