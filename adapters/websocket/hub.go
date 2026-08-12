package websocket

import (
	"ITLAFINAL/pkg/authjwt"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

// Hub keeps authenticated customer connections. A customer may have more than
// one browser tab, so every connection receives status events.
type Hub struct {
	clients  map[string]map[*client]struct{}
	mu       sync.RWMutex
	upgrader websocket.Upgrader
}

type client struct {
	conn    *websocket.Conn
	writeMu sync.Mutex // gorilla permits only one concurrent writer per conn
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]map[*client]struct{}),
		upgrader: websocket.Upgrader{
			CheckOrigin: allowedOrigin,
		},
	}
}

// HandleConnection authenticates the upgrade request. Browser clients should
// use /ws?access_token=<JWT>; non-browser clients can use Authorization: Bearer.
// customer_id is intentionally ignored: the recipient comes only from the JWT.
func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request) {
	claims, err := authenticateRequest(r)
	if err != nil {
		http.Error(w, "no autorizado", http.StatusUnauthorized)
		return
	}
	if claims.Role != "customer" {
		http.Error(w, "solo clientes pueden suscribirse a notificaciones", http.StatusForbidden)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade: %v", err)
		return
	}
	c := &client{conn: conn}
	h.add(claims.UserID, c)
	log.Printf("websocket: customer %s conectado", claims.UserID)

	defer func() {
		h.remove(claims.UserID, c)
		_ = conn.Close()
		log.Printf("websocket: customer %s desconectado", claims.UserID)
	}()

	conn.SetReadLimit(maxMessageSize)
	_ = conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	stopPing := make(chan struct{})
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		defer close(done)
		for {
			select {
			case <-stopPing:
				return
			case <-ticker.C:
				if err := c.write(websocket.PingMessage, nil); err != nil {
					return
				}
			}
		}
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			close(stopPing)
			<-done
			return
		}
	}
}

func authenticateRequest(r *http.Request) (*authjwt.Claims, error) {
	if header := r.Header.Get("Authorization"); header != "" {
		parts := strings.Fields(header)
		if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
			return authjwt.Parse(parts[1])
		}
	}
	return authjwt.Parse(r.URL.Query().Get("access_token"))
}

func allowedOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" { // non-browser WebSocket clients do not send Origin.
		return true
	}
	originURL, err := url.Parse(origin)
	if err != nil || originURL.Scheme == "" || originURL.Host == "" {
		return false
	}
	for _, allowed := range allowedOrigins() {
		if origin == allowed {
			return true
		}
	}
	return false
}

func allowedOrigins() []string {
	configured := strings.TrimSpace(os.Getenv("WS_ALLOWED_ORIGINS"))
	if configured == "" {
		return []string{
			"http://localhost:3000", "http://127.0.0.1:3000",
			"http://localhost:5173", "http://127.0.0.1:5173",
		}
	}
	values := strings.Split(configured, ",")
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			result = append(result, value)
		}
	}
	return result
}

func (h *Hub) add(customerID string, c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[customerID] == nil {
		h.clients[customerID] = make(map[*client]struct{})
	}
	h.clients[customerID][c] = struct{}{}
}

func (h *Hub) remove(customerID string, c *client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	connections := h.clients[customerID]
	delete(connections, c)
	if len(connections) == 0 {
		delete(h.clients, customerID)
	}
}

func (c *client) write(messageType int, data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(messageType, data)
}

func (h *Hub) NotifyOrderReady(customerID, orderID string) error {
	return h.sendMessage(customerID, Message{Type: "ORDER_READY", OrderID: orderID, Message: "Tu orden está lista para recoger."})
}

func (h *Hub) NotifyStatusChange(customerID, orderID, status string) error {
	return h.sendMessage(customerID, Message{Type: "STATUS_CHANGE", OrderID: orderID, Message: "Tu orden cambió de estado: " + status, Status: status})
}

func (h *Hub) sendMessage(customerID string, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	h.mu.RLock()
	connections := make([]*client, 0, len(h.clients[customerID]))
	for c := range h.clients[customerID] {
		connections = append(connections, c)
	}
	h.mu.RUnlock()

	for _, c := range connections {
		if err := c.write(websocket.TextMessage, data); err != nil {
			log.Printf("websocket send to customer %s: %v", customerID, err)
			h.remove(customerID, c)
			_ = c.conn.Close()
		}
	}
	return nil // notifications are best effort; no active connection is valid.
}

type Message struct {
	Type    string `json:"type"`
	OrderID string `json:"order_id"`
	Message string `json:"message"`
	Status  string `json:"status,omitempty"`
}
