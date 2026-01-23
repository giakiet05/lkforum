import { writable, derived } from "svelte/store";
import type { ChannelResponse } from "../dtos/channel-dto";
import type { MessageResponse } from "../dtos/message-dto";

// --- Chat State ---

interface ChatState {
    // Current user ID (to filter own messages from unread count)
    currentUserId: string | null;
    
    // All channels for current user
    channels: ChannelResponse[];
    
    // Currently active channel ID
    activeChannelId: string | null;
    
    // Messages grouped by channel_id
    // Map<channel_id, MessageResponse[]>
    messagesByChannel: Map<string, MessageResponse[]>;
    
    // Typing indicators: Map<channel_id, Set<user_id>>
    typingUsers: Map<string, Set<string>>;
    
    // Online users: Set<user_id>
    onlineUsers: Set<string>;
    
    // Loading states
    isLoadingChannels: boolean;
    isLoadingMessages: boolean;
}

const initialState: ChatState = {
    currentUserId: null,
    channels: [],
    activeChannelId: null,
    messagesByChannel: new Map(),
    typingUsers: new Map(),
    onlineUsers: new Set(),
    isLoadingChannels: false,
    isLoadingMessages: false
};

// Create writable store
export const chatStore = writable<ChatState>(initialState);

// --- Derived Stores ---

/**
 * Get currently active channel
 */
export const activeChannel = derived(
    chatStore,
    $chatStore => {
        if (!$chatStore.activeChannelId) return null;
        return $chatStore.channels.find(c => c.id === $chatStore.activeChannelId) || null;
    }
);

/**
 * Get messages for active channel
 */
export const activeChannelMessages = derived(
    chatStore,
    $chatStore => {
        if (!$chatStore.activeChannelId) return [];
        return $chatStore.messagesByChannel.get($chatStore.activeChannelId) || [];
    }
);

/**
 * Get unread count for each channel (excluding own messages)
 */
export const unreadCounts = derived(
    chatStore,
    $chatStore => {
        const counts = new Map<string, number>();
        const currentUserId = $chatStore.currentUserId;
        
        $chatStore.messagesByChannel.forEach((messages, channelId) => {
            // Only count unread messages that are NOT from the current user
            const unreadCount = messages.filter(m => !m.is_read && m.sender_id !== currentUserId).length;
            counts.set(channelId, unreadCount);
        });
        
        return counts;
    }
);

/**
 * Get total unread messages count across all channels (excluding own messages)
 */
export const totalUnreadCount = derived(
    chatStore,
    $chatStore => {
        let total = 0;
        const currentUserId = $chatStore.currentUserId;
        
        $chatStore.messagesByChannel.forEach((messages, channelId) => {
            // Only count unread messages that are NOT from the current user
            const unreadCount = messages.filter(m => {
                const isUnread = !m.is_read;
                const isFromOther = m.sender_id !== currentUserId;
                return isUnread && isFromOther;
            }).length;
            total += unreadCount;
        });
        console.log("📊 [chat-store] totalUnreadCount:", total, "currentUserId:", currentUserId);
        return total;
    }
);

// --- Actions ---

/**
 * Set current user ID (needed to filter own messages from unread count)
 */
export function setCurrentUserId(userId: string | null) {
    chatStore.update(state => ({
        ...state,
        currentUserId: userId
    }));
}

/**
 * Set all channels
 */
export function setChannels(channels: ChannelResponse[]) {
    chatStore.update(state => ({
        ...state,
        channels,
        isLoadingChannels: false
    }));
}

/**
 * Add a new channel
 */
export function addChannel(channel: ChannelResponse) {
    chatStore.update(state => ({
        ...state,
        channels: [channel, ...state.channels]
    }));
}

/**
 * Update an existing channel
 */
export function updateChannel(channelId: string, updates: Partial<ChannelResponse>) {
    chatStore.update(state => ({
        ...state,
        channels: state.channels.map(c => 
            c.id === channelId ? { ...c, ...updates } : c
        )
    }));
}

/**
 * Remove a channel
 */
export function removeChannel(channelId: string) {
    chatStore.update(state => {
        const newMessagesByChannel = new Map(state.messagesByChannel);
        newMessagesByChannel.delete(channelId);

        return {
            ...state,
            channels: state.channels.filter(c => c.id !== channelId),
            messagesByChannel: newMessagesByChannel,
            activeChannelId: state.activeChannelId === channelId ? null : state.activeChannelId
        };
    });
}

/**
 * Set active channel
 */
export function setActiveChannel(channelId: string | null) {
    chatStore.update(state => ({
        ...state,
        activeChannelId: channelId
    }));
}

