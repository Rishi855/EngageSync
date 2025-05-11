package service

import (
	"encoding/json"
	// Constant "github.com/Rishi855/engagesync/VAR"

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

type ProjectMembers struct {
	ID         int       `json:"id"`
	ProjectId  int       `json:"project_id"`
	UserID     int       `json:"user_id"`
	JoinedAt   time.Time `json:"joined_at"`
	Role       string    `json:"role"`
	Technology string    `json:"technnology"`
}

func GetAllProjectMembersHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}

	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Missing project ID"})
		return
	}

	query := `
		SELECT pm.userid, u.name, u.email, pm.role
		FROM kanaka.projectmembers pm
		JOIN kanaka.users u ON pm.userid = u.userid
		WHERE pm.projectid = $1
	`

	rows, err := db.Query(query, projectID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Database query error"})
		return
	}
	defer rows.Close()

	var members []map[string]interface{}
	for rows.Next() {
		var userID, name, email, role string
		if err := rows.Scan(&userID, &name, &email, &role); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Failed to scan result"})
			return
		}
		members = append(members, map[string]interface{}{
			"user_id": userID,
			"name":    name,
			"email":   email,
			"role":    role,
		})
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Project members retrieved successfully",
		Data:    members,
	})
}

// CreateIdeaHandler creates a new idea
func CreateProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	type ProjectMember struct {
		ProjectID string `json:"project_id"`
		UserID    string `json:"user_id"`
		Role      string `json:"role"`
	}

	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	var member ProjectMember
	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Invalid request body"})
		return
	}

	if member.ProjectID == "" || member.UserID == "" || member.Role == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Missing required fields"})
		return
	}

	query := `
		INSERT INTO kanaka.projectmembers (projectid, userid, role)
		VALUES ($1, $2, $3)
	`

	_, err := db.Exec(query, member.ProjectID, member.UserID, member.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Error inserting project member"})
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusCreated,
		Message: "Project member added successfully",
	})
}

// DeleteIdeaHandler deletes an idea by its ID
func DeleteProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	projectID := r.URL.Query().Get("project_id")
	userID := r.URL.Query().Get("user_id")

	if projectID == "" || userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Missing project_id or user_id"})
		return
	}

	query := `DELETE FROM kanaka.projectmembers WHERE projectid = $1 AND userid = $2`

	result, err := db.Exec(query, projectID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Error deleting project member"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Status: http.StatusNotFound, Message: "Project member not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Status: http.StatusOK, Message: "Project member deleted successfully"})
}
