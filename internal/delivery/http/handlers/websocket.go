package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/collab-platform/backend/internal/delivery/websocket"
	"github.com/collab-platform/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	hub              *websocket.Hub
	authUsecase      *usecase.AuthUsecase
	collabUsecase    *usecase.CollaborationUsecase
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, check origin properly
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWebSocketHandler(
	hub *websocket.Hub,
	authUsecase *usecase.AuthUsecase,
	collabUsecase *usecase.CollaborationUsecase,
) *WebSocketHandler {
	return &WebSocketHandler{
		hub:           hub,
		authUsecase:   authUsecase,
		collabUsecase: collabUsecase,
	}
}

func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	// Get token from query parameter or header
	token := c.Query("token")
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required"})
		return
	}

	// Validate token
	claims, err := h.authUsecase.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	docIDStr := c.Query("document_id")
	if docIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document_id required"})
		return
	}

	docID, err := uuid.Parse(docIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	// Upgrade connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &websocket.Client{
		Hub:        h.hub,
		Conn:       conn,
		Send:       make(chan []byte, 256),
		UserID:     userID,
		DocumentID: docID,
	}

	client.Hub.register <- client

	// Start goroutines
	go client.WritePump()
	go client.ReadPump()
}

