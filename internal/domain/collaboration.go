package domain

import (
	"time"

	"github.com/google/uuid"
)

// Operation represents a single edit operation (CRDT-based)
type Operation struct {
	ID         uuid.UUID `json:"id"`
	DocumentID uuid.UUID `json:"document_id"`
	UserID     uuid.UUID `json:"user_id"`
	Type       string    `json:"type"` // "insert", "delete", "format"
	Position   int       `json:"position"`
	Length     int       `json:"length"`
	Content    string    `json:"content"`
	Timestamp  int64     `json:"timestamp"` // Lamport timestamp
	VectorClock map[uuid.UUID]int64 `json:"vector_clock"` // Vector clock for ordering
	CreatedAt  time.Time `json:"created_at"`
}

// ClientSession represents an active WebSocket connection
type ClientSession struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	DocumentID uuid.UUID
	Conn       interface{} // *websocket.Conn (avoiding import here)
	Send       chan []byte
	LastSeen   time.Time
}

// BroadcastMessage represents a message to be broadcasted
type BroadcastMessage struct {
	DocumentID uuid.UUID `json:"document_id"`
	Operation  Operation `json:"operation"`
	UserID     uuid.UUID `json:"user_id"`
	Timestamp  time.Time `json:"timestamp"`
}

