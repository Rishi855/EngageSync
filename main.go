package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	service "github.com/Rishi855/engagesync/service"
	// quiz "github.com/Rishi855/engagesync/quiz"
)

func init() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Println("Error loading .env file:", err)
    }

    log.Println("JWT_SECRET =", os.Getenv("JWT_SECRET"))

	host := os.Getenv("DB_HOST") // should be 'db' in container
	port := os.Getenv("DB_PORT") // 5432
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
	fmt.Println("Connecting to DB with:", dsn)
}

// Main function
func main() {
	// Create a new router
	r := mux.NewRouter()

	service.InsertInitialUsers()
	// Prefix `/api/` for all API routes
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Define routes under `/api/` and apply CORS middleware
	// apiRouter.HandleFunc("/demo/request", service.EnableCORS(service.DemoRequest)).Methods("POST")

	apiRouter.HandleFunc("/", service.EnableCORS(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Hello, EngageSync API is running!"}`))
	})).Methods("GET")

	apiRouter.HandleFunc("/login", service.EnableCORS(service.LoginHandler)).Methods("GET")

	apiRouter.HandleFunc("/ideas", service.EnableCORS(service.AuthMiddleware(service.GetAllIdeasHandler))).Methods("GET")
	apiRouter.HandleFunc("/idea", service.EnableCORS(service.AuthMiddleware(service.CreateIdeaHandler))).Methods("POST")
	apiRouter.HandleFunc("/idea/delete", service.EnableCORS(service.AuthMiddleware(service.DeleteIdeaHandler))).Methods("DELETE")

	apiRouter.HandleFunc("/comments", service.EnableCORS(service.AuthMiddleware(service.GetAllCommentsByIdeaIDHandler))).Methods("GET")
	apiRouter.HandleFunc("/comment", service.EnableCORS(service.AuthMiddleware(service.CreateCommentHandler))).Methods("POST")
	apiRouter.HandleFunc("/comment/delete", service.EnableCORS(service.AuthMiddleware(service.DeleteCommentHandler))).Methods("DELETE")

	apiRouter.HandleFunc("/projects", service.EnableCORS(service.AuthMiddleware(service.GetAllProjectsHandler))).Methods("GET")
	apiRouter.HandleFunc("/project", service.EnableCORS(service.AuthMiddleware(service.CreateProjectHandler))).Methods("POST")
	apiRouter.HandleFunc("/project/delete", service.EnableCORS(service.AuthMiddleware(service.DeleteProjectHandler))).Methods("DELETE")

	apiRouter.HandleFunc("/project/members", service.EnableCORS(service.AuthMiddleware(service.GetAllProjectMembersHandler))).Methods("GET")
	apiRouter.HandleFunc("/project/member", service.EnableCORS(service.AuthMiddleware(service.CreateProjectMemberHandler))).Methods("POST")
	apiRouter.HandleFunc("/project/member/delete", service.EnableCORS(service.AuthMiddleware(service.DeleteProjectMemberHandler))).Methods("DELETE")

	// apiRouter.HandleFunc("/technologies", service.EnableCORS(service.AuthMiddleware(service.GetAllTechnologiesHandler))).Methods("GET")
	// apiRouter.HandleFunc("/roles", service.EnableCORS(service.AuthMiddleware(service.GetAllRolesHandler))).Methods("GET")
	apiRouter.HandleFunc("/all/users", service.EnableCORS(service.AuthMiddleware(service.GetAllUsersHandler))).Methods("GET")

	apiRouter.HandleFunc("/add/user", service.EnableCORS(service.AuthMiddleware(service.AddUserHandler))).Methods("POST")

	apiRouter.HandleFunc("/add/organization", service.EnableCORS(service.AuthMiddleware(service.AddOrganizationHandler))).Methods("POST")
	apiRouter.HandleFunc("/organizations", service.EnableCORS(service.AuthMiddleware(service.GetOrganizationHandler))).Methods("GET")
	
	apiRouter.HandleFunc("/get/guess/image/{category}",service.EnableCORS(service.AuthMiddleware(service.GetGuessImageHandler))).Methods("GET")
	// go quiz.HubInstance.Run()

	// r.HandleFunc("/ws/quiz", quiz.QuizWebSocketHandler) // WebSocket entry
	// apiRouter.HandleFunc("/send-question", service.EnableCORS(service.AuthMiddleware(quiz.SendQuestionHandler))).Methods("POST")
	// apiRouter.HandleFunc("/start-quiz", service.EnableCORS(service.AuthMiddleware(quiz.StartQuizHandler))).Methods("POST")

	// Start the server
	// Add this just before http.ListenAndServe
	log.Printf("Available routes:")
	err := r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, _ := route.GetMethods()
			log.Printf("Route: %s [%s]", pathTemplate, methods)
		}
		return nil
	})
	if err != nil {
		log.Printf("Error walking routes: %v", err)
	}
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
