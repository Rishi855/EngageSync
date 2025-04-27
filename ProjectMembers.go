package main

import (
	"encoding/json"
	Constant "websocket-demo/VAR"

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
	ProjectId  int    `json:"project_id"`
	UserID     int       `json:"user_id"`
	JoinedAt   time.Time `json:"joined_at"`
	Role       string    `json:"role"`
	Technology string    `json:"technnology"`
}

func GetAllProjectMembersHandler(w http.ResponseWriter, r *http.Request) {
	projectId := r.URL.Query().Get("id")
	if projectId == "" {
		http.Error(w, "Missing project_id query parameter", http.StatusBadRequest)
		return
	}

	query := `
		SELECT id, project_id, user_id, joined_at, role, technology
		FROM ` + SCHEMA + `.project_members
		WHERE project_id = $1
	`

	rows, err := queryRows(query, projectId) // <- pass projectId properly
	if err != nil {
		http.Error(w, "Error fetching project members", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var members []ProjectMembers
	for rows.Next() {
		var id int
		var projId int
		var userId int
		var joinedAt time.Time
		var roleInt int
		var techInt int

		if err := rows.Scan(&id, &projId, &userId, &joinedAt, &roleInt, &techInt); err != nil {
			http.Error(w, "Error scanning project members", http.StatusInternalServerError)
			return
		}

		member := ProjectMembers{
			ID:         id,
			ProjectId:  projId,
			UserID:     userId,
			JoinedAt:   joinedAt,
			Role:       Constant.Position(roleInt).String(),
			Technology: Constant.TechStack(techInt).String(),
		}
		members = append(members, member)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// CreateIdeaHandler creates a new idea
func CreateProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	var member struct {
		ProjectId  int `json:"project_id"`
		UserID     int    `json:"user_id"`
		Role       string `json:"role"`       // will be "Member", "Manager", "TechLead"
		Technology string `json:"technology"` // will be "Golang", "NodeJS", etc.
	}

	if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	roleInt, err := Constant.ParseRole(member.Role)
	if err != nil {
		http.Error(w, "Invalid role value", http.StatusBadRequest)
		return
	}

	techInt, err := Constant.ParseTechnology(member.Technology)
	if err != nil {
		http.Error(w, "Invalid technology value", http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO ` + SCHEMA + `.project_members
	(project_id, user_id, joined_at, role, technology)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	_, err = execQuery(query, member.ProjectId, member.UserID, time.Now(), roleInt, techInt)
	if err != nil {
		http.Error(w, "Error creating project member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Project member created successfully")
}

// DeleteIdeaHandler deletes an idea by its ID
func DeleteProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	ideaID := r.URL.Query().Get("id")
	if ideaID == "" {
		http.Error(w, "Project member ID is required", http.StatusBadRequest)
		return
	}

	query := `
		DELETE FROM ` + SCHEMA + `.project_members  
		WHERE id = $1
	`
	_, err := execQuery(query, ideaID)
	if err != nil {
		http.Error(w, "Error deleting project member", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Project member deleted")
}
