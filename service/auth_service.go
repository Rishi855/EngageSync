package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func SignupHandler(w http.ResponseWriter, r *http.Request) {
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
	json.NewDecoder(r.Body).Decode(&creds)

	var storedUser User
	// Retrieve the user from the database
	err := db.QueryRow(`SELECT id, "name", email, "password", is_admin, created_at FROM users WHERE email=$1`, creds.Email).
		Scan(&storedUser.ID, &storedUser.Name, &storedUser.Email, &storedUser.Password, &storedUser.IsAdmin, &storedUser.CreatedAt)

	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Check password (ensure to uncomment bcrypt logic when implementing hashing)
	// err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Set token expiration time
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create custom claims
	claims := &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   fmt.Sprint(storedUser.ID),
		},
		IsAdmin:   storedUser.IsAdmin,
		CreatedAt: storedUser.CreatedAt, // Keep it as time.Time, not a string
	}

	// Create the token with custom claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	// Send the token as a response
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func AddUserHandler(w http.ResponseWriter, r *http.Request) {

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
		log.Printf("User create error: %v", err)
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "Email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User Added")
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
