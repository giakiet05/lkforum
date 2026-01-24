<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import { push } from "svelte-spa-router";
  import { authStore } from "../stores/auth-store";
  import { toastStore } from "../stores/toast-store";
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
    clearChatStore,
  } from "../stores/chat-store";
  import {
    getChannelsByUser,
    updateChannel,
  } from "../services/channel-service";
  import { getMessages } from "../services/message-service";
  import { websocketService } from "../services/websocket-service";
  import type { ChannelResponse } from "../dtos/channel-dto";
  import type { MessageResponse } from "../dtos/message-dto";

  let messageInput = $state("");
  let showChatMenu = $state(false);
  let searchQuery = $state("");

  // Get current user
  const currentUser = $derived($authStore.user);
  const channels = $derived($chatStore.channels);
  const activeChannelData = $derived($activeChannel);
  const reports = $derived($activeChannelMessages);
  const isLoadingChannels = $derived($chatStore.isLoadingChannels);
  const isLoadingMessages = $derived($chatStore.isLoadingMessages);

  // Filter channels based on search query
  const filteredChannels = $derived(() => {
    if (!searchQuery.trim()) return channels;
    const query = searchQuery.toLowerCase();
    return channels.filter((channel) => {
      const member = channel.members.find((m) => m.user_id !== currentUser?.id);
      if (!member) return false;
      return (
        member.username.toLowerCase().includes(query) ||
        member.full_name?.toLowerCase().includes(query)
      );
    });
  });

  // Get other member info from active channel
  const otherMember = $derived(() => {
    if (!activeChannelData || !currentUser) return null;
    return (
      activeChannelData.members.find((m) => m.user_id !== currentUser.id) ||
      null
    );
  });

  // Check if other member is online
  const isOtherMemberOnline = $derived(
    otherMember() ? $chatStore.onlineUsers.has(otherMember().user_id) : false,
  );

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

  // Format relative time for last message
  function formatRelativeTime(isoString: string): string {
    const now = new Date();
    const date = new Date(isoString);
    const diffMs = now.getTime() - date.getTime();
    const diffMins = Math.floor(diffMs / 60000);

    if (diffMins < 1) return "Vừa xong";
    if (diffMins < 60) return `${diffMins} phút`;

    const diffHours = Math.floor(diffMins / 60);
    if (diffHours < 24) return `${diffHours} giờ`;

    const diffDays = Math.floor(diffHours / 24);
    return `${diffDays} ngày`;
  }

  // Get last message for a channel
  function getLastMessage(channelId: string): string {
    const channelMessages = $chatStore.messagesByChannel.get(channelId) || [];
    if (channelMessages.length === 0) return "Chưa có tin nhắn nào";

    const lastMsg = channelMessages[channelMessages.length - 1];
    const isCurrentUser = lastMsg.sender_id === currentUser?.id;
    const prefix = isCurrentUser ? "Bạn: " : `${lastMsg.sender_username}: `;

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
      (m) => !m.is_read && m.sender_id !== currentUser?.id,
    ).length;
  }

  // Load channels on mount
  async function loadChannels() {
    if (!currentUser) return;

    try {
      console.log("📋 Loading channels for user:", currentUser.id);
      setLoadingChannels(true);
      const response = await getChannelsByUser(currentUser.id);
      console.log("✅ Channels loaded:", response);
      setChannels(response.channels);

      // Auto-select first channel
      if (response.channels.length > 0 && !$chatStore.activeChannelId) {
        console.log("🎯 Auto-selecting first channel:", response.channels[0]);
        await handleSelectChannel(response.channels[0]);
      }
    } catch (error) {
      console.error("❌ Failed to load channels:", error);
    }
  }

  // Load reports for a channel
  async function loadMessagesForChannel(channelId: string) {
    try {
      console.log("💬 Loading messages for channel:", channelId);
      setLoadingMessages(true);
      const channelMessages = await getMessages({ channel_id: channelId });
      console.log("✅ Messages loaded:", channelMessages);
      setMessages(channelId, channelMessages);

      // Mark as read
      markChannelAsRead(channelId);

      // Scroll to bottom
      setTimeout(scrollToBottom, 100);
    } catch (error) {
      console.error("❌ Failed to load messages:", error);
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
        "text",
      );
      messageInput = "";

      // Scroll to bottom after sending
      setTimeout(scrollToBottom, 100);
    } catch (error) {
      console.error("Failed to send message:", error);
      toastStore.error("Không thể gửi tin nhắn. Vui lòng thử lại.");
    }
  }

  // Handle incoming WebSocket reports
  function handleIncomingMessage(message: MessageResponse) {
    addMessage(message.channel_id, message);

    // Scroll to bottom if message is in active channel
    if (message.channel_id === $chatStore.activeChannelId) {
      setTimeout(scrollToBottom, 100);
    }
  }

  // Scroll to bottom of reports
  function scrollToBottom() {
    const messagesArea = document.querySelector(".reports-area");
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

      // Reload channels to reflect changes
      await loadChannels();
    } catch (error) {
      console.error("Failed to toggle mute:", error);
    }
  }

  function toggleChatMenu() {
    showChatMenu = !showChatMenu;
  }

  function handleBack() {
    push("/");
  }

  // Initialize on mount
  onMount(async () => {
    console.log("🚀 Messages page mounted");
    console.log("👤 Current user:", currentUser);

    if (!currentUser) {
      console.log("❌ No current user, redirecting to login");
      push("/login");
      return;
    }

    // Load channels
    await loadChannels();

    // Check for channel query param (from notification click)
    const urlParams = new URLSearchParams(
      window.location.hash.split("?")[1] || "",
    );
    const channelIdFromUrl = urlParams.get("channel");
    if (channelIdFromUrl) {
      console.log("📩 Opening channel from notification:", channelIdFromUrl);
      // Find and select the channel after channels are loaded
      setTimeout(() => {
        const targetChannel = channels.find((c) => c.id === channelIdFromUrl);
        if (targetChannel) {
          handleSelectChannel(targetChannel);
        }
      }, 100);
    }

    // Register message handler (WebSocket already connected in App.svelte)
    if (websocketService.isConnected()) {
      websocketService.onMessage(handleIncomingMessage);
      console.log("✅ WebSocket message handler registered");
    } else {
      console.warn("⚠️ WebSocket not connected yet");
    }
  });

  // Cleanup on destroy
  onDestroy(() => {
    websocketService.offMessage(handleIncomingMessage);
    // Don't disconnect WebSocket here - might be used in other pages
  });
