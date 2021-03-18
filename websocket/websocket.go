package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Client represents a websock connection and pool
type Client struct {
	Conn *websocket.Conn
	Pool *Pool
}

// Pool represents a connection pool of connected clients
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan BroadcastMessage
}

// BroadcastMessage represents a message to be broadcast
type BroadcastMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// ClientMessage repressents a message from a client
type ClientMessage struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Client disconnected from websocket")
			return
		}
	}
}

// Upgrade upgrades client connection to websocket
func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// NewPool creates a pool
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan BroadcastMessage),
	}
}

// Start listens for anything passed in Pool's channels
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			break
		case message := <-pool.Broadcast:
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Println("Failed to broadcast message", err)
					return
				}
			}
		}
	}
}
