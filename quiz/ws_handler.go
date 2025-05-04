package quiz

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/rand"
)

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

func (c *Client) writePump() {
	for {
		select {
		case message := <-c.sendChan:
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("writePump error:", err)
				return
			}
		}
	}
}

func createLobbyHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		AdminID string `json:"admin_id"`
	}
	type Response struct {
		GroupID string `json:"group_id"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.AdminID == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	groupID := generateGroupID() // Generate unique 6-digit or alphanumeric code

	resp := Response{GroupID: groupID}
	json.NewEncoder(w).Encode(resp)
}

func generateGroupID() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func QuizWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	groupID := r.URL.Query().Get("group_id")

	// Check if group exists (optional: track lobby map or validate from DB/cache)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade:", err)
		return
	}

	client := &Client{
		conn:     conn,
		userID:   userID,
		groupID:  groupID,
		sendChan: make(chan []byte, 256),
	}

	HubInstance.register <- client

	go client.readPump()
	go client.writePump()
}

func StartQuizHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		GroupID string `json:"group_id"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.GroupID == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	go StartQuizForGroup(req.GroupID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Quiz started for group " + req.GroupID))
}
