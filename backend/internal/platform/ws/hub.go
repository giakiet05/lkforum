package ws

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/platform/bus"
	"github.com/giakiet05/lkforum/internal/util"
)

// Hub maintains the set of active clients and broadcasts messages to them.
type Hub struct {
	userClients map[string]*Client
	register    chan *Client
	unregister  chan *Client
	incoming    chan []byte
	eventBus    bus.EventBus
}

func NewHub(bus bus.EventBus) *Hub {
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
	h.eventBus.Subscribe(bus.TopicBroadcast, eventChannel)
	h.eventBus.Subscribe(bus.TopicMessageError, eventChannel)

	log.Println("WebSocket Hub started and subscribed to events.")

	go h.run(eventChannel)
}

// RegisterClient sends a client to the register channel.
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// GetOnlineUsers returns a list of currently connected user IDs
func (h *Hub) GetOnlineUsers() []string {
	userIDs := make([]string, 0, len(h.userClients))
	for userID := range h.userClients {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

func (h *Hub) run(eventChannel bus.EventListener) {
	for {
		select {
		case client := <-h.register:
			h.userClients[client.UserID] = client
			log.Printf("WebSocket client registered: %s", client.UserID)

			// Broadcast presence online to all connected users
			onlineUsers := h.GetOnlineUsers()
			presencePayload := dto.PresencePayload{UserID: client.UserID}
			h.broadcastToUsers(onlineUsers, dto.PresenceOnline, presencePayload)
		case client := <-h.unregister:
			if _, ok := h.userClients[client.UserID]; ok {
				delete(h.userClients, client.UserID)
				close(client.send)
				log.Printf("WebSocket client unregistered: %s", client.UserID)

				// Broadcast presence offline to remaining connected users
				onlineUsers := h.GetOnlineUsers()
				presencePayload := dto.PresencePayload{UserID: client.UserID}
				h.broadcastToUsers(onlineUsers, dto.PresenceOffline, presencePayload)
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
						h.sendToUser(recipientID, dto.NewNotification, notification)
					}
				}
			case bus.TopicBroadcast:
				payload := event.Payload()
				recipientIDs, _ := payload["recipient_ids"].([]string)
				eventType, _ := payload["event_type"].(bus.BroadcastEventType)
				tempID, _ := payload["temp_id"].(string)
				data := payload["data"]

				h.handleBroadcast(recipientIDs, string(eventType), tempID, data)
			case bus.TopicMessageError:
				payload := event.Payload()
				senderID, _ := payload["sender_id"].(string)
				tempID, _ := payload["temp_id"].(string)
				errorCode, _ := payload["error_code"].(string)
				errorMsg, _ := payload["error_msg"].(string)

				errPayload := dto.ErrorPayload{
					TempMessageID: &tempID,
					ErrorCode:     &errorCode,
					ErrorMsg:      errorMsg,
				}
				h.sendToUser(senderID, dto.ErrorMessage, errPayload)
			default:
				log.Printf("WebSocket client received unknown event: %s", event.Topic())
			}
		}
	}
}

func (h *Hub) handleIncoming(raw []byte, userID string) {
	var incomingMsg dto.WebSocketMessage
	if err := json.Unmarshal(raw, &incomingMsg); err != nil {
		log.Printf("WebSocket Invalid JSON from client: %v", err)
		return
	}

	switch incomingMsg.Type {
	case dto.NewMessage:
		var payload dto.NewMessagePayload
		if err := util.DecodeJson(incomingMsg.Payload, &payload); err != nil {
			log.Printf("WebSocket invalid new message payload from user %s: %v", userID, err)
			h.sendToUser(userID, dto.ErrorMessage, dto.ErrorPayload{ErrorMsg: err.Error(), TempMessageID: nil})
			return
		}

		h.eventBus.Publish(bus.NewMessageEvent{
			TempMessageID:  payload.TempMessageID,
			ChannelID:      payload.ChannelID,
			SenderID:       userID,
			SenderUsername: payload.SenderUsername,
			Type:           payload.Type,
			Content:        payload.Content,
		})
	case dto.TypingIndicator:
		var payload dto.TypingIndicatorPayload
		if err := util.DecodeJson(incomingMsg.Payload, &payload); err != nil {
			log.Printf("WebSocket invalid typing indicator payload from user %s: %v", userID, err)
			h.sendToUser(userID, dto.ErrorMessage, dto.ErrorPayload{ErrorMsg: err.Error()})
			return
		}

		h.eventBus.Publish(bus.TypingMessageEvent{
			ChannelID: payload.ChannelID,
			SenderID:  userID,
			IsTyping:  payload.IsTyping,
		})
	case dto.InChatIndicator:
		var payload dto.InChatIndicatorPayload
		if err := util.DecodeJson(incomingMsg.Payload, &payload); err != nil {
			log.Printf("WebSocket invalid in-chat indicator payload from user %s: %v", userID, err)
			h.sendToUser(userID, dto.ErrorMessage, dto.ErrorPayload{ErrorMsg: err.Error()})
			return
		}

		h.eventBus.Publish(bus.InChatMessageEvent{
			ChannelID: payload.ChannelID,
			UserID:    userID,
			IsInChat:  payload.IsInChat,
		})
	default:
		log.Printf("Unknown incoming type: %s", incomingMsg.Type)
	}
}

func (h *Hub) handleBroadcast(recipientIDs []string, eventType string, tempMessageID string, data interface{}) {
	switch eventType {
	case string(bus.BroadcastEventMessageCreated):
		// Handle new message
		var messageData dto.MessageResponse
		if err := util.DecodeJson(data, &messageData); err != nil {
			log.Printf("Failed to decode message payload: %v", err)
			return
		}

		// Send ACK back to sender
		ackResponse := dto.ACKMessagePayload{
			TempMessageID: tempMessageID,
			Message:       messageData,
		}
		h.sendToUser(messageData.SenderID, dto.ACKMessage, ackResponse)

		// Send message to recipients
		response := dto.SendMessagePayload{
			Message: messageData,
		}
		h.broadcastToUsers(recipientIDs, dto.SendMessage, response)
	case string(bus.BroadcastEventTypingStart), string(bus.BroadcastEventTypingStop):
		// Handle typing message
		h.broadcastToUsers(recipientIDs, dto.TypingIndicator, data)
	default:
		log.Printf("Unhandled broadcast event type: %s", eventType)
	}
}

// sendToUser is a private method to send a message to a specific user.
func (h *Hub) sendToUser(userID string, messageType dto.WebSocketMessageType, payload interface{}) {
	if client, ok := h.userClients[userID]; ok {
		msg := dto.WebSocketMessage{
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

func (h *Hub) broadcastToUsers(userIDs []string, messageType dto.WebSocketMessageType, payload interface{}) {
	msg := dto.WebSocketMessage{
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