/**
 * Set messages for a channel
 */
export function setMessages(channelId: string, messages: MessageResponse[]) {
    chatStore.update(state => {
        const newMessagesByChannel = new Map(state.messagesByChannel);
        newMessagesByChannel.set(channelId, messages);

        return {
            ...state,
            messagesByChannel: newMessagesByChannel,
            isLoadingMessages: false
        };
    });
}

/**
 * Add a new message to a channel
 */
export function addMessage(channelId: string, message: MessageResponse) {
    console.log("💾 [chat-store] Adding message to channel:", channelId, message);
    chatStore.update(state => {
        const newMessagesByChannel = new Map(state.messagesByChannel);
        const currentMessages = newMessagesByChannel.get(channelId) || [];
        
        console.log("💾 [chat-store] Current messages count:", currentMessages.length);
        
        // Avoid duplicates
        if (!currentMessages.find(m => m.id === message.id)) {
            newMessagesByChannel.set(channelId, [...currentMessages, message]);
            console.log("💾 [chat-store] Message added! New count:", currentMessages.length + 1);
        } else {
            console.log("💾 [chat-store] Message already exists, skipping");
        }

        return {
            ...state,
            messagesByChannel: newMessagesByChannel
        };
    });
}

/**
 * Update a message
 */
export function updateMessage(channelId: string, messageId: string, updates: Partial<MessageResponse>) {
    chatStore.update(state => {
        const newMessagesByChannel = new Map(state.messagesByChannel);
        const messages = newMessagesByChannel.get(channelId) || [];
        
        newMessagesByChannel.set(
            channelId,
            messages.map(m => m.id === messageId ? { ...m, ...updates } : m)
        );

        return {
            ...state,
            messagesByChannel: newMessagesByChannel
        };
    });
}

/**
 * Remove a message
 */
export function removeMessage(channelId: string, messageId: string) {
    chatStore.update(state => {
        const newMessagesByChannel = new Map(state.messagesByChannel);
        const messages = newMessagesByChannel.get(channelId) || [];
        
        newMessagesByChannel.set(
            channelId,
            messages.filter(m => m.id !== messageId)
        );

        return {
            ...state,
            messagesByChannel: newMessagesByChannel
        };
    });
}

/**
 * Mark messages as read for a channel
 */
export function markChannelAsRead(channelId: string) {
    chatStore.update(state => {
        const newMessagesByChannel = new Map(state.messagesByChannel);
        const messages = newMessagesByChannel.get(channelId) || [];
        
        newMessagesByChannel.set(
            channelId,
            messages.map(m => ({ ...m, is_read: true }))
        );

        return {
            ...state,
            messagesByChannel: newMessagesByChannel
        };
    });
}

/**
 * Add typing indicator
 */
export function addTypingUser(channelId: string, userId: string) {
    chatStore.update(state => {
        const newTypingUsers = new Map(state.typingUsers);
        const users = newTypingUsers.get(channelId) || new Set();
        users.add(userId);
        newTypingUsers.set(channelId, users);

        return {
            ...state,
            typingUsers: newTypingUsers
        };
    });
}

/**
 * Remove typing indicator
 */
export function removeTypingUser(channelId: string, userId: string) {
    chatStore.update(state => {
        const newTypingUsers = new Map(state.typingUsers);
        const users = newTypingUsers.get(channelId);
        
        if (users) {
            users.delete(userId);
            if (users.size === 0) {
                newTypingUsers.delete(channelId);
            } else {
                newTypingUsers.set(channelId, users);
            }
        }

        return {
            ...state,
            typingUsers: newTypingUsers
        };
    });
}

/**
 * Set loading states
 */
export function setLoadingChannels(isLoading: boolean) {
    chatStore.update(state => ({
        ...state,
        isLoadingChannels: isLoading
    }));
}

export function setLoadingMessages(isLoading: boolean) {
    chatStore.update(state => ({
        ...state,
        isLoadingMessages: isLoading
    }));
}

/**
 * Add user to online users
 */
export function addOnlineUser(userId: string) {
    chatStore.update(state => {
        const newOnlineUsers = new Set(state.onlineUsers);
        newOnlineUsers.add(userId);
        return {
            ...state,
            onlineUsers: newOnlineUsers
        };
    });
}

/**
 * Remove user from online users
 */
export function removeOnlineUser(userId: string) {
    chatStore.update(state => {
        const newOnlineUsers = new Set(state.onlineUsers);
        newOnlineUsers.delete(userId);
        return {
            ...state,
            onlineUsers: newOnlineUsers
        };
    });
}

/**
 * Clear all chat data
 */
export function clearChatStore() {
    chatStore.set(initialState);
}