</script>

<div class="reports-page">
  <!-- Header -->
  <div class="reports-header">
    <button class="back-btn" onclick={handleBack} title="Về trang chủ">
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
        <path
          d="M15 18L9 12L15 6"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
        />
      </svg>
    </button>
    <h1>Tin nhắn</h1>
  </div>

  <div class="reports-container">
    <!-- Left Side - Conversations List -->
    <div class="conversations-sidebar">
      <div class="conversations-header">
        <div class="user-info">
          <img
            src={currentUser?.profile?.avatar?.url || "/user.jpg"}
            alt="User"
            class="user-avatar"
          />
          <h2>Cuộc trò chuyện</h2>
        </div>
      </div>

      <div class="search-box">
        <img
          src="/search_icon.svg"
          alt="Search"
          class="search-icon-img"
          width="20"
          height="20"
        />
        <input
          type="text"
          placeholder="Tìm kiếm tin nhắn..."
          bind:value={searchQuery}
        />
      </div>

      <div class="conversations-list">
        {#if isLoadingChannels}
          <div class="loading-state">Đang tải...</div>
        {:else if filteredChannels().length === 0}
          <div class="empty-state">
            {searchQuery ? "Không tìm thấy kết quả" : "Chưa có cuộc hội thoại"}
          </div>
        {:else}
          {#each filteredChannels() as channel (channel.id)}
            {@const member = channel.members.find(
              (m) => m.user_id !== currentUser?.id,
            )}
            {@const unreadCount = getUnreadCount(channel.id)}
            {#if member}
              <button
                class="conversation-item"
                class:active={$chatStore.activeChannelId === channel.id}
                class:unread={unreadCount > 0}
                onclick={() => handleSelectChannel(channel)}
              >
                <div class="conversation-avatar-wrapper">
                  <img
                    src={member.avatar || "/user.jpg"}
                    alt={member.username}
                    class="conversation-avatar"
                  />
                  <!-- TODO: Online status needs backend support -->
                </div>
                <div class="conversation-info">
                  <div class="conversation-top">
                    <span class="conversation-name">{member.username}</span>
                    <span class="conversation-time"
                      >{getLastMessageTime(channel.id)}</span
                    >
                  </div>
                  <div class="conversation-bottom">
                    <p class="conversation-preview">
                      {getLastMessage(channel.id)}
                    </p>
                    {#if unreadCount > 0}
                      <span class="unread-badge">{unreadCount}</span>
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
    <div class="chat-detail">
      {#if activeChannelData && otherMember()}
        <!-- Chat Header -->
        <div class="chat-header">
          <div class="chat-user-info">
            <div class="avatar-wrapper">
              <img
                src={otherMember()?.avatar || "/user.jpg"}
                alt={otherMember()?.username || "User"}
                class="chat-avatar"
              />
              {#if isOtherMemberOnline}
                <span class="online-indicator"></span>
              {/if}
            </div>
            <div class="chat-user-details">
              <h3>{otherMember()?.username || "Ẩn danh"}</h3>
              {#if isOtherMemberOnline}
                <span class="status-online">Đang hoạt động</span>
              {:else}
                <span class="status-offline">Ngoại tuyến</span>
              {/if}
            </div>
          </div>
          <div class="chat-actions">
            <div class="chat-menu-wrapper">
              <button
                class="action-icon-btn"
                onclick={toggleChatMenu}
                title="Thêm tùy chọn"
              >
                <svg
                  width="20"
                  height="20"
                  viewBox="0 0 20 20"
                  fill="currentColor"
                >
                  <path
                    d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z"
                  />
                </svg>
              </button>
              {#if showChatMenu}
                <div class="chat-dropdown">
                  <button
                    class="dropdown-item"
                    class:muted={!(channelSettings()?.notification ?? true)}
                    onclick={() => {
                      toggleMute();
                      showChatMenu = false;
                    }}
                  >
                    <img
                      src="/muted_icon.svg"
                      alt="Mute"
                      width="20"
                      height="20"
                    />
                    <span
                      >{(channelSettings()?.notification ?? true)
                        ? "Tắt thông báo"
                        : "Bật thông báo"}</span
                    >
                  </button>
                </div>
              {/if}
            </div>
          </div>
        </div>

        <!-- Messages Area -->
        <div class="reports-area">
          {#if isLoadingMessages}
            <div class="loading-reports">Đang tải tin nhắn...</div>
          {:else}
            <div class="reports-wrapper">
              {#each reports as message (message.id)}
                {@const isSent = message.sender_id === currentUser?.id}
                <div class="message-row" class:sent={isSent}>
                  {#if !isSent}
                    <img
                      src={otherMember()?.avatar || "/user.jpg"}
                      alt={message.sender_username}
                      class="message-avatar"
                    />
                  {/if}
                  <div class="message-bubble" class:sent={isSent}>
                    <p>{message.content}</p>
                    <span class="message-time"
                      >{formatTime(message.created_at)}</span
                    >
                  </div>
                </div>
              {/each}
              <!-- TODO: Typing indicator needs backend support -->
            </div>
          {/if}
        </div>

        <!-- Message Input -->
        <div class="message-input-container">
          <button class="emoji-btn" title="Biểu tượng cảm xúc">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none">
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
            class="send-btn"
            onclick={handleSendMessage}
            disabled={!messageInput.trim() || !websocketService.isConnected()}
            title="Gửi"
          >
            <img src="/send_icon.svg" alt="Gửi" width="24" height="24" />
          </button>
        </div>
      {:else}
        <div class="no-conversation-selected">
          <svg width="96" height="96" viewBox="0 0 96 96" fill="none">
            <path
              d="M48 88C70.0914 88 88 70.0914 88 48C88 25.9086 70.0914 8 48 8C25.9086 8 8 25.9086 8 48C8 70.0914 25.9086 88 48 88Z"
              stroke="#E0E0E0"
              stroke-width="4"
            />
            <path
              d="M32 42L48 58L64 42"
              stroke="#E0E0E0"
              stroke-width="4"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
          <h3>Select a conversation</h3>
          <p>Choose a conversation from the list to start messaging</p>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .reports-page {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background: #f6f7f8;
  }

  /* Header */
  .reports-header {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px 24px;
    background: white;
    border-bottom: 1px solid #edeff1;
  }

  .back-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    color: #1c1c1c;
    transition: background 0.2s;
  }

  .back-btn:hover {
    background: #f6f7f8;
  }

  .reports-header h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
  }

  /* Container */
  .reports-container {
    display: flex;
    flex: 1;
    overflow: hidden;
  }

  /* Conversations Sidebar */
  .conversations-sidebar {
    width: 360px;
    background: white;
    border-right: 1px solid #edeff1;
    display: flex;
    flex-direction: column;
  }

  .conversations-header {
    padding: 20px;
    border-bottom: 1px solid #edeff1;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .user-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    object-fit: cover;
  }

  .user-info h2 {
    margin: 0;
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
  }

  /* Search Box */
  .search-box {
    position: relative;
    padding: 0 20px 20px;
    margin-top: 8px;
  }

  .search-icon-img {
    position: absolute;
    left: 36px;
    top: 10px;
    opacity: 0.6;
    pointer-events: none;
  }

  .search-box input {
    width: 100%;
    padding: 10px 16px 10px 44px;
    border: 1px solid #edeff1;
    border-radius: 20px;
    background: #f6f7f8;
    font-size: 14px;
    color: #1c1c1c;
  }

  .search-box input:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  /* Conversations List */
  .conversations-list {
    flex: 1;
    overflow-y: auto;
  }

  .conversation-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 20px;
    width: 100%;
    border: none;
    background: transparent;
    cursor: pointer;
    text-align: left;
    transition: background 0.2s;
    border-left: 3px solid transparent;
  }

  .conversation-item:hover {
    background: #f6f7f8;
  }

  .conversation-item.active {
    background: #e8f4fd;
    border-left-color: var(--blue--);
  }

  .conversation-avatar-wrapper {
    position: relative;
    flex-shrink: 0;
  }

  .conversation-avatar {
    width: 56px;
    height: 56px;
    border-radius: 50%;
    object-fit: cover;
  }

  .conversation-info {
    flex: 1;
    min-width: 0;
  }

  .conversation-top {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 4px;
  }

  .conversation-name {
    font-size: 15px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .conversation-item.unread .conversation-name {
    font-weight: 700;
  }

  .conversation-item.unread .conversation-preview {
    color: #1c1c1c;
    font-weight: 600;
  }

  .conversation-time {
    font-size: 12px;
    color: #7c7c7c;
  }

  .conversation-bottom {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 8px;
  }

  .conversation-preview {
    flex: 1;
    margin: 0;
    font-size: 13px;
    color: #7c7c7c;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .unread-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 20px;
    height: 20px;
    padding: 0 6px;
    background: var(--error--);
    color: white;
    font-size: 11px;
    font-weight: 700;
    border-radius: 10px;
  }

  /* Chat Detail */
  .chat-detail {
    flex: 1;
    display: flex;
    flex-direction: column;
    background: white;
  }

  .chat-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 24px;
    border-bottom: 1px solid #edeff1;
  }

  .chat-user-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .chat-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    object-fit: cover;
  }

  .chat-user-details h3 {
    margin: 0 0 4px 0;
    font-size: 16px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .status-online,
  .status-offline {
    font-size: 12px;
    color: #7c7c7c;
  }

  .status-online {
    color: #31a24c;
  }

  .chat-actions {
    position: relative;
  }

  .chat-menu-wrapper {
    position: relative;
  }

  .action-icon-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    color: #7c7c7c;
    transition: background 0.2s;
  }

  .action-icon-btn:hover {
    background: #f6f7f8;
  }

  .chat-dropdown {
    position: absolute;
    top: calc(100% + 8px);
    right: 0;
    background: white;
    border: 1px solid #edeff1;
    border-radius: 8px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    min-width: 200px;
    z-index: 100;
    overflow: hidden;
  }

  .dropdown-item {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    padding: 12px 16px;
    border: none;
    background: transparent;
    cursor: pointer;
    text-align: left;
    font-size: 14px;
    color: #1c1c1c;
    transition: background 0.2s;
  }

  .dropdown-item:hover {
    background: #f6f7f8;
  }

  .dropdown-item.muted {
    color: #ff4444;
  }

  .dropdown-item.muted:hover {
    background: #ffe0e0;
  }

  .dropdown-item img {
    opacity: 0.7;
  }

  .dropdown-item span {
    font-weight: 500;
  }

  /* Messages Area */
  .reports-area {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    background: #f6f7f8;
  }

  .reports-wrapper {
    display: flex;
    flex-direction: column;
    gap: 12px;
    max-width: 800px;
    margin: 0 auto;
  }

  .message-row {
    display: flex;
    align-items: flex-end;
    gap: 8px;
  }

  .message-row.sent {
    flex-direction: row-reverse;
  }

  .message-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    object-fit: cover;
    flex-shrink: 0;
  }

  .message-bubble {
    max-width: 70%;
    padding: 10px 16px;
    background: white;
    border-radius: 18px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }

  .message-bubble.sent {
    background: var(--blue--);
    color: white;
  }

  .message-bubble p {
    margin: 0 0 4px 0;
    font-size: 14px;
    line-height: 1.4;
    word-wrap: break-word;
  }

  .message-time {
    font-size: 11px;
    opacity: 0.7;
  }

  /* Message Input */
  .message-input-container {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 16px 24px;
    border-top: 1px solid #edeff1;
    background: white;
  }

  .emoji-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    color: #7c7c7c;
    transition: background 0.2s;
  }

  .emoji-btn:hover {
    background: #f6f7f8;
  }

  .message-input-container input {
    flex: 1;
    padding: 10px 16px;
    border: 1px solid #edeff1;
    border-radius: 20px;
    background: #f6f7f8;
    font-size: 14px;
    color: #1c1c1c;
  }

  .message-input-container input:focus {
    outline: none;
    border-color: var(--blue--);
    background: white;
  }

  .send-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border: none;
    background: transparent;
    border-radius: 50%;
    cursor: pointer;
    transition: background 0.2s;
  }

  .send-btn:hover:not(:disabled) {
    background: #f6f7f8;
  }

  .send-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
  }

  /* No Conversation Selected */
  .no-conversation-selected {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    text-align: center;
    padding: 40px;
  }

  .no-conversation-selected h3 {
    margin: 24px 0 8px 0;
    font-size: 20px;
    font-weight: 600;
    color: #1c1c1c;
  }

  .no-conversation-selected p {
    margin: 0;
    font-size: 14px;
    color: #7c7c7c;
  }

  /* Loading States */
  .loading-state,
  .empty-state,
  .loading-reports {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 40px;
    color: #7c7c7c;
    font-size: 14px;
  }

  .loading-reports {
    height: 100%;
  }
</style>
