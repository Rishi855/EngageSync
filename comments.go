package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	STATIC "websocket-demo/VAR"
)

type Comment struct {
	ID          int       `json:"id"`
	IdeaID      int       `json:"idea_id"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
	CommentText string    `json:"comment_text"`
}

// GetAllCommentsByIdeaID retrieves all comments for a specific idea
func GetAllCommentsByIdeaIDHandler(w http.ResponseWriter, r *http.Request) {
	ideaID := r.URL.Query().Get("idea_id")
	if ideaID == "" {
		http.Error(w, "Missing idea_id", http.StatusBadRequest)
		return
	}

	query := `SELECT id, idea_id, user_id, comment_text, created_at 
			  FROM ` + STATIC.TANENT + `.comments 
			  WHERE idea_id = $1 ORDER BY created_at ASC`

	rows, err := queryRows(query, ideaID)
	if err != nil {
		http.Error(w, "Error fetching comments", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var c Comment
		if err := rows.Scan(&c.ID, &c.IdeaID, &c.UserID, &c.CommentText, &c.CreatedAt); err != nil {
			http.Error(w, "Error scanning comment", http.StatusInternalServerError)
			return
		}
		comments = append(comments, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}

// CreateCommentHandler creates a new comment and increments comment count
func CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var comment Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		log.Println("Invalid request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userID := r.Header.Get("UserID") // Retrieved from JWT or auth middleware
	if userID == "" {
		log.Println("Missing UserID in header")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("STATIC.TANENT: %s", STATIC.TANENT)
	insertQuery := `INSERT INTO ` + STATIC.TANENT + `.comments 
                    (idea_id, user_id, comment_text, created_at) 
                    VALUES ($1, $2, $3, $4)`
	log.Printf("Executing query: %s", insertQuery)
	_, err := execQuery(insertQuery, comment.IdeaID, userID, comment.CommentText, time.Now())
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	updateQuery := `UPDATE ` + STATIC.TANENT + `.ideas 
                    SET comments = comments + 1 
                    WHERE id = $1`
	log.Printf("Executing update query: %s", updateQuery)
	_, err = execQuery(updateQuery, comment.IdeaID)
	if err != nil {
		log.Println("Failed to increment comment count")
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Comment created")
}

// DeleteCommentHandler deletes a comment and decrements comment count
func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentID := r.URL.Query().Get("id")
	if commentID == "" {
		http.Error(w, "Missing comment id", http.StatusBadRequest)
		return	
	}

	var ideaID int
	selectQuery := `SELECT idea_id FROM ` + STATIC.TANENT + `.comments WHERE id = $1`
	err := db.QueryRow(selectQuery, commentID).Scan(&ideaID)
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	deleteQuery := `DELETE FROM ` + STATIC.TANENT + `.comments WHERE id = $1`
	_, err = execQuery(deleteQuery, commentID)
	if err != nil {
		http.Error(w, "Error deleting comment", http.StatusInternalServerError)
		return
	}

	updateQuery := `UPDATE ` + STATIC.TANENT + `.ideas 
					SET comments = comments - 1 
					WHERE id = $1`
	_, err = execQuery(updateQuery, ideaID)
	if err != nil {
		log.Println("Failed to decrement comment count")
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Comment deleted")
}
