package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/collab-platform/backend/internal/domain"
	"github.com/collab-platform/backend/internal/infrastructure/redis"
	"github.com/google/uuid"
)

type Hub struct {
	// Registered clients per document
	documents map[uuid.UUID]map[*Client]bool

	// Inbound messages from clients
	broadcast chan *BroadcastMessage

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Redis client for distributed pub/sub
	redisClient *redis.RedisClient

	// Mutex for thread-safe access
	mu sync.RWMutex
}

type BroadcastMessage struct {
	DocumentID uuid.UUID      `json:"document_id"`
	Operation  domain.Operation `json:"operation"`
	UserID     uuid.UUID      `json:"user_id"`
	Timestamp  time.Time      `json:"timestamp"`
}

func NewHub(redisClient *redis.RedisClient) *Hub {
	hub := &Hub{
		documents:   make(map[uuid.UUID]map[*Client]bool),
		broadcast:   make(chan *BroadcastMessage, 256),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		redisClient: redisClient,
	}
	return hub
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.documents[client.DocumentID] == nil {
				h.documents[client.DocumentID] = make(map[*Client]bool)
			}
			h.documents[client.DocumentID][client] = true
			h.mu.Unlock()

			log.Printf("Client registered: User %s for Document %s", client.UserID, client.DocumentID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.documents[client.DocumentID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.documents, client.DocumentID)
					}
				}
			}
			h.mu.Unlock()

			log.Printf("Client unregistered: User %s for Document %s", client.UserID, client.DocumentID)

		case message := <-h.broadcast:
			// Publish to Redis for distributed broadcasting
			if h.redisClient != nil {
				if err := h.redisClient.PublishOperation(message.DocumentID.String(), *message); err != nil {
					log.Printf("Error publishing to Redis: %v", err)
				}
			}

			// Broadcast to local clients
			h.mu.RLock()
			clients := h.documents[message.DocumentID]
			h.mu.RUnlock()

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Error marshaling broadcast message: %v", err)
				continue
			}

			for client := range clients {
				// Don't send message back to sender
				if client.UserID != message.UserID {
					select {
					case client.send <- data:
					default:
						close(client.send)
						h.mu.Lock()
						delete(clients, client)
						h.mu.Unlock()
					}
				}
			}
		}
	}
}

func (h *Hub) BroadcastToDocument(docID uuid.UUID, message *BroadcastMessage) {
	h.broadcast <- message
}

func (h *Hub) GetDocumentClients(docID uuid.UUID) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := make([]*Client, 0)
	if docClients, ok := h.documents[docID]; ok {
		for client := range docClients {
			clients = append(clients, client)
		}
	}
	return clients
}

// HandleRedisMessage handles messages received from Redis pub/sub
func (h *Hub) HandleRedisMessage(docID uuid.UUID, message *BroadcastMessage) {
	h.mu.RLock()
	clients := h.documents[docID]
	h.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling Redis message: %v", err)
		return
	}

	for client := range clients {
		select {
		case client.send <- data:
		default:
			close(client.send)
			h.mu.Lock()
			delete(clients, client)
			h.mu.Unlock()
		}
	}
}

