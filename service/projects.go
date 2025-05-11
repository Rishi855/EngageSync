package service

import (
	"encoding/json"
	// "strconv"

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
	type Project struct {
		ProjectID   string    `json:"project_id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		ManagerID   string    `json:"manager_id"`
		CreatedAt   time.Time `json:"created_at"`
	}

	type Response struct {
		Status  int       `json:"status"`
		Message string    `json:"message"`
		Data    []Project `json:"data,omitempty"`
	}

	rows, err := db.Query(`SELECT projectid, name, description, managerid, createdat FROM kanaka.projects`)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Failed to fetch projects"})
		return
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var proj Project
		err := rows.Scan(&proj.ProjectID, &proj.Name, &proj.Description, &proj.ManagerID, &proj.CreatedAt)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Error reading project data"})
			return
		}
		projects = append(projects, proj)
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Projects fetched successfully",
		Data:    projects,
	})
}


func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	type Project struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ManagerID   string `json:"manager_id"`
	}

	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	var project Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Invalid request body"})
		return
	}

	if project.Name == "" || project.ManagerID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Missing required fields"})
		return
	}

	projectID := uuid.New()

	query := `
		INSERT INTO kanaka.projects (projectid, name, description, managerid)
		VALUES ($1, $2, $3, $4)
	`

	_, err := db.Exec(query, projectID, project.Name, project.Description, project.ManagerID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Failed to create project"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{Status: http.StatusCreated, Message: "Project created successfully"})
}


// DeleteIdeaHandler deletes an idea by its ID
func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	projectID := r.URL.Query().Get("id")
	if projectID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Missing project ID"})
		return
	}

	query := `DELETE FROM kanaka.projects WHERE projectid = $1`
	result, err := db.Exec(query, projectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Error deleting project"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Status: http.StatusNotFound, Message: "Project not found or already deleted"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: http.StatusOK, Message: "Project deleted successfully"})
}
