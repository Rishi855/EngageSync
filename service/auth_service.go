package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func init() {
	InitConfig()
}

type User struct {
	ID        int       `json:"id"`
	TanentId  string    `json:"tanent_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	CreatedAt time.Time `json:"created_at"`
	Password  string    `json:"password"`
}

type CustomClaims struct {
	jwt.RegisteredClaims
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"createdAt"`
}

var jwtKey = []byte("my_secret_key")

func DemoRequest(w http.ResponseWriter, r *http.Request) {
	var user User
	// Decode the incoming request body
	json.NewDecoder(r.Body).Decode(&user)

	// Set the current time for createdAt
	createdAt := time.Now()

	_, err := db.Exec(
		`INSERT INTO users (tenant_id, name, email, is_admin, created_at, password) 
		VALUES ('kanaka', $1, $2, $3, $4, $5)`,
		user.Name, user.Email, user.IsAdmin, createdAt, user.Password,
	)
	if err != nil {
		log.Printf("Signup error: %v", err)
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User created")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds User
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var storedUser struct {
		ID       string
		Email    string
		Password string
		TenantID string
		Role     string
	}

	// Query GlobalUsers table
	err = db.QueryRow(`
		SELECT globaluserid, email, password, tenantid, role 
		FROM globalusers 
		WHERE email = $1`, creds.Email).
		Scan(&storedUser.ID, &storedUser.Email, &storedUser.Password, &storedUser.TenantID, &storedUser.Role)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Password check — uncomment this once you're hashing passwords
	// err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Define custom claims
	claims := jwt.MapClaims{
		"sub":       storedUser.ID,
		"email":     storedUser.Email,
		"tenant_id": storedUser.TenantID,
		"role":      storedUser.Role,
		"exp":       expirationTime.Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Return the token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
func getSchemaFromTenantID(tenantID string) string {
	var schemaName string
	err := db.QueryRow(`
		SELECT SchemaName 
		FROM TenantRegistry 
		WHERE TenantID = $1 AND IsActive = TRUE
	`, tenantID).Scan(&schemaName)

	if err != nil {
		// Log or handle error as needed (you might return empty string or panic depending on your strategy)
		fmt.Printf("Error fetching schema for tenant %s: %v\n", tenantID, err)
		return ""
	}

	return schemaName
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	userDetails, err := ExtractUserFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if userDetails.Role != "Admin" {
		http.Error(w, "Only admin can add users", http.StatusForbidden)
		return
	}

	tenantID := userDetails.TenantID

	var user struct {
		UserID       string  `json:"user_id"`
		Name         string  `json:"name"`
		Email        string  `json:"email"`
		PasswordHash string  `json:"password"`
		PhotoURL     *string `json:"photo_url,omitempty"`
		BirthDate    *string `json:"birth_date,omitempty"`
		Department   *string `json:"department,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userRole := "User"

	if user.UserID == "" {
		user.UserID = uuid.New().String()
	}

	_, err = db.Exec(`
		INSERT INTO globalusers (globaluserid, email, password, tenantid, role)
		VALUES ($1, $2, $3, $4, $5)`,
		user.UserID, user.Email, user.PasswordHash, tenantID, userRole)
	if err != nil {
		http.Error(w, "Error inserting into globalusers: "+err.Error(), http.StatusInternalServerError)
		return
	}

	schema := getSchemaFromTenantID(tenantID)
	query := fmt.Sprintf(`INSERT INTO "%s".users 
	(userID, tenantID, name, email, passwordhash, photourl, birthdate, department, role) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`, schema)

	_, err = db.Exec(query,
		user.UserID, tenantID, user.Name, user.Email, user.PasswordHash,
		user.PhotoURL, user.BirthDate, user.Department, userRole)
	if err != nil {
		http.Error(w, "Error inserting into tenant users: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User successfully created",
		"user_id": user.UserID,
	})
}

type AuthenticatedUser struct {
	Role     string
	TenantID string
}

func ExtractUserFromToken(r *http.Request) (*AuthenticatedUser, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, fmt.Errorf("missing or invalid Authorization header")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("unauthorized token")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("missing role in token")
	}

	tenantID, ok := claims["tenant_id"].(string)
	if !ok {
		return nil, fmt.Errorf("missing tenant ID in token")
	}

	return &AuthenticatedUser{
		Role:     role,
		TenantID: tenantID,
	}, nil
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]

		// Use CustomClaims to parse the token
		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Accessing custom claims (isAdmin, createdAt)
		isAdmin := claims.IsAdmin
		createdAt := claims.CreatedAt

		// Optionally: set user ID and other claims in context for later use
		r.Header.Set("UserID", claims.Subject)
		r.Header.Set("IsAdmin", fmt.Sprint(isAdmin))              // Optional, if needed in the request context
		r.Header.Set("CreatedAt", createdAt.Format(time.RFC3339)) // Optional, if needed in the request context

		next(w, r)
	}
}

func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Or specify exact origin
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func InsertInitialUsers() {
	query1 := `
INSERT INTO TenantRegistry (
    TenantID, OrgName, SchemaName
) VALUES (
    'a6d892f4-39e8-4f5d-9c2e-d74f4a9be3cf',
    'Kanaka Org',
    'kanaka'
)`

	query2 := `
INSERT INTO GlobalUsers (
    GlobalUserID, Email, Password, TenantID, Role
) VALUES (
    'b3f9d1de-12e2-445c-9c69-7a81c431fd0c',
    'admin@gmail.com',
    'admin',  -- Use hashed version in production
    'a6d892f4-39e8-4f5d-9c2e-d74f4a9be3cf',
    'Admin'
)`

	query3 := `
INSERT INTO kanaka.Users (
    UserID, TenantID, Name, Email, PasswordHash, PhotoURL, BirthDate, Department, Role
) VALUES (
    'f18e304c-2e83-4f67-8f56-6824f579bb8f',
    'a6d892f4-39e8-4f5d-9c2e-d74f4a9be3cf',
    'admin',
    'admin@gmail.com',
    'admin',
    NULL,
    NULL,
    'Administration',
    'Admin'
)`

	if _, err := db.Exec(query1); err != nil {
		log.Println("executing query1:", err)
	}
	if _, err := db.Exec(query2); err != nil {
		log.Println("executing query2:", err)
	}
	if _, err := db.Exec(query3); err != nil {
		log.Println("executing query3:", err)
	}
	log.Println("Initial users inserted successfully\nDon't panic You can start using port given below")
	log.Printf("\n\n#### User default username as '%s' and password as '%s'\n\n", "admin@gmail.com", "admin")
}
