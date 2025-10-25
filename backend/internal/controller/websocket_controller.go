package controller

import (
	"log"
	"net/http"

	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/platform/ws"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// In production, check the request origin to prevent CSRF.
		// e.g., return r.Header.Get("Origin") == config.Cfg.FrontendURL
		return true
	},
}

type WebSocketController struct{}

func NewWebSocketController() *WebSocketController {
	return &WebSocketController{}
}

// HandleConnections handles the WebSocket connection requests.
func (c *WebSocketController) HandleConnections(ctx *gin.Context) {
	authUser, exists := ctx.Get("authUser")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
		return
	}
	userID := authUser.(auth.AuthUser).ID

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection for user %s: %v", userID, err)
		return
	}

	// Create a new client instance.
	client := ws.NewClient(ws.WSHub, conn, userID)

	// Register the client with the hub.
	ws.WSHub.RegisterClient(client)

	// Start the client's processing goroutines.
	client.Serve()
}
