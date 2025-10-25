package ws

import (
	"encoding/json"
	"log"

	"github.com/giakiet05/lkforum/internal/platform/bus"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	userClients map[string]*Client
	register    chan *Client
	unregister  chan *Client
}

// Global Hub instance
var WSHub = NewHub()

func NewHub() *Hub {
	return &Hub{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		userClients: make(map[string]*Client),
	}
}

// Start runs the hub's event loop and subscribes to the event bus.
func (h *Hub) Start() {
	eventChannel := make(bus.EventListener, 100)
	bus.Bus.Subscribe(bus.TopicNotificationCreated, eventChannel)

	log.Println("WebSocket Hub started and subscribed to events.")

	go h.run(eventChannel)
}

// RegisterClient sends a client to the register channel.
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *Hub) run(eventChannel bus.EventListener) {
	for {
		select {
		case client := <-h.register:
			h.userClients[client.UserID] = client
			log.Printf("WebSocket client registered: %s", client.UserID)

		case client := <-h.unregister:
			if _, ok := h.userClients[client.UserID]; ok {
				// Important: only delete from the map. The closing of the send channel
				// and the connection is handled by the client's pumps.
				delete(h.userClients, client.UserID)
				log.Printf("WebSocket client unregistered: %s", client.UserID)
			}

		case event := <-eventChannel:
			if event.Topic() == bus.TopicNotificationCreated {
				payload := event.Payload()
				if recipientID, ok := payload["recipientId"].(string); ok {
					if notification, ok := payload["notification"].(interface{}); ok {
						h.sendToUser(recipientID, "new_notification", notification)
					}
				}
			}
		}
	}
}

type WebSocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

// sendToUser is a private method to send a message to a specific user.
func (h *Hub) sendToUser(userID string, messageType string, payload interface{}) {
	if client, ok := h.userClients[userID]; ok {
		msg := WebSocketMessage{
			Type:    messageType,
			Payload: payload,
		}
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshalling websocket message: %v", err)
			return
	}

		select {
		case client.send <- jsonMsg:
		default:
			// If the client's send buffer is full, we assume the client is lagging
			// and drop the message. The client's own read/write pumps are responsible
			// for detecting a dead connection and unregistering.
			log.Printf("Warning: Client %s channel is full. Message dropped.", userID)
		}
	}
}
