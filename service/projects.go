package service

import (
	"encoding/json"
	"strconv"
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

type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string       `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	AdminId     int       `json:"admin_id"`
}

// GetAllIdeasHandler retrieves all ideas
func GetAllProjectsHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, name, description, created_at, updated_at, admin_id
		FROM ` + SCHEMA + `.projects
	`
	rows, err := queryRows(query)
	if err != nil {
		http.Error(w, "Error fetching projects", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		if err := rows.Scan(&project.ID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt, &project.AdminId); err != nil {
			http.Error(w, "Error scanning projects", http.StatusInternalServerError)
			return
		}
		projects = append(projects, project)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}
func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var project Project
	var user User

	// Decode the incoming request body to the 'project' struct
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract the user ID from the Authorization token (set by the authMiddleware)
	userID := r.Header.Get("UserID")
	if userID == "" {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Convert userID (string) to int if necessary
	// If your user ID is an integer, you may need to parse it
	parsedUserID, err := strconv.Atoi(userID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusUnauthorized)
		return
	}

	// Set the user ID
	user.ID = parsedUserID

	// Prepare SQL query to insert a new project
	query := `
	INSERT INTO ` + SCHEMA + `.projects 
	(name, description, created_at, updated_at, admin_id) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id` 

	// Execute the query using the provided data and get the inserted ID
	var projectID int
	err = db.QueryRow(query, project.Name, project.Description, time.Now(), time.Now(), user.ID).Scan(&projectID)
	if err != nil {
		http.Error(w, "Error creating project", http.StatusInternalServerError)
		return
	}

	// Send success response with the created project ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Project created successfully",
		"id":      projectID, // Return the newly created project ID
	})
}


// DeleteIdeaHandler deletes an idea by its ID
func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	ideaID := r.URL.Query().Get("id")
	if ideaID == "" {
		http.Error(w, "Projects ID is required", http.StatusBadRequest)
		return
	}

	query := `
		DELETE FROM ` + SCHEMA + `.projects 
		WHERE id = $1
	`
	_, err := execQuery(query, ideaID)
	if err != nil {
		http.Error(w, "Error deleting projects", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Projects deleted")
}
