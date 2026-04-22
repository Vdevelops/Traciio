package hub

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10

	// Time allowed to read the next pong message from the peer
	// Reduced to 30s for better compatibility with Cloudflare Proxied (timeout ~100s)
	pongWait = 30

	// Send pings to peer with this period (must be less than pongWait)
	// More aggressive ping (25s) to keep connection alive through Cloudflare proxy
	pingPeriod = 25

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// getAllowedOrigins returns the list of allowed origins from environment variable
func getAllowedOrigins() []string {
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"http://localhost:8081", // Frontend dev server
		// Production origins
		"https://api.gilabs.id",
		"https://crm-demo.gilabs.id",
	}

	// Add origins from environment variable
	// Format: comma-separated list, e.g., "https://gilabs.id,https://www.gilabs.id"
	if envOrigins := os.Getenv("CORS_ALLOWED_ORIGINS"); envOrigins != "" {
		origins := strings.Split(envOrigins, ",")
		for _, origin := range origins {
			trimmed := strings.TrimSpace(origin)
			if trimmed != "" {
				// Check if already exists to avoid duplicates
				exists := false
				for _, existing := range allowedOrigins {
					if existing == trimmed {
						exists = true
						break
					}
				}
				if !exists {
					allowedOrigins = append(allowedOrigins, trimmed)
				}
			}
		}
	}

	return allowedOrigins
}

// GetUpgrader returns the WebSocket upgrader with origin validation
func GetUpgrader() websocket.Upgrader {
	allowedOrigins := getAllowedOrigins()
	isDevelopment := os.Getenv("ENV") != "production"
	
	// Log allowed origins for debugging (only in development)
	if isDevelopment {
		log.Printf("WebSocket allowed origins: %v", allowedOrigins)
	}

	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				// Allow requests without Origin header (same-origin or non-browser clients)
				if isDevelopment {
					log.Printf("WebSocket connection: no Origin header, allowing")
				}
				return true
			}

			// Check if origin is in allowed list
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					if isDevelopment {
						log.Printf("WebSocket connection allowed: origin %s", origin)
					}
					return true
				}
			}

			// In development, allow localhost origins even if not explicitly listed
			if isDevelopment {
				if strings.HasPrefix(origin, "http://localhost:") || strings.HasPrefix(origin, "https://localhost:") {
					log.Printf("WebSocket connection allowed (dev mode): origin %s", origin)
					return true
				}
			}

			log.Printf("WebSocket connection rejected: origin %s not in allowed list %v", origin, allowedOrigins)
			return false
		},
	}
}

// Client represents a WebSocket client connection
type Client struct {
	hub    *NotificationHub
	conn   *websocket.Conn
	userID string
	send   chan []byte
}

// NotificationHub manages WebSocket connections
type NotificationHub struct {
	// Registered clients per user
	clients map[string]map[*Client]bool

	// Broadcast channel
	broadcast chan *NotificationMessage

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NotificationMessage represents a notification message
type NotificationMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// NewNotificationHub creates a new notification hub
func NewNotificationHub() *NotificationHub {
	return &NotificationHub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan *NotificationMessage, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *NotificationHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]bool)
			}
			h.clients[client.userID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered for user: %s", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.userID)
					}
				}
			}
			h.mu.Unlock()
			log.Printf("Client unregistered for user: %s", client.userID)

		case message := <-h.broadcast:
			h.mu.RLock()
			// Extract userID from message data
			var userID string
			if dataMap, ok := message.Data.(map[string]interface{}); ok {
				if uid, ok := dataMap["user_id"].(string); ok {
					userID = uid
				}
			}
			
			if userID != "" {
				if userClients, ok := h.clients[userID]; ok {
					data, err := json.Marshal(message)
					if err != nil {
						log.Printf("Error marshaling message: %v", err)
						h.mu.RUnlock()
						continue
					}
					for client := range userClients {
						select {
						case client.send <- data:
						default:
							close(client.send)
							delete(userClients, client)
						}
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastNotification broadcasts a notification to a specific user
func (h *NotificationHub) BroadcastNotification(userID string, notification interface{}) {
	// Ensure notification has user_id
	var data map[string]interface{}
	if notifMap, ok := notification.(map[string]interface{}); ok {
		data = notifMap
		data["user_id"] = userID
	} else {
		// If notification is a struct, convert to map
		notifBytes, _ := json.Marshal(notification)
		json.Unmarshal(notifBytes, &data)
		data["user_id"] = userID
	}
	
	message := &NotificationMessage{
		Type: "notification.created",
		Data: data,
	}
	h.broadcast <- message
}

// BroadcastNotificationUpdate broadcasts a notification update
func (h *NotificationHub) BroadcastNotificationUpdate(userID string, notification interface{}) {
	var data map[string]interface{}
	if notifMap, ok := notification.(map[string]interface{}); ok {
		data = notifMap
		data["user_id"] = userID
	} else {
		notifBytes, _ := json.Marshal(notification)
		json.Unmarshal(notifBytes, &data)
		data["user_id"] = userID
	}
	
	message := &NotificationMessage{
		Type: "notification.updated",
		Data: data,
	}
	h.broadcast <- message
}

// BroadcastNotificationDelete broadcasts a notification deletion
func (h *NotificationHub) BroadcastNotificationDelete(userID string, notificationID string) {
	message := &NotificationMessage{
		Type: "notification.deleted",
		Data: map[string]interface{}{
			"user_id":         userID,
			"notification_id": notificationID,
		},
	}
	h.broadcast <- message
}

// BroadcastPermissionsUpdate broadcasts a permissions update event to a specific user
// This notifies the frontend to invalidate and refetch the user's permissions
func (h *NotificationHub) BroadcastPermissionsUpdate(userID string, roleID string) {
	message := &NotificationMessage{
		Type: "permissions.updated",
		Data: map[string]interface{}{
			"user_id": userID,
			"role_id": roleID,
		},
	}
	h.broadcast <- message
}

// BroadcastPermissionsUpdateToMultipleUsers broadcasts permissions update to multiple users
func (h *NotificationHub) BroadcastPermissionsUpdateToMultipleUsers(userIDs []string, roleID string) {
	for _, userID := range userIDs {
		h.BroadcastPermissionsUpdate(userID, roleID)
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWS handles WebSocket requests from clients
func (h *NotificationHub) ServeWS(conn *websocket.Conn, userID string) {
	client := &Client{
		hub:    h,
		conn:   conn,
		userID: userID,
		send:   make(chan []byte, 256),
	}

	client.hub.register <- client

	// Start goroutines for read and write pumps
	go client.writePump()
	go client.readPump()
}

