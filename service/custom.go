package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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
	schema, err := GetSchemaByToken(w, r)
	query := fmt.Sprintf(`SELECT userid, tenantid, name, email, photourl, birthdate, department, role FROM "%s".users`, schema)

	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Error fetching users: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}

	for rows.Next() {
		var (
			userID     string
			tenantID   string
			name       string
			email      string
			photoURL   *string
			birthDate  *time.Time
			department *string
			role       string
		)

		if err := rows.Scan(&userID, &tenantID, &name, &email, &photoURL, &birthDate, &department, &role); err != nil {
			http.Error(w, "Error scanning users: "+err.Error(), http.StatusInternalServerError)
			return
		}

		user := map[string]interface{}{
			"user_id":    userID,
			"tenant_id":  tenantID,
			"name":       name,
			"email":      email,
			"photo_url":  photoURL,
			"birth_date": birthDate,
			"department": department,
			"role":       role,
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetSchemaByToken(w http.ResponseWriter, r *http.Request) (string, error) {
	userDetails, err := ExtractUserFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return "", err
	}
	schema := getSchemaFromTenantID(userDetails.TenantID)
	return schema,nil
}