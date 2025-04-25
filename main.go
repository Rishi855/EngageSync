package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
	// "golang.org/x/crypto/bcrypt"
)

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=engagesync sslmode=disable password=root")
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
    defer func() {
        if r := recover(); r != nil {
            log.Fatalf("Server crashed with error: %v", r)
        }
    }()

    mux := http.NewServeMux()

    // Serve Description.md at root "/"
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "Description.md")
    })

    // API routes under /api/
    apiMux := http.NewServeMux()

    apiMux.HandleFunc("/signup", enableCORS(signupHandler))
    apiMux.HandleFunc("/login", enableCORS(loginHandler))

    apiMux.HandleFunc("/ideas", enableCORS(authMiddleware(GetAllIdeasHandler)))
    apiMux.HandleFunc("/idea", enableCORS(authMiddleware(CreateIdeaHandler)))
    apiMux.HandleFunc("/idea/delete", enableCORS(authMiddleware(DeleteIdeaHandler)))

    apiMux.HandleFunc("/comments", enableCORS(authMiddleware(GetAllCommentsByIdeaIDHandler)))
    apiMux.HandleFunc("/comment", enableCORS(authMiddleware(CreateCommentHandler)))
    apiMux.HandleFunc("/comment/delete", enableCORS(authMiddleware(DeleteCommentHandler)))

    // Prefix /api routes
    mux.Handle("/api/", http.StripPrefix("/api", apiMux))

    log.Println("Server started at :8000")
    if err := http.ListenAndServe(":8000", mux); err != nil {
        log.Fatalf("Server stopped with error: %v", err)
    }
}
