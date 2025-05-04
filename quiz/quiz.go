package quiz

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Question struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type Message struct {
	UserID   string `json:"user_id"`
	Answer   string `json:"answer"`
	Question int    `json:"question_id"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var questions = []Question{
	{1, "What is the capital of India?"},
	{2, "2 + 2 = ?"},
	{3, "Which planet is known as the Red Planet?"},
	{4, "Go is developed by which company?"},
	{5, "What is 5 x 6?"},
	{6, "Which is the largest ocean?"},
	{7, "Sun rises in the ___?"},
	{8, "Who wrote the Ramayana?"},
	{9, "What color is the sky?"},
	{10, "Is Go statically typed?"},
}

var clients = make(map[*websocket.Conn]bool)

func quizHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer ws.Close()

	clients[ws] = true
	log.Println("New client connected")

	// Receive responses
	go func() {
		for {
			var msg Message
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Println("Read error:", err)
				delete(clients, ws)
				break
			}
			log.Printf("Answer from User %s for Q%d: %s\n", msg.UserID, msg.Question, msg.Answer)
		}
	}()

	// Send 10 questions, one every 10 seconds
	for i, q := range questions {
		broadcastQuestion(q)
		log.Printf("Sent question %d", i+1)
		time.Sleep(10 * time.Second)
	}
}

func broadcastQuestion(q Question) {
	for client := range clients {
		err := client.WriteJSON(q)
		if err != nil {
			log.Println("Write error:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

func SendQuestionHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		GroupID  string `json:"group_id"`
		Question string `json:"question"`
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	msg := map[string]string{
		"type":     "question",
		"question": req.Question,
	}
	bytes, _ := json.Marshal(msg)

	HubInstance.broadcast <- GroupMessage{
		GroupID: req.GroupID,
		Message: bytes,
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Question sent"))
}
func StartQuizForGroup(groupID string) {
	go func() {
		for i, q := range questions {
			msg := map[string]interface{}{
				"type":     "question",
				"question": q.Content,
				"id":       q.ID,
			}
			data, err := json.Marshal(msg)
			if err != nil {
				log.Println("Error marshalling question:", err)
				continue
			}
			HubInstance.broadcast <- GroupMessage{
				GroupID: groupID,
				Message: data,
			}
			log.Printf("Sent question %d to group %s", i+1, groupID)
			time.Sleep(10 * time.Second)
		}
	}()
}
