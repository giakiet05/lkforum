package ws

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/platform/bus"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	userClients map[string]*Client
	register    chan *Client
	unregister  chan *Client
	incoming    chan []byte
	eventBus    *bus.EventBus
}

func NewHub(bus *bus.EventBus) *Hub {
	return &Hub{
		incoming:    make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		userClients: make(map[string]*Client),
		eventBus:    bus,
	}
}

// Start runs the hub's event loop and subscribes to the event eventBus.
func (h *Hub) Start() {
	eventChannel := make(bus.EventListener, 100)
	h.eventBus.Subscribe(bus.TopicNotificationCreated, eventChannel)
	h.eventBus.Subscribe(bus.TopicMessageCreated, eventChannel)
	h.eventBus.Subscribe(bus.TopicMessageError, eventChannel)

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
				delete(h.userClients, client.UserID)
				close(client.send)
				log.Printf("WebSocket client unregistered: %s", client.UserID)
			}

		case data := <-h.incoming:
			//Handle message receive from client
			parts := bytes.SplitN(data, []byte("|"), 2)
			if len(parts) != 2 {
				log.Println("Invalid incoming message format")
				continue
			}

			userID := string(parts[0])
			message := parts[1]
			h.handleIncoming(message, userID)

		case event := <-eventChannel:
			//Handle event
			switch event.Topic() {
			case bus.TopicNotificationCreated:
				payload := event.Payload()
				if recipientID, ok := payload["recipientId"].(string); ok {
					if notification, ok := payload["notification"].(interface{}); ok {
						h.sendToUser(recipientID, NewNotification, notification)
					}
				}

			case bus.TopicMessageCreated:
				payload := event.Payload()
				recipientIDs, _ := payload["recipient_ids"].([]string)
				tempMessageID, _ := payload["temp_message_id"].(string)
				messageData, _ := payload["message"].(dto.MessageResponse)

				ackResponse := ACKMessagePayload{
					TempMessageID: tempMessageID,
					Message:       messageData,
				}
				h.sendToUser(messageData.SenderID, ACKMessage, ackResponse)

				response := SendMessagePayload{
					Message: messageData,
				}
				h.broadcastToUsers(recipientIDs, SendMessage, response)

			case bus.TopicMessageError:
				payload := event.Payload()
				senderID, _ := payload["sender_id"].(string)
				tempID, _ := payload["temp_id"].(string)
				errorCode, _ := payload["error_code"].(string)
				errorMsg, _ := payload["error_msg"].(string)

				errPayload := ErrorPayload{
					TempMessageID: &tempID,
					ErrorCode:     &errorCode,
					ErrorMsg:      errorMsg,
				}
				h.sendToUser(senderID, ErrorMessage, errPayload)

			default:
				log.Printf("WebSocket client received unknown event: %s", event.Topic())
			}
		}
	}
}

func (h *Hub) handleIncoming(raw []byte, userID string) {
	var incomingMsg WebSocketMessage
	if err := json.Unmarshal(raw, &incomingMsg); err != nil {
		log.Printf("WebSocket Invalid JSON from client: %v", err)
		return
	}

	switch incomingMsg.Type {
	case NewMessage:
		var payload NewMessagePayload
		if err := decodePayload(incomingMsg.Payload, &payload); err != nil {
			log.Printf("WebSocket invalid new message payload from user %s: %v", userID, err)
			h.sendToUser(userID, ErrorMessage, ErrorPayload{ErrorMsg: err.Error(), TempMessageID: nil})
			return
		}

		h.eventBus.Publish(bus.NewMessageEvent{
			TempMessageID: payload.TempMessageID,
			ChannelID:     payload.ChannelID,
			SenderID:      userID,
			Type:          payload.Type,
			Content:       payload.Content,
		})

	default:
		log.Printf("Unknown incoming type: %s", incomingMsg.Type)
	}
}

func decodePayload[T any](data interface{}, out *T) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(marshal, out)
}

// sendToUser is a private method to send a message to a specific user.
func (h *Hub) sendToUser(userID string, messageType SocketMessageType, payload interface{}) {
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
			log.Printf("Warning: Client %s channel is full. Message dropped.", userID)
		}
	}
}

func (h *Hub) broadcastToUsers(userIDs []string, messageType SocketMessageType, payload interface{}) {
	msg := WebSocketMessage{
		Type:    messageType,
		Payload: payload,
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshalling websocket broadcast message: %v", err)
		return
	}

	for _, userID := range userIDs {
		if client, ok := h.userClients[userID]; ok {
			select {
			case client.send <- jsonMsg:
			default:
				log.Printf("Warning: Client %s channel is full. Broadcast message dropped.", userID)
			}
		}
	}
}
