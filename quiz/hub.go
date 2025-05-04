package quiz

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn     *websocket.Conn
	userID   string
	groupID  string
	sendChan chan []byte
}

func (c *Client) readPump() {
	panic("unimplemented")
}

type Hub struct {
	clients    map[string][]*Client // groupID -> clients
	register   chan *Client
	unregister chan *Client
	broadcast  chan GroupMessage
	mutex      sync.Mutex
}

type GroupMessage struct {
	GroupID string
	Message []byte
}

var HubInstance = Hub{
	clients:    make(map[string][]*Client),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	broadcast:  make(chan GroupMessage),
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client.groupID] = append(h.clients[client.groupID], client)
			h.mutex.Unlock()

		case client := <-h.unregister:
			h.mutex.Lock()
			clients := h.clients[client.groupID]
			for i, c := range clients {
				if c == client {
					h.clients[client.groupID] = append(clients[:i], clients[i+1:]...)
					break
				}
			}
			h.mutex.Unlock()

		case message := <-h.broadcast:
			h.mutex.Lock()
			for _, client := range h.clients[message.GroupID] {
				select {
				case client.sendChan <- message.Message:
				default:
					log.Println("Dropping message to slow client")
				}
			}
			h.mutex.Unlock()
		}
	}
}
