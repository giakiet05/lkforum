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

type WebSocketController struct {
	wsHub *ws.Hub
}

func NewWebSocketController(hub *ws.Hub) *WebSocketController {
	return &WebSocketController{wsHub: hub}
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade WebSocket"})
		return
	}

	// Create a new client instance.
	client := ws.NewClient(c.wsHub, conn, userID)

	// Register the client with the hub.
	c.wsHub.RegisterClient(client)

	// Start the client's processing goroutines.
	client.Serve()
}
