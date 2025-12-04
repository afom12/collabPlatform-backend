package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/collab-platform/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 1024 // 512KB
)

type Client struct {
	Hub        *Hub
	Conn       *websocket.Conn
	Send       chan []byte
	UserID     uuid.UUID
	DocumentID uuid.UUID
}

type ClientMessage struct {
	Type      string          `json:"type"`
	Operation domain.Operation `json:"operation"`
	DocumentID string         `json:"document_id"`
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var clientMsg ClientMessage
		if err := json.Unmarshal(message, &clientMsg); err != nil {
			log.Printf("Error unmarshaling client message: %v", err)
			continue
		}

		docID, err := uuid.Parse(clientMsg.DocumentID)
		if err != nil {
			log.Printf("Invalid document ID: %v", err)
			continue
		}

		// Set timestamps and IDs
		clientMsg.Operation.ID = uuid.New()
		clientMsg.Operation.DocumentID = docID
		clientMsg.Operation.UserID = c.UserID
		clientMsg.Operation.Timestamp = time.Now().UnixNano()
		clientMsg.Operation.CreatedAt = time.Now()

		// Broadcast operation
		broadcastMsg := &BroadcastMessage{
			DocumentID: docID,
			Operation:  clientMsg.Operation,
			UserID:     c.UserID,
			Timestamp:  time.Now(),
		}

		c.Hub.BroadcastToDocument(docID, broadcastMsg)
	}
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

