import type { MessageResponse, MessageType } from "../dtos/message-dto";
import { getValidAccessToken } from "../auth/token";
import { API_BASE_URL } from "./api";

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
                throw new Error("No valid access token");
            }

            // Convert http/https to ws/wss
            const wsUrl = API_BASE_URL.replace("http://", "ws://").replace("https://", "wss://");
            const url = `${wsUrl}/api/ws?token=${token}`;

            this.ws = new WebSocket(url);

            this.ws.onopen = () => {
                console.log("WebSocket connected");
                this.reconnectAttempts = 0;
                this.reconnectDelay = 1000;
            };

            this.ws.onmessage = (event) => {
                try {
                    const message: MessageResponse = JSON.parse(event.data);
                    this.messageCallbacks.forEach(callback => callback(message));
                } catch (error) {
                    console.error("Failed to parse WebSocket message:", error);
                }
            };

            this.ws.onerror = (event) => {
                console.error("WebSocket error:", event);
                this.errorCallbacks.forEach(callback => callback(event));
            };

            this.ws.onclose = (event) => {
                console.log("WebSocket closed:", event.code, event.reason);
                this.closeCallbacks.forEach(callback => callback(event));

                // Auto-reconnect if not a normal closure
                if (event.code !== 1000 && this.reconnectAttempts < this.maxReconnectAttempts) {
                    setTimeout(() => {
                        this.reconnectAttempts++;
                        console.log(`Reconnecting... Attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}`);
                        this.connect();
                    }, this.reconnectDelay);

                    // Exponential backoff
                    this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000); // Max 30 seconds
                }
            };
        } catch (error) {
            console.error("Failed to connect WebSocket:", error);
            throw error;
        }
    }

    /**
     * Send message through WebSocket
     */
    sendMessage(channelId: string, content: string, type: MessageType = "text"): void {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            throw new Error("WebSocket is not connected");
        }

        // Get current user info
        const userStr = localStorage.getItem("user_key");
        if (!userStr) {
            throw new Error("User not authenticated");
        }
        const currentUser = JSON.parse(userStr);

        const message = {
            channel_id: channelId,
            sender_id: currentUser.id,
            type,
            content
        };

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
