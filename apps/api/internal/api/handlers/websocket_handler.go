package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/hub"
	"github.com/gilabs/crm-healthcare/api/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type WebSocketHandler struct {
	hub       *hub.NotificationHub
	jwtManager *jwt.JWTManager
}

func NewWebSocketHandler(hub *hub.NotificationHub, jwtManager *jwt.JWTManager) *WebSocketHandler {
	return &WebSocketHandler{
		hub:        hub,
		jwtManager: jwtManager,
	}
}

// HandleWebSocket handles WebSocket connections for notifications
func (h *WebSocketHandler) HandleWebSocket(c *gin.Context) {
	var tokenString string

	// Priority 1: Get token from cookie (most secure for production)
	tokenString, _ = c.Cookie("token")

	// Priority 2: Fallback to query parameter (for backward compatibility)
	if tokenString == "" {
		tokenString = c.Query("token")
	}

	// Priority 3: Fallback to Authorization header
	if tokenString == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}
	}

	if tokenString == "" {
		log.Printf("WebSocket connection rejected: no token provided (Origin: %s)", c.GetHeader("Origin"))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token required"})
		return
	}

	// Validate token
	claims, err := h.jwtManager.ValidateToken(tokenString)
	if err != nil {
		log.Printf("WebSocket connection rejected: invalid token - %v (Origin: %s)", err, c.GetHeader("Origin"))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Upgrade connection to WebSocket
	upgrader := hub.GetUpgrader()
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v (Origin: %s, UserID: %s)", err, c.GetHeader("Origin"), claims.UserID)
		return
	}

	log.Printf("WebSocket connection established for user: %s (Origin: %s)", claims.UserID, c.GetHeader("Origin"))
	
	// Register client with hub
	h.hub.ServeWS(conn, claims.UserID)
}

