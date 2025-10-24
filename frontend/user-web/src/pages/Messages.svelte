<script lang="ts">
  import { push } from "svelte-spa-router";
  import { mockConversations } from "../mocks/conversations.mock";
  import type { Conversation, Message } from "../mocks/conversations.mock";

  let conversations = $state(mockConversations);
  let selectedConversation = $state<Conversation | null>(conversations[0]);
  let messageInput = $state("");
  let isTyping = $state(false);
  let showChatMenu = $state(false);

  function handleSelectConversation(conversation: Conversation) {
    selectedConversation = conversation;
    // Mark as read
    if (selectedConversation) {
      selectedConversation.unreadCount = 0;
    }
  }

  function handleSendMessage() {
    if (!messageInput.trim() || !selectedConversation) return;

    const newMessage: Message = {
      id: `m${Date.now()}`,
      senderId: "me",
      senderName: "You",
      senderAvatar: "/avatar.jpg",
      content: messageInput,
      timestamp: new Date().toLocaleTimeString("en-US", {
        hour: "2-digit",
        minute: "2-digit",
      }),
      isRead: false,
      isSent: true,
    };

    selectedConversation.messages.push(newMessage);
    selectedConversation.lastMessage = `You: ${messageInput}`;
    selectedConversation.lastMessageTime = "Just now";

    messageInput = "";

    // Scroll to bottom
    setTimeout(() => {
      const messagesArea = document.querySelector(".messages-area");
      if (messagesArea) {
        messagesArea.scrollTop = messagesArea.scrollHeight;
      }
    }, 0);
  }

  function handleBack() {
    push("/");
  }

  function toggleMute() {
    if (selectedConversation) {
      selectedConversation.isMuted = !selectedConversation.isMuted;
    }
  }

  function toggleChatMenu() {
    showChatMenu = !showChatMenu;
  }

  function formatTime(time: string): string {
    return time;
  }
</script>

<div class="messages-page">
  <!-- Header -->
  <div class="messages-header">
    <button class="back-btn" onclick={handleBack} title="Back to home">
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

  <div class="messages-container">
    <!-- Left Side - Conversations List -->
    <div class="conversations-sidebar">
      <div class="conversations-header">
        <div class="user-info">
          <img src="/avatar.jpg" alt="User" class="user-avatar" />
          <h2>Chats</h2>
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
        <input type="text" placeholder="Search messages..." />
      </div>

      <div class="conversations-list">
        {#each conversations as conversation}
          <button
            class="conversation-item"
            class:active={selectedConversation?.id === conversation.id}
            class:unread={conversation.unreadCount > 0}
            onclick={() => handleSelectConversation(conversation)}
          >
            <div class="conversation-avatar-wrapper">
              <img
                src={conversation.userAvatar}
                alt={conversation.userName}
                class="conversation-avatar"
              />
              {#if conversation.isOnline}
                <span class="online-indicator"></span>
              {/if}
            </div>
            <div class="conversation-info">
              <div class="conversation-top">
                <span class="conversation-name">{conversation.userName}</span>
                <span class="conversation-time"
                  >{conversation.lastMessageTime}</span
                >
              </div>
              <div class="conversation-bottom">
                <p class="conversation-preview">{conversation.lastMessage}</p>
                {#if conversation.unreadCount > 0}
                  <span class="unread-badge">{conversation.unreadCount}</span>
                {/if}
              </div>
            </div>
          </button>
        {/each}
      </div>
    </div>

    <!-- Right Side - Chat Detail -->
    <div class="chat-detail">
      {#if selectedConversation}
        <!-- Chat Header -->
        <div class="chat-header">
          <div class="chat-user-info">
            <img
              src={selectedConversation.userAvatar}
              alt={selectedConversation.userName}
              class="chat-avatar"
            />
            <div class="chat-user-details">
              <h3>{selectedConversation.userName}</h3>
              {#if selectedConversation.isOnline}
                <span class="status-online">Active now</span>
              {:else}
                <span class="status-offline">Offline</span>
              {/if}
            </div>
          </div>
          <div class="chat-actions">
            <div class="chat-menu-wrapper">
              <button
                class="action-icon-btn"
                onclick={toggleChatMenu}
                title="More options"
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
                    class:muted={selectedConversation.isMuted}
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
                      >{selectedConversation.isMuted ? "Unmute" : "Mute"}</span
                    >
                  </button>
                  <button
                    class="dropdown-item"
                    onclick={() => {
                      showChatMenu = false;
                    }}
                  >
                    <img
                      src="/notification_icon.svg"
                      alt="Notification"
                      width="20"
                      height="20"
                    />
                    <span>Notifications</span>
                  </button>
                </div>
              {/if}
            </div>
          </div>
        </div>

        <!-- Messages Area -->
        <div class="messages-area">
          <div class="messages-wrapper">
            {#each selectedConversation.messages as message}
              <div class="message-row" class:sent={message.isSent}>
                {#if !message.isSent}
                  <img
                    src={message.senderAvatar}
                    alt={message.senderName}
                    class="message-avatar"
                  />
                {/if}
                <div class="message-bubble" class:sent={message.isSent}>
                  <p>{message.content}</p>
                  <span class="message-time">{message.timestamp}</span>
                </div>
              </div>
            {/each}
            {#if isTyping}
              <div class="message-row">
                <img
                  src={selectedConversation.userAvatar}
                  alt={selectedConversation.userName}
                  class="message-avatar"
                />
                <div class="typing-indicator">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
              </div>
            {/if}
          </div>
        </div>

        <!-- Message Input -->
        <div class="message-input-container">
          <button class="emoji-btn" title="Emoji">
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
            disabled={!messageInput.trim()}
            title="Send"
          >
            <img src="/send_icon.svg" alt="Send" width="24" height="24" />
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
  .messages-page {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background: #f6f7f8;
  }

  /* Header */
  .messages-header {
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

  .messages-header h1 {
    margin: 0;
    font-size: 24px;
    font-weight: 700;
    color: #1c1c1c;
  }

  /* Container */
  .messages-container {
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

  .online-indicator {
    position: absolute;
    bottom: 2px;
    right: 2px;
    width: 14px;
    height: 14px;
    background: #31a24c;
    border: 2px solid white;
    border-radius: 50%;
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
  .messages-area {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    background: #f6f7f8;
  }

  .messages-wrapper {
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

  /* Typing Indicator */
  .typing-indicator {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 12px 16px;
    background: white;
    border-radius: 18px;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
  }

  .typing-indicator span {
    width: 8px;
    height: 8px;
    background: #7c7c7c;
    border-radius: 50%;
    animation: typing 1.4s infinite;
  }

  .typing-indicator span:nth-child(2) {
    animation-delay: 0.2s;
  }

  .typing-indicator span:nth-child(3) {
    animation-delay: 0.4s;
  }

  @keyframes typing {
    0%,
    60%,
    100% {
      transform: translateY(0);
      opacity: 0.7;
    }
    30% {
      transform: translateY(-10px);
      opacity: 1;
    }
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
</style>
