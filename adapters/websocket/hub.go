package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub mantiene todas las conexiones activas
// Cada cliente se conecta con su customerID
type Hub struct {
	clients  map[string]*websocket.Conn // customerID → conexión
	mu       sync.RWMutex
	upgrader websocket.Upgrader
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*websocket.Conn),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // en producción validar el origin
			},
		},
	}
}

// HandleConnection registra al cliente cuando se conecta
// GET /ws?customer_id=c1
func (h *Hub) HandleConnection(w http.ResponseWriter, r *http.Request) {
	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		http.Error(w, "customer_id requerido", http.StatusBadRequest)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	h.mu.Lock()
	h.clients[customerID] = conn
	h.mu.Unlock()

	log.Printf("🔌 Cliente %s conectado via WebSocket", customerID)

	// Mantener la conexión viva y limpiar cuando cierre
	defer func() {
		h.mu.Lock()
		delete(h.clients, customerID)
		h.mu.Unlock()
		conn.Close()
		log.Printf("🔌 Cliente %s desconectado", customerID)
	}()

	// Escuchar mensajes (mantiene la goroutine activa)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

// NotifyOrderReady implementa el port Notifier
func (h *Hub) NotifyOrderReady(customerID, orderID string) error {
	return h.sendMessage(customerID, Message{
		Type:    "ORDER_READY",
		OrderID: orderID,
		Message: "¡Tu orden está lista para recoger!",
	})
}

// NotifyStatusChange implementa el port Notifier
func (h *Hub) NotifyStatusChange(customerID, orderID, status string) error {
	return h.sendMessage(customerID, Message{
		Type:    "STATUS_CHANGE",
		OrderID: orderID,
		Message: "Tu orden cambió de estado: " + status,
		Status:  status,
	})
}

func (h *Hub) sendMessage(customerID string, msg Message) error {
	h.mu.RLock()
	conn, ok := h.clients[customerID]
	h.mu.RUnlock()

	if !ok {
		// Cliente no conectado — guardar para después (opcional)
		log.Printf("⚠️  Cliente %s no está conectado", customerID)
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, data)
}

type Message struct {
	Type    string `json:"type"`
	OrderID string `json:"order_id"`
	Message string `json:"message"`
	Status  string `json:"status,omitempty"`
}
