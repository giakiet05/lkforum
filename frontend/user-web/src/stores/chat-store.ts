import { writable, derived } from "svelte/store";
import type { ChannelResponse } from "../dtos/channel-dto";
import type { MessageResponse } from "../dtos/message-dto";

// --- Chat State ---

interface ChatState {
    // All channels for current user
    channels: ChannelResponse[];
    
    // Currently active channel ID
    activeChannelId: string | null;
    
    // Messages grouped by channel_id
    // Map<channel_id, MessageResponse[]>
    messagesByChannel: Map<string, MessageResponse[]>;
    
    // Typing indicators: Map<channel_id, Set<user_id>>
    typingUsers: Map<string, Set<string>>;
    
    // Loading states
    isLoadingChannels: boolean;
    isLoadingMessages: boolean;
}

const initialState: ChatState = {
    channels: [],
    activeChannelId: null,
    messagesByChannel: new Map(),
    typingUsers: new Map(),
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
 * Get unread count for each channel
 */
export const unreadCounts = derived(
    chatStore,
    $chatStore => {
        const counts = new Map<string, number>();
        
        $chatStore.messagesByChannel.forEach((messages, channelId) => {
            const unreadCount = messages.filter(m => !m.is_read).length;
            counts.set(channelId, unreadCount);
        });
        
        return counts;
    }
);

// --- Actions ---

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
    chatStore.update(state => {
        const newMessagesByChannel = new Map(state.messagesByChannel);
        const currentMessages = newMessagesByChannel.get(channelId) || [];
        
        // Avoid duplicates
        if (!currentMessages.find(m => m.id === message.id)) {
            newMessagesByChannel.set(channelId, [...currentMessages, message]);
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
 * Clear all chat data
 */
export function clearChatStore() {
    chatStore.set(initialState);
}
