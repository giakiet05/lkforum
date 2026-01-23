import type { MessageResponse, MessageType } from "../dtos/message-dto";
import { getValidAccessToken } from "../auth/token";
import { API_BASE_URL } from "./api";
import { USER_KEY } from "../constants/auth-constants";

export type WebSocketMessageType =
  | 'new_message'
  | 'send_message'
  | 'ack_message'
  | 'typing'
  | 'in_chat'
  | 'new_notification'
  | 'error'
  | 'presence_online'
  | 'presence_offline';

export interface WebSocketMessage {
  type: WebSocketMessageType;
  payload: any;
}

export interface PresencePayload {
  user_id: string;
}

// Handler types
type MessageHandler = (payload: any) => void;
type ErrorCallback = (error: Event) => void;
type CloseCallback = (event: CloseEvent) => void;

/**
 * WebSocket Service for real-time messaging
 */
class WebSocketService {
    private ws: WebSocket | null = null;
    private handlers: Map<WebSocketMessageType, MessageHandler[]> = new Map();
    private errorCallbacks: ErrorCallback[] = [];
    private closeCallbacks: CloseCallback[] = [];
    private reconnectAttempts = 0;
    private maxReconnectAttempts = 5;
    private reconnectDelay = 1000;
    private isConnecting = false;
    private shouldReconnect = true;

    /**
     * Connect to WebSocket server
     */
    async connect(): Promise<void> {
        if (this.ws?.readyState === WebSocket.OPEN) return Promise.resolve();
        if (this.isConnecting) return Promise.resolve();

        return new Promise(async (resolve, reject) => {
            try {
                this.isConnecting = true;
                const token = await getValidAccessToken();
                if (!token) {
                    this.isConnecting = false;
                    console.error("❌ WebSocket: No valid access token");
                    throw new Error("No valid access token");
                }

                // Convert http/https to ws/wss
                const wsUrl = API_BASE_URL.replace("http://", "ws://").replace("https://", "wss://");
                const url = `${wsUrl}/api/ws?token=${token}`;
                
                console.log("🔌 WebSocket: Connecting to:", wsUrl + "/api/ws");

                this.ws = new WebSocket(url);

                this.ws.onopen = () => {
                    console.log("✅ WebSocket Connected");
                    this.isConnecting = false;
                    this.reconnectAttempts = 0;
                    this.reconnectDelay = 1000;
                    resolve();
                };

                this.ws.onmessage = (event) => {
                    try {
                        const message: WebSocketMessage = JSON.parse(event.data);
                        console.log("📨 WebSocket message:", message.type, message.payload);
                        this.handleMessage(message.type, message.payload);
                    } catch (error) {
                        console.error("❌ Failed to parse WS message:", error);
                    }
                };

                this.ws.onerror = (event) => {
                    console.error("❌ WebSocket error:", event);
                    this.isConnecting = false;
                    this.errorCallbacks.forEach(callback => callback(event));
                    reject(event);
                };

                this.ws.onclose = (event) => {
                    console.log("🔌 WebSocket closed:", event.code, event.reason);
                    this.isConnecting = false;
                    this.ws = null;
                    this.closeCallbacks.forEach(callback => callback(event));

                    // Auto-reconnect
                    if (this.shouldReconnect && this.reconnectAttempts < this.maxReconnectAttempts) {
                        this.reconnectAttempts++;
                        console.log(`🔄 Reconnecting... Attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts}`);
                        setTimeout(() => this.connect(), this.reconnectDelay);
                        this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
                    }
                };
            } catch (error) {
                this.isConnecting = false;
                console.error("❌ Failed to connect WebSocket:", error);
                reject(error);
            }
        });
    }

    /**
     * Send generic WebSocket message
     */
    send(type: WebSocketMessageType, payload: any): boolean {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            console.error("❌ WebSocket not connected");
            return false;
        }
        try {
            this.ws.send(JSON.stringify({ type, payload }));
            return true;
        } catch (error) {
            console.error("❌ Failed to send message:", error);
            return false;
        }
    }

    /**
     * Send message through WebSocket
     */
    sendMessage(channelId: string, content: string, type: MessageType = "text"): void {
        const userStr = localStorage.getItem(USER_KEY);
        if (!userStr) throw new Error("User not authenticated");
        const currentUser = JSON.parse(userStr);

        const tempMessageId = `temp-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

        const payload = {
            temp_message_id: tempMessageId,
            channel_id: channelId,
            sender_username: currentUser.username,
            type,
            content
        };

        console.log("📤 Sending message:", payload);
        this.send("new_message", payload);
    }

    /**
     * Send typing indicator
     */
    sendTyping(channelId: string, isTyping: boolean): void {
        this.send("typing", {
            channel_id: channelId,
            is_typing: isTyping
        });
    }

    /**
     * Send in-chat indicator
     */
    sendInChat(channelId: string, isInChat: boolean): void {
        this.send("in_chat", {
            channel_id: channelId,
            is_in_chat: isInChat
        });
    }

    /**
     * Register handler for specific message type
     */
    on(type: WebSocketMessageType, handler: MessageHandler): void {
        if (!this.handlers.has(type)) this.handlers.set(type, []);
        this.handlers.get(type)!.push(handler);
    }

    /**
     * Remove handler for specific message type
     */
    off(type: WebSocketMessageType, handler: MessageHandler): void {
        const handlers = this.handlers.get(type);
        if (handlers) {
            const index = handlers.indexOf(handler);
            if (index > -1) handlers.splice(index, 1);
        }
    }

    /**
     * Handle incoming message
     */
    private handleMessage(type: WebSocketMessageType, payload: any): void {
        const handlers = this.handlers.get(type);
        if (handlers) {
            handlers.forEach(handler => handler(payload));
        }
    }

    /**
     * Legacy: Register callback for incoming messages (backward compatibility)
     */
    onMessage(callback: (message: MessageResponse) => void): void {
        // Wrapper to handle both send_message and ack_message
        const handler = (payload: any) => {
            const message = payload.message || payload;
            callback(message);
        };
        this.on("send_message", handler);
        this.on("ack_message", handler);
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
     * Legacy: Remove message callback (backward compatibility)
     */
    offMessage(callback: (message: MessageResponse) => void): void {
        // Cannot easily remove without keeping references, so we'll clear all
        // In practice, this should use the new off() method
        console.warn("offMessage is deprecated, use off() method instead");
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
        this.shouldReconnect = false;
        if (this.ws) {
            this.ws.close(1000, "User disconnected");
            this.ws = null;
        }
        this.handlers.clear();
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
