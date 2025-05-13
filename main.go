package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	service "github.com/Rishi855/engagesync/service"
	// quiz "github.com/Rishi855/engagesync/quiz"
)

// var SCHEMA string

// func init() {
// 	service.InitConfig()
// }

// Main function
func main() {
	// Create a new router
	r := mux.NewRouter()

	service.InsertInitialUsers()
	// Prefix `/api/` for all API routes
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Define routes under `/api/` and apply CORS middleware
	// apiRouter.HandleFunc("/demo/request", service.EnableCORS(service.DemoRequest)).Methods("POST")
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
	log.Println("Server started at :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
