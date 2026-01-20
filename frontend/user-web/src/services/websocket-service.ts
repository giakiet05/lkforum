import type { MessageResponse, MessageType } from "../dtos/message-dto";
import { getValidAccessToken } from "../auth/token";
import { API_BASE_URL } from "./api";
import { USER_KEY } from "../constants/auth-constants";

type MessageCallback = (message: MessageResponse) => void;
type ErrorCallback = (error: Event) => void;
type CloseCallback = (event: CloseEvent) => void;

/**
 * WebSocket Service for real-time messaging
 */
class WebSocketService {
    private ws: WebSocket | null = null;
    private messageCallbacks: MessageCallback[] = [];
    private errorCallbacks: ErrorCallback[] = [];
    private closeCallbacks: CloseCallback[] = [];
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectDelay = 1000; // Start with 1 second

    /**
     * Connect to WebSocket server
     */
    async connect(): Promise<void> {
        try {
            const token = await getValidAccessToken();
            if (!token) {
                console.error("❌ WebSocket: No valid access token");
                throw new Error("No valid access token");
            }

            // Convert http/https to ws/wss
            const wsUrl = API_BASE_URL.replace("http://", "ws://").replace("https://", "wss://");
            const url = `${wsUrl}/api/ws?token=${token}`;
            
            console.log("🔌 WebSocket: Connecting to:", wsUrl + "/api/ws");
            console.log("🔑 WebSocket: Using token:", token.substring(0, 20) + "...");

            this.ws = new WebSocket(url);

            this.ws.onopen = () => {
                console.log("✅ WebSocket connected successfully");
                this.reconnectAttempts = 0;
                this.reconnectDelay = 1000;
            };

            this.ws.onmessage = (event) => {
                console.log("📨 WebSocket: Raw message received:", event.data);
                try {
                    const wsMessage: any = JSON.parse(event.data);
                    console.log("📨 WebSocket: Parsed websocket message:", wsMessage);
                    
                    // Backend sends: { type: "send_message|ack_message|...", payload: {...} }
                    const messageType = wsMessage.type;
                    const payload = wsMessage.payload;

                    if (messageType === "send_message" || messageType === "ack_message") {
                        // Extract the actual message from payload
                        const message: MessageResponse = messageType === "send_message" 
                            ? payload.message 
                            : payload.message;
                        
                        console.log("📨 WebSocket: Extracted message:", message);
                        console.log("📨 WebSocket: Notifying", this.messageCallbacks.length, "callbacks");
                        this.messageCallbacks.forEach(callback => callback(message));
                    } else if (messageType === "typing") {
                        console.log("📨 WebSocket: Typing indicator:", payload);
                        // Handle typing indicator if needed
                    } else {
                        console.log("📨 WebSocket: Unhandled message type:", messageType, payload);
                    }
                } catch (error) {
                    console.error("❌ WebSocket: Failed to parse message:", error, "Raw:", event.data);
                }
            };

            this.ws.onerror = (event) => {
                console.error("❌ WebSocket error:", event);
                this.errorCallbacks.forEach(callback => callback(event));
            };

            this.ws.onclose = (event) => {
                console.log("🔌 WebSocket closed:", event.code, event.reason);
                this.closeCallbacks.forEach(callback => callback(event));

                // Auto-reconnect if not a normal closure
                if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
                    setTimeout(() => {
                        this.reconnectAttempts++;
                        console.log(`🔄 Reconnecting... Attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}`);
                        this.connect();
                    }, this.reconnectDelay);

                    // Exponential backoff
                    this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000); // Max 30 seconds
                }
            };
        } catch (error) {
            console.error("❌ Failed to connect WebSocket:", error);
            throw error;
        }
    }

    /**
     * Send message through WebSocket
     */
    sendMessage(channelId: string, content: string, type: MessageType = "text"): void {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            console.error("❌ WebSocket: Cannot send message - not connected. State:", this.ws?.readyState);
            throw new Error("WebSocket is not connected");
        }

        // Get current user info
        const userStr = localStorage.getItem(USER_KEY);
        if (!userStr) {
            console.error("❌ WebSocket: User not authenticated");
            throw new Error("User not authenticated");
        }
        const currentUser = JSON.parse(userStr);

        // Generate temp message ID
        const tempMessageId = `temp-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

        // Backend expects: { type: "new_message", payload: { ... } }
        const message = {
            type: "new_message",
            payload: {
                temp_message_id: tempMessageId,
                channel_id: channelId,
                sender_username: currentUser.username,
                type,
                content
            }
        };

        console.log("📤 WebSocket: Sending message:", message);
        this.ws.send(JSON.stringify(message));
    }

    /**
     * Register callback for incoming messages
     */
    onMessage(callback: MessageCallback): void {
        this.messageCallbacks.push(callback);
    }

    /**
     * Register callback for errors
     */
    onError(callback: ErrorCallback): void {
        this.errorCallbacks.push(callback);
    }

    /**
     * Register callback for connection close
     */
    onClose(callback: CloseCallback): void {
        this.closeCallbacks.push(callback);
    }

    /**
     * Remove message callback
     */
    offMessage(callback: MessageCallback): void {
        this.messageCallbacks = this.messageCallbacks.filter(cb => cb !== callback);
    }

    /**
     * Remove error callback
     */
    offError(callback: ErrorCallback): void {
        this.errorCallbacks = this.errorCallbacks.filter(cb => cb !== callback);
    }

    /**
     * Remove close callback
     */
    offClose(callback: CloseCallback): void {
        this.closeCallbacks = this.closeCallbacks.filter(cb => cb !== callback);
    }

    /**
     * Disconnect WebSocket
     */
    disconnect(): void {
        if (this.ws) {
            this.ws.close(1000, "User disconnected");
            this.ws = null;
        }
        this.messageCallbacks = [];
        this.errorCallbacks = [];
        this.closeCallbacks = [];
        this.reconnectAttempts = 0;
    }

    /**
     * Check if WebSocket is connected
     */
    isConnected(): boolean {
        return this.ws !== null && this.ws.readyState === WebSocket.OPEN;
    }

    /**
     * Get current WebSocket connection state
     */
    getState(): number {
        return this.ws?.readyState ?? WebSocket.CLOSED;
    }
}

// Export singleton instance
export const websocketService = new WebSocketService();
