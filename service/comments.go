package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	// "os"
	// "github.com/joho/godotenv"
	// STATIC "websocket-demo/VAR"
)

// var SCHEMA string

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}
// 	SCHEMA = os.Getenv("DB_SCHEMA")
// }

type Comment struct {
	ID          int       `json:"id"`
	IdeaID      int       `json:"idea_id"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	CommentText string    `json:"comment_text"`
}

// GetAllCommentsByIdeaID retrieves all comments for a specific idea
func GetAllCommentsByIdeaIDHandler(w http.ResponseWriter, r *http.Request) {
	// schema, err := GetSchemaByToken(w, r)

	type Comment struct {
		CommentID   string    `json:"comment_id"`
		IdeaID      string    `json:"idea_id"`
		CommentedBy string    `json:"commented_by"`
		CommentText string    `json:"comment_text"`
		CommentedAt time.Time `json:"commented_at"`
	}

	type Response struct {
		Status  int       `json:"status"`
		Message string    `json:"message"`
		Data    []Comment `json:"data,omitempty"`
	}

	ideaID := r.URL.Query().Get("id")
	if ideaID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Missing idea ID"})
		return
	}

	rows, err := db.Query(`
		SELECT commentid, ideaid, commentedby, commenttext, commentedat
		FROM kanaka.ideacomments
		WHERE ideaid = $1
		ORDER BY commentedat ASC
	`, ideaID)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Failed to retrieve comments"})
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.CommentID, &c.IdeaID, &c.CommentedBy, &c.CommentText, &c.CommentedAt); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Error scanning comment row"})
			return
		}
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Comments retrieved successfully",
		Data:    comments,
	})
}

// CreateCommentHandler creates a new comment and increments comment count
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	type CommentInput struct {
		IdeaID      string `json:"idea_id"`
		CommentedBy string `json:"commented_by"`
		CommentText string `json:"comment_text"`
	}

	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	var input CommentInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Invalid request body"})
		return
	}

	if input.IdeaID == "" || input.CommentedBy == "" || input.CommentText == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Missing required fields"})
		return
	}

	_, err := db.Exec(`
		INSERT INTO kanaka.ideacomments (commentid, ideaid, commentedby, commenttext)
		VALUES ($1, $2, $3, $4)
	`, uuid.New().String(), input.IdeaID, input.CommentedBy, input.CommentText)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Failed to create comment"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusCreated,
		Message: "Comment created successfully",
	})
}

// DeleteCommentHandler deletes a comment and decrements comment count
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}

	commentID := r.URL.Query().Get("id")
	if commentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Status: http.StatusBadRequest, Message: "Missing comment ID"})
		return
	}

	var ideaID string
	selectQuery := `SELECT ideaid FROM kanaka.ideacomments WHERE commentid = $1`
	err := db.QueryRow(selectQuery, commentID).Scan(&ideaID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(Response{Status: http.StatusNotFound, Message: "Comment not found"})
		return
	}

	deleteQuery := `DELETE FROM kanaka.ideacomments WHERE commentid = $1`
	_, err = db.Exec(deleteQuery, commentID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{Status: http.StatusInternalServerError, Message: "Failed to delete comment"})
		return
	}

	// Optional: If you have a separate comment count in Ideas table, update it here
	// Note: This logic assumes you store a count. If not, you can omit this block.
	/*
		updateQuery := `UPDATE kanaka.ideas SET comment_count = comment_count - 1 WHERE ideaid = $1`
		_, err = db.Exec(updateQuery, ideaID)
		if err != nil {
			log.Println("Failed to update comment count:", err)
		}
	*/
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusOK,
		Message: "Comment deleted successfully",
	})
}
