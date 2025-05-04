package service

import (
	"encoding/json"
	"net/http"
	Constant "github.com/Rishi855/engagesync/VAR"
)

func GetAllTechnologiesHandler(w http.ResponseWriter, r *http.Request) {
	technologies := []string{
		Constant.Golang.String(),
		Constant.NodeJS.String(),
		Constant.Python.String(),
		Constant.Java.String(),
		Constant.React.String(),
		Constant.Angular.String(),
		Constant.Vue.String(),
		Constant.Postgres.String(),
		Constant.MongoDB.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(technologies)
}

func GetAllRolesHandler(w http.ResponseWriter, r *http.Request) {
	roles := []string{
		Constant.Member.String(),
		Constant.Manager.String(),
		Constant.TechLead.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, tenant_id, name, email, is_admin, created_at, password
		FROM users
	`
	rows, err := queryRows(query)
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.TanentId,
			&user.Name,
			&user.Email,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.Password,
		); err != nil {
			http.Error(w, "Error scanning user data", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
