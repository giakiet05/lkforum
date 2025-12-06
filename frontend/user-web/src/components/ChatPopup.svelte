<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { push } from "svelte-spa-router";
  import { authStore } from "../stores/auth-store";
  import {
    chatStore,
    activeChannel,
    activeChannelMessages,
    setChannels,
    setActiveChannel,
    setMessages,
    addMessage,
    markChannelAsRead,
    setLoadingChannels,
    setLoadingMessages,
  } from "../stores/chat-store";
  import {
    getChannelsByUser,
    updateChannel,
    createChannel,
  } from "../services/channel-service";
  import { getMessages } from "../services/message-service";
  import { websocketService } from "../services/websocket-service";
  import type { ChannelResponse } from "../dtos/channel-dto";
  import type { MessageResponse } from "../dtos/message-dto";

  type ChatPopupProps = {
    show: boolean;
    onClose: () => void;
  };

  let { show, onClose }: ChatPopupProps = $props();

  let messageInput = $state("");
  let showChatMenu = $state(false);

  // Get current user
  const currentUser = $derived($authStore.user);
  const channels = $derived($chatStore.channels);
  const activeChannelData = $derived($activeChannel);
  const reports = $derived($activeChannelMessages);
  const isLoadingChannels = $derived($chatStore.isLoadingChannels);
  const isLoadingMessages = $derived($chatStore.isLoadingMessages);

  // Get other member info from active channel
  const otherMember = $derived(() => {
    if (!activeChannelData || !currentUser) return null;
    return (
      activeChannelData.members.find((m) => m.user_id !== currentUser.id) ||
      null
    );
  });

  // Get channel settings for current user
  const channelSettings = $derived(() => {
    if (!activeChannelData || !currentUser) return null;
    return (
      activeChannelData.settings.find((s) => s.user_id === currentUser.id) ||
      null
    );
  });

  // Format timestamp
  function formatTime(isoString: string): string {
    const date = new Date(isoString);
    return date.toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  // Format relative time
  function formatRelativeTime(isoString: string): string {
    const now = new Date();
    const date = new Date(isoString);
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);

    if (diffMins < 1) return "Just now";
    if (diffMins < 60) return `${diffMins}m`;

    const diffHours = Math.floor(diffMins / 60);
    if (diffHours < 24) return `${diffHours}h`;

    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays}d`;
  }

  // Get last message for a channel
  function getLastMessage(channelId: string): string {
    const channelMessages = $chatStore.messagesByChannel.get(channelId) || [];
    if (channelMessages.length === 0) return "No reports yet";

    const lastMsg = channelMessages[channelMessages.length - 1];
    const isCurrentUser = lastMsg.sender_id === currentUser?.id;
    const prefix = isCurrentUser ? "You: " : `${lastMsg.sender_username}: `;

    return prefix + lastMsg.content;
  }

  // Get last message time for a channel
  function getLastMessageTime(channelId: string): string {
    const channelMessages = $chatStore.messagesByChannel.get(channelId) || [];
    if (channelMessages.length === 0) return "";

    const lastMsg = channelMessages[channelMessages.length - 1];
    return formatRelativeTime(lastMsg.created_at);
  }

  // Get unread count for a channel
  function getUnreadCount(channelId: string): number {
    const channelMessages = $chatStore.messagesByChannel.get(channelId) || [];
    return channelMessages.filter(
      (m) => !m.is_read && m.sender_id !== currentUser?.id
    ).length;
  }

  // Load channels
  async function loadChannels() {
    if (!currentUser) {
      console.log("[ChatPopup] No current user, skipping load");
      return;
    }

    try {
      console.log("[ChatPopup] Loading channels for user:", currentUser.id);
      setLoadingChannels(true);
      const response = await getChannelsByUser(currentUser.id);
      console.log("[ChatPopup] Loaded channels:", response);
      setChannels(response.channels);

      // Auto-select first channel
      if (response.channels.length > 0 && !$chatStore.activeChannelId) {
        await handleSelectChannel(response.channels[0]);
      }
    } catch (error) {
      console.error("[ChatPopup] Failed to load channels:", error);
    } finally {
      setLoadingChannels(false);
    }
  }

  // Load reports for a channel
  async function loadMessagesForChannel(channelId: string) {
    try {
      setLoadingMessages(true);
      const channelMessages = await getMessages({ channel_id: channelId });
      setMessages(channelId, channelMessages);
      console.log(
        `✅ Loaded ${channelMessages.length} messages for channel ${channelId}`
      );

      // Mark as read
      markChannelAsRead(channelId);

      // Scroll to bottom
      setTimeout(scrollToBottom, 100);
    } catch (error) {
      console.error("Failed to load messages:", error);
      console.log(
        `⚠️ API error - keeping any existing messages from WebSocket`
      );
      setLoadingMessages(false);
    }
  }

  // Handle channel selection
  async function handleSelectChannel(channel: ChannelResponse) {
    setActiveChannel(channel.id);
    await loadMessagesForChannel(channel.id);
  }

  // Handle send message
  async function handleSendMessage() {
    if (!messageInput.trim() || !$chatStore.activeChannelId) return;

    try {
      websocketService.sendMessage(
        $chatStore.activeChannelId,
        messageInput.trim(),
        "text"
      );
      messageInput = "";

      setTimeout(scrollToBottom, 100);
    } catch (error) {
      console.error("Failed to send message:", error);
    }
  }

  // Handle incoming WebSocket reports
  function handleIncomingMessage(message: MessageResponse) {
    console.log("📨 [ChatPopup] Received WebSocket message:", message);
    addMessage(message.channel_id, message);

    if (message.channel_id === $chatStore.activeChannelId) {
      setTimeout(scrollToBottom, 100);
    }
  }

  // Scroll to bottom
  function scrollToBottom() {
    const messagesArea = document.querySelector(".popup-reports-area");
    if (messagesArea) {
      messagesArea.scrollTop = messagesArea.scrollHeight;
    }
  }

  // Toggle mute
  async function toggleMute() {
    if (!activeChannelData || !currentUser) return;

    const currentSettings = channelSettings();
    const newNotificationState = !(currentSettings?.notification ?? true);

    try {
      await updateChannel({
        channel_id: activeChannelData.id,
        notification: newNotificationState,
      });

      await loadChannels();
    } catch (error) {
      console.error("Failed to toggle mute:", error);
    }
  }

  function toggleChatMenu() {
    showChatMenu = !showChatMenu;
  }

  function handleExpand() {
    push("/reports");
    onClose();
  }

  // DEBUG: Create test channel with hardcoded user
  async function handleDebugCreateChannel() {
    if (!currentUser) {
      alert("Not logged in!");
      return;
    }

    console.log("=== YOUR USER INFO ===");
    console.log("Your User ID:", currentUser.id);
    console.log("Your Username:", currentUser.username);
    console.log("Full User Object:", currentUser);
    console.log("====================");

    // Quick test: Use hardcoded user IDs from your communities
    const testUserIds = [
      { id: "6915936f2ae7e4bba023dad1", username: "sample0_creator" },
      { id: "6905e0bcbaa66c4cb5effe24", username: "golang_enthusiast" },
    ];

    // Show selection
    const choice = prompt(
      `YOUR USER ID: ${currentUser.id}\n\nChoose test user to chat with:\n1. ${testUserIds[0].username} (${testUserIds[0].id})\n2. ${testUserIds[1].username} (${testUserIds[1].id})\n\nOr enter custom User ID:`
    );

    if (!choice) return;

    let targetUserId: string;
    let targetUsername: string;

    if (choice === "1") {
      targetUserId = testUserIds[0].id;
      targetUsername = testUserIds[0].username;
    } else if (choice === "2") {
      targetUserId = testUserIds[1].id;
      targetUsername = testUserIds[1].username;
    } else {
      targetUserId = choice;
      targetUsername = prompt("Enter their username:") || "Unknown";
    }

    try {
      console.log("[ChatPopup] Creating test channel with:", targetUserId);
      const newChannel = await createChannel(
        targetUserId,
        targetUsername,
        "" // Empty avatar for now
      );
      console.log("[ChatPopup] Channel created:", newChannel);

      // Reload channels
      await loadChannels();
    } catch (error) {
      console.error("[ChatPopup] Failed to create channel:", error);
      alert("Failed to create channel: " + error);
    }
  }

  // Initialize when popup shows
  // WebSocket lifecycle
  $effect(() => {
    if (show && currentUser) {
      // Connect WebSocket and register message handler
      if (!websocketService.isConnected()) {
        websocketService
          .connect()
          .then(() => {
            console.log("📞 Registering WebSocket message handler");
            websocketService.onMessage(handleIncomingMessage);
          })
          .catch((error) => {
            console.error("Failed to connect WebSocket:", error);
          });
      } else {
        console.log(
          "📞 Registering WebSocket message handler (already connected)"
        );
        websocketService.onMessage(handleIncomingMessage);
      }
    }
  });

  // Load channels on open
  $effect(() => {
    if (show && currentUser && channels.length === 0) {
      loadChannels();
    }
  });

  // Debug: Log channels changes
  $effect(() => {
    console.log("[ChatPopup] Channels changed:", channels.length, channels);
  });

  // Cleanup
  onDestroy(() => {
    websocketService.offMessage(handleIncomingMessage);
  });
</script>

{#if show}
  <div class="chat-popup">
    <!-- Popup Header -->
    <div class="popup-header">
      <h2>Tin nhắn</h2>
      <div class="popup-header-actions">
        <button
          class="header-icon-btn debug-btn"
          onclick={handleDebugCreateChannel}
          title="Create Test Channel"
        >
          🔧
        </button>
        <button class="header-icon-btn" onclick={handleExpand} title="Expand">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M13 3h4v4M7 17H3v-4M17 3l-7 7M3 17l7-7"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </button>
        <button class="header-icon-btn" onclick={onClose} title="Close">
          <svg width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path
              d="M15 5L5 15M5 5l10 10"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
            />
          </svg>
        </button>
      </div>
    </div>

    <div class="popup-container">
      <!-- Left Side - Conversations List -->
      <div class="popup-conversations">
        <div class="popup-conversations-list">
          {#if isLoadingChannels}
            <div class="popup-loading">Loading...</div>
          {:else if channels.length === 0}
            <div class="popup-empty">
              <p>No conversations</p>
              <button
                class="debug-create-channel-btn"
                onclick={handleDebugCreateChannel}
              >
                🔧 Debug: Create Test Channel
              </button>
            </div>
          {:else}
            {#each channels as channel (channel.id)}
              {@const member = channel.members.find(
                (m) => m.user_id !== currentUser?.id
              )}
              {@const unreadCount = getUnreadCount(channel.id)}
              {#if member}
                <button
                  class="popup-conversation-item"
                  class:active={$chatStore.activeChannelId === channel.id}
                  class:unread={unreadCount > 0}
                  onclick={() => handleSelectChannel(channel)}
                >
                  <div class="popup-avatar-wrapper">
                    <img
                      src={member.avatar || "/avatar.jpg"}
                      alt={member.username}
                      class="popup-avatar"
                    />
                  </div>
                  <div class="popup-conversation-info">
                    <div class="popup-conversation-top">
                      <span class="popup-conversation-name"
                        >{member.username}</span
                      >
                      <span class="popup-conversation-time"
                        >{getLastMessageTime(channel.id)}</span
                      >
                    </div>
                    <div class="popup-conversation-bottom">
                      <p class="popup-conversation-preview">
                        {getLastMessage(channel.id)}
                      </p>
                      {#if unreadCount > 0}
                        <span class="popup-unread-badge">{unreadCount}</span>
                      {/if}
                    </div>
                  </div>
                </button>
              {/if}
            {/each}
          {/if}
        </div>
      </div>

      <!-- Right Side - Chat Detail -->
      <div class="popup-chat-detail">
        {#if activeChannelData && otherMember()}
          <!-- Chat Header -->
          <div class="popup-chat-header">
            <div class="popup-chat-user-info">
              <img
                src={otherMember()?.avatar || "/avatar.jpg"}
                alt={otherMember()?.username || "User"}
                class="popup-chat-avatar"
              />
              <div class="popup-chat-user-details">
                <h3>{otherMember()?.username || "Unknown"}</h3>
              </div>
            </div>
            <div class="popup-chat-actions">
              <div class="popup-chat-menu-wrapper">
                <button
                  class="popup-action-icon-btn"
                  onclick={toggleChatMenu}
                  title="More options"
                >
                  <svg
                    width="18"
                    height="18"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                  >
                    <path
                      d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z"
                    />
                  </svg>
                </button>
                {#if showChatMenu}
                  <div class="popup-chat-dropdown">
                    <button
                      class="popup-dropdown-item"
                      class:muted={!(channelSettings()?.notification ?? true)}
                      onclick={() => {
                        toggleMute();
                        showChatMenu = false;
                      }}
                    >
                      <img
                        src="/muted_icon.svg"
                        alt="Mute"
                        width="18"
                        height="18"
                      />
                      <span
                        >{(channelSettings()?.notification ?? true)
                          ? "Mute"
                          : "Unmute"}</span
                      >
                    </button>
                  </div>
                {/if}
              </div>
            </div>
          </div>

          <!-- Messages Area -->
          <div class="popup-reports-area">
            {#if isLoadingMessages}
              <div class="popup-loading">Loading reports...</div>
            {:else}
              <div class="popup-reports-wrapper">
                {#each reports as message (message.id)}
                  {@const isSent = message.sender_id === currentUser?.id}
                  <div class="popup-message-row" class:sent={isSent}>
                    {#if !isSent}
                      <img
                        src={otherMember()?.avatar || "/avatar.jpg"}
                        alt={message.sender_username}
                        class="popup-message-avatar"
                      />
                    {/if}
                    <div class="popup-message-bubble" class:sent={isSent}>
                      <p>{message.content}</p>
                      <span class="popup-message-time"
                        >{formatTime(message.created_at)}</span
                      >
                    </div>
                  </div>
                {/each}
              </div>
            {/if}
          </div>

          <!-- Message Input -->
          <div class="popup-message-input">
            <button class="popup-emoji-btn" title="Emoji">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none">
                <circle
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="2"
                />
                <path
                  d="M8 14s1.5 2 4 2 4-2 4-2"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                />
                <circle cx="9" cy="9" r="1" fill="currentColor" />
                <circle cx="15" cy="9" r="1" fill="currentColor" />
              </svg>
            </button>
            <input
              type="text"
              placeholder="Nhắn tin..."
              bind:value={messageInput}
              disabled={!websocketService.isConnected()}
              onkeydown={(e) => {
                if (e.key === "Enter" && !e.shiftKey) {
                  e.preventDefault();
                  handleSendMessage();
                }
              }}
            />
            <button
              class="popup-send-btn"
              onclick={handleSendMessage}
              disabled={!messageInput.trim() || !websocketService.isConnected()}
              title="Send"
            >
              <img src="/send_icon.svg" alt="Send" width="20" height="20" />
            </button>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<style>
  .chat-popup {
    position: fixed;
    bottom: 0;
    right: 80px;
    width: 680px;
    height: 520px;
    background: white;
    border-radius: 12px 12px 0 0;
    display: flex;
    flex-direction: column;
    box-shadow: 0 -2px 16px rgba(0, 0, 0, 0.15);
    z-index: 9999;
    overflow: hidden;
  }

  /* Popup Header */
  .popup-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid #edeff1;
    background: white;
  }

  .popup-header h2 {
    margin: 0;
    font-size: 16px;
    font-weight: 700;
    color: #1c1c1c;
  }

  .popup-header-actions {
    display: flex;
    gap: 8px;
  }

  .header-icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    color: #7c7c7c;
    transition: background 0.2s;
  }

  .header-icon-btn:hover {
    background: #f6f7f8;
  }

  .header-icon-btn.debug-btn {
    font-size: 16px;
    color: var(--blue--);
  }

  .header-icon-btn.debug-btn:hover {
    background: #e8f4fd;
  }

  /* Container */
  .popup-container {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  /* Conversations List */
  .popup-conversations {
    width: 240px;
    border-right: 1px solid #edeff1;
    display: flex;
    flex-direction: column;
    background: white;
  }

  .popup-conversations-list {
    flex: 1;
    overflow-y: auto;
  }

  .popup-conversation-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 16px;
    width: 100%;
    border: none;
    background: transparent;
    cursor: pointer;
    text-align: left;
    transition: background 0.2s;
    border-left: 3px solid transparent;
  }

  .popup-conversation-item:hover {
    background: #f6f7f8;
  }

  .popup-conversation-item.active {
    background: #e8f4fd;
    border-left-color: var(--blue--);
  }

  .popup-avatar-wrapper {
    position: relative;
    flex-shrink: 0;
  }

  .popup-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    object-fit: cover;
  }

  .popup-conversation-info {
    flex: 1;
    min-width: 0;
  }

  .popup-conversation-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 4px;
  }

  .popup-conversation-name {
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .popup-conversation-item.unread .popup-conversation-name {
    font-weight: 700;
  }

  .popup-conversation-item.unread .popup-conversation-preview {
    color: #1c1c1c;
    font-weight: 600;
  }

  .popup-conversation-time {
    font-size: 11px;
    color: #7c7c7c;
  }

  .popup-conversation-bottom {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 8px;
  }

  .popup-conversation-preview {
    flex: 1;
    margin: 0;
    font-size: 12px;
    color: #7c7c7c;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .popup-unread-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 18px;
    height: 18px;
    padding: 0 5px;
    background: var(--error--);
    color: white;
    font-size: 10px;
    font-weight: 700;
    border-radius: 9px;
  }

  /* Chat Detail */
  .popup-chat-detail {
    flex: 1;
    display: flex;
    flex-direction: column;
    background: white;
  }

  .popup-chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    border-bottom: 1px solid #edeff1;
  }

  .popup-chat-user-info {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .popup-chat-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
  }

  .popup-chat-user-details h3 {
    margin: 0 0 2px 0;
    font-size: 14px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .popup-chat-actions {
    position: relative;
  }

  .popup-chat-menu-wrapper {
    position: relative;
  }

  .popup-action-icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    color: #7c7c7c;
    transition: background 0.2s;
  }

  .popup-action-icon-btn:hover {
    background: #f6f7f8;
  }

  .popup-chat-dropdown {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    background: white;
    border: 1px solid #edeff1;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    min-width: 180px;
    z-index: 100;
    overflow: hidden;
  }

  .popup-dropdown-item {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding: 10px 14px;
    border: none;
    background: transparent;
    cursor: pointer;
    text-align: left;
    font-size: 13px;
    color: #1c1c1c;
    transition: background 0.2s;
  }

  .popup-dropdown-item:hover {
    background: #f6f7f8;
  }

  .popup-dropdown-item.muted {
    color: #ff4444;
  }

  .popup-dropdown-item.muted:hover {
    background: #ffe0e0;
  }

  .popup-dropdown-item img {
    opacity: 0.7;
  }

  .popup-dropdown-item span {
    font-weight: 500;
  }

  /* Messages Area */
  .popup-reports-area {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    background: #f6f7f8;
  }

  .popup-reports-wrapper {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .popup-message-row {
    display: flex;
    align-items: flex-end;
    gap: 6px;
  }

  .popup-message-row.sent {
    flex-direction: row-reverse;
  }

  .popup-message-avatar {
    width: 28px;
    height: 28px;
    border-radius: 50%;
    object-fit: cover;
    flex-shrink: 0;
  }

  .popup-message-bubble {
    max-width: 70%;
    padding: 8px 12px;
    background: white;
    border-radius: 16px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }

  .popup-message-bubble.sent {
    background: var(--blue--);
    color: white;
  }

  .popup-message-bubble p {
    margin: 0 0 3px 0;
    font-size: 13px;
    line-height: 1.4;
    word-wrap: break-word;
  }

  .popup-message-time {
    font-size: 10px;
    opacity: 0.7;
  }

  /* Message Input */
  .popup-message-input {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 16px;
    border-top: 1px solid #edeff1;
    background: white;
  }

  .popup-emoji-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    color: #7c7c7c;
    transition: background 0.2s;
  }

  .popup-emoji-btn:hover {
    background: #f6f7f8;
  }

  .popup-message-input input {
    flex: 1;
    padding: 8px 14px;
    border: 1px solid #edeff1;
    border-radius: 18px;
    background: #f6f7f8;
    font-size: 13px;
    color: #1c1c1c;
  }

  .popup-message-input input:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  .popup-send-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    transition: background 0.2s;
  }

  .popup-send-btn:hover:not(:disabled) {
    background: #f6f7f8;
  }

  .popup-send-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  /* Loading & Empty States */
  .popup-loading,
  .popup-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 12px;
    padding: 20px;
    color: #7c7c7c;
    font-size: 13px;
  }

  .popup-empty p {
    margin: 0;
  }

  .debug-create-channel-btn {
    padding: 8px 16px;
    background: var(--blue--);
    color: white;
    border: none;
    border-radius: 6px;
    font-size: 12px;
    cursor: pointer;
    transition: background 0.2s;
  }

  .debug-create-channel-btn:hover {
    background: #0056b3;
  }
</style>
